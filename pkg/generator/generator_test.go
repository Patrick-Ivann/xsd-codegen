package generator

import (
	"os"
	"testing"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

func TestGenerateGoCode_WithRestrictionTags(t *testing.T) {
	min := "1"
	max := "10"
	pattern := "^[a-z]+$"
	schema := &model.Schema{
		Types: []model.XSDType{
			{
				Name: "TestType",
				Fields: []model.XSDField{
					{
						Name: "Field1",
						Type: "string",
						Restriction: &model.Restriction{
							Enumeration:  []string{"A", "B"},
							Pattern:      &pattern,
							MinInclusive: &min,
							MaxInclusive: &max,
						},
					},
				},
			},
		},
	}
	tmpl := "type {{.Name}} struct {\n{{range .Fields}}    {{.Name}} {{.Type}} `xml:\"{{.Name}}\" validate:\"{{restrictionTag .Restriction}}\"`\n{{end}}}\n"
	tmplPath := "test_struct.tmpl"
	outputPath := "test_output.go"
	if err := os.WriteFile(tmplPath, []byte(tmpl), 0644); err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}
	defer os.Remove(tmplPath)
	defer os.Remove(outputPath)
	err := GenerateGoCode(schema, tmplPath, outputPath)
	if err != nil {
		t.Errorf("GenerateGoCode failed: %v", err)
	}
	data, err := os.ReadFile(outputPath)
	if err != nil || len(data) == 0 {
		t.Errorf("Output file not written or empty")
	}
}
