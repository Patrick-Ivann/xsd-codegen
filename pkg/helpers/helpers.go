package helpers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

// Constants for default values and magic numbers.
const (
	defaultMinInt       = 1
	defaultMaxInt       = 100
	defaultStringLength = 5
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
		return "string"
	default:
		return "string"
	}
}

// EscapeQuotes replaces double quotes with escaped quotes.
func EscapeQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

// GenerateValue returns a sample value for a given XSD type and restriction.
// It handles enumerations, patterns, and type-specific value generation.
func GenerateValue(xsdType string, restriction *model.XSDRestriction) string {
	// Handle enumerations if present
	if restriction != nil && len(restriction.Enumerations) > 0 {
		return pickRandomEnumeration(restriction)
	}

	// Dispatch based on the XSD type
	switch xsdType {
	case "xsd:string":
		return generateStringValue(restriction)
	case "xsd:decimal", "xsd:double", "xsd:float":
		return generateFloatValue()
	case "xsd:positiveInteger", "xsd:nonNegativeInteger":
		return generatePositiveIntegerValue(restriction)
	case "xsd:integer", "xsd:int", "xsd:long", "xsd:short", "xsd:byte":
		return generateIntegerValue()
	case "xsd:NMTOKEN":
		return RandomIdentifier()
	case "xsd:date":
		return RandomDate()
	case "xsd:time":
		return RandomTime()
	case "xsd:dateTime":
		return randomDateTime()
	case "xsd:duration":
		return randomDuration()
	case "xsd:boolean":
		return randomBoolean()
	default:
		return generateDefaultValue(restriction)
	}
}

func secureIntn(n int) int {
	if n <= 0 {
		return 0
	}
	bn := big.NewInt(int64(n))
	r, err := rand.Int(rand.Reader, bn)
	if err != nil {
		return 0
	}
	return int(r.Int64())
}

func secureIntBetween(minVal, maxVal int) int {
	if minVal >= maxVal {
		return minVal
	}
	return secureIntn(maxVal-minVal+1) + minVal
}

func secureFloatBetween(minVal, maxVal float64) float64 {
	r, err := rand.Int(rand.Reader, big.NewInt(1e9))
	if err != nil {
		return minVal
	}
	f := float64(r.Int64()) / 1e9
	return minVal + f*(maxVal-minVal)
}

// pickRandomEnumeration selects a random value from enumeration restrictions.
func pickRandomEnumeration(restriction *model.XSDRestriction) string {
	return restriction.Enumerations[secureIntn(len(restriction.Enumerations))].Value
}

// generateStringValue generates a string value, using pattern if present.
func generateStringValue(restriction *model.XSDRestriction) string {
	if restriction != nil && restriction.Pattern != nil {
		return generateStringFromPattern(restriction.Pattern.Value)
	}
	return RandomString(secureIntn(10) + defaultStringLength)
}

// generateFloatValue generates a random floating point number as a string.
func generateFloatValue() string {
	return strconv.FormatFloat(RandomFloat(-1e6, 1e6), 'f', 6, 64)
}

// generatePositiveIntegerValue generates a positive integer value considering restrictions.
func generatePositiveIntegerValue(restriction *model.XSDRestriction) string {
	return strconv.Itoa(randomIntFromRestriction(restriction, defaultMinInt, defaultMaxInt))
}

// generateIntegerValue generates an integer value in a fixed range.
func generateIntegerValue() string {
	return strconv.Itoa(RandomInt(-10000, 10000))
}

// randomBoolean randomly returns "true" or "false".
func randomBoolean() string {
	n, err := rand.Int(rand.Reader, big.NewInt(2))
	if err != nil {
		// fallback or panic if crypto/rand fails, depending on your error policy
		panic("crypto/rand failed: " + err.Error())
	}
	if n.Int64() == 0 {
		return "false"
	}
	return "true"
}

