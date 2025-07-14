package xmlgen

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

// TestMarshalDummyXML checks that XML is produced and contains expected elements.
func TestMarshalDummyXML(t *testing.T) {
	typ := model.XSDType{
		Name: "Person",
		Fields: []model.XSDField{
			{Name: "Name", Type: "string"},
			{Name: "Age", Type: "int"},
		},
		Attributes: []model.XSDAttribute{
			{Name: "id", Type: "string"},
		},
	}
	xmlBytes, err := MarshalDummyXML("Person", typ)
	if err != nil {
		t.Fatalf("MarshalDummyXML failed: %v", err)
	}
	out := string(xmlBytes)
	if !strings.Contains(out, "<Person>") || !strings.Contains(out, "<Name>") {
		t.Errorf("Expected XML root and field tags, got: %s", out)
	}
	// Validate XML structure
	var v interface{}
	if err := xml.Unmarshal(xmlBytes, &v); err != nil {
		t.Errorf("Unmarshal failed: %v", err)
	}
}

// TestGenerateDummyStruct ensures the dummy struct contains expected keys and types.
func TestGenerateDummyStruct(t *testing.T) {
	typ := model.XSDType{
		Name: "TestType",
		Fields: []model.XSDField{
			{Name: "Field1", Type: "string"},
			{Name: "Field2", Type: "int"},
		},
	}
	dummy := GenerateDummyStruct(typ)
	m, ok := dummy.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", dummy)
	}
	if _, ok := m["Field1"]; !ok {
		t.Errorf("Expected key Field1")
	}
	if _, ok := m["Field2"]; !ok {
		t.Errorf("Expected key Field2")
	}
}
