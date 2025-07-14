package main

import (
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

	schema, err := parser.ParseXSD(xsdPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parser error: %v\n", err)
		os.Exit(1)
	}

	basePath := filepath.Dir(xsdPath)
	if err := parser.ResolveImportsAndIncludes(schema, basePath); err != nil {
		fmt.Fprintf(os.Stderr, "resolve error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Schema parsed and resolved successfully.")
}
