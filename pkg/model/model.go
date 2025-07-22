// // Package model defines internal representations for parsed XSD components.
// package model

// // Schema represents an XSD schema.
// type Schema struct {
//     Elements []Element
//     Types    []Type
//     Includes []Directive
//     Imports  []Directive
// }

// // Element maps to an <xs:element>.
// type Element struct {
//     Name       string
//     Type       string
//     Attributes []Attribute
// }

// // Type maps to <xs:complexType> or <xs:simpleType>.
// type Type struct {
//     Name        string
//     Base        string
//     Constraints map[string]string
// }

// // Attribute represents an XSD attribute for an element.
// type Attribute struct {
//     Name string
//     Type string
// }

// // Directive represents <import> or <include>.
// type Directive struct {
//     SchemaLocation string
//     Namespace      string
// }

// Package model provides internal representations of XSD schema constructs.
package model

// Schema represents a parsed XSD schema with its types and elements.
type Schema struct {
	Types         []XSDType    // List of complex/simple types
	Elements      []XSDElement // Top-level elements
	Includes      []Directive  // <xs:include> directives
	Imports       []Directive  // <xs:import> directives
	Documentation string       // Schema-level documentation
}

// XSDType represents a complex or simple type in the schema.
type XSDType struct {
	Name          string         // Name of the type
	Fields        []XSDField     // Fields for complex types
	Attributes    []XSDAttribute // Attributes for the type
	Documentation string         // Type-level comment from <annotation><documentation>
}

// XSDField represents a field (element) within a complex type.
type XSDField struct {
	Name          string       // Field name
	Type          string       // Field type (XSD or Go type) XSD type string (e.g., xs:string)
	MinOccurs     int          // Minimum occurrences
	MaxOccurs     int          // Maximum occurrences
	Restriction   *Restriction // Optional restrictions for the field
	Documentation string       // Optional documentation from <annotation><documentation>
	Default       string
	Fixed         string
}

// XSDElement represents a top-level element in the schema.
type XSDElement struct {
	Name          string       // Element name
	Type          string       // Element type
	Restriction   *Restriction // Optional inline simpleType restrictions
	Documentation string       // Optional documentation from <annotation><documentation>
	MinOccurs     int          // Minimum occurrences
	MaxOccurs     int          // Maximum occurrences
}

// XSDAttribute represents an <xs:attribute> inside a complexType.
type XSDAttribute struct {
	Name          string // Attribute name
	Type          string // Attribute type (e.g., xs:string)
	Documentation string // Optional documentation from <annotation><documentation>
	Default       string
	Fixed         string
	MinOccurs     int // Minimum occurrences
	MaxOccurs     int // Maximum occurrences
}

// Directive models an <xs:import> or <xs:include>.
type Directive struct {
	SchemaLocation string // Path to external schema file
	Namespace      string // Optional target namespace
}
