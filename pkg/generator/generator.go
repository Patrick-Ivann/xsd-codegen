package generator

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

// GenerateGoCode generates Go structs from the schema model and writes them to the output file.
func GenerateGoCode(schema *model.Schema, tmplPath, outputPath string) error {
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}
	tmpl, err := template.New("struct").Funcs(template.FuncMap{
		"restrictionTag": restrictionTag,
	}).Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	for _, t := range schema.Types {
		err := tmpl.Execute(&buf, t)
		if err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}
		buf.WriteString("\n")
	}
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	return nil
}

// restrictionTag generates a struct tag string for XSD restrictions.
func restrictionTag(r *model.Restriction) string {
	if r == nil {
		return ""
	}
	tags := []string{}
	if len(r.Enumeration) > 0 {
		tags = append(tags, fmt.Sprintf("enum=%s", strings.Join(r.Enumeration, "|")))
	}
	if r.Pattern != nil {
		tags = append(tags, fmt.Sprintf("pattern=%s", *r.Pattern))
	}
	if r.MinLength != nil {
		tags = append(tags, fmt.Sprintf("minlen=%d", *r.MinLength))
	}
	if r.MaxLength != nil {
		tags = append(tags, fmt.Sprintf("maxlen=%d", *r.MaxLength))
	}
	if r.Length != nil {
		tags = append(tags, fmt.Sprintf("len=%d", *r.Length))
	}
	if r.MinInclusive != nil {
		tags = append(tags, fmt.Sprintf("min=%s", *r.MinInclusive))
	}
	if r.MaxInclusive != nil {
		tags = append(tags, fmt.Sprintf("max=%s", *r.MaxInclusive))
	}
	if r.MinExclusive != nil {
		tags = append(tags, fmt.Sprintf("minex=%s", *r.MinExclusive))
	}
	if r.MaxExclusive != nil {
		tags = append(tags, fmt.Sprintf("maxex=%s", *r.MaxExclusive))
	}
	if r.FractionDigits != nil {
		tags = append(tags, fmt.Sprintf("frac=%d", *r.FractionDigits))
	}
	if r.TotalDigits != nil {
		tags = append(tags, fmt.Sprintf("digits=%d", *r.TotalDigits))
	}
	if r.WhiteSpace != nil {
		tags = append(tags, fmt.Sprintf("ws=%s", *r.WhiteSpace))
	}
	return strings.Join(tags, ",")
}
