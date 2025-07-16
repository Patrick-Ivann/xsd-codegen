package generator

import (
	"fmt"
	"strings"
)

// XSDToGoType maps an XSD type to its Go equivalent.
// Returns "string" by default for unknown types.
func XSDToGoType(xsdType string, maxOccurs int) string {
	xsdTypeLower := strings.ToLower(xsdType)
	fmt.Printf("xsdTypeLower: %v\n", xsdTypeLower)
	switch xsdTypeLower {
	case "string", "xs:string", "xsd:string":
		return wrapSlice("string", maxOccurs)
	case "int", "xs:int", "xsd:int", "integer", "xs:integer", "xsd:integer":
		return wrapSlice("int", maxOccurs)
	case "float", "xs:float", "xsd:float", "double", "xs:double", "xsd:double":
		return wrapSlice("float64", maxOccurs)
	case "boolean", "xs:boolean", "xsd:boolean":
		return wrapSlice("bool", maxOccurs)
	case "date", "xs:date", "xsd:date", "datetime", "xs:datetime", "xsd:datetime":
		return wrapSlice("time.Time", maxOccurs)
	default:
		// Handle custom types by stripping prefix and Title-casing
		parts := strings.Split(xsdType, ":")
		fmt.Printf("parts[len(parts)-1]: %v\n", parts[len(parts)-1])
		customType := Title(parts[len(parts)-1])
		return wrapSlice(customType, maxOccurs)
	}
}

func wrapSlice(base string, maxOccurs int) string {
	if maxOccurs > 1 || maxOccurs == -1 {
		return "[]" + base
	}
	return base
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
