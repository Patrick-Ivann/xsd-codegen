package helpers

import "github.com/Patrick-Ivann/xsd-codegen/pkg/model"

// ValueGenerator defines an interface for generating values from XSD types.
type ValueGenerator interface {
	Generate(xsdType string, restriction *model.XSDRestriction) string
}
