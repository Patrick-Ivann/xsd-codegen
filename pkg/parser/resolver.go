package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

// schemaCacheType is a simple in-memory cache for resolved schemas.
type schemaCacheType map[string]*model.Schema

// ResolveImportsAndIncludes processes <import> and <include> directives in an XSD file.
// It recursively loads referenced schemas and merges their types/elements into the main schema.
// This version uses a local cache to avoid global state and deadlocks.
func ResolveImportsAndIncludes(schema *model.Schema, basePath string) error {
	cache := make(schemaCacheType)
	return resolveImportsAndIncludes(schema, basePath, cache)
}

// resolveImportsAndIncludes is the internal recursive resolver using a local cache.
func resolveImportsAndIncludes(schema *model.Schema, basePath string, cache schemaCacheType) error {
	// Open the XSD file for parsing.
	file, err := os.Open(basePath)
	if err != nil {
		return fmt.Errorf("failed to open XSD file: %w", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to decode XML: %w", err)
		}
		switch se := token.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "import", "include":
				var schemaLocation string
				for _, attr := range se.Attr {
					if attr.Name.Local == "schemaLocation" {
						schemaLocation = attr.Value
						break
					}
				}
				if schemaLocation == "" {
					continue
				}
				// Resolve the referenced schema path.
				refPath := schemaLocation
				if !filepath.IsAbs(schemaLocation) {
					refPath = filepath.Join(filepath.Dir(basePath), schemaLocation)
				}
				absRefPath, err := filepath.Abs(refPath)
				if err != nil {
					return fmt.Errorf("failed to resolve schema path: %w", err)
				}
				// Avoid duplicate parsing.
				if cachedSchema, ok := cache[absRefPath]; ok {
					mergeSchemas(schema, cachedSchema)
					continue
				}
				// Recursively parse the referenced schema.
				refSchema, err := ParseXSD(absRefPath)
				if err != nil {
					return fmt.Errorf("failed to parse imported/included schema: %w", err)
				}
				// Recursively resolve imports/includes in the referenced schema.
				if err := resolveImportsAndIncludes(refSchema, absRefPath, cache); err != nil {
					return fmt.Errorf("failed to resolve imports/includes for %s: %w", absRefPath, err)
				}
				// Cache and merge the referenced schema.
				cache[absRefPath] = refSchema
				mergeSchemas(schema, refSchema)
			}
		}
	}
	return nil
}

// mergeSchemas merges types and elements from src into dst, avoiding duplicates.
func mergeSchemas(dst, src *model.Schema) {
	typeExists := func(types []model.XSDType, name string) bool {
		for _, t := range types {
			if t.Name == name {
				return true
			}
		}
		return false
	}
	elemExists := func(elems []model.XSDElement, name string) bool {
		for _, e := range elems {
			if e.Name == name {
				return true
			}
		}
		return false
	}
	for _, t := range src.Types {
		if !typeExists(dst.Types, t.Name) {
			dst.Types = append(dst.Types, t)
		}
	}
	for _, e := range src.Elements {
		if !elemExists(dst.Elements, e.Name) {
			dst.Elements = append(dst.Elements, e)
		}
	}
}
