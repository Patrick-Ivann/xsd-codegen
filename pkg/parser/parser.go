package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

// ParseXSD orchestrates the loading, parsing, and recursive inclusion/import handling
func ParseXSD(filePath string, loadedSchemas map[string]*model.XSDSchema) (*model.XSDSchema, error) {
	absPath, _ := filepath.Abs(filePath)
	if loadedSchemas == nil {
		loadedSchemas = make(map[string]*model.XSDSchema)
	}
	// If schema is already loaded, return it to avoid reprocessing (handles include/import cycles!)
	if s, exists := loadedSchemas[absPath]; exists {
		return s, nil
	}

	// Read and unmarshal schema XML
	schema, err := readAndUnmarshalSchema(filePath)
	if err != nil {
		return nil, err
	}
	loadedSchemas[absPath] = schema

	// Handle <xs:include> elements
	if err := processIncludes(schema, loadedSchemas, filePath); err != nil {
		return nil, err
	}

	// Handle <xs:import> elements
	if err := processImports(schema, loadedSchemas, filePath); err != nil {
		return nil, err
	}
	return schema, nil
}

// readAndUnmarshalSchema securely reads the XSD file from a safe location and unmarshals its XML.
func readAndUnmarshalSchema(filePath string) (*model.XSDSchema, error) {
	cleanedPath, err := sanitizeAndVerifyPath(filePath)
	if err != nil {
		return nil, fmt.Errorf("unsafe file path: %w", err)
	}
	//nolint:gosec // potential file inclusion expected as this is made for a CLI tool and avoiding new flag on command usage
	data, err := os.ReadFile(cleanedPath)
	if err != nil {
		return nil, err
	}

	var schema model.XSDSchema
	if err := xml.Unmarshal(data, &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}

// sanitizeAndVerifyPath cleans the given path and ensures it resides inside the trusted base directory.
func sanitizeAndVerifyPath(path string) (string, error) {
	// Clean path for lexical normalization
	cleaned := filepath.Clean(path)

	// Get absolute path to handle relative inputs robustly
	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// processIncludes recursively loads and merges schemas from <xs:include> elements
func processIncludes(schema *model.XSDSchema, loadedSchemas map[string]*model.XSDSchema, filePath string) error {
	dir := filepath.Dir(filePath)
	for _, inc := range schema.Includes {
		incPath := filepath.Join(dir, inc.SchemaLocation)
		incSchema, err := ParseXSD(incPath, loadedSchemas)
		if err != nil {
			return err
		}
		// Merge elements and types from included schema into current schema
		schema.Elements = append(schema.Elements, incSchema.Elements...)
		schema.ComplexTypes = append(schema.ComplexTypes, incSchema.ComplexTypes...)
		schema.SimpleTypes = append(schema.SimpleTypes, incSchema.SimpleTypes...)
	}
	return nil
}

// processImports recursively loads and merges schemas from <xs:import> elements
func processImports(schema *model.XSDSchema, loadedSchemas map[string]*model.XSDSchema, filePath string) error {
	dir := filepath.Dir(filePath)
	for _, imp := range schema.Imports {
		impPath := filepath.Join(dir, imp.SchemaLocation)
		impSchema, err := ParseXSD(impPath, loadedSchemas)
		if err != nil {
			return err
		}
		// Merge elements and types from imported schema into current schema
		schema.Elements = append(schema.Elements, impSchema.Elements...)
		schema.ComplexTypes = append(schema.ComplexTypes, impSchema.ComplexTypes...)
		schema.SimpleTypes = append(schema.SimpleTypes, impSchema.SimpleTypes...)
	}
	return nil
}
