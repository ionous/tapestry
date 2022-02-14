package js

import (
	"strconv"
	"strings"
)

const Space = ' '
const Comma = ','
const Colon = ':'
const Quote = '"'
const Score = '_'
const Percent = '%'
const True = "true"
const False = "false"

var Obj = [2]rune{'{', '}'}
var Array = [2]rune{'[', ']'}
var Quotes = [2]rune{Quote, Quote}

// wraps strings builder with some json helper functions.
type Builder struct{ strings.Builder }

func (out *Builder) R(els ...rune) *Builder {
	for _, el := range els {
		if el >= 0 {
			out.WriteRune(el)
		}
	}
	return out
}

func (out *Builder) N(n int) *Builder {
	out.WriteString(strconv.Itoa(n))
	return out
}

func (out *Builder) S(el string) *Builder {
	out.WriteString(el)
	return out
}

// write the inner bits of a json friend string
// see: https://cs.opensource.google/go/go/+/refs/tags/go1.17.6:src/encoding/json/encode.go;l=1036
func (out *Builder) X(s string) *Builder {
	var start int
	for i, c := range s {
		// is this rune a control character or other special value?
		if esc, special := escapes[c]; special || c < 32 {
			if start < i { // flush
				out.WriteString(s[start:i])
			}
			start = i + 1 // skip the character we're escaping
			if special {
				out.WriteString(esc)
			} else {
				// write in the default \uxxxx format
				const hex = "0123456789abcdef"
				out.WriteString(`u00`)
				out.WriteByte(hex[c>>4])
				out.WriteByte(hex[c&0xF])
			}
		}
	}
	if start < len(s) { // flush
		out.WriteString(s[start:])
	}
	return out
}

// special escapes for json strings.
// any codepoint except " or \ or a control character is fine
// ( re: https://www.json.org/json-en.html, plus or minus some separators. )
// all control characters *can* be written in \uxxxx format
// while some also have their own special syntax ( ex. \n instead of \u000a )
var escapes = map[rune]string{
	'"':  `\"`,
	'\\': `\\`,
	'\b': `\b`,
	'\f': `\f`,
	'\n': `\n`,
	'\r': `\r`,
	'\t': `\t`,
	// U+2028 is LINE SEPARATOR.
	// U+2029 is PARAGRAPH SEPARATOR.
	'\u2028': `\u2028`,
	'\u2029': `\u2029`,
}

func (out *Builder) Q(el string) *Builder {
	return out.Brace(Quotes, func(out *Builder) {
		out.X(el)
	})
}

func (out *Builder) Kv(k, v string) *Builder {
	return out.Q(k).R(Colon).Q(v)
}

func (out *Builder) If(b bool, cb func(*Builder)) *Builder {
	if b {
		cb(out)
	}
	return out
}

func (out *Builder) Brace(style [2]rune, cb func(*Builder)) *Builder {
	out.WriteRune(style[0])
	cb(out)
	out.WriteRune(style[1])
	return out
}
