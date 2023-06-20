package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/cli"
	"github.com/vk-cs/terraform-provider-vkcs/helpers/changelog"
	"github.com/vk-cs/terraform-provider-vkcs/helpers/tfsdk2fw/naming"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/provider"
	"golang.org/x/exp/slices"
)

var (
	dataSourceType = flag.String("data-source", "", "Data Source type")
	resourceType   = flag.String("resource", "", "Resource type")
	changelogPath  = flag.String("changelog", "CHANGELOG.md", "Path to changelog file")

	packageServiceMap = map[string]string{
		"regions":          "Regions",
		"compute":          "Virtual Machines",
		"blockstorage":     "Disks",
		"image":            "Images",
		"sharedfilesystem": "File Share",
		"network":          "Network",
		"firewall":         "Firewall",
		"lb":               "Load Balancers",
		"vpnaas":           "VPN",
		"publicdns":        "DNS",
		"kubernetes":       "Kubernetes",
		"db":               "Databases",
		"keymanager":       "Key Manager",
	}
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "\ttfsdk2fw [-resource <resource-type>|-data-source <data-source-type>] <package-name> <name> <generated-file>\n\n")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 3 || (*dataSourceType == "" && *resourceType == "") {
		flag.Usage()
		os.Exit(2)
	}

	packageName := args[0]
	name := args[1]
	outputFilename := args[2]

	g := &generator{
		ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	migrator := &migrator{
		Generator:     g,
		Name:          name,
		PackageName:   packageName,
		ChangelogPath: *changelogPath,
	}

	p := provider.SDKProvider()

	if v := *dataSourceType; v != "" {
		resource, ok := p.DataSourcesMap[v]

		if !ok {
			g.fatalf("data source type %s not found", v)
		}

		migrator.IsDataSource = true
		migrator.Resource = resource
		migrator.Template = datasourceImpl
		migrator.TestTemplate = datasourceTestImpl
		migrator.TFTypeName = v
	} else if v := *resourceType; v != "" {
		resource, ok := p.ResourcesMap[v]

		if !ok {
			g.fatalf("resource type %s not found", v)
		}

		migrator.Resource = resource
		migrator.Template = resourceImpl
		migrator.TestTemplate = resourceTestImpl
		migrator.TFTypeName = v
	}

	if err := migrator.migrate(outputFilename); err != nil {
		g.fatalf("error migrating Terraform %s schema: %s", *resourceType, err)
	}
}

type generator struct {
	ui cli.Ui
}

type destination interface {
	write() error
	writeBytes(body []byte) error
	writeTemplate(templateName, templateBody string, templateData any) error
}

func (g *generator) newGoFileDestination(filename string) destination {
	return &fileDestination{
		filename:  filename,
		formatter: format.Source,
	}
}

type fileDestination struct {
	append    bool
	filename  string
	formatter func([]byte) ([]byte, error)
	buffer    strings.Builder
}

func (d *fileDestination) write() error {
	var flags int
	if d.append {
		flags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	} else {
		flags = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	}

	f, err := os.OpenFile(d.filename, flags, 0644) //nolint:gomnd
	if err != nil {
		return fmt.Errorf("opening file (%s): %w", d.filename, err)
	}
	defer f.Close()

	_, err = f.WriteString(d.buffer.String())
	if err != nil {
		return fmt.Errorf("writing to file (%s): %w", d.filename, err)
	}

	return nil
}

func (d *fileDestination) writeBytes(body []byte) error {
	_, err := d.buffer.Write(body)
	return err
}

func (d *fileDestination) writeTemplate(templateName, templateBody string, templateData any) error {
	unformattedBody, err := parseTemplate(templateName, templateBody, templateData)
	if err != nil {
		return err
	}

	body, err := d.formatter(unformattedBody)
	if err != nil {
		return fmt.Errorf("formatting parsed template:\n%s\n%w", unformattedBody, err)
	}

	_, err = d.buffer.Write(body)
	return err
}

