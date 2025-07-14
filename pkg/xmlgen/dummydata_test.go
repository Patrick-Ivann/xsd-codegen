package xmlgen

import (
	"testing"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

// TestGenerateDummyValue_Enumeration checks that the first enumeration value is returned.
func TestGenerateDummyValue_Enumeration(t *testing.T) {
	f := model.XSDField{
		Name: "Color",
		Type: "string",
		Restriction: &model.Restriction{
			Enumeration: []string{"Red", "Green", "Blue"},
		},
	}
	val := GenerateDummyValue(f)
	if val != "Red" {
		t.Errorf("expected 'Red', got '%s'", val)
	}
}

// TestGenerateDummyValue_Length checks that the generated value matches the specified length.
func TestGenerateDummyValue_Length(t *testing.T) {
	l := 6
	f := model.XSDField{
		Name: "Code",
		Type: "string",
		Restriction: &model.Restriction{
			Length: &l,
		},
	}
	val := GenerateDummyValue(f)
	if len(val) != 6 {
		t.Errorf("expected length 6, got %d", len(val))
	}
}

// TestGenerateDummyValue_Numeric checks that a numeric value is generated within the correct range.
func TestGenerateDummyValue_Numeric(t *testing.T) {
	min := "10"
	max := "12"
	f := model.XSDField{
		Name: "Age",
		Type: "int",
		Restriction: &model.Restriction{
			MinInclusive: &min,
			MaxInclusive: &max,
		},
	}
	val := GenerateDummyValue(f)
	// Should be "10", "11", or "12"
	if val != "10" && val != "11" && val != "12" {
		t.Errorf("expected value between 10 and 12, got %s", val)
	}
}

// TestGenerateDummyValue_Default checks that a default value is returned for unknown types.
func TestGenerateDummyValue_Default(t *testing.T) {
	f := model.XSDField{
		Name: "Something",
		Type: "unknownType",
	}
	val := GenerateDummyValue(f)
	if val != "value" {
		t.Errorf("expected 'value', got '%s'", val)
	}
}
