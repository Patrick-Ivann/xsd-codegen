package helpers_test

import (
	"fmt"
	"testing"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/helpers"
)

func TestNewXeger_ValidPattern(t *testing.T) {
	x, err := helpers.NewXeger(`\d{3}\w{2}`)
	if err != nil {
		t.Fatalf("NewXeger failed: %v", err)
	}
	if x == nil {
		t.Fatal("NewXeger returned nil")
	}
}

func TestNewXeger_InvalidPattern(t *testing.T) {
	_, err := helpers.NewXeger(`[`)
	if err == nil {
		t.Error("Expected error for invalid pattern")
	}
}

func TestXeger_GenerateLiteral(t *testing.T) {
	x, err := helpers.NewXeger(`abc`)
	if err != nil {
		t.Fatal(err)
	}
	val := x.Generate()
	if val != "abc" {
		t.Errorf("Expected 'abc', got %q", val)
	}
}

func TestXeger_GenerateCharClass(t *testing.T) {
	x, err := helpers.NewXeger(`[a-z]`)
	if err != nil {
		t.Fatal(err)
	}
	val := x.Generate()
	if len(val) != 1 || val[0] < 'a' || val[0] > 'z' {
		t.Errorf("Expected lowercase letter, got %q", val)
	}
}

func TestXeger_GenerateAnyCharAlt(t *testing.T) {
	x, err := helpers.NewXeger(`\s`)
	if err != nil {
		t.Fatal(err)
	}
	val := x.Generate()
	if len(val) != 1 {
		t.Errorf("Expected single character, got %q", val)
	}
}
func TestXeger_GenerateAnyChar(t *testing.T) {
	x, err := helpers.NewXeger(`.`)
	if err != nil {
		t.Fatal(err)
	}
	val := x.Generate()
	if len(val) != 1 {
		t.Errorf("Expected single character, got %q", val)
	}
}
func TestXeger_GenerateAnyCharOneOrMore(t *testing.T) {
	x, err := helpers.NewXeger(`a*`)
	if err != nil {
		t.Fatal(err)
	}
	val := x.Generate()
	if len(val) < +1 {
		t.Errorf("Expected one or more characters, got %d", len(val))
	}
}

func TestXeger_GenerateAnyCharZeroOrOne(t *testing.T) {
	x, err := helpers.NewXeger(`a?`)
	if err != nil {
		t.Fatal(err)
	}
	val := x.Generate()
	if len(val) != 1 && val != "" {
		fmt.Printf("val: %v\n", len(val))
		t.Errorf("Expected zero or one character, got %d", len(val))
	}
}

func TestXeger_GenerateRepeat(t *testing.T) {
	x, err := helpers.NewXeger(`a{2,4}`)
	if err != nil {
		t.Fatal(err)
	}
	val := x.Generate()
	if len(val) < 2 || len(val) > 4 {
		t.Errorf("Expected 2-4 'a's, got %q", val)
	}
}

func TestXeger_GenerateAlternate(t *testing.T) {
	x, err := helpers.NewXeger(`foo|bar`)
	if err != nil {
		t.Fatal(err)
	}
	val := x.Generate()
	if val != "foo" && val != "bar" {
		t.Errorf("Expected 'foo' or 'bar', got %q", val)
	}
}
