package xmlgen

import (
	"strings"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/helpers"
	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
	"github.com/beevik/etree"
)

// GenerateElement creates an XML element (etree.Element) based on the provided XSD schema definition.
// It evaluates whether the element is of a simple type, complex type, reference, or inline definition.
// gen is used as a value generator for populating element content.
func GenerateElement(schema *model.XSDSchema, element model.XSDElement, gen helpers.ValueGenerator) *etree.Element {
	// Create a new XML element with the given element name from the schema
	elem := etree.NewElement(element.Name)

	// Handle explicitly defined types
	if element.Type != "" {
		// Case: type refers to a schema-defined type (tns: prefix indicates "this namespace")
		if strings.HasPrefix(element.Type, "tns:") {
			typeName := strings.TrimPrefix(element.Type, "tns:")

			// Look for a complex type with that name in the schema
			for _, ct := range schema.ComplexTypes {
				if ct.Name == typeName {
					// Append content for the complex type
					appendComplexContent(schema, elem, ct, gen)
					return elem
				}
			}

			// Look for a simple type with that name in the schema
			for _, st := range schema.SimpleTypes {
				if st.Name == typeName {
					// Generate a value for this type restriction (like enumeration, pattern, etc.)
					elem.SetText(helpers.GenerateValue(st.Restriction.Base, st.Restriction))
					return elem
				}
			}
		} else {
			// Case: type is a built-in primitive type (e.g., xs:string, xs:int, etc.)
			elem.SetText(helpers.GenerateValue(element.Type, nil))
		}

	} else if element.ComplexType != nil {
		// Inline-defined complex type
		appendComplexContent(schema, elem, *element.ComplexType, gen)

	} else if element.SimpleType != nil {
		// Inline-defined simple type
		elem.SetText(helpers.GenerateValue(element.SimpleType.Restriction.Base, element.SimpleType.Restriction))

	} else if element.Ref != "" {
		// Case: this element is a reference to another schema element
		refName := strings.Split(element.Ref, ":")[1] // Extract referenced element name
		for _, el := range schema.Elements {
			if el.Name == refName {
				// Generate the referenced element recursively
				refElem := GenerateElement(schema, el, gen)
				elem = refElem
			}
		}
	}

	return elem
}

// appendComplexContent populates a complexType into the target XML element.
// It handles sequences, choices, and attributes as defined in the XSD.
func appendComplexContent(schema *model.XSDSchema, elem *etree.Element, ct model.XSDComplexType, gen helpers.ValueGenerator) {
	// Handle <xs:sequence> — ordered elements
	if ct.Sequence != nil {
		for _, child := range ct.Sequence.Elements {
			// Parse occurrence constraints (default to 1 if not specified)
			minOccurs := helpers.ParseOccurs(child.MinOccurs, 1)
			maxOccurs := helpers.ParseOccurs(child.MaxOccurs, 1)

			// Randomly pick how many times to repeat this element (within min/max)
			count := helpers.RandomBetween(minOccurs, maxOccurs)
			for i := 0; i < count; i++ {
				// Recursively generate child elements
				childXML := GenerateElement(schema, child, gen)
				elem.AddChild(childXML)
			}
		}
	}

	// Handle <xs:choice> — only one of the listed elements should be chosen
	if ct.Choice != nil && len(ct.Choice.Elements) > 0 {
		// Randomly pick one child element from the choice
		choice := ct.Choice.Elements[helpers.RandomBetween(0, len(ct.Choice.Elements)-1)]
		elem.AddChild(GenerateElement(schema, choice, gen))
	}

	// Handle attributes defined in the complex type
	for _, attr := range ct.Attrs {
		// Generate a value according to type (unless a 'fixed' value is provided)
		val := helpers.GenerateValue(attr.Type, nil)
		if attr.Fixed != "" {
			val = attr.Fixed
		}
		// Add the attribute to the element
		elem.CreateAttr(attr.Name, val)
	}
}
