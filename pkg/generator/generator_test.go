package generator

import (
	"os"
	"path/filepath"
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

func TestGenerateGoCode_WithDocumentation(t *testing.T) {
	// Simulate schema with one documented field
	schema := &model.Schema{
		Types: []model.XSDType{
			{
				Name: "MM_THENAMEIWANT",
				Fields: []model.XSDField{
					{
						Name:          "testtId",
						Type:          "string",
						Documentation: "this is the object identifier",
					},
				},
			},
		},
	}

	// Temporary template with documentation support
	tmpl := `{{- if .Documentation}}// {{.Documentation}}{{end}}
type {{title .Name}} struct {
{{- range .Fields }}
    {{- if .Documentation}}// {{.Documentation}}{{end}}
    {{title .Name}} {{goType .Type}} ` + "`xml:\"{{.Name}}\"`" + `
{{- end }}
}
`

	tmpDir := t.TempDir()
	tmplPath := filepath.Join(tmpDir, "struct.tmpl")
	outputPath := filepath.Join(tmpDir, "generated.go")

	if err := os.WriteFile(tmplPath, []byte(tmpl), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	err := GenerateGoCode(schema, tmplPath, outputPath)
	if err != nil {
		t.Fatalf("GenerateGoCode failed: %v", err)
	}

	contents, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	text := string(contents)
	if !contains(text, "this is the object identifier") {
		t.Errorf("Documentation comment not found in output:\n%s", text)
	}
	if !contains(text, "type MM_THENAMEIWANT") {
		t.Errorf("Struct name missing: %s", text)
	}
}

func contains(text, substr string) bool {
	return len(text) >= len(substr) && (text[0:len(substr)] == substr || text[len(text)-len(substr):] == substr || (len(text) > len(substr) && text[len(text)/2:len(text)/2+len(substr)] == substr))
}
