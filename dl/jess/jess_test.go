package jess_test

import (
	"errors"
	"fmt"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/support/groktest"
)

func TestTraits(t *testing.T) {
	groktest.RunTraitTests(t, func(testPhrase string) (ret grok.TraitSet, err error) {
		t.Log("testing:", testPhrase)
		if ws, e := grok.MakeSpan(testPhrase); e != nil {
			err = e
		} else {
			var t jess.TraitsKind //
			input := jess.InputState(ws)
			if !t.Match(jess.MakeQuery(&known), &input) {
				err = errors.New("failed to match traits")
			} else if cnt := len(input); cnt != 0 {
				err = fmt.Errorf("partially matched %d words", len(ws)-cnt)
			} else {
				ret = t.GetTraitSet()
			}
		}
		return
	})
}

func TestPhrases(t *testing.T) {
	groktest.RunPhraseTests(t, func(testPhrase string) (ret grok.Results, err error) {
		t.Log("testing:", testPhrase)
		if ws, e := grok.MakeSpan(testPhrase); e != nil {
			err = e
		} else {
			ret, err = jess.Match(&known, ws)
		}
		return
	})
}

type Word = grok.Word
type Span = grok.Span
type MacroType = grok.MacroType
type Macro = grok.Macro
type Match = grok.Matched

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

// "kinds of", "inherit", grok.Macro_PrimaryOnly, false
// "a kind of", "inherit", grok.Macro_PrimaryOnly, false

var known = info{
	macros: groktest.PanicMacros(
		// source carries/ is carrying the targets
		// reverse would be: targets are carried by the source.
		"carried by", "carry", grok.Macro_ManySecondary, true,
		"carrying", "carry", grok.Macro_ManySecondary, false,
		// source contains the targets
		// the targets are in the source ( rhs macro )
		// in the source are the targets ( lhs macro; re-reversed )
		"in", "contain", grok.Macro_ManySecondary, true,
		// source supports/is supporting the targets
		// so, "targets are on source" is reversed ( rhs macro )
		// and, "on source are targets" ( lhs macro; re-reversed )
		"on", "support", grok.Macro_ManySecondary, true,
		//
		"suspicious of", "suspect", grok.Macro_ManyMany, false,
		// kinds: should these really be a macro?
		"kinds of", "inherit", grok.Macro_PrimaryOnly, false, // for "are kinds of containers"
		"a kind of", "inherit", grok.Macro_PrimaryOnly, false, // for "a kind of container"
		// fix: should these really be a macro?
		"usually", "implies", grok.Macro_PrimaryOnly, false,
	),
	kinds: grok.PanicSpans(
		"kind", "kinds",
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
