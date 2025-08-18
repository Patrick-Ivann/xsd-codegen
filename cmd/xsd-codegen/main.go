package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/helpers"
	"github.com/Patrick-Ivann/xsd-codegen/pkg/parser"
	"github.com/Patrick-Ivann/xsd-codegen/pkg/xmlgen"
	"github.com/beevik/etree"
)

func main() {
	xsdPath := flag.String("xsd", "", "Path to XSD file")
	outPath := flag.String("out", "", "Output XML file path (default stdout)")
	flag.Parse()

	if *xsdPath == "" {
		log.Fatal("XSD file path is required. Use -xsd flag.")
	}

	schema, err := parser.ParseXSD(*xsdPath, nil)
	if err != nil {
		log.Fatalf("Failed to parse XSD: %v", err)
	}

	doc := etree.NewDocument()
	for _, el := range schema.Elements {
		if el.Name == "purchaseOrder" { // entrypoint element
			root := xmlgen.GenerateElement(schema, el, helpers.DefaultValueGenerator{})
			root.CreateAttr("xmlns", schema.TargetNamespace)
			root.CreateAttr("xsi:schemaLocation", schema.TargetNamespace+" schema.xsd")
			root.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
			doc.SetRoot(root)
			break
		}
	}

	doc.Indent(2)
	if *outPath == "" {
		doc.WriteTo(os.Stdout)
	} else {
		file, err := os.Create(*outPath)
		if err != nil {
			log.Fatalf("Unable to create output file: %v", err)
		}
		defer file.Close()
		doc.WriteTo(file)
		fmt.Println("XML generated at", *outPath)
	}
}
