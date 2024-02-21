package jess_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/support/jesstest"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func TestPhrases(t *testing.T) {
	jesstest.RunPhraseTests(t, func(testPhrase string) (ret jess.Generator, err error) {
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
	kinds, traits, fields match.SpanList
	macros                jesstest.MacroList
}

func (n *info) GetContext() int {
	return 0
}
func (n *info) FindKind(ws match.Span) (string, int) {
	m, cnt := n.kinds.FindMatch(ws)
	return m.String(), cnt
}
func (n *info) FindTrait(ws match.Span) (string, int) {
	m, cnt := n.traits.FindMatch(ws)
	return m.String(), cnt
}
func (n *info) FindField(ws match.Span) (string, int) {
	m, cnt := n.fields.FindMatch(ws)
	return m.String(), cnt
}
func (n *info) FindMacro(ws match.Span) (mdl.Macro, int) {
	return n.macros.FindMacro(ws)
}

var known = info{
	macros: jesstest.PanicMacros(
		// source carries/ is carrying the targets
		// reverse would be: targets are carried by the source.
		"carried by", "carry", mdl.Macro_ManySecondary, true,
		"carrying", "carry", mdl.Macro_ManySecondary, false,
		// source contains the targets
		// the targets are in the source ( rhs macro )
		// in the source are the targets ( lhs macro; re-reversed )
		"in", "contain", mdl.Macro_ManySecondary, true,
		// source supports/is supporting the targets
		// so, "targets are on source" is reversed ( rhs macro )
		// and, "on source are targets" ( lhs macro; re-reversed )
		"on", "support", mdl.Macro_ManySecondary, true,
		//
		"suspicious of", "suspect", mdl.Macro_ManyMany, false,
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
	fields: match.PanicSpans(
		"description",
		"title",
		"age",
	),
}
