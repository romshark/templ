package fmtstr

import (
	"errors"
	"strings"
	"unicode"

	"github.com/a-h/parse"
)

// PartType is either a literal (PartTypeLiteral) or any of:
// https://pkg.go.dev/fmt#hdr-Printing
type PartType int8

const (
	_ PartType = iota

	PartTypeLiteral // plain text not part of a format specifier
	PartTypeBase10  // integer formats: %d, %b
	PartTypeString  // string formats: %s, %q
	PartTypeFloat   // floating-point formats: %f, %e, %E, %g, %G
	PartTypeBool    // boolean format: %t
	PartTypePointer // pointer format: %p
	PartTypeChar    // character format: %c
	PartTypeHex     // hexadecimal formats: %x, %X
	PartTypeOctal   // octal format: %o
	PartTypeGeneric // generic formats: %v, %+v etc.
)

type Part struct {
	Type  PartType
	Value string
}

type FormatString struct {
	Raw   string
	Parts []Part
}

var (
	ErrUnterm          = errors.New("unterminated quoted string")
	ErrTrailingPercent = errors.New(`trailing '%' at end of string`)
	ErrUntermVerb      = errors.New("unterminated format verb")
)

// Parse parses a quoted format string like `"foo-bar-%d (%s)"` into a FormatString.
func Parse(pi *parse.Input) (FormatString, error) {
	start := pi.Index()

	// Opening quote
	chunk, ok := pi.Peek(1)
	if !ok || chunk[0] != '"' {
		return FormatString{}, nil // not a quoted string at this position
	}
	_, _ = pi.Take(1) // consume opening "

	var parts []Part
	var lit strings.Builder

	flushLit := func() {
		if lit.Len() > 0 {
			parts = append(parts, Part{
				Type:  PartTypeLiteral,
				Value: lit.String(),
			})
			lit.Reset()
		}
	}

	for {
		c, ok := pi.Peek(1)
		if !ok {
			pi.Seek(start)
			return FormatString{}, ErrUnterm
		}
		if c[0] == '"' {
			// reached closing quote
			_, _ = pi.Take(1)
			break
		}

		if c[0] != '%' {
			// literal character
			lit.WriteByte(c[0])
			_, _ = pi.Take(1)
			continue
		}

		// found '%'
		_, _ = pi.Take(1) // consume '%'
		next, ok := pi.Peek(1)
		if !ok {
			return FormatString{}, ErrTrailingPercent
		}
		if next == "%" {
			// literal %%
			lit.WriteByte('%')
			_, _ = pi.Take(1)
			continue
		}

		flushLit()

		// parse format specifier
		spec := "%"
		for {
			c2, ok := pi.Peek(1)
			if !ok {
				return FormatString{}, ErrUntermVerb
			}
			verb := rune(c2[0])
			spec += c2
			_, _ = pi.Take(1)

			if unicode.IsLetter(verb) {
				parts = append(parts, Part{
					Type:  classifyVerb(verb),
					Value: spec,
				})
				break
			}
			// else keep consuming flags/width/precision
		}
	}

	end := pi.Index()
	pi.Seek(start)
	raw, _ := pi.Peek(end - start)
	pi.Seek(end)

	flushLit()
	return FormatString{Raw: raw, Parts: parts}, nil
}

func classifyVerb(r rune) PartType {
	switch r {
	case 'd', 'b':
		return PartTypeBase10
	case 's', 'q':
		return PartTypeString
	case 'f', 'e', 'E', 'g', 'G':
		return PartTypeFloat
	case 't':
		return PartTypeBool
	case 'p':
		return PartTypePointer
	case 'c':
		return PartTypeChar
	case 'x', 'X':
		return PartTypeHex
	case 'o':
		return PartTypeOctal
	case 'v': // generic
	}
	return PartTypeGeneric
}
