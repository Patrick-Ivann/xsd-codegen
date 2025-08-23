package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/helpers"
	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
	"github.com/Patrick-Ivann/xsd-codegen/pkg/parser"
	"github.com/Patrick-Ivann/xsd-codegen/pkg/xmlgen"
	"github.com/beevik/etree"
)

func main() {
	xsdPath, outPath := parseFlags()
	schema := mustParseSchema(xsdPath)

	doc, found := generateXMLDocument(schema)
	if !found {
		log.Fatal("Entrypoint element 'purchaseOrder' not found in schema.")
	}
	writeOutput(doc, outPath)
}

// parseFlags handles command-line flag parsing and validation.
// Uses flag.StringVar to avoid immediate dereference issues.
func parseFlags() (pathToXSDFile, pathToOutputXML string) {
	var xsdPath, outPath string
	flag.StringVar(&xsdPath, "xsd", "", "Path to XSD file")
	flag.StringVar(&outPath, "out", "", "Output XML file path (default stdout)")
	flag.Parse()

	if xsdPath == "" {
		log.Fatal("XSD file path is required. Use -xsd flag.")
	}
	return xsdPath, outPath
}

// mustParseSchema parses the XSD file and exits the program on error.
// Pass schema elements by pointer for efficiency.
func mustParseSchema(xsdPath string) *model.XSDSchema {
	schema, err := parser.ParseXSD(xsdPath, nil)
	if err != nil {
		log.Fatalf("Failed to parse XSD: %v", err)
	}
	return schema
}

// generateXMLDocument creates the XML document with root populated from the entrypoint element.
func generateXMLDocument(schema *model.XSDSchema) (*etree.Document, bool) {
	doc := etree.NewDocument()
	for i := range schema.Elements {
		el := &schema.Elements[i]       // pass element by pointer to avoid copying large struct
		if el.Name != "purchaseOrder" { // invert if to reduce nesting and continue early
			continue
		}
		root := xmlgen.GenerateElement(schema, el, helpers.DefaultValueGenerator{})
		// Add schema namespace attributes
		root.CreateAttr("xmlns", schema.TargetNamespace)
		root.CreateAttr("xsi:schemaLocation", schema.TargetNamespace+" schema.xsd")
		root.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
		doc.SetRoot(root)
		return doc, true
	}
	return nil, false
}

// writeOutput outputs the XML document either to stdout or a file, with indenting for readability.
// Checks file.Close error explicitly without using fatal inside defer.
// Ensures output path does not allow directory traversal for safer file creation.
func writeOutput(doc *etree.Document, outPath string) {
	doc.Indent(2)
	if outPath == "" {
		_, err := doc.WriteTo(os.Stdout)
		if err != nil {
			log.Fatalf("Unable to write output: %v", err)
		}
		return
	}

	// Sanitize output path: e.g., prevent path traversal attacks or unsupported characters.
	cleanPath := filepath.Clean(outPath)
	if strings.HasPrefix(cleanPath, "..") || strings.Contains(cleanPath, string(os.PathSeparator)+"..") {
		log.Fatalf("Invalid output file path: %s", outPath)
	}

	file, err := os.Create(cleanPath)
	if err != nil {
		log.Fatalf("Unable to create output file: %v", err)
	}

	_, err = doc.WriteTo(file)
	if err != nil {
		fileCloseErr := file.Close()
		if fileCloseErr != nil {
			log.Fatalf("Unable to write output file: %v", fileCloseErr)
		}
		log.Fatalf("Unable to write output file: %v", err)
	}

	// Explicitly check Close error here, after successful write
	if err = file.Close(); err != nil {
		log.Fatalf("Failed to close output file: %v", err)
	}
}
