package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	f := flag.NewFlagSet("generate-providerschema", flag.ExitOnError)

	path := f.String("path", "./schema_base.go", "export the schema to the given path/filename.")
	packageName := f.String("package", "provider", "the package where the schema file lays in")
	jsonSchemaPath := f.String("schemajson", ".release/provider-schema.json", "path to provider schema json file")

	if err := f.Parse(os.Args[1:]); err != nil {
		fmt.Printf("error parsing args: %+v", err)
		os.Exit(1)
	}

	run(*path, *packageName, *jsonSchemaPath)
}

func run(path, packageName, jsonSchemaPath string) {
	generator := ProviderSchemaJSONGenerator{
		PackageName: packageName,
		SchemaPath:  jsonSchemaPath,
	}
	if err := goFmtAndWriteToFile(path, generator.Code()); err != nil {
		panic(err)
	}
}

type ProviderSchemaJSONGenerator struct {
	PackageName string
	Version     string
	SchemaPath  string
}

func (g ProviderSchemaJSONGenerator) Code() string {
	return fmt.Sprintf(`package %s

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

%s 
`, g.PackageName, g.codeForSchema())
}

func (g ProviderSchemaJSONGenerator) codeForSchema() string {
	pSchemaJSON, err := os.ReadFile(g.SchemaPath)
	if err != nil {
		panic(err)
	}

	compactSchemaJSON := new(bytes.Buffer)
	err = json.Compact(compactSchemaJSON, pSchemaJSON)
	if err != nil {
		panic(err)
	}

	pSchema := strings.ReplaceAll(compactSchemaJSON.String(), "`", "`+\"`\"+`")

	return fmt.Sprintf("const ProviderSchemaJSON = `%s`", pSchema)
}

func goFmtAndWriteToFile(filePath, fileContents string) error {
	fmt, err := GolangCodeFormatter{}.Format(fileContents)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, []byte(*fmt), 0o644); err != nil {
		return err
	}

	return nil
}

type GolangCodeFormatter struct{}

func (f GolangCodeFormatter) Format(input string) (*string, error) {
	tmpfile, err := os.CreateTemp("", "temp-*.go")
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %+v", err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	filePath := tmpfile.Name()

	if _, err := tmpfile.WriteString(input); err != nil {
		return nil, fmt.Errorf("writing contents to %q: %+v", filePath, err)
	}

	f.runGoFmt(filePath)
	f.runGoImports(filePath)

	contents, err := f.readFileContents(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading contents from %q: %+v", filePath, err)
	}

	return contents, nil
}

func (f GolangCodeFormatter) runGoFmt(filePath string) {
	cmd := exec.Command("gofmt", "-w", filePath)
	// intentionally not using these errors since the exit codes are kinda uninteresting
	_ = cmd.Start()
	_ = cmd.Wait()
}

func (f GolangCodeFormatter) runGoImports(filePath string) {
	cmd := exec.Command("goimports", "-w", filePath)
	// intentionally not using these errors since the exit codes are kinda uninteresting
	_ = cmd.Start()
	_ = cmd.Wait()
}

func (f GolangCodeFormatter) readFileContents(filePath string) (*string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	contents := string(data)
	return &contents, nil
}
