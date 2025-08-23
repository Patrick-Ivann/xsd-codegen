package parser

import (
	"path/filepath"
	"testing"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

func TestParseXSD(t *testing.T) {
	schemaPath := filepath.Join("testdata", "sample.xsd")
	schema, err := ParseXSD(schemaPath, nil)
	if err != nil {
		t.Fatalf("Failed to parse XSD: %v", err)
	}
	if len(schema.Elements) == 0 {
		t.Error("Expected schema elements, got none")
	}
}

func TestParseXSD_WithInclude(t *testing.T) {
	schemaPath := filepath.Join("testdata", "with_include.xsd")
	schema, err := ParseXSD(schemaPath, nil)
	if err != nil {
		t.Fatalf("Failed to parse XSD with include: %v", err)
	}
	if len(schema.ComplexTypes) == 0 {
		t.Error("Expected complex types from included schema")
	}
}

func TestParseXSD_WithImport(t *testing.T) {
	schemaPath := filepath.Join("testdata", "with_import.xsd")
	schema, err := ParseXSD(schemaPath, nil)
	if err != nil {
		t.Fatalf("Failed to parse XSD with import: %v", err)
	}
	if len(schema.SimpleTypes) == 0 {
		t.Error("Expected simple types from imported schema")
	}
}

func TestParseXSD_CyclicReference(t *testing.T) {
	schemaPath := filepath.Join("testdata", "cyclic_a.xsd")
	schema, err := ParseXSD(schemaPath, nil)
	if err != nil {
		t.Fatalf("Failed to parse cyclic schema: %v", err)
	}
	// Should not crash or loop infinitely
	if len(schema.Elements) == 0 {
		t.Error("Expected elements despite cyclic reference")
	}
}

func TestParseXSD_InvalidPath(t *testing.T) {
	schemaPath := filepath.Join("testdata", "nonexistent.xsd")
	_, err := ParseXSD(schemaPath, nil)
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestParseXSD_MalformedXML(t *testing.T) {
	schemaPath := filepath.Join("testdata", "malformed.xsd")
	_, err := ParseXSD(schemaPath, nil)
	if err == nil {
		t.Error("Expected error for malformed XML")
	}
}

func TestParseXSD_AlreadyLoaded(t *testing.T) {
	schemaPath := filepath.Join("testdata", "sample.xsd")
	loaded := make(map[string]*model.XSDSchema)
	schema1, err := ParseXSD(schemaPath, loaded)
	if err != nil {
		t.Fatalf("First parse failed: %v", err)
	}
	schema2, err := ParseXSD(schemaPath, loaded)
	if err != nil {
		t.Fatalf("Second parse failed: %v", err)
	}
	if schema1 != schema2 {
		t.Error("Expected cached schema to be reused")
	}
}
