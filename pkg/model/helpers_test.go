package model

import "testing"

func TestNormalizeType(t *testing.T) {
    if got := NormalizeType("decimal"); got != "float64" {
        t.Errorf("expected float64, got %s", got)
    }
    if got := NormalizeType("DATE"); got != "string" {
        t.Errorf("expected string, got %s", got)
    }
}
