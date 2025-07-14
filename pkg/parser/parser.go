package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

// ParseXSD parses an XSD file and returns a Schema model, handling all restrictions.
// including inline complex/simple types, attributes, and restrictions.
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

	var inAnnotation bool
	var inDocumentation bool
	var docBuffer strings.Builder

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("XML decode error: %w", err)
		}

		switch se := token.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "element":
				field := model.XSDField{}
				elem := model.XSDElement{}
				for _, attr := range se.Attr {
					switch attr.Name.Local {
					case "name":
						field.Name = attr.Value
						elem.Name = attr.Value
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

			case "complexType":
				currType = &model.XSDType{}
				for _, attr := range se.Attr {
					if attr.Name.Local == "name" {
						currType.Name = attr.Value
					}
				}
				if currType.Name == "" && currElem != nil {
					currType.Name = currElem.Name
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
				if currType != nil && currType.Name == "" && currElem != nil {
					currType.Name = currElem.Name
				}

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

			case "annotation":
				inAnnotation = true
			case "documentation":
				if inAnnotation {
					inDocumentation = true
					docBuffer.Reset()
				}
			}

		case xml.CharData:
			if inDocumentation {
				docBuffer.Write([]byte(se))
			}

		case xml.EndElement:
			switch se.Name.Local {

			case "documentation":
				inDocumentation = false
				doc := strings.TrimSpace(docBuffer.String())
				if currField != nil {
					currField.Documentation = doc
				} else if currElem != nil {
					currElem.Documentation = doc
				} else if currType != nil {
					currType.Documentation = doc
				}
				docBuffer.Reset()
			case "annotation":
				inAnnotation = false

			case "complexType", "simpleType":

				if currType != nil {
					currType.Documentation = strings.TrimSpace(docBuffer.String())
					schema.Types = append(schema.Types, *currType)
					currType = nil
					docBuffer.Reset()
				}
				inSimpleType = false
				fmt.Printf("schema: %v\n", schema)
			case "element":
				if currElem != nil {
					currElem.Documentation = strings.TrimSpace(docBuffer.String())
					docBuffer.Reset()
				}
				if currType != nil && currField != nil {
					currType.Fields = append(currType.Fields, *currField)
					currType.Documentation = strings.TrimSpace(docBuffer.String())

				} else if currElem != nil {
					schema.Elements = append(schema.Elements, *currElem)
				}
				currField = nil
				currElem = nil
			}
		}
	}

	fmt.Printf("schema: %v\n", schema)

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
