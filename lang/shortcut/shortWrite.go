package shortcut

import (
	"math"
	"strconv"
	"strings"
	"unicode"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/rt"
)

// this can only write addresses and dots composed of literals
func WriteDots(a rt.Address) (ret string, okay bool) {
	var out strings.Builder
	switch a := a.(type) {
	case *object.VariableDot:
		if n, ok := getLiteralString(a.VariableName); ok && len(n) > 0 {
			out.WriteRune(AtSign)
			out.WriteString(formatName(n))
			if writeDots(&out, a.Dot) {
				ret, okay = out.String(), true
			}
		}
	case *object.ObjectDot:
		if n, ok := getLiteralString(a.NounName); ok && len(n) > 0 {
			out.WriteRune(HashMark)
			out.WriteString(formatName(n))
			if writeDots(&out, a.Dot) {
				ret, okay = out.String(), true
			}
		}
	}
	return
}

func writeDots(out *strings.Builder, dots []object.Dot) (okay bool) {
	okay = true // provisionally
Loop:
	for _, d := range dots {
		switch d := d.(type) {
		case *object.AtField:
			if field, ok := getLiteralString(d.FieldName); !ok {
				okay = false
				break Loop
			} else {
				out.WriteRune(FieldSeparator)
				out.WriteString(formatName(field))
			}
		case *object.AtIndex:
			if idx, ok := getLiteralNumber(d.Index); !ok {
				okay = false
				break Loop
			} else {
				out.WriteRune(FieldSeparator)
				out.WriteString(strconv.Itoa(idx + 1))
			}
		default:
			okay = false
			break Loop
		}
	}
	return
}

func formatName(str string) (ret string) {
	var b strings.Builder
	var spacing bool
Loop:
	for _, q := range str {
		switch {
		case unicode.IsLetter(q) || (b.Len() > 0 && unicode.IsDigit(q)):
			b.WriteRune(q)
			spacing = false
		case q == Space || q == Underscore:
			if !spacing {
				b.WriteRune(Underscore)
				spacing = true
			}
		default:
			ret = "`" + str + "`"
			b.Reset()
			break Loop
		}
	}
	if b.Len() > 0 {
		ret = b.String()
	}
	return
}

func getLiteralString(t rt.TextEval) (ret string, okay bool) {
	if a, ok := t.(*literal.TextValue); ok && len(a.Kind) == 0 {
		ret, okay = a.Value, true
	}
	return
}

func getLiteralNumber(t rt.NumEval) (ret int, okay bool) {
	if a, ok := t.(*literal.NumValue); ok {
		if v := math.Floor(a.Value); v == a.Value {
			ret, okay = int(v), true
		}
	}
	return
}
