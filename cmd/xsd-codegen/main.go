package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
func parseFlags() (string, string) {
	xsdPath := flag.String("xsd", "", "Path to XSD file")
	outPath := flag.String("out", "", "Output XML file path (default stdout)")
	flag.Parse()

	if *xsdPath == "" {
		log.Fatal("XSD file path is required. Use -xsd flag.")
	}
	return *xsdPath, *outPath
}

// mustParseSchema parses the XSD file and exits the program on error.
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
	for _, el := range schema.Elements {
		if el.Name == "purchaseOrder" { // entrypoint element
			root := xmlgen.GenerateElement(schema, el, helpers.DefaultValueGenerator{})
			// Add schema namespace attributes
			root.CreateAttr("xmlns", schema.TargetNamespace)
			root.CreateAttr("xsi:schemaLocation", schema.TargetNamespace+" schema.xsd")
			root.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
			doc.SetRoot(root)
			return doc, true
		}
	}
	return nil, false
}

// writeOutput outputs the XML document either to stdout or a file, with indenting for readability.
func writeOutput(doc *etree.Document, outPath string) {
	doc.Indent(2)
	if outPath == "" {
		_, err := doc.WriteTo(os.Stdout)
		if err != nil {
			log.Fatalf("Unable to write output file: %v", err)
		}
	} else {
		file, err := os.Create(outPath)
		if err != nil {
			log.Fatalf("Unable to create output file: %v", err)
		}
		defer file.Close()
		_, err = doc.WriteTo(file)
		if err != nil {
			log.Fatalf("Unable to write output file: %v", err)
		}
		fmt.Println("XML generated at", outPath)
	}
}
