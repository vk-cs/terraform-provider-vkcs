// Wrapper around tfplugingen-framework that works around issue #20
// (collision-prone naming of nested types).
//
// What it does:
//  1. Reads spec.json.
//  2. Walks every nested attribute depth-first with alphabetical ordering
//     (list_nested, single_nested, map_nested, set_nested, object).
//  3. For each such attribute, replaces its "name" field with a path-qualified
//     form ("configs_maintenance_backup_full" instead of just "full"). Thanks to
//     this, the generator emits unique type names (ConfigsMaintenanceBackupFullType
//     instead of the colliding FullType).
//  4. Runs tfplugingen-framework on the rewritten spec.
//  5. Post-processes the result: restores the short form wherever the name is
//     user-facing or referenced by other Go code - tfsdk tags, parent struct
//     field names, map keys in Attributes()/AttributeTypes().
//  6. Writes the final file to the requested path.
//
// Idea: TYPE names (XxxType, XxxValue) stay qualified - the compiler needs that.
// Everything user-facing (HCL names, Go struct fields) stays short.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func main() {
	var (
		specPath      = flag.String("spec", "", "path to spec.json (required)")
		outPath       = flag.String("out", "", "path to write final cluster_resource_gen.go (required)")
		pkgName       = flag.String("package", "resource_cluster", "Go package name")
		generator     = flag.String("generator", "tfplugingen-framework", "path to tfplugingen-framework binary")
		overridesPath = flag.String("overrides", "", "path to dpgen-overrides.json (optional)")
		keepTmpDir    = flag.Bool("keep-tmp", false, "do not remove temp dir on exit (debug)")
	)
	flag.Parse()
	if *specPath == "" || *outPath == "" {
		fmt.Fprintln(os.Stderr, "usage: dpgen -spec <spec.json> -out <out.go> [-package <name>]")
		os.Exit(2)
	}

	specBytes, err := os.ReadFile(*specPath)
	if err != nil {
		die("read spec: %v", err)
	}

	// Spec is walked loosely-typed through generic maps; full schema knowledge is unnecessary.
	var spec map[string]any
	if err := json.Unmarshal(specBytes, &spec); err != nil {
		die("parse spec: %v", err)
	}

	// renames collects pairs used during post-processing.
	// snake is the original name (last path segment); pascalShort is its pascal form.
	// snakeQualified is what we temporarily substitute into the spec (full path joined by underscores).
	// pascalQualified is what the generator emits as a chunk of the type/field name.
	type rename struct {
		snake           string
		pascalShort     string
		snakeQualified  string
		pascalQualified string
	}
	var renames []rename
	seenQualified := map[string]bool{}

	// Spec has the shape {"resources":[{"name":"cluster","schema":{"attributes":[...]}}], ...}.
	// Resource attributes live under resources[].schema.attributes; inside, walk depth-first
	// with alphabetical ordering (matching the generator's traversal order).
	resources, _ := spec["resources"].([]any)
	for _, r := range resources {
		rmap, _ := r.(map[string]any)
		if rmap == nil {
			continue
		}
		schema, _ := rmap["schema"].(map[string]any)
		if schema == nil {
			continue
		}
		attrs, _ := schema["attributes"].([]any)
		walkAttributeList(attrs, nil, func(path []string, attr map[string]any) {
			if os.Getenv("DPGEN_DEBUG") != "" {
				fmt.Fprintf(os.Stderr, "visit: %s (kinds: %v)\n", strings.Join(path, "."), kindsOf(attr))
			}
			if !generatesNestedType(attr) {
				return
			}
			short, _ := attr["name"].(string)
			if short == "" {
				return
			}
			snakeQualified := strings.Join(path, "_")
			pascalQualified := snakeToPascal(snakeQualified)
			pascalShort := snakeToPascal(short)
			if seenQualified[pascalQualified] {
				// Unique by construction; guard against double visits.
				return
			}
			seenQualified[pascalQualified] = true
			renames = append(renames, rename{
				snake:           short,
				pascalShort:     pascalShort,
				snakeQualified:  snakeQualified,
				pascalQualified: pascalQualified,
			})
			// Rename in place; after this the spec goes to the generator with unique names.
			attr["name"] = snakeQualified
		})
	}

	// Sort renames so longer matches come first - guards against the case where one
	// qualified prefix is a prefix of another.
	sort.SliceStable(renames, func(i, j int) bool {
		return len(renames[i].pascalQualified) > len(renames[j].pascalQualified)
	})

	// Write the modified spec to a temp file.
	tmpDir, err := os.MkdirTemp("", "dpgen-*")
	if err != nil {
		die("mktemp: %v", err)
	}
	if !*keepTmpDir {
		defer os.RemoveAll(tmpDir)
	} else {
		fmt.Fprintf(os.Stderr, "tmp dir: %s\n", tmpDir)
	}

	tmpSpec := filepath.Join(tmpDir, "spec.json")
	out, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		die("marshal spec: %v", err)
	}
	if err := os.WriteFile(tmpSpec, out, 0o644); err != nil {
		die("write tmp spec: %v", err)
	}

	// Run the generator.
	tmpOut := filepath.Join(tmpDir, "out")
	if err := os.MkdirAll(tmpOut, 0o755); err != nil {
		die("mkdir tmp out: %v", err)
	}
	cmd := exec.Command(*generator, "generate", "resources",
		"--input", tmpSpec, "--output", tmpOut, "--package", *pkgName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		die("run generator: %v", err)
	}

	genFile := filepath.Join(tmpOut, "cluster_resource_gen.go")
	body, err := os.ReadFile(genFile)
	if err != nil {
		die("read gen: %v", err)
	}
	src := string(body)

	// Post-processing: for each rename, restore the user-facing form to the short one.
	// Strict replacements in two contexts. Type names (XxxType, XxxValue) are NOT touched -
	// they must stay qualified, otherwise the collisions return.
	for _, r := range renames {
		// 1) Any string literal in snake-qualified form → short form.
		// Covers tfsdk tags, map keys in Attributes()/AttributeTypes(),
		// bracket assignments like attrTypes["..."] = ..., etc. Snake names are
		// specific enough (with underscores and full path) that no other text collides.
		src = strings.ReplaceAll(src,
			"\""+r.snakeQualified+"\"",
			"\""+r.snake+"\"",
		)
		// 2) Field name in the parent struct and its use sites: PascalQualified → PascalShort,
		// but NOT inside type-name identifiers (PascalQualifiedType/Value).
		// The \b on both sides guarantees we won't match inside XxxType.
		fieldRe := regexp.MustCompile(`\b` + regexp.QuoteMeta(r.pascalQualified) + `\b`)
		src = fieldRe.ReplaceAllString(src, r.pascalShort)
	}

	// Apply manual override names. These touch ONLY field contexts (field declaration,
	// selector, struct-literal, String() return); type contexts (type X struct,
	// func (t X), X{}, X{...}) are left alone.
	if *overridesPath != "" {
		appliedOverrides, err := applyFieldOverrides(&src, *overridesPath)
		if err != nil {
			die("apply overrides: %v", err)
		}
		fmt.Fprintf(os.Stderr, "applied %d override(s) from %s\n", appliedOverrides, *overridesPath)
	}

	if err := os.WriteFile(*outPath, []byte(src), 0o644); err != nil {
		die("write out: %v", err)
	}
	fmt.Fprintf(os.Stderr, "wrote %s (%d nested types qualified)\n", *outPath, len(renames))
}

