package parser

import (
	"encoding/xml"
	"os"
	"path/filepath"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

func ParseXSD(filePath string, loadedSchemas map[string]*model.XSDSchema) (*model.XSDSchema, error) {
	absPath, _ := filepath.Abs(filePath)
	if loadedSchemas == nil {
		loadedSchemas = make(map[string]*model.XSDSchema)
	}
	if s, exists := loadedSchemas[absPath]; exists {
		return s, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var schema model.XSDSchema
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
