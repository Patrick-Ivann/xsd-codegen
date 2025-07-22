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
	var currAttr *model.XSDAttribute
	var currRestriction *model.Restriction
	var inDocumentation bool
	var docBuffer strings.Builder

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("XML decoding failed: %w", err)
		}

		switch tok := token.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "include":
				for _, a := range tok.Attr {
					if a.Name.Local == "schemaLocation" {
						schema.Includes = append(schema.Includes, model.Directive{
							SchemaLocation: a.Value,
						})
					}
				}

			case "import":
				for _, a := range tok.Attr {
					if a.Name.Local == "namespace" {
						schema.Imports = append(schema.Imports, model.Directive{
							SchemaLocation: a.Value,
						})
					}
				}

			case "element":
				field := model.XSDField{}
				elem := model.XSDElement{}
				for _, a := range tok.Attr {
					switch a.Name.Local {
					case "name":
						field.Name = a.Value
						elem.Name = a.Value
					case "ref":
						refName := strings.TrimPrefix(a.Value, "tns:")
						field.Name = refName
						field.Type = refName // assuming global element has type = name
					case "type":
						field.Type = a.Value
						elem.Type = a.Value
					case "minOccurs":
						field.MinOccurs = atoi(a.Value)
						elem.MinOccurs = field.MinOccurs
					case "maxOccurs":
						field.MaxOccurs = atoi(a.Value)
						elem.MaxOccurs = field.MaxOccurs
					}
				}
				currField = &field
				currElem = &elem

			case "complexType":
				currType = &model.XSDType{}
				for _, a := range tok.Attr {
					if a.Name.Local == "name" {
						currType.Name = a.Value
					}
				}
				if currType.Name == "" && currElem != nil {
					currType.Name = currElem.Name
					currType.Documentation = currElem.Documentation
				}

			case "attribute":
				attr := model.XSDAttribute{}
				for _, a := range tok.Attr {
					switch a.Name.Local {
					case "name":
						attr.Name = a.Value
					case "type":
						attr.Type = a.Value
					}
				}
				currAttr = &attr

			case "simpleType":
				if currRestriction == nil {
					currRestriction = &model.Restriction{}
				}

			case "restriction":
				if currRestriction == nil {
					currRestriction = &model.Restriction{}
				}
				for _, a := range tok.Attr {
					if a.Name.Local == "base" {
						if currField != nil && currField.Type == "" {
							currField.Type = a.Value
						}
					}
				}

			case "annotation":
				// pass

			case "documentation":
				inDocumentation = true
				docBuffer.Reset()

			// ➕ Restriction facets
			case "enumeration":
				for _, a := range tok.Attr {
					if a.Name.Local == "value" {
						currRestriction.Enumeration = append(currRestriction.Enumeration, a.Value)
					}
				}

			case "length":
				for _, a := range tok.Attr {
					if a.Name.Local == "value" {
						currRestriction.Length = parseIntPtr(a.Value)
					}
				}

			case "minLength":
				for _, a := range tok.Attr {
					if a.Name.Local == "value" {
						currRestriction.MinLength = parseIntPtr(a.Value)
					}
				}

			case "maxLength":
				for _, a := range tok.Attr {
					if a.Name.Local == "value" {
						currRestriction.MaxLength = parseIntPtr(a.Value)
					}
				}

			case "pattern":
				for _, a := range tok.Attr {
					if a.Name.Local == "value" {
						currRestriction.Pattern = parseStrPtr(a.Value)
					}
				}

			case "whiteSpace":
				for _, a := range tok.Attr {
					if a.Name.Local == "value" {
						currRestriction.WhiteSpace = parseStrPtr(a.Value)
					}
				}

			case "minInclusive":
				for _, a := range tok.Attr {
					if a.Name.Local == "value" {
						currRestriction.MinInclusive = parseStrPtr(a.Value)
					}
				}

			case "maxInclusive":
				for _, a := range tok.Attr {
					if a.Name.Local == "value" {
						currRestriction.MaxInclusive = parseStrPtr(a.Value)
					}
				}

			case "minExclusive":
				for _, a := range tok.Attr {
					if a.Name.Local == "value" {
						currRestriction.MinExclusive = parseStrPtr(a.Value)
					}
				}

			case "maxExclusive":
				for _, a := range tok.Attr {
					if a.Name.Local == "value" {
						currRestriction.MaxExclusive = parseStrPtr(a.Value)
					}
				}

			case "totalDigits":
				for _, a := range tok.Attr {
					if a.Name.Local == "value" {
						currRestriction.TotalDigits = parseIntPtr(a.Value)
					}
				}

			case "fractionDigits":
				for _, a := range tok.Attr {
					if a.Name.Local == "value" {
						currRestriction.FractionDigits = parseIntPtr(a.Value)
					}
				}
			}

		case xml.CharData:
			if inDocumentation {
				docBuffer.Write([]byte(tok))
			}

		case xml.EndElement:
			switch tok.Name.Local {
			case "documentation":
				inDocumentation = false
				doc := strings.TrimSpace(docBuffer.String())
				if doc != "" {
					if currField != nil {
						currField.Documentation = doc
						fmt.Printf("Assigned documentation to field: %s\n", doc)
					} else if currElem != nil {
						currElem.Documentation = doc
						fmt.Printf("Assigned documentation to element: %s\n", doc)
					} else if currType != nil {
						currType.Documentation = doc
						fmt.Printf("Assigned documentation to type: %s\n", doc)
					} else {
						schema.Documentation = doc
						fmt.Printf("Assigned documentation to schema: %s\n", doc)
					}
				}
				docBuffer.Reset()

			case "restriction":
				if currField != nil {
					currField.Restriction = currRestriction
				}
				if currElem != nil {
					currElem.Restriction = currRestriction
				}
				currRestriction = nil

			case "element":
				if currType != nil && currField != nil {
					currType.Fields = append(currType.Fields, *currField)
				} else if currElem != nil {
					schema.Elements = append(schema.Elements, *currElem)
				}
				currField = nil
				currElem = nil

			case "complexType":
				if currType != nil {
					fmt.Printf("currType: %v\n", currType)
					fmt.Printf("currField: %v\n", currField)
					fmt.Printf("currElem: %v\n", currElem)
					if currElem != nil {
						fmt.Printf("currElem.Documentation: %v\n", currElem.Documentation)
					}
					schema.Types = append(schema.Types, *currType)
					currType = nil
				}

			case "attribute":
				if currType != nil && currAttr != nil {
					currType.Attributes = append(currType.Attributes, *currAttr)
				}
				currAttr = nil
			}
		}
	}

	return &schema, nil
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func parseIntPtr(s string) *int {
	if s == "" {
		return nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &i
}

func parseStrPtr(s string) *string {
	if s == "" {
		return nil
	}
	str := strings.TrimSpace(s)
	return &str
}
