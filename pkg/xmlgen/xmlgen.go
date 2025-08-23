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
func GenerateElement(schema *model.XSDSchema, element *model.XSDElement, gen helpers.ValueGenerator) *etree.Element {
	elem := etree.NewElement(element.Name)

	switch {
	case element.Type != "":
		handleType(schema, elem, element, gen)
	case element.ComplexType != nil:
		appendComplexContent(schema, elem, *element.ComplexType, gen)
	case element.SimpleType != nil:
		elem.SetText(helpers.GenerateValue(element.SimpleType.Restriction.Base, element.SimpleType.Restriction))
	case element.Ref != "":
		elem = handleRef(schema, element.Ref, gen)
	}

	return elem
}

func handleType(schema *model.XSDSchema, elem *etree.Element, element *model.XSDElement, gen helpers.ValueGenerator) {
	if strings.HasPrefix(element.Type, "tns:") {
		typeName := strings.TrimPrefix(element.Type, "tns:")
		if tryAppendComplexType(schema, elem, typeName, gen) {
			return
		}
		if trySetSimpleType(schema, elem, typeName) {
			return
		}
	} else {
		elem.SetText(helpers.GenerateValue(element.Type, nil))
	}
}

func tryAppendComplexType(schema *model.XSDSchema, elem *etree.Element, typeName string, gen helpers.ValueGenerator) bool {
	for _, ct := range schema.ComplexTypes {
		if ct.Name == typeName {
			appendComplexContent(schema, elem, ct, gen)
			return true
		}
	}
	return false
}

func trySetSimpleType(schema *model.XSDSchema, elem *etree.Element, typeName string) bool {
	for _, st := range schema.SimpleTypes {
		if st.Name == typeName {
			elem.SetText(helpers.GenerateValue(st.Restriction.Base, st.Restriction))
			return true
		}
	}
	return false
}

func handleRef(schema *model.XSDSchema, ref string, gen helpers.ValueGenerator) *etree.Element {
	refName := strings.Split(ref, ":")[1]
	for _, el := range schema.Elements {
		if el.Name == refName {
			return GenerateElement(schema, &el, gen)
		}
	}
	return nil
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
				childXML := GenerateElement(schema, &child, gen)
				elem.AddChild(childXML)
			}
		}
	}

	// Handle <xs:choice> — only one of the listed elements should be chosen
	if ct.Choice != nil && len(ct.Choice.Elements) > 0 {
		// Randomly pick one child element from the choice
		choice := ct.Choice.Elements[helpers.RandomBetween(0, len(ct.Choice.Elements)-1)]
		elem.AddChild(GenerateElement(schema, &choice, gen))
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
