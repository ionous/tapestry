package grok_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/support/groktest"
)

type Word = grok.Word
type Span = grok.Span
type MacroType = grok.MacroType
type Macro = grok.Macro
type Match = grok.Match

func TestPhrases(t *testing.T) {
	groktest.Phrases(t, &known)
}

func TestTraits(t *testing.T) {
	groktest.Phrases(t, &known)
}

type info struct {
	kinds, traits grok.SpanList
	macros        groktest.MacroList
}

func (n *info) FindArticle(ws Span) (ret grok.Article, err error) {
	return grok.FindArticle(ws)
}

func (n *info) FindKind(ws Span) (Match, error) {
	return n.kinds.FindMatch(ws)
}

func (n *info) FindTrait(ws Span) (Match, error) {
	return n.traits.FindMatch(ws)
}

func (n *info) FindMacro(ws Span) (Macro, error) {
	return n.macros.FindMacro(ws)
}

var known = info{
	macros: groktest.PanicMacros(
		// source carries/ is carrying the targets
		// reverse would be: targets are carried by the source.
		"carried by", "carry", grok.Macro_ManyTargets, true, //
		"carrying", "carry", grok.Macro_ManyTargets, false, //
		// source contains the targets
		// the targets are in the source ( rhs macro )
		// in the source are the targets ( lhs macro; re-reversed )
		"in", "contain", grok.Macro_ManyTargets, true,
		// kinds:
		"kind of", "inherit", grok.Macro_SourcesOnly, false, // for "a closed kind of container"
		"kinds of", "inherit", grok.Macro_SourcesOnly, false, // for "are closed containers"
		"a kind of", "inherit", grok.Macro_SourcesOnly, false, // for "a kind of container"
		// source supports/is supporting the targets
		// so, "targets are on source" is reversed ( rhs macro )
		// and, "on source are targets" ( lhs macro; re-reversed )
		"on", "support", grok.Macro_ManyTargets, true,
		//
		"suspicious of", "suspect", grok.Macro_ManyMany, false,
	),
	kinds: grok.PanicSpans(
		"thing", "things",
		"container", "containers",
		"supporter", "supporters",
	),
	traits: grok.PanicSpans(
		"closed",
		"open",
		"openable",
		"transparent",
		"fixed in place",
	),
}
