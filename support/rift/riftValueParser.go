package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/charmed"
	"github.com/ionous/errutil"
)

// parses the "right hand side" of a collection or map
// assumes the next rune is at the start of the value: no leading whitespace.
type ValueParser struct {
	hist   *History
	indent int
	inner  valueGetter
}

type valueGetter interface {
	GetValue() (ret any, err error)
}

func NewValue(hist *History, indent int, writeBack func(v any) error) charm.State {
	p := &ValueParser{hist: hist, indent: indent}
	return hist.PushIndent(indent, p, func() (err error) {
		if p.inner == nil {
			err = errutil.New("no value found") // is this an error?
		} else if v, e := p.inner.GetValue(); e != nil {
			err = e
		} else {
			err = writeBack(v)
		}
		return
	})
}

func (p *ValueParser) NewRune(r rune) (ret charm.State) {
	const dashOrMinus = SequenceDash
	switch {
	case r == InterpretedQuotes:
		next := new(interpretedString)
		ret = p.runInner(r, next, next)
	case charmed.IsNumber(r) || r == '+':
		next := new(numValue)
		ret = p.runInner(r, next, next)
	case r == 't' || r == 'f':
		next := new(boolValue)
		ret = p.runInner(r, next, next)
	case r == dashOrMinus:
		// ahh the pain of negative numbers and sequences
		// no space indicates a number `-5`
		// otherwise, its a sequence `- 5`
		ret = charm.Statement("dashing", func(r rune) (ret charm.State) {
			if r != Space && r != Newline {
				next := new(numValue)
				ret = p.runInner(dashOrMinus, next, next).NewRune(r)
			} else {
				next := NewSequence(p.hist, p.indent, func(vs []any) (_ error) {
					p.inner = computedValue{vs}
					return
				})
				ret = next.NewRune(dashOrMinus).NewRune(r)
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

func (p *ValueParser) runInner(r rune, i valueGetter, c charm.State) charm.State {
	p.inner = i
	// on the first unhandled rune; pop ourselves
	// we could wait till the parent is ready
	// return charm.RunStep(r, c, charm.Statement("post value", func(next rune) charm.State {
	// 	return charm.RunState(next, p.hist.PopIndent(p.indent))
	// }))
	return charm.RunState(r, c)
}

type computedValue struct{ v any }

func (p computedValue) GetValue() (ret any, err error) {
	ret = p.v
	return
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
