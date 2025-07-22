package model

import "strings"

// NormalizeType returns a Go-friendly version of the type string.
func NormalizeType(xsdType string) string {
	switch strings.ToLower(xsdType) {
	case "string":
		return "string"
	case "int", "integer":
		return "int"
	case "float", "decimal":
		return "float64"
	case "boolean":
		return "bool"
	case "date":
		return "string" // time.Time with custom unmarshaler?
	default:
		return "string"
	}
}

func EscapeQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
