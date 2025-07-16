package generator

import "strings"

// XSDToGoType maps an XSD type to its Go equivalent.
// Returns "string" by default for unknown types.
func XSDToGoType(xsdType string) string {
	switch strings.ToLower(xsdType) {
	case "string", "xs:string":
		return "string"
	case "int", "xs:int", "integer", "xs:integer":
		return "int"
	case "float", "xs:float", "double", "xs:double":
		return "float64"
	case "boolean", "xs:boolean":
		return "bool"
	case "date", "xs:date", "datetime", "xs:datetime":
		return "time.Time"
	default:
		return "string"
	}
}

func Title(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func omitTag(min int) string {
	if min == 0 {
		return ",omitempty"
	}
	return ""
}
