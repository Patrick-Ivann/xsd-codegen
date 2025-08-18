package helpers

import (
	"fmt"
	"math/rand"
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
// This function handles enumerations, patterns, and type-specific value generation.
func GenerateValue(xsdType string, restriction *model.XSDRestriction) string {
	// If there are enumerations, pick one at random.
	if restriction != nil && len(restriction.Enumerations) > 0 {
		return restriction.Enumerations[rand.Intn(len(restriction.Enumerations))].Value
	}

	switch xsdType {
	case "xsd:string":
		if restriction != nil && restriction.Pattern != nil {
			return generateStringFromPattern(restriction.Pattern.Value)
		}
		return RandomString(rand.Intn(10) + defaultStringLength)

	case "xsd:decimal", "xsd:double", "xsd:float":
		return strconv.FormatFloat(RandomFloat(-1e6, 1e6), 'f', 6, 64)

	case "xsd:positiveInteger", "xsd:nonNegativeInteger":
		return strconv.Itoa(randomIntFromRestriction(restriction, defaultMinInt, defaultMaxInt))

	case "xsd:integer", "xsd:int", "xsd:long", "xsd:short", "xsd:byte":
		return strconv.Itoa(RandomInt(-10000, 10000))

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
		return []string{"true", "false"}[rand.Intn(2)]

	default:
		if restriction != nil {
			if restriction.Pattern != nil {
				return generateStringFromPattern(restriction.Pattern.Value)
			}
			if restriction.MinIncl != nil && restriction.MaxExcl != nil {
				min, _ := strconv.Atoi(restriction.MinIncl.Value)
				max, _ := strconv.Atoi(restriction.MaxExcl.Value)
				return strconv.Itoa(RandomInt(min, max-defaultMinInt))
			}
		}
		return "default"
	}
}

// randomIntFromRestriction returns a random integer within the given restriction bounds.
// If no restriction is provided, it uses the provided default values.
func randomIntFromRestriction(r *model.XSDRestriction, defMin, defMax int) int {
	min, max := defMin, defMax
	if r != nil {
		if r.MinIncl != nil {
			if v, err := strconv.Atoi(r.MinIncl.Value); err == nil {
				min = v
			}
		}
		if r.MaxExcl != nil {
			if v, err := strconv.Atoi(r.MaxExcl.Value); err == nil {
				max = v
			}
		}
	}
	return RandomInt(min, max)
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
		sb.WriteByte(chars[rand.Intn(len(chars))])
	}
	return sb.String()
}

// RandomIdentifier generates a random identifier string.
// The first character is always a letter, followed by letters, numbers, or certain symbols.
func RandomIdentifier() string {
	const start = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const rest = start + "0123456789_-."
	var sb strings.Builder
	sb.WriteByte(start[rand.Intn(len(start))])
	for i := 0; i < rand.Intn(8)+2; i++ {
		sb.WriteByte(rest[rand.Intn(len(rest))])
	}
	return sb.String()
}

// RandomFloat generates a random float64 between min and max.
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandomInt generates a random integer between min and max, inclusive.
func RandomInt(min, max int) int {
	if min >= max {
		return min
	}
	return rand.Intn(max-min+defaultMinInt) + min
}

// RandomDate generates a random date string in YYYY-MM-DD format.
func RandomDate() string {
	start := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2030, 12, 31, 0, 0, 0, 0, time.UTC)
	delta := end.Sub(start)
	randomTime := start.Add(time.Duration(rand.Int63n(int64(delta))))
	return randomTime.Format("2006-01-02")
}

// RandomTime generates a random time string in HH:MM:SS format.
func RandomTime() string {
	return fmtTime(rand.Intn(24), rand.Intn(60), rand.Intn(60))
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
		rand.Intn(10), rand.Intn(12), rand.Intn(28),
		rand.Intn(24), rand.Intn(60), rand.Intn(60))
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
		return rand.Intn(3) + 1
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
func RandomBetween(min, max int) int {
	if min >= max {
		return min
	}
	return rand.Intn(max-min+1) + min
}
