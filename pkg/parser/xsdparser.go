package parser

import (
	"encoding/xml"
	"os"
	"path/filepath"
)

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
	Base    string      `xml:"base,attr"`
	MinIncl *XSDValue   `xml:"minInclusive"`
	MaxExcl *XSDValue   `xml:"maxExclusive"`
	Pattern *XSDPattern `xml:"pattern"`
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

func ParseXSD(filePath string, loadedSchemas map[string]*XSDSchema) (*XSDSchema, error) {
	absPath, _ := filepath.Abs(filePath)
	if loadedSchemas == nil {
		loadedSchemas = make(map[string]*XSDSchema)
	}
	if s, exists := loadedSchemas[absPath]; exists {
		return s, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var schema XSDSchema
	if err := xml.Unmarshal(data, &schema); err != nil {
		return nil, err
	}
	loadedSchemas[absPath] = &schema
	dir := filepath.Dir(filePath)

	for _, inc := range schema.Includes {
		incPath := filepath.Join(dir, inc.SchemaLocation)
		incSchema, err := ParseXSD(incPath, loadedSchemas)
		if err != nil {
			return nil, err
		}
		schema.Elements = append(schema.Elements, incSchema.Elements...)
		schema.ComplexTypes = append(schema.ComplexTypes, incSchema.ComplexTypes...)
		schema.SimpleTypes = append(schema.SimpleTypes, incSchema.SimpleTypes...)
	}

	for _, imp := range schema.Imports {
		impPath := filepath.Join(dir, imp.SchemaLocation)
		impSchema, err := ParseXSD(impPath, loadedSchemas)
		if err != nil {
			return nil, err
		}
		schema.Elements = append(schema.Elements, impSchema.Elements...)
		schema.ComplexTypes = append(schema.ComplexTypes, impSchema.ComplexTypes...)
		schema.SimpleTypes = append(schema.SimpleTypes, impSchema.SimpleTypes...)
	}

	return &schema, nil
}