// applyFieldOverrides reads a JSON file with {from, to} renames and applies them
// to src in FIELD contexts. Returns the number of applied rules.
func applyFieldOverrides(src *string, path string) (int, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read overrides: %w", err)
	}
	var cfg struct {
		FieldRenames []struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"field_renames"`
	}
	if err := json.Unmarshal(body, &cfg); err != nil {
		return 0, fmt.Errorf("parse overrides: %w", err)
	}
	for _, r := range cfg.FieldRenames {
		if r.From == "" || r.To == "" {
			continue
		}
		from := regexp.QuoteMeta(r.From)
		// 1) Struct field declaration: "\n\t...<from>\s+basetypes.X" -> "\n\t...<to>\s+basetypes.X".
		decl := regexp.MustCompile(`(\n[\t ]+)` + from + `(\s+basetypes\.)`)
		*src = decl.ReplaceAllString(*src, `${1}`+r.To+`${2}`)
		// 2) Selector: ".<from>\b" -> ".<to>".
		sel := regexp.MustCompile(`\.` + from + `\b`)
		*src = sel.ReplaceAllString(*src, `.`+r.To)
		// 3) Field in struct literal: "\n\t...<from>:\s" -> "\n\t...<to>:\s".
		lit := regexp.MustCompile(`(\n[\t ]+)` + from + `(:\s)`)
		*src = lit.ReplaceAllString(*src, `${1}`+r.To+`${2}`)
		// 4) String() method on the wrapper type returns the type name as a debug label.
		// When a rename is applied, the label is worth keeping in sync.
		// This matches `func (t <from>) String() string { return "<from>" }`.
		strRet := regexp.MustCompile(
			`(func \(t ` + from + `\) String\(\) string \{\n[\t ]*return ")` + from + `("\n)`,
		)
		*src = strRet.ReplaceAllString(*src, `${1}`+r.To+`${2}`)
	}
	return len(cfg.FieldRenames), nil
}

// walkAttributeList traverses an attribute array ([{"name": "...", "list_nested": ...}, ...])
// in alphabetical order (matching the generator's SortedKeys()). For each attribute it
// calls cb(path, attr), then recurses into its nested variants.
func walkAttributeList(attrs []any, path []string, cb func(path []string, attr map[string]any)) {
	sortAttrsByName(attrs)
	for _, a := range attrs {
		am, _ := a.(map[string]any)
		if am == nil {
			continue
		}
		name, _ := am["name"].(string)
		if name == "" {
			continue
		}
		fullPath := append(append([]string{}, path...), name)
		cb(fullPath, am)

		// Recurse into the attribute's type variants. For list_nested/map_nested/set_nested
		// the children live inside nested_object.attributes; for single_nested they live
		// directly under attributes.
		for _, kind := range []string{"list_nested", "map_nested", "set_nested"} {
			holder, ok := am[kind].(map[string]any)
			if !ok {
				continue
			}
			nested, _ := holder["nested_object"].(map[string]any)
			if nested == nil {
				continue
			}
			children, _ := nested["attributes"].([]any)
			walkAttributeList(children, fullPath, cb)
		}
		if sn, ok := am["single_nested"].(map[string]any); ok {
			children, _ := sn["attributes"].([]any)
			walkAttributeList(children, fullPath, cb)
		}
		if obj, ok := am["object"].(map[string]any); ok {
			children, _ := obj["attributes"].([]any)
			walkAttributeList(children, fullPath, cb)
		}
	}
}

func sortAttrsByName(attrs []any) {
	sort.SliceStable(attrs, func(i, j int) bool {
		ni, _ := attrs[i].(map[string]any)["name"].(string)
		nj, _ := attrs[j].(map[string]any)["name"].(string)
		return ni < nj
	})
}

// generatesNestedType returns true if the attribute produces its own *Type/*Value
// (i.e. a compileable type that can collide with same-named siblings).
func generatesNestedType(attr map[string]any) bool {
	for _, key := range []string{"list_nested", "single_nested", "map_nested", "set_nested"} {
		if _, ok := attr[key]; ok {
			return true
		}
	}
	return false
}

func kindsOf(attr map[string]any) []string {
	var ks []string
	for _, k := range []string{"string", "int64", "bool", "list", "map", "set", "object", "list_nested", "single_nested", "map_nested", "set_nested", "float64", "number"} {
		if _, ok := attr[k]; ok {
			ks = append(ks, k)
		}
	}
	return ks
}

func snakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "dpgen: "+format+"\n", args...)
	os.Exit(1)
}
