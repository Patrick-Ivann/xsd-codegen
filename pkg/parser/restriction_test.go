package parser

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

// helper to parse a restriction XML fragment using parseRestriction
func parseRestrictionFromString(t *testing.T, xmlFragment string) *model.Restriction {
	decoder := xml.NewDecoder(strings.NewReader(xmlFragment))
	// Advance to <restriction> start
	for {
		token, err := decoder.Token()
		if err != nil {
			t.Fatalf("Failed to parse XML: %v", err)
		}
		if se, ok := token.(xml.StartElement); ok && se.Name.Local == "restriction" {
			return parseRestriction(decoder)
		}
	}
}

func TestParseRestriction_Enumeration(t *testing.T) {
	xmlFragment := `
    <restriction base="xs:string">
        <enumeration value="A"/>
        <enumeration value="B"/>
        <enumeration value="C"/>
    </restriction>`
	r := parseRestrictionFromString(t, xmlFragment)
	if len(r.Enumeration) != 3 || r.Enumeration[0] != "A" || r.Enumeration[2] != "C" {
		t.Errorf("Enumeration not parsed correctly: %+v", r.Enumeration)
	}
}

func TestParseRestriction_Pattern(t *testing.T) {
	xmlFragment := `
    <restriction base="xs:string">
        <pattern value="[A-Z]{3}"/>
    </restriction>`
	r := parseRestrictionFromString(t, xmlFragment)
	if r.Pattern == nil || *r.Pattern != "[A-Z]{3}" {
		t.Errorf("Pattern not parsed: %+v", r.Pattern)
	}
}

func TestParseRestriction_Lengths(t *testing.T) {
	xmlFragment := `
    <restriction base="xs:string">
        <length value="5"/>
        <minLength value="2"/>
        <maxLength value="10"/>
    </restriction>`
	r := parseRestrictionFromString(t, xmlFragment)
	if r.Length == nil || *r.Length != 5 {
		t.Errorf("Length not parsed: %+v", r.Length)
	}
	if r.MinLength == nil || *r.MinLength != 2 {
		t.Errorf("MinLength not parsed: %+v", r.MinLength)
	}
	if r.MaxLength == nil || *r.MaxLength != 10 {
		t.Errorf("MaxLength not parsed: %+v", r.MaxLength)
	}
}

func TestParseRestriction_InclusiveExclusive(t *testing.T) {
	xmlFragment := `
    <restriction base="xs:int">
        <minInclusive value="1"/>
        <maxInclusive value="10"/>
        <minExclusive value="2"/>
        <maxExclusive value="9"/>
    </restriction>`
	r := parseRestrictionFromString(t, xmlFragment)
	if r.MinInclusive == nil || *r.MinInclusive != "1" {
		t.Errorf("MinInclusive not parsed: %+v", r.MinInclusive)
	}
	if r.MaxInclusive == nil || *r.MaxInclusive != "10" {
		t.Errorf("MaxInclusive not parsed: %+v", r.MaxInclusive)
	}
	if r.MinExclusive == nil || *r.MinExclusive != "2" {
		t.Errorf("MinExclusive not parsed: %+v", r.MinExclusive)
	}
	if r.MaxExclusive == nil || *r.MaxExclusive != "9" {
		t.Errorf("MaxExclusive not parsed: %+v", r.MaxExclusive)
	}
}

func TestParseRestriction_DigitsAndWhiteSpace(t *testing.T) {
	xmlFragment := `
    <restriction base="xs:decimal">
        <fractionDigits value="3"/>
        <totalDigits value="7"/>
        <whiteSpace value="collapse"/>
    </restriction>`
	r := parseRestrictionFromString(t, xmlFragment)
	if r.FractionDigits == nil || *r.FractionDigits != 3 {
		t.Errorf("FractionDigits not parsed: %+v", r.FractionDigits)
	}
	if r.TotalDigits == nil || *r.TotalDigits != 7 {
		t.Errorf("TotalDigits not parsed: %+v", r.TotalDigits)
	}
	if r.WhiteSpace == nil || *r.WhiteSpace != "collapse" {
		t.Errorf("WhiteSpace not parsed: %+v", r.WhiteSpace)
	}
}
