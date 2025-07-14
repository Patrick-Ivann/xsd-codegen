package generator

import "testing"

// TestXSDToGoType checks correct mapping of XSD types to Go types.
func TestXSDToGoType(t *testing.T) {
	cases := map[string]string{
		"xs:string":  "string",
		"xs:int":     "int",
		"xs:float":   "float64",
		"xs:boolean": "bool",
		"xs:date":    "time.Time",
		"unknown":    "string",
	}
	for xsd, want := range cases {
		got := XSDToGoType(xsd)
		if got != want {
			t.Errorf("XSDToGoType(%q) = %q; want %q", xsd, got, want)
		}
	}
}
