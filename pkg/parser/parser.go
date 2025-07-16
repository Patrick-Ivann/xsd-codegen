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
	var inSimpleType bool

	var inAnnotation bool
	var inDocumentation bool
	var docBuffer strings.Builder

	var elementStack []string

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("XML decode error: %w", err)
		}

		switch tok := token.(type) {
		case xml.StartElement:
			elementStack = append(elementStack, tok.Name.Local)

			switch tok.Name.Local {
			case "element":
				isInsideType := currType != nil && len(elementStack) >= 2 && elementStack[len(elementStack)-2] == "sequence"
				if isInsideType {
					field := model.XSDField{}
					for _, attr := range tok.Attr {
						switch attr.Name.Local {
						case "name":
							field.Name = attr.Value
						case "type":
							field.Type = attr.Value
						case "minOccurs":
							field.MinOccurs, _ = strconv.Atoi(attr.Value)
						case "maxOccurs":
							field.MaxOccurs, _ = strconv.Atoi(attr.Value)
						}
					}
					currField = &field
				} else {
					elem := model.XSDElement{}
					for _, attr := range tok.Attr {
						switch attr.Name.Local {
						case "name":
							elem.Name = attr.Value
						case "type":
							elem.Type = attr.Value
						}
					}
					currElem = &elem
				}

			case "complexType":
				currType = &model.XSDType{}
				for _, attr := range tok.Attr {
					if attr.Name.Local == "name" {
						currType.Name = attr.Value
					}
				}
				if currType.Name == "" && currElem != nil {
					// currType.Name = currElem.Name
					currType.Name = currElem.Name + "Type"
					currType.Documentation = currElem.Documentation
				}

			case "simpleType":
				inSimpleType = true
				if currField != nil || currElem != nil {
					// Could handle inline restriction parsing here
				} else {
					currType = &model.XSDType{}
					for _, attr := range tok.Attr {
						if attr.Name.Local == "name" {
							currType.Name = attr.Value
						}
					}
				}

			case "attribute":
				if currType != nil {
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
					// currType.Attributes = append(currType.Attributes, attr)
				}

			case "restriction":
				if inSimpleType {
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
				docBuffer.Write(tok)
			}

		case xml.EndElement:
			switch tok.Name.Local {
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
				if currAttr != nil {
					currAttr.Documentation = doc
				}
				docBuffer.Reset()

			case "annotation":
				inAnnotation = false

			case "complexType", "simpleType":
				if currType != nil {
					// Ensure inline element-level documentation isn't lost
					if currType.Documentation == "" && currElem != nil && currElem.Documentation != "" {
						currType.Documentation = currElem.Documentation
					}
					schema.Types = append(schema.Types, *currType)
				}
				currType = nil
				inSimpleType = false
			case "element":
				if currField != nil && currType != nil {
					currType.Fields = append(currType.Fields, *currField)
				} else if currElem != nil {
					schema.Elements = append(schema.Elements, *currElem)
				}
				currField = nil
				currElem = nil
				docBuffer.Reset()
			}

			if len(elementStack) > 0 {
				elementStack = elementStack[:len(elementStack)-1]
			}
		}
	}

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
