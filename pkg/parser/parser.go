package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

// ParseXSD parses an XSD file and returns a Schema model, handling all restrictions.
func ParseXSD(path string) (*model.Schema, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open XSD: %w", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	var schema model.Schema

	var currType *model.XSDType
	var currField *model.XSDField
	var currElem *model.XSDElement
	var inSimpleType bool

	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("XML decode error: %w", err)
		}
		switch se := t.(type) {
		case xml.StartElement:
			fmt.Printf("se: %v\n", se)
			switch se.Name.Local {
			case "complexType":
				currType = &model.XSDType{}
				for _, attr := range se.Attr {
					if attr.Name.Local == "name" {
						currType.Name = attr.Value
					}
				}
			case "simpleType":
				if currField != nil || currElem != nil {
					inSimpleType = true
				} else {
					currType = &model.XSDType{}
					for _, attr := range se.Attr {
						if attr.Name.Local == "name" {
							currType.Name = attr.Value
						}
					}
				}
			case "element":
				field := model.XSDField{}
				elem := model.XSDElement{}
				for _, attr := range se.Attr {
					switch attr.Name.Local {
					case "name":
						field.Name = attr.Value
						elem.Name = attr.Value
						fmt.Printf("parser.go field.Name: %v\n", field.Name)
						fmt.Printf("parser.go elem.Name: %v\n", elem.Name)
					case "type":
						field.Type = attr.Value
						elem.Type = attr.Value
					case "minOccurs":
						field.MinOccurs, _ = strconv.Atoi(attr.Value)
					case "maxOccurs":
						field.MaxOccurs, _ = strconv.Atoi(attr.Value)
					}
				}
				currField = &field
				currElem = &elem
				fmt.Printf("parser.go currField: %v\n", currField)
				fmt.Printf("parser.go currElem: %v\n", currElem)
			case "attribute":
				if currType != nil {
					attr := model.XSDAttribute{}
					for _, a := range se.Attr {
						switch a.Name.Local {
						case "name":
							attr.Name = a.Value
						case "type":
							attr.Type = a.Value
						}
					}
					currType.Attributes = append(currType.Attributes, attr)
				}
			case "restriction":
				if inSimpleType && (currField != nil || currElem != nil) {
					restriction := parseRestriction(decoder)
					if currField != nil {
						currField.Restriction = restriction
					}
					if currElem != nil {
						currElem.Restriction = restriction
					}
				}
			}
			fmt.Printf("parser.go out switch schema: %v\n", schema)
			fmt.Printf("parser.go se.Name: %v\n", se.Name)
			fmt.Printf("parser.go se: %v\n", se)
			fmt.Printf("parser.go startElement currElem: %v\n", currElem)

		case xml.EndElement:
			fmt.Printf("se.Name.Local: %v\n", se.Name.Local)
			switch se.Name.Local {
			case "complexType", "simpleType":
				if currType != nil {
					schema.Types = append(schema.Types, *currType)
					currType = nil
				}
				inSimpleType = false
			case "element":
				fmt.Printf("currType: %v\n", currType)
				fmt.Printf("currType: %v\n", currField)
				fmt.Printf("currType: %v\n", currElem)
				if currType != nil && currField != nil {
					currType.Fields = append(currType.Fields, *currField)
				} else if currElem != nil {
					schema.Elements = append(schema.Elements, *currElem)
				}
				currField = nil
				currElem = nil
			}

		}
	}
	fmt.Printf("parser.go currElem: %v\n", currElem)
	fmt.Printf("parser.go schema: %v\n", schema)
	fmt.Printf("Final schema.Elements: %+v\n", schema.Elements)
	return &schema, nil
}

// parseRestriction parses an <xs:restriction> element and its children.
func parseRestriction(decoder *xml.Decoder) *model.Restriction {
	r := &model.Restriction{}
	depth := 1
	for depth > 0 {
		t, err := decoder.Token()
		if err != nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "enumeration":
				for _, attr := range se.Attr {
					if attr.Name.Local == "value" {
						r.Enumeration = append(r.Enumeration, attr.Value)
					}
				}
			case "pattern":
				for _, attr := range se.Attr {
					if attr.Name.Local == "value" {
						v := attr.Value
						r.Pattern = &v
					}
				}
			case "minLength":
				for _, attr := range se.Attr {
					if attr.Name.Local == "value" {
						v, _ := strconv.Atoi(attr.Value)
						r.MinLength = &v
					}
				}
			case "maxLength":
				for _, attr := range se.Attr {
					if attr.Name.Local == "value" {
						v, _ := strconv.Atoi(attr.Value)
						r.MaxLength = &v
					}
				}
			case "length":
				for _, attr := range se.Attr {
					if attr.Name.Local == "value" {
						v, _ := strconv.Atoi(attr.Value)
						r.Length = &v
					}
				}
			case "minInclusive":
				for _, attr := range se.Attr {
					if attr.Name.Local == "value" {
						v := attr.Value
						r.MinInclusive = &v
					}
				}
			case "maxInclusive":
				for _, attr := range se.Attr {
					if attr.Name.Local == "value" {
						v := attr.Value
						r.MaxInclusive = &v
					}
				}
			case "minExclusive":
				for _, attr := range se.Attr {
					if attr.Name.Local == "value" {
						v := attr.Value
						r.MinExclusive = &v
					}
				}
			case "maxExclusive":
				for _, attr := range se.Attr {
					if attr.Name.Local == "value" {
						v := attr.Value
						r.MaxExclusive = &v
					}
				}
			case "fractionDigits":
				for _, attr := range se.Attr {
					if attr.Name.Local == "value" {
						v, _ := strconv.Atoi(attr.Value)
						r.FractionDigits = &v
					}
				}
			case "totalDigits":
				for _, attr := range se.Attr {
					if attr.Name.Local == "value" {
						v, _ := strconv.Atoi(attr.Value)
						r.TotalDigits = &v
					}
				}
			case "whiteSpace":
				for _, attr := range se.Attr {
					if attr.Name.Local == "value" {
						v := attr.Value
						r.WhiteSpace = &v
					}
				}
			}
			depth++
		case xml.EndElement:
			depth--
			if se.Name.Local == "restriction" {
				return r
			}
		}
	}
	return r
}
