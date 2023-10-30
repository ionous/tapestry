package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/charmed"
	"github.com/ionous/errutil"
)

// parses the "right hand side" of a collection or map
type ValueParser struct {
	inner valueGetter
	charm.State
}

type Value struct {
	Result any
	// comments
	// line number
	// etc.
}

type valueGetter interface {
	GetValue() (any, error)
}

func (p *ValueParser) StateName() string {
	return "values"
}

func (p *ValueParser) GetValue() (ret Value, err error) {
	if p.inner == nil {
		err = errutil.New("unknown value")
	} else if a, e := p.inner.GetValue(); e != nil {
		err = e
	} else {
		ret = Value{Result: a}
	}
	return
}

// fix? see operand
func (p *ValueParser) NewRune(r rune) (ret charm.State) {
	// could do a parallel, but the different types are distinguishable based on the first character
	// heredocs might actually break that, but all the string types could be housed under a single "string" parser
	var b boolValue
	var n numValue
	var q interpretedString
	// fix? could make this into a "ParseChoice" if you wanted
	// var args an arbitrary number of states.
	if next := b.NewRune(r); next != nil {
		p.inner = &b
		ret = next
	} else if next := n.NewRune(r); next != nil {
		p.inner = &n
		ret = next
	} else if next := q.NewRune(r); next != nil {
		p.inner = &q
		ret = next
	}
	return
}

// func (p *numValue) GetValue() (ret any, err error) {
// 	ret, err = p.GetNumber()
// 	return
// }

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
