package grok_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/support/groktest"
)

type Word = grok.Word
type Span = grok.Span
type MacroType = grok.MacroType
type MacroInfo = grok.MacroInfo
type Match = grok.Match

func TestPhrases(t *testing.T) {
	groktest.Phrases(t, &known)
}

func TestTraits(t *testing.T) {
	groktest.Phrases(t, &known)
}

type info struct {
	determiners, kinds, traits groktest.SpanList
	macros                     groktest.MacroList
}

func (n *info) FindDeterminer(ws []Word) Match {
	return n.determiners.FindMatch(ws)
}

func (n *info) FindKind(ws []Word) Match {
	return n.kinds.FindMatch(ws)
}

func (n *info) FindTrait(ws []Word) Match {
	return n.traits.FindMatch(ws)
}

func (n *info) FindMacro(ws []Word) (MacroInfo, bool) {
	return n.macros.FindMacro(ws)
}

var known = info{
	determiners: groktest.PanicSpans(
		"the", "a", "an", "some", "our",
		// ex. kettle of fish
		"a kettle of",
	),
	macros: groktest.PanicMacros(
		grok.ManyToOne, "kind of", // for "a closed kind of container"
		grok.ManyToOne, "kinds of", // for "are closed containers"
		grok.ManyToOne, "a kind of", // for "a kind of container"
		// tbd: flags need more thought.
		grok.OneToMany, "carrying", //
		// other macros
		grok.OneToMany, "on", // on the x are the w,y,z
		grok.OneToMany, "in",
		//
		grok.ManyToMany, "suspicious of",
	),
	kinds: groktest.PanicSpans(
		"thing", "things",
		"container", "containers",
		"supporter", "supporters",
	),
	traits: groktest.PanicSpans(
		"closed",
		"open",
		"openable",
		"transparent",
		"fixed in place",
	),
}
