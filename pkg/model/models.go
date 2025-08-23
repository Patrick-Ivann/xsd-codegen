package model

import "encoding/xml"

type XSDSchema struct {
	XMLName         xml.Name         `xml:"schema"`
	TargetNamespace string           `xml:"targetNamespace,attr"`
	ElementForm     string           `xml:"elementFormDefault,attr"`
	Includes        []XSDInclude     `xml:"include"`
	Imports         []XSDImport      `xml:"import"`
	Elements        []XSDElement     `xml:"element"`
	ComplexTypes    []XSDComplexType `xml:"complexType"`
	SimpleTypes     []XSDSimpleType  `xml:"simpleType"`
}

type XSDInclude struct {
	SchemaLocation string `xml:"schemaLocation,attr"`
}

type XSDImport struct {
	SchemaLocation string `xml:"schemaLocation,attr"`
	Namespace      string `xml:"namespace,attr"`
}

type XSDElement struct {
	Name        string          `xml:"name,attr"`
	Type        string          `xml:"type,attr,omitempty"`
	Ref         string          `xml:"ref,attr,omitempty"`
	MinOccurs   string          `xml:"minOccurs,attr,omitempty"`
	MaxOccurs   string          `xml:"maxOccurs,attr,omitempty"`
	ComplexType *XSDComplexType `xml:"complexType"`
	SimpleType  *XSDSimpleType  `xml:"simpleType"`
}

type XSDComplexType struct {
	Name     string         `xml:"name,attr,omitempty"`
	Sequence *XSDSequence   `xml:"sequence"`
	Choice   *XSDChoice     `xml:"choice"`
	Attrs    []XSDAttribute `xml:"attribute"`
}

type XSDSimpleType struct {
	Name        string          `xml:"name,attr,omitempty"`
	Restriction *XSDRestriction `xml:"restriction"`
}

type XSDRestriction struct {
	Base         string      `xml:"base,attr"`
	MinIncl      *XSDValue   `xml:"minInclusive"`
	MaxExcl      *XSDValue   `xml:"maxExclusive"`
	Pattern      *XSDPattern `xml:"pattern"`
	Enumerations []XSDValue  `xml:"enumeration"`
}

type XSDPattern struct {
	Value string `xml:"value,attr"`
}

type XSDValue struct {
	Value string `xml:"value,attr"`
}

type XSDSequence struct {
	Elements []XSDElement `xml:"element"`
}

type XSDChoice struct {
	Elements []XSDElement `xml:"element"`
}

type XSDAttribute struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Use   string `xml:"use,attr,omitempty"`
	Fixed string `xml:"fixed,attr,omitempty"`
}
