package blocks

import (
	"strings"
)

const space = ' '
const comma = ','
const colon = ':'

var obj = [2]rune{'{', '}'}
var array = [2]rune{'[', ']'}
var quote = [2]rune{'"', '"'}

func Embrace(style [2]rune, cb func(*Js)) string {
	var out Js
	return out.Brace(style, cb).String()
}

// an array of comma separated quoted strings
func QuotedStrings(els []string) string {
	var out Js
	return out.Brace(array, func(out *Js) {
		var sep Commas
		for _, el := range els {
			out.R(sep.Next()).Q(el)
		}
	}).String()
}

type Commas int

func (c *Commas) Next() (ret rune) {
	n := *c
	if n == 0 {
		ret = -1
	} else {
		ret = comma
	}
	*c = n + 1
	return
}

type Js struct {
	strings.Builder
}

func (out *Js) R(el rune) *Js {
	if el >= 0 {
		out.WriteRune(el)
	}
	return out
}

func (out *Js) S(el string) *Js {
	out.WriteString(el)
	return out
}

func (out *Js) Q(el string) *Js {
	return out.Brace(quote, func(out *Js) {
		out.S(el)
	})
}

func Q(el string) *Js {
	var out Js
	return out.Q(el)
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
