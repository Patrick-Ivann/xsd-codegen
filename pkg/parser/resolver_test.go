package parser

import (
	"os"
	"testing"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

// TestResolveImportsAndIncludes verifies that imported and included schemas are merged correctly.
func TestResolveImportsAndIncludes(t *testing.T) {
	// Create a base schema file that imports another schema.
	baseXSD := `<?xml version="1.0"?>
    <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
        <xs:import schemaLocation="imported.xsd"/>
        <xs:element name="Root" type="xs:string"/>
    </xs:schema>`
	importedXSD := `<?xml version="1.0"?>
    <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
        <xs:element name="ImportedElem" type="xs:string"/>
        <xs:complexType name="ImportedType">
            <xs:sequence>
                <xs:element name="Field" type="xs:string"/>
            </xs:sequence>
        </xs:complexType>
    </xs:schema>`

	basePath := "base.xsd"
	importedPath := "imported.xsd"
	if err := os.WriteFile(basePath, []byte(baseXSD), 0644); err != nil {
		t.Fatalf("Failed to write base XSD: %v", err)
	}
	if err := os.WriteFile(importedPath, []byte(importedXSD), 0644); err != nil {
		t.Fatalf("Failed to write imported XSD: %v", err)
	}
	defer os.Remove(basePath)
	defer os.Remove(importedPath)

	schema, err := ParseXSD(basePath)
	if err != nil {
		t.Fatalf("ParseXSD failed: %v", err)
	}
	err = ResolveImportsAndIncludes(schema, basePath)
	if err != nil {
		t.Fatalf("ResolveImportsAndIncludes failed: %v", err)
	}
	// Should contain both 'Root' and 'ImportedElem' elements, and 'ImportedType' type.
	foundRoot := false
	foundImportedElem := false
	foundImportedType := false
	for _, e := range schema.Elements {
		if e.Name == "Root" {
			foundRoot = true
		}
		if e.Name == "ImportedElem" {
			foundImportedElem = true
		}
	}
	for _, tpe := range schema.Types {
		if tpe.Name == "ImportedType" {
			foundImportedType = true
		}
	}
	if !foundRoot || !foundImportedElem || !foundImportedType {
		t.Errorf("Schema did not merge imported elements/types correctly: Root=%v ImportedElem=%v ImportedType=%v",
			foundRoot, foundImportedElem, foundImportedType)
	}
}

// TestResolverAvoidsDuplicates ensures duplicate includes/imports are not merged multiple times.
func TestResolverAvoidsDuplicates(t *testing.T) {
	// Two schemas, both include the same third schema.
	mainXSD := `<?xml version="1.0"?>
    <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
        <xs:include schemaLocation="shared.xsd"/>
        <xs:import schemaLocation="other.xsd"/>
        <xs:element name="MainElem" type="xs:string"/>
    </xs:schema>`
	otherXSD := `<?xml version="1.0"?>
    <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
        <xs:include schemaLocation="shared.xsd"/>
        <xs:element name="OtherElem" type="xs:string"/>
    </xs:schema>`
	sharedXSD := `<?xml version="1.0"?>
    <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
        <xs:element name="SharedElem" type="xs:string"/>
    </xs:schema>`

	files := map[string]string{
		"main.xsd":   mainXSD,
		"other.xsd":  otherXSD,
		"shared.xsd": sharedXSD,
	}
	for name, content := range files {
		if err := os.WriteFile(name, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", name, err)
		}
		defer os.Remove(name)
	}

	schema, err := ParseXSD("main.xsd")
	if err != nil {
		t.Fatalf("ParseXSD failed: %v", err)
	}
	err = ResolveImportsAndIncludes(schema, "main.xsd")
	if err != nil {
		t.Fatalf("ResolveImportsAndIncludes failed: %v", err)
	}
	// Should contain MainElem, OtherElem, SharedElem only once each.
	count := make(map[string]int)
	for _, e := range schema.Elements {
		count[e.Name]++
	}
	if count["SharedElem"] != 1 {
		t.Errorf("SharedElem should appear once, got %d", count["SharedElem"])
	}
	if count["OtherElem"] != 1 {
		t.Errorf("OtherElem should appear once, got %d", count["OtherElem"])
	}
	if count["MainElem"] != 1 {
		t.Errorf("MainElem should appear once, got %d", count["MainElem"])
	}
}

func TestResolveImportsAndIncludes_InvalidAbsPath(t *testing.T) {
	schema := &model.Schema{}
	invalidPath := string([]byte{0}) // Invalid path on most OS
	cache := make(schemaCacheType)
	err := resolveImportsAndIncludes(schema, invalidPath, cache)
	if err == nil {
		t.Errorf("Expected error for invalid abs path")
	}
}
