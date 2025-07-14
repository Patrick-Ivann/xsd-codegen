// Package xmlgen provides utilities for generating dummy XML data from Go structs.
package xmlgen

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

// GenerateDummyValue produces a dummy value for a given field, respecting XSD restrictions.
func GenerateDummyValue(field model.XSDField) string {
	r := field.Restriction
	if r != nil {
		// Handle enumeration: pick the first allowed value.
		if len(r.Enumeration) > 0 {
			return r.Enumeration[0]
		}
		// Handle pattern: return a placeholder matching the pattern.
		if r.Pattern != nil {
			return "patternmatch"
		}
		// Handle minLength/maxLength/length for strings.
		if r.Length != nil {
			return strings.Repeat("x", *r.Length)
		}
		if r.MinLength != nil {
			n := *r.MinLength
			if r.MaxLength != nil && *r.MaxLength > n {
				n = *r.MaxLength
			}
			return strings.Repeat("x", n)
		}
		if r.MaxLength != nil {
			return strings.Repeat("x", *r.MaxLength)
		}
		// Handle numeric restrictions.
		if isNumericType(field.Type) {
			return generateNumericDummy(r)
		}
	}
	// Default dummy values by type.
	switch field.Type {
	case "string":
		return "example"
	case "int":
		return "42"
	case "float64":
		return "3.14"
	case "bool":
		return "true"
	case "time.Time":
		return time.Now().Format(time.RFC3339)
	default:
		return "value"
	}
}

// isNumericType determines if the field type is numeric.
func isNumericType(typ string) bool {
	switch typ {
	case "int", "float64", "decimal":
		return true
	}
	return false
}

// generateNumericDummy creates a dummy number as a string, respecting restrictions.
func generateNumericDummy(r *model.Restriction) string {
	min := 1
	max := 100
	// Handle min/max inclusive/exclusive.
	if r.MinInclusive != nil {
		if v, err := strconv.Atoi(*r.MinInclusive); err == nil {
			min = v
		}
	}
	if r.MinExclusive != nil {
		if v, err := strconv.Atoi(*r.MinExclusive); err == nil {
			min = v + 1
		}
	}
	if r.MaxInclusive != nil {
		if v, err := strconv.Atoi(*r.MaxInclusive); err == nil {
			max = v
		}
	}
	if r.MaxExclusive != nil {
		if v, err := strconv.Atoi(*r.MaxExclusive); err == nil {
			max = v - 1
		}
	}
	if max < min {
		max = min
	}
	val := rand.Intn(max-min+1) + min

	// Handle totalDigits/fractionDigits.
	if r.TotalDigits != nil && r.FractionDigits != nil {
		// Compose a float with the correct number of digits.
		base := float64(val)
		return strconv.FormatFloat(base, 'f', *r.FractionDigits, 64)
	}
	if r.TotalDigits != nil {
		// Generate a number with the exact number of digits.
		digits := *r.TotalDigits
		s := "1" + strings.Repeat("0", digits-1)
		return s
	}
	return strconv.Itoa(val)
}
