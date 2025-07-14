package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/generator"
	"github.com/Patrick-Ivann/xsd-codegen/pkg/parser"
	"github.com/Patrick-Ivann/xsd-codegen/pkg/xmlgen"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: xsd-codegen <path/to/schema.xsd>")
		os.Exit(1)
	}

	xsdPath := os.Args[1]

	// Check that the path exists and is a file
	info, err := os.Stat(xsdPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if info.IsDir() {
		fmt.Fprintf(os.Stderr, "error: %s is a directory, not a file\n", xsdPath)
		os.Exit(1)
	}

	absPath, err := filepath.Abs(xsdPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving absolute path: %v\n", err)
		os.Exit(1)
	}
	baseDir := filepath.Dir(absPath)

	// Parse the XSD file
	schema, err := parser.ParseXSD(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parser error: %v\n", err)
		os.Exit(1)
	}

	// Resolve imports/includes
	if err := parser.ResolveImportsAndIncludes(schema, absPath); err != nil {
		fmt.Fprintf(os.Stderr, "resolve error: %v\n", err)
		os.Exit(1)
	}

	// Generate Go code
	tmplPath := filepath.Join(baseDir, "structs.tmpl") // or use a fixed template path
	goFile := filepath.Join(baseDir, "schema_gen.go")
	if err := generator.GenerateGoCode(schema, tmplPath, goFile); err != nil {
		fmt.Fprintf(os.Stderr, "Go code generation error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Go code generated: %s\n", goFile)

	// Generate dummy XML for the first type (as an example)
	if len(schema.Types) > 0 {
		rootType := schema.Types[0]
		xmlBytes, err := xmlgen.MarshalDummyXML(rootType.Name, rootType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Dummy XML generation error: %v\n", err)
			os.Exit(1)
		}
		xmlFile := filepath.Join(baseDir, "dummy.xml")
		if err := os.WriteFile(xmlFile, xmlBytes, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write dummy XML: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Dummy XML generated: %s\n", xmlFile)
	} else {
		fmt.Println("No types found in schema; skipping dummy XML generation.")
	}

	fmt.Println("Schema parsed, resolved, and code generated successfully.")
	jsonBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal schema to JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonBytes))
}
