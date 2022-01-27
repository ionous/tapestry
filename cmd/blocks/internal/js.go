package blocks

import (
	"strings"
)

const space = ' '
const comma = ','
const colon = ':'
const quote = '"'
const score = '_'
const percent = '%'
const unchecked = ""

var obj = [2]rune{'{', '}'}
var array = [2]rune{'[', ']'}
var quotes = [2]rune{quote, quote}

func Embrace(style [2]rune, cb func(*Js)) string {
	var out Js
	return out.Brace(style, cb).String()
}

type Js struct {
	strings.Builder
}

func (out *Js) R(els ...rune) *Js {
	for _, el := range els {
		if el >= 0 {
			out.WriteRune(el)
		}
	}
	return out
}

func (out *Js) S(el string) *Js {
	out.WriteString(el)
	return out
}

func (out *Js) Q(el string) *Js {
	return out.Brace(quotes, func(out *Js) {
		out.S(el)
	})
}

func (out *Js) Kv(k, v string) *Js {
	return out.Q(k).R(colon).Q(v)
}

func (out *Js) If(b bool, cb func(*Js)) *Js {
	if b {
		cb(out)
	}
	return out
}

func (out *Js) Brace(style [2]rune, cb func(*Js)) *Js {
	out.WriteRune(style[0])
	cb(out)
	out.WriteRune(style[1])
	return out
}

func quotedStrings(values []string) string {
	var out Js
	for i, el := range values {
		if i > 0 {
			out.R(comma)
		}
		out.Q(el)
	}
	return out.String()
}
