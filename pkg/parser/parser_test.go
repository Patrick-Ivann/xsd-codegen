package parser

import (
	"os"
	"testing"
)

// TestParseXSD_SimpleRestriction parses a simpleType with enumeration restriction.
func TestParseXSD_SimpleRestriction(t *testing.T) {
	xsd := `<?xml version="1.0"?>
    <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
        <xs:simpleType name="ColorType">
            <xs:restriction base="xs:string">
                <xs:enumeration value="Red"/>
                <xs:enumeration value="Green"/>
                <xs:enumeration value="Blue"/>
            </xs:restriction>
        </xs:simpleType>
        <xs:element name="Color" type="ColorType"/>
    </xs:schema>`
	path := "test_enum.xsd"
	if err := os.WriteFile(path, []byte(xsd), 0644); err != nil {
		t.Fatalf("Failed to write XSD: %v", err)
	}
	defer os.Remove(path)
	schema, err := ParseXSD(path)
	if err != nil {
		t.Fatalf("ParseXSD failed: %v", err)
	}
	foundType := false
	for _, typ := range schema.Types {
		if typ.Name == "ColorType" {
			foundType = true
			// In your parser, attach restriction to the type or its field as needed.
		}
	}
	if !foundType {
		t.Errorf("ColorType not found or restriction not parsed")
	}
}

// TestParseXSD_ComplexTypeWithSequence parses a complexType with a sequence of elements.
func TestParseXSD_ComplexTypeWithSequence(t *testing.T) {
	xsd := `<?xml version="1.0"?>
    <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
      <xs:complexType name="Person">
        <xs:sequence>
          <xs:element name="FirstName" type="xs:string"/>
          <xs:element name="LastName" type="xs:string"/>
        </xs:sequence>
      </xs:complexType>
      <xs:element name="Employee" type="Person"/>
    </xs:schema>`
	path := "test_complex.xsd"
	if err := os.WriteFile(path, []byte(xsd), 0644); err != nil {
		t.Fatalf("Failed to write XSD: %v", err)
	}
	defer os.Remove(path)
	schema, err := ParseXSD(path)
	if err != nil {
		t.Fatalf("ParseXSD failed: %v", err)
	}
	found := false
	for _, typ := range schema.Types {
		if typ.Name == "Person" && len(typ.Fields) == 2 {
			found = true
		}
	}
	if !found {
		t.Errorf("Person type with sequence not parsed correctly")
	}
}

// TestParseXSD_Attribute parses a complexType with an attribute.
func TestParseXSD_Attribute(t *testing.T) {
	xsd := `<?xml version="1.0"?>
    <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
      <xs:complexType name="Book">
        <xs:sequence>
          <xs:element name="Title" type="xs:string"/>
        </xs:sequence>
        <xs:attribute name="id" type="xs:string" use="required"/>
      </xs:complexType>
      <xs:element name="BookElement" type="Book"/>
    </xs:schema>`
	path := "test_attr.xsd"
	if err := os.WriteFile(path, []byte(xsd), 0644); err != nil {
		t.Fatalf("Failed to write XSD: %v", err)
	}
	defer os.Remove(path)
	schema, err := ParseXSD(path)
	if err != nil {
		t.Fatalf("ParseXSD failed: %v", err)
	}
	found := false
	for _, typ := range schema.Types {
		if typ.Name == "Book" && len(typ.Attributes) == 1 && typ.Attributes[0].Name == "id" {
			found = true
		}
	}
	if !found {
		t.Errorf("Book type with attribute not parsed correctly")
	}
}

// TestParseXSD_MinMaxOccurs parses minOccurs and maxOccurs attributes.
func TestParseXSD_MinMaxOccurs(t *testing.T) {
	xsd := `<?xml version="1.0"?>
    <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
      <xs:complexType name="Group">
        <xs:sequence>
          <xs:element name="Member" type="xs:string" minOccurs="1" maxOccurs="5"/>
        </xs:sequence>
      </xs:complexType>
    </xs:schema>`
	path := "test_occurs.xsd"
	if err := os.WriteFile(path, []byte(xsd), 0644); err != nil {
		t.Fatalf("Failed to write XSD: %v", err)
	}
	defer os.Remove(path)
	schema, err := ParseXSD(path)
	if err != nil {
		t.Fatalf("ParseXSD failed: %v", err)
	}
	found := false
	for _, typ := range schema.Types {
		if typ.Name == "Group" && len(typ.Fields) == 1 &&
			typ.Fields[0].MinOccurs == 1 && typ.Fields[0].MaxOccurs == 5 {
			found = true
		}
	}
	if !found {
		t.Errorf("Group type with minOccurs/maxOccurs not parsed correctly")
	}
}

// TestParseXSD_InvalidFile checks error handling for non-existent files.
func TestParseXSD_InvalidFile(t *testing.T) {
	_, err := ParseXSD("nonexistent.xsd")
	if err == nil {
		t.Errorf("Expected error for nonexistent file, got nil")
	}
}

// TestParseXSD_RestrictionFacets parses various restriction facets.
func TestParseXSD_RestrictionFacets(t *testing.T) {
	xsd := `<?xml version="1.0"?>
    <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
      <xs:simpleType name="LimitedString">
        <xs:restriction base="xs:string">
          <xs:minLength value="2"/>
          <xs:maxLength value="10"/>
          <xs:pattern value="[A-Z]+"/>
        </xs:restriction>
      </xs:simpleType>
    </xs:schema>`
	path := "test_facets.xsd"
	if err := os.WriteFile(path, []byte(xsd), 0644); err != nil {
		t.Fatalf("Failed to write XSD: %v", err)
	}
	defer os.Remove(path)
	schema, err := ParseXSD(path)
	if err != nil {
		t.Fatalf("ParseXSD failed: %v", err)
	}
	found := false
	for _, typ := range schema.Types {
		if typ.Name == "LimitedString" && typ.Fields == nil && len(typ.Attributes) == 0 {
			found = true
		}
	}
	if !found {
		t.Errorf("LimitedString type with restriction facets not parsed correctly")
	}
}

func TestParseXSD_ElementInlineSimpleTypeRestriction(t *testing.T) {
	xsd := `<?xml version="1.0"?>
    <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
      <xs:element name="Code">
        <xs:simpleType>
          <xs:restriction base="xs:string">
            <xs:pattern value="[A-Z]{3}"/>
          </xs:restriction>
        </xs:simpleType>
      </xs:element>
    </xs:schema>`
	path := "test_inline_restriction.xsd"
	os.WriteFile(path, []byte(xsd), 0644)
	defer os.Remove(path)
	schema, err := ParseXSD(path)
	if err != nil {
		t.Fatalf("ParseXSD failed: %v", err)
	}
	if len(schema.Elements) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(schema.Elements))
	}
	elem := schema.Elements[0]
	if elem.Name != "Code" {
		t.Fatalf("Expected element 'Code', got '%s'", elem.Name)
	}
	if elem.Restriction == nil || elem.Restriction.Pattern == nil || *elem.Restriction.Pattern != "[A-Z]{3}" {
		t.Fatalf("Restriction or pattern not parsed correctly: %+v", elem.Restriction)
	}

}
