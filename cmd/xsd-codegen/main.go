package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/parser"
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

	// Use absolute path for base directory
	absPath, err := filepath.Abs(xsdPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving absolute path: %v\n", err)
		os.Exit(1)
	}

	// Parse the XSD file
	schema, err := parser.ParseXSD(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parser error: %v\n", err)
		os.Exit(1)
	}

	// Resolve imports and includes using the base directory
	if err := parser.ResolveImportsAndIncludes(schema, xsdPath); err != nil {
		fmt.Fprintf(os.Stderr, "resolve error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Schema parsed and resolved successfully.")

	jsonBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal schema to JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonBytes))
}
