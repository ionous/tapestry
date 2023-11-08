package rift

import (
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/charmed"
)

// parses the "right hand side" of a collection or map
// assumes the next rune is at the start of the value: no leading whitespace.
type valueState struct {
	entry        *collectionEntry // values get written to a spot in a collection
	delayedValue valueFinalizer
}

type valueWriter interface {
	WriteValue(any) error
}

// helper to adapt the number parser to rift
type valueFinalizer interface {
	finalizeValue(valueWriter) error
}

func (p *valueState) finalizeValue() (err error) {
	if p.delayedValue != nil {
		err = p.delayedValue.finalizeValue(p.entry)
	}
	return
}

func (p *valueState) NewRune(r rune) (ret charm.State) {
	const dashOrMinus = Dash
	switch {
	case r == InterpretedString:
		ret = p.tryString(r, true)

	case r == RawString:
		ret = p.tryString(r, false)

	case charmed.IsNumber(r) || r == '+':
		ret = p.tryNum().NewRune(r)

	case unicode.IsLetter(r):
		// might be a mapping, or might be a bool literal.
		ret = p.tryBool(r, func(partial string) charm.State {
			next := p.tryMapping()
			for _, r := range partial {
				next = next.NewRune(r)
			}
			return next
		})

	case r == dashOrMinus:
		// negative numbers or sequences
		ret = charm.Statement("dashing", func(r rune) (ret charm.State) {
			// no space indicates a number `-5`
			if r != Space && r != Newline {
				ret = p.tryNum()
			} else {
				// a space indicates a sequence `- 5`
				doc := p.entry.Document()
				ret = p.trySequence(doc.Col - 1)
			}
			// send the dash and the new character along to the num or seq
			return ret.NewRune(dashOrMinus).NewRune(r)
		})
	}
	return
}

func (p *valueState) tryString(r rune, interpreted bool) charm.State {
	return charmed.ScanQuote(r, interpreted, func(res string) {
		p.entry.WriteValue(res)
	})
}

func (p *valueState) tryNum() charm.State {
	num := new(numValue)
	p.delayedValue = num
	return num
}

func (p *valueState) trySequence(depth int) (ret charm.State) {
	if header, e := p.entry.flushHeader(); e != nil {
		ret = charm.Error(e)
	} else {
		p.delayedValue = nil
		ret = NewSequence(p.entry, header, depth)
	}
	return
}

func (p *valueState) tryMapping() (ret charm.State) {
	if header, e := p.entry.flushHeader(); e != nil {
		ret = charm.Error(e)
	} else {
		doc := p.entry.Document()
		// hack: starting a map in a collection,
		// we expect the indent should be greater than the parent.
		//   - Field:
		//     Next: 5
		//
		// ... or...
		// Field:
		//
		//	Next: 5
		//
		// but we don't want that extra indent for reading documents containing a single value
		// tbd: after a key we should expect some amount of spaces,
		// maybe we can generate / pass in padding from the parent
		hack := doc.Col
		if len(doc.els) > 2 {
			hack += 1
		}
		p.delayedValue = nil
		ret = NewMapping(p.entry, header, hack)
	}
	return
}

// if the passed rune might be start a bool value
// for example, `- trouble:` would match `- true` temporarily
// and `- false:` would match `- false` runtil the colon.
func (p *valueState) tryBool(r rune, makeNext func(str string) charm.State) (ret charm.State) {
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
			// fix: for this and number we should require a space or eol after
			// to complete the value; *clear* some value on success rather setting some ambient state
			// and error on pop if the value wasnt cleared
			p.delayedValue = computedValue{res}
			ret = charm.Statement("post bool", func(r rune) (ret charm.State) {
				// the word "true" or "false" needs to check the rune after it
				// if it doesn't then jump to the next state
				if r != Space && r != Newline {
					ret = makeNext(partial).NewRune(r)
				}
				return
			})
		}
		return
	}).NewRune(r)
}

// a final value, ex. from a boolean.
type computedValue struct{ v any }

func (p computedValue) finalizeValue(tgt valueWriter) error {
	return tgt.WriteValue(p.v)
}

// a number --
// note this is a little different than the other types
// because there's no terminal value for it.
type numValue struct{ charmed.NumParser }

// fix? returns float64 because json does
// could also return int64 when its int like
func (p *numValue) finalizeValue(tgt valueWriter) (err error) {
	if n, e := p.GetFloat(); e != nil {
		err = e
	} else {
		err = tgt.WriteValue(n)
	}
	return
}
