package parser

import (
	"path/filepath"
	"testing"
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