// generateDefaultValue generates a default value when no specific type matched.
func generateDefaultValue(restriction *model.XSDRestriction) string {
	if restriction != nil {
		if restriction.Pattern != nil {
			return generateStringFromPattern(restriction.Pattern.Value)
		}
		if restriction.MinIncl != nil && restriction.MaxExcl != nil {
			minVal, _ := strconv.Atoi(restriction.MinIncl.Value)
			maxVal, _ := strconv.Atoi(restriction.MaxExcl.Value)
			return strconv.Itoa(RandomInt(minVal, maxVal-defaultMinInt))
		}
	}
	return "default"
}

// randomIntFromRestriction returns a random integer within the given restriction bounds.
// If no restriction is provided, it uses the provided default values.
func randomIntFromRestriction(r *model.XSDRestriction, defMin, defMax int) int {
	minVal, maxVal := defMin, defMax
	if r != nil {
		if r.MinIncl != nil {
			if v, err := strconv.Atoi(r.MinIncl.Value); err == nil {
				minVal = v
			}
		}
		if r.MaxExcl != nil {
			if v, err := strconv.Atoi(r.MaxExcl.Value); err == nil {
				maxVal = v
			}
		}
	}
	return RandomInt(minVal, maxVal)
}

// generateStringFromPattern generates a string that matches the given regex pattern.
// If pattern generation fails, it returns a default value.
func generateStringFromPattern(pattern string) string {
	pattern = strings.ReplaceAll(pattern, `\\`, `\`)
	xeger, err := NewXeger(pattern)
	if err != nil {
		return "match123"
	}
	return xeger.Generate()
}

// RandomString generates a random string of the specified length.
func RandomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(chars[secureIntn(len(chars))])
	}
	return sb.String()
}

// RandomIdentifier generates a random identifier string.
// The first character is always a letter, followed by letters, numbers, or certain symbols.
func RandomIdentifier() string {
	const start = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const rest = start + "0123456789_-."
	var sb strings.Builder
	sb.WriteByte(start[secureIntn(len(start))])
	for i := 0; i < secureIntn(8)+2; i++ {
		sb.WriteByte(rest[secureIntn(len(rest))])
	}
	return sb.String()
}

// RandomFloat generates a random float64 between min and max.
func RandomFloat(minVal, maxVal float64) float64 {
	return secureFloatBetween(minVal, maxVal)
}

// RandomInt generates a random integer between min and max, inclusive.
func RandomInt(minVal, maxVal int) int {
	return secureIntBetween(minVal, maxVal)
}

// RandomDate generates a random date string in YYYY-MM-DD format.
func RandomDate() string {
	start := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2030, 12, 31, 0, 0, 0, 0, time.UTC)
	delta := end.Sub(start)
	seconds := secureIntn(int(delta.Seconds()))
	randomTime := start.Add(time.Duration(seconds) * time.Second)
	return randomTime.Format("2006-01-02")
}

// RandomTime generates a random time string in HH:MM:SS format.
func RandomTime() string {
	return fmtTime(secureIntn(24), secureIntn(60), secureIntn(60))
}

// fmtTime formats hours, minutes, and seconds into a time string.
func fmtTime(h, m, s int) string {
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// randomDateTime generates a random dateTime string in ISO 8601 format.
func randomDateTime() string {
	return RandomDate() + "T" + RandomTime()
}

// randomDuration generates a random duration string in ISO 8601 format.
func randomDuration() string {
	return fmt.Sprintf("P%dY%dM%dDT%dH%dM%dS",
		secureIntn(10), secureIntn(12), secureIntn(28),
		secureIntn(24), secureIntn(60), secureIntn(60))
}

// ParseOccurs parses an XSD occurrence value into an integer.
// If the value is "unbounded", it returns a random value between 1 and 3.
// If parsing fails, it returns the default value.
// ParseOccurs parses an XSD occurrence value.
func ParseOccurs(val string, defaultVal int) int {
	switch val {
	case "":
		return defaultVal
	case "unbounded":
		return secureIntn(3) + 1
	default:
		i, err := strconv.Atoi(val)
		if err != nil {
			return defaultVal
		}
		return i
	}
}

// RandomBetween returns a random integer between min and max, inclusive.
// If min >= max, it returns min.
func RandomBetween(minVal, maxVal int) int {
	if minVal >= maxVal {
		return minVal
	}
	return secureIntn(maxVal-minVal+1) + minVal
}
