package helpers

import (
	"math/rand"
	"regexp/syntax"
	"time"
)

const (
	limit          = 10
	printable      = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~\t\v\f\r\n"
	printableNotNL = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~\t\v\f\r"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type Xeger struct {
	re *syntax.Regexp
}

// NewXeger compiles a regex pattern into a syntax tree.
func NewXeger(pattern string) (*Xeger, error) {
	re, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		return nil, err
	}
	return &Xeger{re: re}, nil
}

// Generate returns a random string that matches the compiled regex.
func (x *Xeger) Generate() string {
	return x.gen(x.re)
}

// gen recursively generates a string that matches the given regular expression syntax.
// It supports all standard regex operations: literals, char classes, wildcards, groups, repetition, alternation, etc.
func (x *Xeger) gen(re *syntax.Regexp) string {
	switch re.Op {
	case syntax.OpLiteral:
		// Returns a string with the literal runes stored in the regex node
		return string(re.Rune)
	case syntax.OpCharClass:
		// Handles character classes ([abc], [a-z], etc.) by generating a single random rune from the class
		return x.genCharClass(re)
	case syntax.OpAnyCharNotNL:
		// Generates one random printable character except for newlines
		return string(printableNotNL[randInt(len(printableNotNL))])
	case syntax.OpAnyChar:
		// Generates one random printable character (may include newline)
		return string(printable[randInt(len(printable))])
	case syntax.OpCapture, syntax.OpConcat:
		// Handles grouping (parentheses) and concatenation by generating the sequence once
		return x.genSub(re, 1)
	case syntax.OpStar:
		// Handles the '*' repetition: zero or more times (up to a random limit)
		return x.genSub(re, randInt(limit+1))
	case syntax.OpPlus:
		// Handles the '+' repetition: one or more times (at least one, up to a random limit)
		return x.genSub(re, randInt(limit)+1)
	case syntax.OpQuest:
		// Handles the '?' repetition: zero or one time (randomly chosen)
		return x.genSub(re, randInt(2))
	case syntax.OpRepeat:
		// Handles explicit repetition: {min,max}
		max := re.Max
		if max < 0 { // No explicit max: use a global limit
			max = limit
		}
		// Generate between Min and Max repetitions randomly
		return x.genSub(re, randInt(max-re.Min+1)+re.Min)
	case syntax.OpAlternate:
		// Handles alternation (a|b): randomly pick one of the subexpressions
		return x.gen(re.Sub[randInt(len(re.Sub))])
	}
	// Fallback for unsupported/empty ops
	return ""
}

// genCharClass randomly selects and returns a rune from the character class specified by re.Rune.
// Handles compact character ranges ([a-z] is represented by pairs in re.Rune).
func (x *Xeger) genCharClass(re *syntax.Regexp) string {
	total := 0
	// Count total number of possible runes in all ranges
	for i := 0; i < len(re.Rune); i += 2 {
		total += int(re.Rune[i+1]-re.Rune[i]) + 1
	}
	idx := randInt(total) // Choose a random position in the total range
	for i := 0; i < len(re.Rune); i += 2 {
		delta := int(re.Rune[i+1] - re.Rune[i])
		// If idx falls into this range, compute the rune
		if idx <= delta {
			return string(rune(re.Rune[i] + rune(idx)))
		}
		idx -= delta + 1
	}
	// Fallback: empty string
	return ""
}

// genSub generates the subexpressions of a regex 'count' times, concatenating the results.
// Used for sequences (groups, concat, repetition).
func (x *Xeger) genSub(re *syntax.Regexp, count int) string {
	var out []byte
	for i := 0; i < count; i++ {
		for _, sub := range re.Sub {
			out = append(out, x.gen(sub)...)
		}
	}
	return string(out)
}

// randInt returns a uniform random integer in [0, n).
// Returns 0 if n <= 0 to avoid panics or negative indexes.
func randInt(n int) int {
	if n <= 0 {
		return 0
	}
	return rng.Intn(n)
}
