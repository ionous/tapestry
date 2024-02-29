package jess_test

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/jesstest"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func TestPhrases(t *testing.T) {
	jesstest.TestPhrases(t, &known)
}

type info struct {
	kinds                 []string
	traits, fields, nouns match.SpanList
	macros                jesstest.MacroList
}

func (n *info) GetContext() int {
	return 0
}
func (n *info) FindKind(ws match.Span, out *kindsOf.Kinds) (ret string, width int) {
	str := strings.ToLower(ws.String())
	for i, k := range n.kinds {
		if strings.HasPrefix(str, k) {
			if i&1 == 0 { // singular are the even numbers
				k = n.kinds[i+1]
			}
			ret, width = k, 1 // always in the tests
			if out != nil {
				// hacks for testing
				k := kindsOf.Kind
				switch ret {
				case "color":
					k = kindsOf.Aspect
				case "groups":
					k = kindsOf.Record
				case "storing":
					k = kindsOf.Action
				}
				*out = k
			}
			break
		}
	}
	return
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
func (n *info) FindNoun(ws match.Span) (string, int) {
	m, cnt := n.nouns.FindMatch(ws)
	return m.String(), cnt
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
	kinds: []string{
		"kind", "kinds",
		"thing", "things",
		"container", "containers",
		"supporter", "supporters",
		"aspect", "aspects",
		"color", "color", // aspect; uses the singular
		"group", "groups", // record
		"storing", "storing", //
	},
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
	nouns: match.PanicSpans(
		"message",
		"missive",
	),
}
