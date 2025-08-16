package helpers

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/parser"
)

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

func GenerateValue(xsdType string, restriction *parser.XSDRestriction) string {
	switch xsdType {
	case "xsd:string":
		return "string"
	case "xsd:decimal":
		return strconv.FormatFloat(rand.Float64()*1e7*(randSign()), 'f', 7, 64)
	case "xsd:positiveInteger":
		min := 1
		max := 100
		if restriction != nil {
			if restriction.MinIncl != nil {
				min, _ = strconv.Atoi(restriction.MinIncl.Value)
			}
			if restriction.MaxExcl != nil {
				max, _ = strconv.Atoi(restriction.MaxExcl.Value)
			}
		}
		return strconv.Itoa(rand.Intn(max-min) + min)
	case "xsd:date":
		return "2006-03-07"
	case "xsd:NMTOKEN":
		return "US"
	default:
		if restriction != nil && restriction.Pattern != nil {
			return "match123" // Stub: Implement go-xeger for real regex
		}
		return "default"
	}
}

func randSign() float64 {
	if rand.Intn(2) == 0 {
		return -1
	}
	return 1
}

func ParseOccurs(val string, defaultVal int) int {
	if val == "" {
		return defaultVal
	}
	if val == "unbounded" {
		return rand.Intn(2) + 1
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

func RandomBetween(min, max int) int {
	if min > max {
		return min
	}
	return rand.Intn(max-min+1) + min
}
