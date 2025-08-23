package helpers_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/helpers"
	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

func TestNormalizeType(t *testing.T) {
	cases := map[string]string{
		"string":  "string",
		"int":     "int",
		"float":   "float64",
		"decimal": "float64",
		"boolean": "bool",
		"date":    "string",
		"custom":  "string",
	}
	for input, expected := range cases {
		if got := helpers.NormalizeType(input); got != expected {
			t.Errorf("NormalizeType(%q) = %q; want %q", input, got, expected)
		}
	}
}

func TestEscapeQuotes(t *testing.T) {
	in := `He said "hello"`
	want := `He said \"hello\"`
	got := helpers.EscapeQuotes(in)
	if got != want {
		t.Errorf("EscapeQuotes = %q; want %q", got, want)
	}
}

func TestRandomIntFromRestriction_MinOnly(t *testing.T) {
	r := &model.XSDRestriction{
		MinIncl: &model.XSDValue{Value: "5"},
	}
	fmt.Printf("r: %v\n", r)
	n := helpers.RandomBetween(5, 100)
	if n < 5 || n > 100 {
		t.Errorf("Expected value between 5 and 100, got %d", n)
	}
}

func TestRandomIntFromRestriction_MaxOnly(t *testing.T) {
	r := &model.XSDRestriction{
		MaxExcl: &model.XSDValue{Value: "10"},
	}

	fmt.Printf("r: %v\n", r)
	n := helpers.RandomBetween(1, 10)
	if n < 1 || n >= 10 {
		t.Errorf("Expected value between 1 and 9, got %d", n)
	}
}
func TestGenerateValue_Enum(t *testing.T) {
	r := &model.XSDRestriction{
		Enumerations: []model.XSDValue{
			{Value: "A"}, {Value: "B"}, {Value: "C"},
		},
	}
	val := helpers.GenerateValue("xsd:string", r)
	if val != "A" && val != "B" && val != "C" {
		t.Errorf("unexpected enum value: %q", val)
	}
}

func TestGenerateValue_Pattern(t *testing.T) {
	r := &model.XSDRestriction{
		Pattern: &model.XSDPattern{Value: `\d{3}\w{3}`},
	}
	val := helpers.GenerateValue("xsd:string", r)
	if len(val) != 6 {
		t.Errorf("pattern value length = %d; want 6", len(val))
	}
}

func TestGenerateValue_Types(t *testing.T) {
	types := []string{
		"xsd:decimal", "xsd:positiveInteger", "xsd:int",
		"xsd:NMTOKEN", "xsd:date", "xsd:time", "xsd:dateTime",
		"xsd:boolean", "xsd:duration", "xsd:custom",
	}
	for _, typ := range types {
		val := helpers.GenerateValue(typ, nil)
		if val == "" {
			t.Errorf("GenerateValue(%q) returned empty string", typ)
		}
	}
}

func TestRandomString(t *testing.T) {
	s := helpers.RandomBetween(5, 10)
	val := helpers.RandomString(s)
	if len(val) != s {
		t.Errorf("RandomString length = %d; want %d", len(val), s)
	}
}

func TestRandomIdentifier(t *testing.T) {
	id := helpers.RandomIdentifier()
	if len(id) < 3 {
		t.Errorf("RandomIdentifier too short: %q", id)
	}
	if !strings.ContainsAny(string(id[0]), "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		t.Errorf("RandomIdentifier does not start with letter: %q", id)
	}
}

func TestRandomDate(t *testing.T) {
	d := helpers.RandomDate()
	if _, err := time.Parse("2006-01-02", d); err != nil {
		t.Errorf("RandomDate invalid format: %q", d)
	}
}

func TestRandomTime(t *testing.T) {
	tm := helpers.RandomTime()
	if _, err := time.Parse("15:04:05", tm); err != nil {
		t.Errorf("RandomTime invalid format: %q", tm)
	}
}

func TestParseOccurs(t *testing.T) {
	cases := map[string]int{
		"":          2,
		"unbounded": 1,
		"3":         3,
		"invalid":   2,
	}
	for input, expected := range cases {
		got := helpers.ParseOccurs(input, 2)
		if got < 0 || got > 3 {
			t.Errorf("ParseOccurs(%q) = %d; want ~%d", input, got, expected)
		}
	}
}

func TestRandomBetween(t *testing.T) {
	got := helpers.RandomBetween(5, 10)
	if got < 5 || got > 10 {
		t.Errorf("RandomBetween out of range: %d", got)
	}
}

func TestRandomInt(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "get 1",
			args: args{min: 1, max: 1},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helpers.RandomInt(tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("RandomInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
