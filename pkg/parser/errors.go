// // errors.go
// // Defines custom error types for parsing and resolution errors.

// package parser

// import "fmt"

// // ErrParseFailure represents a general parsing failure.
// type ErrParseFailure struct {
//     File string
//     Msg  string
// }

// func (e *ErrParseFailure) Error() string {
//     return fmt.Sprintf("failed to parse XSD file %s: %s", e.File, e.Msg)
// }

// // ErrImportResolution represents an error resolving an <import> or <include>.
// type ErrImportResolution struct {
//     Ref string
//     Msg string
// }

// func (e *ErrImportResolution) Error() string {
//     return fmt.Sprintf("failed to resolve import/include %s: %s", e.Ref, e.Msg)
// }


// Package parser defines reusable parser-related error types.
package parser

import "errors"

// ErrEmptySchema indicates no content in schema file.
var ErrEmptySchema = errors.New("schema is empty")

// ErrInvalidXSD indicates the XSD couldn't be parsed.
var ErrInvalidXSD = errors.New("invalid XSD format")
