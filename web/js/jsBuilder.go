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

func (out *Builder) Q(el string) *Builder {
	return out.Brace(Quotes, func(out *Builder) {
		out.S(el)
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
