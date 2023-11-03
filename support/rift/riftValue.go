package rift

import (
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/charmed"
	"github.com/ionous/errutil"
)

// parses the "right hand side" of a collection or map
// assumes the next rune is at the start of the value: no leading whitespace.
type Value struct {
	hist  *History
	inner valueGetter
}

type valueGetter interface {
	GetValue() (ret any, err error)
}

func NewValue(hist *History, indent int, writeBack func(v any) error) charm.State {
	p := &Value{hist: hist}
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

func (p *Value) NewRune(r rune) (ret charm.State) {
	const dashOrMinus = SequenceDash
	switch {
	case r == InterpretedQuotes:
		next := new(interpretedString)
		ret = p.runInner(r, next, next)

	case charmed.IsNumber(r) || r == '+':
		next := new(numValue)
		ret = p.runInner(r, next, next)

	case unicode.IsLetter(r):
		// handle keys that might look like bools.
		mapIndent := p.mapIndent()
		ret = p.tryBool(r, func(partial string) charm.State {
			next := NewMapping(p.hist, mapIndent, func(vs MapValues) (_ error) {
				p.inner = computedValue{vs}
				return
			})
			for _, r := range partial {
				next = next.NewRune(r)
			}
			return next
		})

	case r == dashOrMinus:
		// ahh the pain of negative numbers and sequences
		// no space indicates a number `-5`
		// otherwise, its a sequence `- 5`
		ret = charm.Statement("dashing", func(r rune) (ret charm.State) {
			if r != Space && r != Newline {
				next := new(numValue)
				ret = p.runInner(dashOrMinus, next, next).NewRune(r)
			} else {
				next := NewSequence(p.hist, p.hist.CurrentIndent(), func(vs []any) (_ error) {
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
		e := errutil.New("unexpected value, maybe you're missing string quotes?")
		ret = charm.Error(e)
	}
	return
}

// hack: starting a map in a collection,
// we expect the indent should be one more than the parent:
//   - Field:
//     Next: 5
//
// ... or...
// Field:
//
//	Next: 5
//
// but we don't want that extra indent for reading documents containing a single m
func (p *Value) mapIndent() int {
	var hack int
	if len(p.hist.els) > 1 {

		hack = 1
	}
	return p.hist.CurrentIndent() + hack
}

func (p *Value) runInner(r rune, i valueGetter, c charm.State) charm.State {
	p.inner = i
	return charm.RunState(r, c)
}

// if the passed rune might be start a bool value
// for example, `- trouble:` would match `- true` temporarily
// and `- false:` would match `- false` runtil the colon.
func (p *Value) tryBool(r rune, makeNext func(str string) charm.State) (ret charm.State) {
	var match string
	var res bool
	if r == 't' {
		match, res = "true", true
	} else if r == 'f' {
		match, res = "false", false
	}
	// a true parallel state would have been simpler in concept
	// but would need signaling out of the parallel to interrupt it
	// instead, this matches as much as it can and later re-runs whatever didnt match
	return charmed.MatchString(match, func(r rune, at int) (ret charm.State) {
		partial, matched := match[:at], at == len(match) && len(match) > 0
		if !matched {
			ret = makeNext(partial).NewRune(r)
		} else {
			// store the result early, in case we're at the end of the document.
			p.inner = computedValue{res}
			ret = charm.Statement("post bool", func(r rune) (ret charm.State) {
				// the word "true" or "false" needs to check the rune after it
				if r != Space && r != Newline {
					ret = makeNext(partial).NewRune(r)
				}
				return
			})
		}
		return
	}).NewRune(r)
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

type numValue struct{ charmed.NumParser }

// fix? returns float64 because json does
// could also return int64 when its int like
func (p *numValue) GetValue() (ret any, err error) {
	ret, err = p.GetFloat()
	return
}