func parseTemplate(templateName, templateBody string, templateData any) ([]byte, error) {
	funcs := template.FuncMap{
		"lower_first": lowerFirst,
		"upper_first": upperFirst,
	}
	tmpl, err := template.New(templateName).Funcs(funcs).Parse(templateBody)
	if err != nil {
		return nil, fmt.Errorf("parsing function template: %w", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)
	if err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return buffer.Bytes(), nil
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func (g *generator) infof(format string, a ...interface{}) {
	g.ui.Info(fmt.Sprintf(format, a...))
}

func (g *generator) warnf(format string, a ...interface{}) {
	g.ui.Warn(fmt.Sprintf(format, a...))
}

func (g *generator) errorf(format string, a ...interface{}) {
	g.ui.Error(fmt.Sprintf(format, a...))
}

func (g *generator) fatalf(format string, a ...interface{}) {
	g.errorf(format, a...)
	os.Exit(1)
}

type migrator struct {
	Generator     *generator
	IsDataSource  bool
	Name          string
	PackageName   string
	ChangelogPath string
	Resource      *schema.Resource
	Template      string
	TestTemplate  string
	TFTypeName    string
}

// migrate generates an identical schema into the specified output file.
func (m *migrator) migrate(outputFilename string) error {
	m.infof("generating into %[1]q", outputFilename)

	// Create target directory.
	dirname := path.Dir(outputFilename)
	err := os.MkdirAll(dirname, 0755)
	if err != nil {
		return fmt.Errorf("creating target directory %s: %w", dirname, err)
	}

	templateData, err := m.generateTemplateData()
	if err != nil {
		return err
	}

	d := m.Generator.newGoFileDestination(outputFilename)
	if err := d.writeTemplate("schema", m.Template, templateData); err != nil {
		return err
	}
	if err := d.write(); err != nil {
		return err
	}

	testTemplateData, err := m.generateTestTemplateData()
	if err != nil {
		return err
	}

	d = m.Generator.newGoFileDestination(strings.Replace(outputFilename, ".go", "_test.go", 1))
	if err := d.writeTemplate("migrate_test", m.TestTemplate, testTemplateData); err != nil {
		return err
	}
	return d.write()
}

func (m *migrator) generateTemplateData() (*templateData, error) {
	sbSchema := strings.Builder{}
	sbModel := strings.Builder{}
	emitter := &emitter{
		Generator:    m.Generator,
		IsDataSource: m.IsDataSource,
		SchemaWriter: &sbSchema,
		ModelWriter:  &sbModel,
	}

	if err := emitter.emitSchemaForResource(m.Resource); err != nil {
		return nil, fmt.Errorf("emitting schema code: %w", err)
	}

	serviceName := "UNKNOWN"
	if v, ok := packageServiceMap[m.PackageName]; ok {
		serviceName = v
	}

	sort.SliceStable(emitter.NestedModels, func(i, j int) bool {
		return emitter.NestedModels[i].Name < emitter.NestedModels[j].Name
	})

	templateData := &templateData{
		DefaultCreateTimeout:         emitter.DefaultCreateTimeout,
		DefaultReadTimeout:           emitter.DefaultReadTimeout,
		DefaultUpdateTimeout:         emitter.DefaultUpdateTimeout,
		DefaultDeleteTimeout:         emitter.DefaultDeleteTimeout,
		EmitResourceImportState:      m.Resource.Importer != nil,
		EmitResourceUpdateSkeleton:   m.Resource.UpdateContext != nil || m.Resource.UpdateWithoutTimeout != nil,
		HasTimeouts:                  emitter.HasTimeouts,
		ImportFrameworkAttr:          emitter.ImportFrameworkAttr,
		ImportProviderFrameworkTypes: emitter.ImportProviderFrameworkTypes,
		Name:                         m.Name,
		PackageName:                  m.PackageName,
		ServiceName:                  serviceName,
		Schema:                       sbSchema.String(),
		Model:                        sbModel.String(),
		NestedModels:                 emitter.NestedModels,
		TFTypeName:                   m.TFTypeName,
	}

	for _, v := range emitter.FrameworkPlanModifierPackages {
		if !slices.Contains(templateData.FrameworkPlanModifierPackages, v) {
			templateData.FrameworkPlanModifierPackages = append(templateData.FrameworkPlanModifierPackages, v)
		}
	}
	for _, v := range emitter.FrameworkValidatorsPackages {
		if !slices.Contains(templateData.FrameworkValidatorsPackages, v) {
			templateData.FrameworkValidatorsPackages = append(templateData.FrameworkValidatorsPackages, v)
		}
	}
	for _, v := range emitter.ProviderPlanModifierPackages {
		if !slices.Contains(templateData.ProviderPlanModifierPackages, v) {
			templateData.ProviderPlanModifierPackages = append(templateData.ProviderPlanModifierPackages, v)
		}
	}

	return templateData, nil
}

func (m *migrator) generateTestTemplateData() (*testTemplateData, error) {
	cl, err := changelog.NewChangelogFromFile(*changelogPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing changelog: %w", err)
	}

	if len(cl.Versions) < 2 {
		return nil, fmt.Errorf("error parsing changelog: unable to find last released version")
	}

	return &testTemplateData{
		Name:            m.Name,
		PackageName:     m.PackageName,
		ReleasedVersion: strings.TrimPrefix(cl.Versions[1].Version, "v"),
	}, nil
}

func (m *migrator) infof(format string, a ...interface{}) {
	m.Generator.infof(format, a...)
}

type emitter struct {
	DefaultCreateTimeout          int64
	DefaultReadTimeout            int64
	DefaultUpdateTimeout          int64
	DefaultDeleteTimeout          int64
	Generator                     *generator
	FrameworkPlanModifierPackages []string // Package names for any terraform-plugin-framework plan modifiers. May contain duplicates.
	FrameworkValidatorsPackages   []string // Package names for any terraform-plugin-framework-validators validators. May contain duplicates.
	HasTimeouts                   bool
	ImportFrameworkAttr           bool
	ImportProviderFrameworkTypes  bool
	IsDataSource                  bool
	ProviderPlanModifierPackages  []string // Package names for any provider plan modifiers. May contain duplicates.
	SchemaWriter                  io.Writer
	ModelWriter                   io.Writer
	NestedModels                  []nestedModelData
}

// emitSchemaForResource generates the Plugin Framework code for a Plugin SDK Resource and emits the generated code to the emitter's Writer.
func (e *emitter) emitSchemaForResource(resource *schema.Resource) error {
	if _, ok := resource.Schema["id"]; ok {
		e.warnf("Explicit `id` attribute defined")
	} else {
		resource.Schema["id"] = &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the resource.",
		}
	}

	if v := resource.Timeouts; v != nil {
		e.HasTimeouts = true

		if v := v.Create; v != nil {
			e.DefaultCreateTimeout = int64(*v)
		}
		if v := v.Read; v != nil {
			e.DefaultReadTimeout = int64(*v)
		}
		if v := v.Update; v != nil {
			e.DefaultUpdateTimeout = int64(*v)
		}
		if v := v.Delete; v != nil {
			e.DefaultDeleteTimeout = int64(*v)
		}
	}

	e.fprintf(e.SchemaWriter, "schema.Schema{\n")

	if err := e.emitAttributesAndBlocks(nil, resource.Schema, true); err != nil {
		return err
	}

	if version := resource.SchemaVersion; version > 0 {
		e.fprintf(e.SchemaWriter, "Version:%d,\n", version)
	}

	if description := resource.Description; description != "" {
		e.fprintf(e.SchemaWriter, "Description:%q,\n", description)
	}

	if deprecationMessage := resource.DeprecationMessage; deprecationMessage != "" {
		e.fprintf(e.SchemaWriter, "DeprecationMessage:%q,\n", deprecationMessage)
	}

	e.fprintf(e.SchemaWriter, "}")

	return nil
}

// emitAttributesAndBlocks generates the Plugin Framework code for a set of Plugin SDK Attributes and Blocks
// and emits the generated code to the emitter's Writer.
// Property names are sorted prior to code generation to reduce diffs.
func (e *emitter) emitAttributesAndBlocks(path []string, schema map[string]*schema.Schema, emitModel bool) error {
	isTopLevel := len(path) == 0

	// At this point we are emitting code for a schema.Block or Schema.
	names := make([]string, 0, len(schema))
	utilityNames := make([]string, 0)
	for name, attr := range schema {
		if isTopLevel && name == "id" && attr.Description == "ID of the resource." {
			utilityNames = append(utilityNames, name)
			continue
		}
		if isTopLevel && name == "region" {
			utilityNames = append(utilityNames, name)
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)

	if isTopLevel {
		sort.Strings(utilityNames)
		names = append(utilityNames, names...)
	}

	emittedFieldName := false

	for i, name := range names {
		property := schema[name]

		if !isAttribute(property) {
			continue
		}

		if isTopLevel && i == len(utilityNames) {
			e.fprintf(e.ModelWriter, "\n")
		}

		if !emittedFieldName {
			e.fprintf(e.SchemaWriter, "Attributes: map[string]schema.Attribute{")
			emittedFieldName = true
		}

		e.fprintf(e.SchemaWriter, "\n%q:", name)
		if emitModel {
			e.fprintf(e.ModelWriter, "%s ", naming.ToCamelCase(name))
		}

		if err := e.emitAttributeProperty(append(path, name), property, emitModel); err != nil {
			return err
		}

		if emitModel {
			e.fprintf(e.ModelWriter, " `tfsdk:%q`\n", name)
		}
		e.fprintf(e.SchemaWriter, ",\n")
	}
	if emittedFieldName {
		e.fprintf(e.SchemaWriter, "},\n")
	}

	emittedFieldName = false
	for _, name := range names {
		property := schema[name]
		if isAttribute(property) {
			continue
		}

		if !emittedFieldName {
			e.fprintf(e.SchemaWriter, "Blocks: map[string]schema.Block{")
			emittedFieldName = true
		}

		e.fprintf(e.SchemaWriter, "\n%q:", name)
		if emitModel {
			e.fprintf(e.ModelWriter, "%s ", naming.ToCamelCase(name))
		}

		if err := e.emitBlockProperty(append(path, name), property); err != nil {
			return err
		}

		if emitModel {
			e.fprintf(e.ModelWriter, " `tfsdk:%q`\n", name)
		}
		e.fprintf(e.SchemaWriter, ",\n")
	}
	if emittedFieldName {
		e.fprintf(e.SchemaWriter, "},\n")
	}

	return nil
}

// emitAttributeProperty generates the Plugin Framework code for a Plugin SDK Attribute's property
// and emits the generated code to the emitter's Writer.
func (e *emitter) emitAttributeProperty(path []string, property *schema.Schema, emitModel bool) error {
	name := strings.Join(path, "_")
	var planModifiers []string
	var fwPlanModifierPackage, fwPlanModifierType, fwValidatorsPackage, fwValidatorType string

	// At this point we are emitting code for the values of a schema.Schema's Attributes (map[string]schema.Attribute).
	switch v := property.Type; v {
	//
	// Primitive types.
	//
	case schema.TypeBool:
		e.fprintf(e.SchemaWriter, "schema.BoolAttribute{\n")
		if emitModel {
			e.fprintf(e.ModelWriter, "types.Bool")
		}

		fwPlanModifierPackage = "boolplanmodifier"
		fwPlanModifierType = "Bool"

	case schema.TypeFloat:
		e.fprintf(e.SchemaWriter, "schema.Float64Attribute{\n")
		if emitModel {
			e.fprintf(e.ModelWriter, "types.Float64")
		}

		fwPlanModifierPackage = "float64planmodifier"
		fwPlanModifierType = "Float64"

	case schema.TypeInt:
		e.fprintf(e.SchemaWriter, "schema.Int64Attribute{\n")
		if emitModel {
			e.fprintf(e.ModelWriter, "types.Int64")
		}

		fwPlanModifierPackage = "int64planmodifier"
		fwPlanModifierType = "Int64"

	case schema.TypeString:
		e.fprintf(e.SchemaWriter, "schema.StringAttribute{\n")
		if emitModel {
			e.fprintf(e.ModelWriter, "types.String")
		}

		fwPlanModifierPackage = "stringplanmodifier"
		fwPlanModifierType = "String"

	//
	// Complex types.
	//
	case schema.TypeList, schema.TypeMap, schema.TypeSet:
		var aggregateSchemaFactory, typeName string

		switch v := property.Elem.(type) {
		case *schema.Schema:
			switch v := property.Type; v {
			case schema.TypeList:
				aggregateSchemaFactory = "schema.ListAttribute{"
				typeName = "list"

				if emitModel {
					e.fprintf(e.ModelWriter, "types.List")
				}

				fwPlanModifierPackage = "listplanmodifier"
				fwPlanModifierType = "List"
				fwValidatorsPackage = "listvalidator"
				fwValidatorType = "List"

			case schema.TypeMap:
				aggregateSchemaFactory = "schema.MapAttribute{"
				typeName = "map"

				if emitModel {
					e.fprintf(e.ModelWriter, "types.Map")
				}

				fwPlanModifierPackage = "mapplanmodifier"
				fwPlanModifierType = "Map"
				fwValidatorsPackage = "mapvalidator"
				fwValidatorType = "Map"

			case schema.TypeSet:
				aggregateSchemaFactory = "schema.SetAttribute{"
				typeName = "set"

				if emitModel {
					e.fprintf(e.ModelWriter, "types.Set")
				}

				fwPlanModifierPackage = "setplanmodifier"
				fwPlanModifierType = "Set"
				fwValidatorsPackage = "setvalidator"
				fwValidatorType = "Set"
			}

			var elementType string

			switch v := v.Type; v {
			case schema.TypeBool:
				elementType = "types.BoolType"

			case schema.TypeFloat:
				elementType = "types.Float64Type"

			case schema.TypeInt:
				elementType = "types.Int64Type"

			case schema.TypeString:
				elementType = "types.StringType"

			default:
				return unsupportedTypeError(path, fmt.Sprintf("(Attribute) %s of %s", typeName, v.String()))
			}

			e.fprintf(e.SchemaWriter, "%s\n", aggregateSchemaFactory)
			e.fprintf(e.SchemaWriter, "ElementType:%s,\n", elementType)

		case *schema.Resource:
			// We get here for Computed-only nested blocks or when ConfigMode is SchemaConfigModeBlock.
			var fwType string

			switch v := property.Type; v {
			case schema.TypeList:
				fwType = "types.List"
				aggregateSchemaFactory = "schema.ListNestedAttribute{"
				fwPlanModifierPackage = "listplanmodifier"
				fwPlanModifierType = "List"
				fwValidatorsPackage = "listvalidator"
				fwValidatorType = "List"

			case schema.TypeMap:
				fwType = "types.Map"
				aggregateSchemaFactory = "schema.MapNestedAttribute{"
				fwPlanModifierPackage = "mapplanmodifier"
				fwPlanModifierType = "Map"
				fwValidatorsPackage = "mapvalidator"
				fwValidatorType = "Map"

			case schema.TypeSet:
				fwType = "types.Set"
				aggregateSchemaFactory = "schema.SetNestedAttribute{"
				fwPlanModifierPackage = "setplanmodifier"
				fwPlanModifierType = "Set"
				fwValidatorsPackage = "setvalidator"
				fwValidatorType = "Set"

			default:
				return unsupportedTypeError(path, v.String())
			}

			e.fprintf(e.SchemaWriter, "%s\n", aggregateSchemaFactory)
			e.fprintf(e.SchemaWriter, "NestedObject:")
			e.fprintf(e.SchemaWriter, "schema.NestedAttributeObject{\n")

			if e.IsDataSource || isRequiredOnly(property) {
				e.fprintf(e.ModelWriter, "[]struct{\n")
				if err := e.emitAttributesAndBlocks(path, v.Schema, true); err != nil {
					return err
				}
				e.fprintf(e.ModelWriter, "\n}")
			} else {
				e.fprintf(e.ModelWriter, fwType)

				rootModelWriter := e.ModelWriter
				sbNestedModel := strings.Builder{}
				e.ModelWriter = &sbNestedModel
				if err := e.emitAttributesAndBlocks(path, v.Schema, true); err != nil {
					return err
				}

				sbNestedModelAttrTypes := strings.Builder{}
				e.ModelWriter = &sbNestedModelAttrTypes
				e.ImportFrameworkAttr = true

				e.fprintf(e.ModelWriter, "return map[string]attr.Type{\n")
				if err := e.emitAttributeTypes(path, v); err != nil {
					return err
				}
				e.fprintf(e.ModelWriter, "}\n")

				e.ModelWriter = rootModelWriter

				e.NestedModels = append(e.NestedModels, nestedModelData{
					Name:               naming.ToCamelCase(name),
					Model:              sbNestedModel.String(),
					AttributeTypesFunc: sbNestedModelAttrTypes.String(),
				})
			}

			e.fprintf(e.SchemaWriter, "},\n")

		default:
			return unsupportedTypeError(path, fmt.Sprintf("(Attribute) %s of %T", typeName, v))
		}

	default:
		return unsupportedTypeError(path, v.String())
	}

	if property.Required {
		e.fprintf(e.SchemaWriter, "Required:true,\n")
	}

	if property.Optional {
		e.fprintf(e.SchemaWriter, "Optional:true,\n")
	}

	if property.Computed {
		e.fprintf(e.SchemaWriter, "Computed:true,\n")
	}

	if property.Sensitive {
		e.fprintf(e.SchemaWriter, "Sensitive:true,\n")
	}

	var fwDefaultsPackage string

	if def := property.Default; def != nil {
		switch v := def.(type) {
		case bool:
			fwDefaultsPackage = "booldefault"
			e.fprintf(e.SchemaWriter, fmt.Sprintf("Default:%s.StaticBool(%t),\n", fwDefaultsPackage, v))
			e.FrameworkPlanModifierPackages = append(e.FrameworkPlanModifierPackages, fwDefaultsPackage)
		case int:
			fwDefaultsPackage = "int64default"
			e.fprintf(e.SchemaWriter, fmt.Sprintf("Default:%s.StaticInt64(%d),\n", fwDefaultsPackage, v))
			e.FrameworkPlanModifierPackages = append(e.FrameworkPlanModifierPackages, fwDefaultsPackage)
		case float64:
			fwDefaultsPackage = "float64default"
			e.fprintf(e.SchemaWriter, fmt.Sprintf("Default:%s.StaticFloat64(%f),\n", fwDefaultsPackage, v))
			e.FrameworkPlanModifierPackages = append(e.FrameworkPlanModifierPackages, fwDefaultsPackage)
		case string:
			fwDefaultsPackage = "stringdefault"
			e.fprintf(e.SchemaWriter, fmt.Sprintf("Default:%s.StaticString(%q),\n", fwDefaultsPackage, v))
			e.FrameworkPlanModifierPackages = append(e.FrameworkPlanModifierPackages, fwDefaultsPackage)
		default:
			e.fprintf(e.SchemaWriter, "// TODO Default:%#v,\n", def)
		}
	}

	if description := property.Description; description != "" {
		e.fprintf(e.SchemaWriter, "Description:%q,\n", description)
	}

	if deprecationMessage := property.Deprecated; deprecationMessage != "" {
		e.fprintf(e.SchemaWriter, "DeprecationMessage:%q,\n", deprecationMessage)
	}

	if maxItems, minItems := property.MaxItems, property.MinItems; maxItems > 0 || minItems > 0 && fwValidatorsPackage != "" && fwValidatorType != "" {
		e.FrameworkValidatorsPackages = append(e.FrameworkValidatorsPackages, fwValidatorsPackage)

		e.fprintf(e.SchemaWriter, "Validators:[]validator.%s{\n", fwValidatorType)
		if minItems > 0 {
			e.fprintf(e.SchemaWriter, "%s.SizeAtLeast(%d),\n", fwValidatorsPackage, minItems)
		}
		if maxItems > 0 {
			e.fprintf(e.SchemaWriter, "%s.SizeAtMost(%d),\n", fwValidatorsPackage, maxItems)
		}
		e.fprintf(e.SchemaWriter, "},\n")
	}

	if property.ForceNew {
		planModifiers = append(planModifiers, fmt.Sprintf("%s.RequiresReplace()", fwPlanModifierPackage))
		e.FrameworkPlanModifierPackages = append(e.FrameworkPlanModifierPackages, fwPlanModifierPackage)
	}

	if len(planModifiers) > 0 {
		e.fprintf(e.SchemaWriter, "PlanModifiers:[]planmodifier.%s{\n", fwPlanModifierType)
		for _, planModifier := range planModifiers {
			e.fprintf(e.SchemaWriter, "%s,\n", planModifier)
		}
		e.fprintf(e.SchemaWriter, "},\n")
	}

	// Features that we can't (yet) migrate:

	if property.ValidateFunc != nil || property.ValidateDiagFunc != nil {
		e.fprintf(e.SchemaWriter, "// TODO Validate,\n")
	}

	e.fprintf(e.SchemaWriter, "}")

	return nil
}

// emitBlockProperty generates the Plugin Framework code for a Plugin SDK Block's property
// and emits the generated code to the emitter's Writer.
func (e *emitter) emitBlockProperty(path []string, property *schema.Schema) error {
	name := strings.Join(path, "_")
	var planModifiers []string
	var fwPlanModifierPackage, fwPlanModifierType, fwValidatorsPackage, fwValidatorType string

	// At this point we are emitting code for the values of a schema.Block or Schema's Blocks (map[string]schema.Block).
	switch v := property.Type; v {
	//
	// Complex types.
	//
	case schema.TypeList:
		switch v := property.Elem.(type) {
		case *schema.Resource:
			fwPlanModifierPackage = "listplanmodifier"
			fwPlanModifierType = "List"
			fwValidatorsPackage = "listvalidator"
			fwValidatorType = "List"

			e.fprintf(e.SchemaWriter, "schema.ListNestedBlock{\n")
			e.fprintf(e.SchemaWriter, "NestedObject:schema.NestedBlockObject{\n")
			e.fprintf(e.ModelWriter, "types.List")

			rootModelWriter := e.ModelWriter
			sbNestedModel := strings.Builder{}
			e.ModelWriter = &sbNestedModel
			if err := e.emitAttributesAndBlocks(path, v.Schema, true); err != nil {
				return err
			}

			sbNestedModelAttrTypes := strings.Builder{}
			e.ModelWriter = &sbNestedModelAttrTypes
			e.ImportFrameworkAttr = true

			e.fprintf(e.ModelWriter, "return map[string]attr.Type{\n")
			if err := e.emitAttributeTypes(path, v); err != nil {
				return err
			}
			e.fprintf(e.ModelWriter, "}\n")

			e.ModelWriter = rootModelWriter

			e.NestedModels = append(e.NestedModels, nestedModelData{
				Name:               naming.ToCamelCase(name),
				Model:              sbNestedModel.String(),
				AttributeTypesFunc: sbNestedModelAttrTypes.String(),
			})

			e.fprintf(e.SchemaWriter, "},\n")

		default:
			return unsupportedTypeError(path, fmt.Sprintf("(Block) list of %T", v))
		}

	case schema.TypeSet:
		switch v := property.Elem.(type) {
		case *schema.Resource:
			fwPlanModifierPackage = "setplanmodifier"
			fwPlanModifierType = "Set"
			fwValidatorsPackage = "setvalidator"
			fwValidatorType = "Set"

			e.fprintf(e.SchemaWriter, "schema.SetNestedBlock{\n")
			e.fprintf(e.SchemaWriter, "NestedObject:schema.NestedBlockObject{\n")
			e.fprintf(e.ModelWriter, "types.Set")

			rootModelWriter := e.ModelWriter
			sbNestedModel := strings.Builder{}
			e.ModelWriter = &sbNestedModel
			if err := e.emitAttributesAndBlocks(path, v.Schema, true); err != nil {
				return err
			}

			sbNestedModelAttrTypes := strings.Builder{}
			e.ModelWriter = &sbNestedModelAttrTypes
			e.ImportFrameworkAttr = true

			e.fprintf(e.ModelWriter, "return map[string]attr.Type{\n")
			if err := e.emitAttributeTypes(path, v); err != nil {
				return err
			}
			e.fprintf(e.ModelWriter, "}\n")

			e.ModelWriter = rootModelWriter

			e.NestedModels = append(e.NestedModels, nestedModelData{
				Name:               naming.ToCamelCase(name),
				Model:              sbNestedModel.String(),
				AttributeTypesFunc: sbNestedModelAttrTypes.String(),
			})

			e.fprintf(e.SchemaWriter, "},\n")
		default:
			return unsupportedTypeError(path, fmt.Sprintf("(Block) set of %T", v))
		}

	default:
		return unsupportedTypeError(path, v.String())
	}

	// Compatibility hacks.
	// See Schema::coreConfigSchemaBlock.
	if property.Required && property.MinItems == 0 {
		property.MinItems = 1
	}
	if property.Optional && property.MinItems > 0 {
		property.MinItems = 0
	}
	if property.Computed && !property.Optional {
		property.MaxItems = 0
		property.MinItems = 0
	}

	if description := property.Description; description != "" {
		e.fprintf(e.SchemaWriter, "Description:%q,\n", description)
	}

	if deprecationMessage := property.Deprecated; deprecationMessage != "" {
		e.fprintf(e.SchemaWriter, "DeprecationMessage:%q,\n", deprecationMessage)
	}

	if maxItems, minItems := property.MaxItems, property.MinItems; maxItems > 0 || minItems > 0 && fwValidatorsPackage != "" && fwValidatorType != "" {
		e.FrameworkValidatorsPackages = append(e.FrameworkValidatorsPackages, fwValidatorsPackage)

		e.fprintf(e.SchemaWriter, "Validators:[]validator.%s{\n", fwValidatorType)
		if minItems > 0 {
			e.fprintf(e.SchemaWriter, "%s.SizeAtLeast(%d),\n", fwValidatorsPackage, minItems)
		}
		if maxItems > 0 {
			e.fprintf(e.SchemaWriter, "%s.SizeAtMost(%d),\n", fwValidatorsPackage, maxItems)
		}
		e.fprintf(e.SchemaWriter, "},\n")
	}

	if property.ForceNew {
		planModifiers = append(planModifiers, fmt.Sprintf("%s.RequiresReplace()", fwPlanModifierPackage))
		e.FrameworkPlanModifierPackages = append(e.FrameworkPlanModifierPackages, fwPlanModifierPackage)
	}

	if len(planModifiers) > 0 {
		e.fprintf(e.SchemaWriter, "PlanModifiers:[]planmodifier.%s{\n", fwPlanModifierType)
		for _, planModifier := range planModifiers {
			e.fprintf(e.SchemaWriter, "%s,\n", planModifier)
		}
		e.fprintf(e.SchemaWriter, "},\n")
	}

	if def := property.Default; def != nil {
		e.warnf("Block %s has non-nil Default: %v", strings.Join(path, "/"), def)
	}

	e.fprintf(e.SchemaWriter, "}")

	return nil
}

func (e *emitter) emitAttributeTypes(path []string, rs *schema.Resource) error {
	var names []string
	for name := range rs.Schema {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		property := rs.Schema[name]
		var fwType string
		var elementType string

		switch v := property.Type; v {
		case schema.TypeBool:
			fwType = "types.BoolType"
		case schema.TypeInt:
			fwType = "types.Int64Type"
		case schema.TypeFloat:
			fwType = "types.Float64Type"
		case schema.TypeString:
			fwType = "types.StringType"
		case schema.TypeList, schema.TypeMap, schema.TypeSet:
			switch v := property.Elem.(type) {
			case *schema.Schema:
				switch v := property.Type; v {
				case schema.TypeList:
					fwType = "types.ListType"
				case schema.TypeMap:
					fwType = "types.MapType"
				case schema.TypeSet:
					fwType = "types.SetType"
				}

				switch v := v.Type; v {
				case schema.TypeBool:
					elementType = "types.BoolType"
				case schema.TypeFloat:
					elementType = "types.Float64Type"
				case schema.TypeInt:
					elementType = "types.Int64Type"
				case schema.TypeString:
					elementType = "types.StringType"
				default:
					return unsupportedTypeError(path, fmt.Sprintf("path: %s", v.String()))
				}

				e.fprintf(e.ModelWriter, "%q:%s{ElemType:%s},\n", name, fwType, elementType)
				continue

			case *schema.Resource:
				switch property.Type {
				case schema.TypeList:
					fwType = "types.ListType{"
				case schema.TypeMap:
					fwType = "types.MapType{"
				case schema.TypeSet:
					fwType = "types.SetType{"
				}

				e.fprintf(e.ModelWriter, "%q:%s\n", name, fwType)
				e.fprintf(e.ModelWriter, "ElemType:types.ObjectType{\nAttrTypes:map[string]attr.Type{\n")
				if err := e.emitAttributeTypes(append(path, name), v); err != nil {
					return err
				}
				e.fprintf(e.ModelWriter, "},\n},\n}, // Consider using corresponding .AttrTypes() result \n")
				continue
			}
		default:
			return unsupportedTypeError(path, v.String())
		}

		e.fprintf(e.ModelWriter, "%q:%s,\n", name, fwType)
	}

	return nil
}

func isRequiredOnly(s *schema.Schema) bool {
	if s == nil {
		return false
	}
	return s.Required && !s.Computed
}

// warnf emits a formatted warning message to the UI.
func (e *emitter) warnf(format string, a ...interface{}) {
	e.Generator.warnf(format, a...)
}

// errof emits a formatted error message to the UI.
func (e *emitter) errorf(format string, a ...interface{}) {
	e.Generator.errorf(format, a...)
}

// e.fprintf writes a formatted string to a Writer.
func (e *emitter) fprintf(w io.Writer, format string, a ...interface{}) {
	_, err := io.WriteString(w, fmt.Sprintf(format, a...))
	if err != nil {
		e.errorf("error writing: %s", err)
	}
}

// isAttribute returns whether or not the specified property should be emitted as an Attribute (vs. a Block).
// See https://github.com/hashicorp/terraform-plugin-sdk/blob/6ffc92796f0716c07502e4d36aaafa5fd85e94cf/helper/schema/core_schema.go#L57.
func isAttribute(property *schema.Schema) bool {
	if property.Elem == nil {
		return true
	}

	if property.Type == schema.TypeMap {
		return true
	}

	switch property.ConfigMode {
	case schema.SchemaConfigModeAttr:
		return true

	case schema.SchemaConfigModeBlock:
		return false

	default:
		if property.Computed && !property.Optional {
			// Computed-only schemas are always handled as attributes because they never appear in configuration.
			return true
		}

		if _, ok := property.Elem.(*schema.Schema); ok {
			return true
		}
	}

	return false
}

func unsupportedTypeError(path []string, typ string) error {
	return fmt.Errorf("%s is of unsupported type: %s", strings.Join(path, "/"), typ)
}

type templateData struct {
	DefaultCreateTimeout          int64
	DefaultReadTimeout            int64
	DefaultUpdateTimeout          int64
	DefaultDeleteTimeout          int64
	EmitResourceImportState       bool
	EmitResourceModifyPlan        bool
	EmitResourceUpdateSkeleton    bool
	FrameworkPlanModifierPackages []string
	FrameworkValidatorsPackages   []string
	HasTimeouts                   bool
	ImportFrameworkAttr           bool
	ImportProviderFrameworkTypes  bool
	Name                          string // e.g. Instance
	PackageName                   string // e.g. db
	ServiceName                   string // e.g. Databases
	ProviderPlanModifierPackages  []string
	Schema                        string
	Model                         string
	NestedModels                  []nestedModelData
	TFTypeName                    string // e.g. vkcs_db_instance
}

type nestedModelData struct {
	Name               string
	Model              string
	AttributeTypesFunc string
}

type testTemplateData struct {
	Name            string
	PackageName     string
	ReleasedVersion string
}

//go:embed datasource.tmpl
var datasourceImpl string

//go:embed datasource_test.tmpl
var datasourceTestImpl string

//go:embed resource.tmpl
var resourceImpl string

//go:embed resource_test.tmpl
var resourceTestImpl string
