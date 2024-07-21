package shortcut

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/ionous/tell/charm"
)

type Type int

//go:generate stringer -type=Type
const (
	InvalidToken Type = iota
	ObjMarker         // rune
	VarMarker         // rune
	PlainName         // string
	QuotedName        // string
	Number            // number
)

func (t Type) IsMarker() bool {
	return t == ObjMarker || t == VarMarker
}
func (t Type) IsName() bool {
	return t == PlainName || t == QuotedName
}
func (t Type) IsNumber() bool {
	return t == Number
}

// tbd: maybe a channel instead?
type Notifier interface {
	Decoded(Type, any) error
}

const (
	HashMark       = '#'
	AtSign         = '@'
	FieldSeparator = '.'
	RawQuote       = '`'
	DoubleQuote    = '"'
	Underscore     = '_'
	Eof            = charm.Eof
	Space          = ' '
)

// this isnt a shortcut string
// the value of the error is 0 for normal strings
// or 1 if there is a double at or double hash
// ( so you can slice with this error code ]
type NotShort int

func (e NotShort) Error() string {
	return "not a shortcut string"
}

func invalidRune(q rune, at int) charm.State {
	err := fmt.Errorf("invalid rune %v at %d", q, at)
	return charm.Error(err)
}

func errorAt(e error, at int) charm.State {
	err := fmt.Errorf("%w at %d", e, at)
	return charm.Error(err)
}

func shortless(ofs int) charm.State {
	return charm.Error(NotShort(ofs))
}

func Tokenize(str string, n Notifier) error {
	return charm.ParseEof(str, NewTokenizer(n))
}

func NewTokenizer(n Notifier) charm.State {
	return charm.Statement("tokenizer", func(q rune) (ret charm.State) {
		switch q {
		default:
			ret = shortless(0)
		case HashMark:
			if e := n.Decoded(ObjMarker, q); e != nil {
				ret = errorAt(e, 0)
			} else {
				ret = readName(n, q, 1)
			}
		case AtSign:
			if e := n.Decoded(VarMarker, q); e != nil {
				ret = errorAt(e, 0)
			} else {
				ret = readName(n, q, 1)
			}
		}
		return
	})
}

func readName(n Notifier, prev rune, pos int) charm.State {
	var name strings.Builder
	var spacing bool
	return charm.Self("readName", func(self charm.State, q rune) (ret charm.State) {
		at := pos
		pos++
		switch {
		// quotes at the start; read the whole quote
		case name.Len() == 0 && (q == RawQuote || q == DoubleQuote):
			ret = readQuote(n, q, pos)

		// a dot after a name?
		case name.Len() > 0 && q == FieldSeparator:
			res := name.String()
			if e := n.Decoded(PlainName, res); e != nil {
				ret = errorAt(e, at)
			} else {
				name.Reset()
				ret, spacing = self, false
			}

		// the end of the name after text.
		case name.Len() > 0 && q == Eof:
			res := name.String()
			if e := n.Decoded(PlainName, res); e != nil {
				ret = errorAt(e, at)
			} else {
				ret = charm.Finished() // done!
			}

		// building a name
		case name.Len() > 0 && q == Underscore:
			if !spacing {
				name.WriteRune(Space)
				spacing = true
			}
			ret = self

		case unicode.IsLetter(q):
			name.WriteRune(q)
			ret, spacing = self, false

		case unicode.IsDigit(q):
			// after a letter, handle numbers like letters.
			if name.Len() > 0 {
				name.WriteRune(q)
				ret, spacing = self, false
			} else {
				next := readNumber(n, at)
				ret = next.NewRune(q)
			}

		case name.Len() == 0 && q == prev && prev != 0:
			ret = shortless(1)

		default:
			// eof is invalid if we are expecting a name.
			ret = invalidRune(q, at)
		}
		return
	})
}

// expects to be called with at least one number
func readNumber(n Notifier, pos int) charm.State {
	var accum int
	return charm.Self("readNumber", func(self charm.State, q rune) (ret charm.State) {
		at := pos
		pos++
		switch v := int(q - '0'); {
		case v >= 0 && v <= 9:
			accum = (accum * 10) + v
			ret = self

		case q == Eof:
			if e := n.Decoded(Number, accum); e != nil {
				ret = errorAt(e, pos)
			} else {
				ret = charm.Finished() // done!
			}

		case q == FieldSeparator:
			if e := n.Decoded(Number, accum); e != nil {
				ret = errorAt(e, pos)
			} else {
				ret = readName(n, 0, pos)
			}

		default:
			ret = invalidRune(q, at)
		}
		return
	})
}

// all that can appear is a field separator or end of string.
func readPath(n Notifier, ofs int) charm.State {
	return charm.Statement("readPath", func(q rune) (ret charm.State) {
		switch q {
		case FieldSeparator:
			ret = readName(n, 0, ofs+1)
		case Eof:
			ret = charm.Finished() // done
		default:
			ret = invalidRune(q, ofs)
		}
		return
	})
}

// reads without escaping.
// ( assumes that the string was descaped when reading from disk )
func readQuote(n Notifier, match rune, ofs int) charm.State {
	var name strings.Builder
	return charm.Self("readQuote", func(self charm.State, q rune) (ret charm.State) {
		switch q {
		default:
			name.WriteRune(q)
			ret, ofs = self, ofs+1

		case match:
			if e := n.Decoded(QuotedName, name.String()); e != nil {
				ret = errorAt(e, ofs)
			} else {
				ret = readPath(n, ofs+name.Len()+1)
			}

		case Eof:
			ret = invalidRune(q, ofs)
		}
		return
	})
}
