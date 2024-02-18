package jess_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/support/jesstest"
	"git.sr.ht/~ionous/tapestry/support/match"
)

func TestPhrases(t *testing.T) {
	jesstest.RunPhraseTests(t, func(testPhrase string) (ret jess.Applicant, err error) {
		t.Log("testing:", testPhrase)
		if ws, e := match.MakeSpan(testPhrase); e != nil {
			err = e
		} else {
			ret, err = jess.Match(&known, ws)
		}
		return
	})
}

type info struct {
	kinds, traits match.SpanList
	macros        jesstest.MacroList
}

func (n *info) FindKind(ws match.Span) (jess.Matched, int) {
	return n.kinds.FindMatch(ws)
}

func (n *info) FindTrait(ws match.Span) (jess.Matched, int) {
	return n.traits.FindMatch(ws)
}

func (n *info) FindMacro(ws match.Span) (jess.Macro, int) {
	return n.macros.FindMacro(ws)
}

var known = info{
	macros: jesstest.PanicMacros(
		// source carries/ is carrying the targets
		// reverse would be: targets are carried by the source.
		"carried by", "carry", jess.Macro_ManySecondary, true,
		"carrying", "carry", jess.Macro_ManySecondary, false,
		// source contains the targets
		// the targets are in the source ( rhs macro )
		// in the source are the targets ( lhs macro; re-reversed )
		"in", "contain", jess.Macro_ManySecondary, true,
		// source supports/is supporting the targets
		// so, "targets are on source" is reversed ( rhs macro )
		// and, "on source are targets" ( lhs macro; re-reversed )
		"on", "support", jess.Macro_ManySecondary, true,
		//
		"suspicious of", "suspect", jess.Macro_ManyMany, false,
	),
	kinds: match.PanicSpans(
		"kind", "kinds",
		"thing", "things",
		"container", "containers",
		"supporter", "supporters",
	),
	traits: match.PanicSpans(
		"closed",
		"open",
		"openable",
		"transparent",
		"fixed in place",
	),
}
