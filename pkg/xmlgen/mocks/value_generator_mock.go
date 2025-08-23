package mocks

import (
	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
	"github.com/stretchr/testify/mock"
)

// MockValueGenerator is a mock implementation of ValueGenerator.
type MockValueGenerator struct {
	mock.Mock
}

func (m *MockValueGenerator) Generate(xsdType string, restriction *model.XSDRestriction) string {
	args := m.Called(xsdType, restriction)
	return args.String(0)
}
