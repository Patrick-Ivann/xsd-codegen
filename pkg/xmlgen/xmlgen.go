// Package xmlgen provides utilities to generate XML output from Go structs with dummy data.
package xmlgen

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/Patrick-Ivann/xsd-codegen/pkg/model"
)

// GenerateDummyStruct recursively creates a Go struct instance with dummy data for XML marshaling.
func GenerateDummyStruct(typ model.XSDType) interface{} {
	// Use a map to represent the struct fields dynamically.
	m := make(map[string]interface{})
	for _, f := range typ.Fields {
		m[f.Name] = dummyValueForType(f)
	}
	for _, attr := range typ.Attributes {
		m[attr.Name] = dummyValueForType(model.XSDField{Name: attr.Name, Type: attr.Type})
	}
	return m
}

// dummyValueForType returns a dummy value for a given XSDField.
func dummyValueForType(f model.XSDField) interface{} {
	val := GenerateDummyValue(f)
	switch f.Type {
	case "int":
		return parseInt(val)
	case "float64":
		return parseFloat(val)
	case "bool":
		return val == "true"
	case "time.Time":
		t, _ := time.Parse(time.RFC3339, val)
		return t
	default:
		return val
	}
}

// parseInt safely converts a string to int.
func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

// parseFloat safely converts a string to float64.
func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

// MarshalDummyXML marshals the dummy struct into XML with the given root element name.
func MarshalDummyXML(rootName string, typ model.XSDType) ([]byte, error) {
	dummy := GenerateDummyStruct(typ)
	// Create a wrapper struct for the root element.
	root := map[string]interface{}{rootName: dummy}
	buf := &bytes.Buffer{}
	enc := xml.NewEncoder(buf)
	enc.Indent("", "  ")
	// Use a helper to convert map to XML tokens.
	if err := encodeMapToXML(enc, root, ""); err != nil {
		return nil, err
	}
	enc.Flush()
	return buf.Bytes(), nil
}

// encodeMapToXML recursively encodes a map[string]interface{} as XML tokens.
func encodeMapToXML(enc *xml.Encoder, m map[string]interface{}, parent string) error {
	for k, v := range m {
		start := xml.StartElement{Name: xml.Name{Local: k}}
		if err := enc.EncodeToken(start); err != nil {
			return err
		}
		switch val := v.(type) {
		case map[string]interface{}:
			if err := encodeMapToXML(enc, val, k); err != nil {
				return err
			}
		case []interface{}:
			for _, item := range val {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if err := encodeMapToXML(enc, itemMap, k); err != nil {
						return err
					}
				} else {
					if err := enc.EncodeElement(item, xml.StartElement{Name: xml.Name{Local: k}}); err != nil {
						return err
					}
				}
			}
		default:
			if err := enc.EncodeToken(xml.CharData([]byte(fmt.Sprint(val)))); err != nil {
				return err
			}
		}
		if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: k}}); err != nil {
			return err
		}
	}
	return nil
}
