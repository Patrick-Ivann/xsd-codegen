// model_test.go tests the model package structures.
package model

import "testing"

// TestSchemaStructs ensures Schema and related types can be instantiated and linked.
func TestSchemaStructs(t *testing.T) {
    attr := XSDAttribute{Name: "id", Type: "string"}
    field := XSDField{Name: "Name", Type: "string", MinOccurs: 1, MaxOccurs: 1}
    typ := XSDType{Name: "Person", Fields: []XSDField{field}, Attributes: []XSDAttribute{attr}}
    elem := XSDElement{Name: "Person", Type: "Person"}
    schema := Schema{Types: []XSDType{typ}, Elements: []XSDElement{elem}}
    if len(schema.Types) != 1 || len(schema.Elements) != 1 {
        t.Errorf("Schema not initialized correctly")
    }
}
