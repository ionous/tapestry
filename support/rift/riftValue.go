package rift

import (
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/charmed"
)

// parses the "right hand side" of a collection or map
// assumes the next rune is at the start of the value: no leading whitespace.
func NewValue(ent *riftEntry) charm.State {
	val := riftValue{ent}
	return charm.Self("value", func(self charm.State, r rune) (ret charm.State) {
		const dashOrMinus = Dash
		switch {
		case r == InterpretedString:
			ret = val.newString(r, true)

		case r == RawString:
			ret = val.newString(r, false)

		case charmed.IsNumber(r) || r == '+':
			ret = val.newNum().NewRune(r)

		case unicode.IsLetter(r):
			// might be a mapping, or might be a bool literal.
			ret = val.tryBool(r, func(partial string) charm.State {
				next := val.newMapping()
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
					ret = val.newNum()
				} else {
					// a space indicates a sequence `- 5`
					ret = val.newSequence(ent.doc.Col - 1)
				}
				// send the dash and the new character along to the num or seq
				return ret.NewRune(dashOrMinus).NewRune(r)
			})
		}
		return
	})
}

type riftValue struct {
	entry *riftEntry // values get written into the entry.
}

// fix?: can be set multiple times
// ex. "true" as a bool might be set until "true:" as a mapping is detected
// would be nicer to require whitespace after values.
// ( and require newline at the end of a file )
func (val *riftValue) setPendingValue(p pendingValue) {
	val.entry.pendingValue = p
}

func (val *riftValue) setComputedValue(newValue any) {
	val.entry.pendingValue = computedValue{newValue}
}

func (val *riftValue) newString(r rune, interpreted bool) charm.State {
	return charmed.ScanQuote(r, interpreted, func(res string) {
		val.setComputedValue(res)
	})
}

func (val *riftValue) newNum() charm.State {
	num := new(numValue)
	val.setPendingValue(num)
	return num
}

func (val *riftValue) newSequence(depth int) (ret charm.State) {
	if header, e := val.entry.writeHeader(); e != nil {
		ret = charm.Error(e)
	} else {
		sub := NewSequence(val.entry.doc, header, depth)
		val.setPendingValue(sub)
		ret = StartSequence(sub)
	}
	return
}

func (val *riftValue) newMapping() (ret charm.State) {
	if header, e := val.entry.writeHeader(); e != nil {
		ret = charm.Error(e)
	} else {
		doc := val.entry.doc
		sub := NewMapping(doc, header, doc.Col)
		val.setPendingValue(sub)
		ret = StartMapping(sub)
	}
	return
}

// if the passed rune might be start a bool value
// for example, `- trouble:` would match `- true` temporarily
// and `- false:` would match `- false` runtil the colon.
func (val *riftValue) tryBool(r rune, makeNext func(str string) charm.State) (ret charm.State) {
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
			val.setComputedValue(res)
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
