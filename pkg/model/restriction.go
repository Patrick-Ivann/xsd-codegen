package model

// Restriction represents XSD simple type restrictions/facets.
type Restriction struct {
	Enumeration    []string // Allowed values
	FractionDigits *int     // Max decimal places
	Length         *int     // Exact length
	MaxExclusive   *string  // Value must be less than this
	MaxInclusive   *string  // Value must be <= this
	MaxLength      *int     // Max length
	MinExclusive   *string  // Value must be greater than this
	MinInclusive   *string  // Value must be >= this
	MinLength      *int     // Min length
	Pattern        *string  // Regex pattern
	TotalDigits    *int     // Exact number of digits
	WhiteSpace     *string  // Whitespace handling
}
