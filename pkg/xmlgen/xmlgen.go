package xmlgen

import (
	"strings"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/helpers"
	"github.com/Patrick-Ivann/xsd-codegen/pkg/parser"
	"github.com/beevik/etree"
)

func GenerateElement(schema *parser.XSDSchema, element parser.XSDElement) *etree.Element {
	elem := etree.NewElement(element.Name)

	if element.Type != "" {
		if strings.HasPrefix(element.Type, "tns:") {
			typeName := strings.TrimPrefix(element.Type, "tns:")
			for _, ct := range schema.ComplexTypes {
				if ct.Name == typeName {
					appendComplexContent(schema, elem, ct)
					return elem
				}
			}
			for _, st := range schema.SimpleTypes {
				if st.Name == typeName {
					elem.SetText(helpers.GenerateValue(st.Restriction.Base, st.Restriction))
					return elem
				}
			}
		} else {
			elem.SetText(helpers.GenerateValue(element.Type, nil))
		}
	} else if element.ComplexType != nil {
		appendComplexContent(schema, elem, *element.ComplexType)
	} else if element.SimpleType != nil {
		elem.SetText(helpers.GenerateValue(element.SimpleType.Restriction.Base, element.SimpleType.Restriction))
	} else if element.Ref != "" {
		refName := strings.Split(element.Ref, ":")[1]
		for _, el := range schema.Elements {
			if el.Name == refName {
				refElem := GenerateElement(schema, el)
				elem = refElem
			}
		}
	}

	return elem
}

func appendComplexContent(schema *parser.XSDSchema, elem *etree.Element, ct parser.XSDComplexType) {
	if ct.Sequence != nil {
		for _, child := range ct.Sequence.Elements {
			minOccurs := helpers.ParseOccurs(child.MinOccurs, 1)
			maxOccurs := helpers.ParseOccurs(child.MaxOccurs, 1)
			count := helpers.RandomBetween(minOccurs, maxOccurs)
			for i := 0; i < count; i++ {
				childXML := GenerateElement(schema, child)
				elem.AddChild(childXML)
			}
		}
	}

	if ct.Choice != nil && len(ct.Choice.Elements) > 0 {
		choice := ct.Choice.Elements[helpers.RandomBetween(0, len(ct.Choice.Elements)-1)]
		elem.AddChild(GenerateElement(schema, choice))
	}

	for _, attr := range ct.Attrs {
		val := helpers.GenerateValue(attr.Type, nil)
		if attr.Fixed != "" {
			val = attr.Fixed
		}
		elem.CreateAttr(attr.Name, val)
	}
}
