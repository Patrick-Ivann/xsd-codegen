package helpers

import "github.com/Patrick-Ivann/xsd-codegen/pkg/model"

// DefaultValueGenerator is the production implementation.
type DefaultValueGenerator struct{}

// Generate returns a random value based on type and restriction.
func (d DefaultValueGenerator) Generate(xsdType string, restriction *model.XSDRestriction) string {
	return GenerateValue(xsdType, restriction)
}
