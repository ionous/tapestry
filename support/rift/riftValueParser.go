package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/charmed"
)

// parses the "right hand side" of a collection or map
// assumes the next rune is at the start of the value: no leading whitespace.
type ValueParser struct {
	indent int
	inner  valueState
}

type Value struct {
	Result any // can be nil
	// comments
	// line number
	// etc.
}

func (v *Value) Empty() bool {
	// if we ever decided to allow an explicit 'null'
	// empty would still mean: "nothing at all was specified"
	return v.Result == nil
}

type valueState interface {
	charm.State
	GetValue() (any, error)
}

func (p *ValueParser) StateName() string {
	return "values"
}

// returns the parsed result, or a nil result if no value was found
func (p *ValueParser) GetValue() (ret Value, err error) {
	if p.inner != nil {
		if a, e := p.inner.GetValue(); e != nil {
			err = e
		} else {
			ret = Value{Result: a}
		}
	}
	return
}

// fix? see operand
func (p *ValueParser) NewRune(r rune) (ret charm.State) {
	const dashOrMinus = SequenceDash
	switch {
	case r == InterpretedQuotes:
		ret = runInner(r, p, &interpretedString{})
	case charmed.IsNumber(r) || r == '+':
		ret = runInner(r, p, &numValue{})
	case r == 't' || r == 'f':
		ret = runInner(r, p, &boolValue{})
	case r == dashOrMinus:
		// ahh the pain of negative numbers and sequences
		ret = charm.Statement("dashing", func(next rune) (ret charm.State) {
			// no space, then it must be a number `-5`
			if next != Space {
				ret = runInner(dashOrMinus, p, &numValue{}).NewRune(next)
			} else {
				// space, then a sub sequence `- 5`
				seq := Sequence{indent: p.indent}
				ret = runInner(dashOrMinus, p, &seq).NewRune(next)
			}
			return
		})

	default:
		// some other rune indicates either a nil value
		// ( or some parsing error if no other state can handle the rune )
	}
	return
}

func runInner(r rune, p *ValueParser, inner valueState) (ret charm.State) {
	p.inner = inner
	return inner.NewRune(r)
}

type interpretedString struct{ charmed.QuoteParser }

func (p *interpretedString) GetValue() (ret any, err error) {
	ret, err = p.GetString()
	return
}

// NewRune starts with the leading quote mark; it finishes just after the matching quote mark.
func (p *interpretedString) NewRune(r rune) (ret charm.State) {
	if r == InterpretedQuotes {
		ret = p.ScanQuote(r)
	}
	return
}

type boolValue struct{ charmed.BoolParser }

func (p *boolValue) GetValue() (ret any, err error) {
	ret, err = p.GetBool()
	return
}

type numValue struct{ charmed.NumParser }

// fix? returns float64 because json does
// could also return int64 when its int like
func (p *numValue) GetValue() (ret any, err error) {
	ret, err = p.GetFloat()
	return
}
