package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/charmed"
	"github.com/ionous/errutil"
)

// parses the "right hand side" of a collection or map
// assumes the next rune is at the start of the value: no leading whitespace.
type ValueParser struct {
	target CollectionTarget
	indent int
}

func NewValue(tgt CollectionTarget, indent int) charm.State {
	return &ValueParser{target: tgt, indent: indent}
}

func (p *ValueParser) NewRune(r rune) (ret charm.State) {
	const dashOrMinus = SequenceDash
	switch {
	case r == InterpretedQuotes:
		ret = p.runInner(r, &interpretedString{})
	case charmed.IsNumber(r) || r == '+':
		ret = p.runInner(r, &numValue{})
	case r == 't' || r == 'f':
		ret = p.runInner(r, &boolValue{})
	case r == dashOrMinus:
		// ahh the pain of negative numbers and sequences
		ret = charm.Statement("dashing", func(next rune) (ret charm.State) {
			// no space after the dash indicates a number `-5`
			if next != Space && next != Newline && next != charm.Eof {
				ret = p.runInner(dashOrMinus, &numValue{}).NewRune(next)
			} else {
				// a space after the dash is a subsequence `- 5`
				ret = NewSequence(p.target, p.indent).
					NewRune(dashOrMinus).
					NewRune(next)
			}
			return
		})
	default:
		// note: implicit nil values dont reach here
		// ex. for sequences, the sequence hits the indent of the next value first.
		e := errutil.New("unexpected value")
		ret = charm.Error(e)
	}
	return
}

// after finishing the inner state, write the value
// ( alt: the value parsers could write to the target themselves:
//   but this allows reusing the existing code )
func (p *ValueParser) runInner(r rune, inner interface {
	charm.State
	GetValue() (any, error)
}) charm.State {
	return charm.RunStep(r, inner,
		charm.Statement("write value", func(_ rune) (ret charm.State) {
			if v, e := inner.GetValue(); e != nil {
				ret = charm.Error(e)
			} else if e := p.target.WriteValue(v); e != nil {
				ret = charm.Error(e)
			}
			return
		}))
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
