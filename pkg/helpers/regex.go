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

func (x *Xeger) gen(re *syntax.Regexp) string {
	switch re.Op {
	case syntax.OpLiteral:
		return string(re.Rune)

	case syntax.OpCharClass:
		return x.genCharClass(re)

	case syntax.OpAnyCharNotNL:
		return string(printableNotNL[randInt(len(printableNotNL))])

	case syntax.OpAnyChar:
		return string(printable[randInt(len(printable))])

	case syntax.OpCapture, syntax.OpConcat:
		return x.genSub(re, 1)

	case syntax.OpStar:
		return x.genSub(re, randInt(limit+1))

	case syntax.OpPlus:
		return x.genSub(re, randInt(limit)+1)

	case syntax.OpQuest:
		return x.genSub(re, randInt(2))

	case syntax.OpRepeat:
		max := re.Max
		if max < 0 {
			max = limit
		}
		return x.genSub(re, randInt(max-re.Min+1)+re.Min)

	case syntax.OpAlternate:
		return x.gen(re.Sub[randInt(len(re.Sub))])
	}
	return ""
}

func (x *Xeger) genCharClass(re *syntax.Regexp) string {
	total := 0
	for i := 0; i < len(re.Rune); i += 2 {
		total += int(re.Rune[i+1]-re.Rune[i]) + 1
	}
	idx := randInt(total)
	for i := 0; i < len(re.Rune); i += 2 {
		delta := int(re.Rune[i+1] - re.Rune[i])
		if idx <= delta {
			return string(rune(re.Rune[i] + rune(idx)))
		}
		idx -= delta + 1
	}
	return ""
}

func (x *Xeger) genSub(re *syntax.Regexp, count int) string {
	var out []byte
	for i := 0; i < count; i++ {
		for _, sub := range re.Sub {
			out = append(out, x.gen(sub)...)
		}
	}
	return string(out)
}

func randInt(n int) int {
	if n <= 0 {
		return 0
	}
	return rng.Intn(n)
}
