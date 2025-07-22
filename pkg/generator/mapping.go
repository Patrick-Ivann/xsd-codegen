package generator

import (
	"strings"
)

// XSDToGoType maps an XSD type to its Go equivalent.
// Returns "string" by default for unknown types.
func XSDToGoType(xsdType string, maxOccurs int) string {
	// Normalize input
	t := strings.ToLower(strings.TrimSpace(stripPrefix(xsdType)))

	base := map[string]string{
		// String types
		"string":           "string",
		"token":            "string",
		"normalizedstring": "string",
		"language":         "string",
		"name":             "string",
		"ncname":           "string",
		"id":               "string",
		"idref":            "string",
		"nmtoken":          "string",
		"notation":         "string",
		"anyuri":           "string",

		// Boolean
		"boolean": "bool",

		// Numbers
		"decimal": "float64",
		"float":   "float64",
		"double":  "float64",

		"int":                "int",
		"integer":            "int",
		"long":               "int",
		"short":              "int",
		"byte":               "int",
		"nonnegativeinteger": "int",
		"positiveinteger":    "int",
		"nonpositiveinteger": "int",
		"negativeinteger":    "int",
		"unsignedint":        "int",
		"unsignedlong":       "int",
		"unsignedshort":      "int",
		"unsignedbyte":       "int",

		// Dates & time
		"date":       "time.Time",
		"datetime":   "time.Time",
		"time":       "time.Time",
		"gyear":      "time.Time",
		"gmonth":     "time.Time",
		"gday":       "time.Time",
		"gyearmonth": "time.Time",
		"gmonthday":  "time.Time",
		"duration":   "string", // optionally time.Duration with custom parsing

		// Binary & misc
		"hexbinary":    "[]byte",
		"base64binary": "[]byte",
		"qname":        "xml.Name",
		"idrefs":       "[]string",
		"nmtokens":     "[]string",
		"entities":     "[]string",
		"anytype":      "interface{}",
	}

	goType := base[t]
	if goType == "" {
		// Treat custom type
		goType = title(stripPrefix(xsdType))
	}

	if maxOccurs > 1 || maxOccurs == -1 {
		return "[]" + goType
	}
	return goType
}

func stripPrefix(s string) string {
	// Remove known prefixes like xs:, xsd:, etc.
	if strings.Contains(s, ":") {
		parts := strings.Split(s, ":")
		return parts[len(parts)-1]
	}
	return s
}

func title(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
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
