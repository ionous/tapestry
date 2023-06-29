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
	macros                     macroList
}

type macroList struct {
	groktest.SpanList
	types []MacroType
}

func (ml macroList) get(i int) ([]Word, MacroType) {
	return ml.SpanList[i], ml.types[i]
}

func panicMacros(pairs ...any) (out macroList) {
	cnt := len(pairs) / 2
	out.SpanList = make(groktest.SpanList, cnt)
	out.types = make([]MacroType, cnt)
	for i := 0; i < cnt; i++ {
		out.types[i] = pairs[i*2+0].(MacroType)
		out.SpanList[i] = groktest.PanicSpan(pairs[i*2+1].(string))
	}
	return
}

func (n *info) FindDeterminer(ws []Word) (ret Match) {
	return n.determiners.FindMatch(ws)
}

func (n *info) FindKind(ws []Word) (ret Match) {
	return n.kinds.FindMatch(ws)
}

func (n *info) FindTrait(ws []Word) (ret Match) {
	return n.traits.FindMatch(ws)
}

func (n *info) FindMacro(ws []Word) (ret MacroInfo, okay bool) {
	if at, skip := n.macros.FindPrefix(ws); skip > 0 {
		w, t := n.macros.get(at)
		ret = MacroInfo{
			Match: Span(w[:skip]),
			Type:  t,
		}
		okay = true
	}
	return
}

var known = info{
	determiners: groktest.PanicSpans(
		"the", "a", "an", "some", "our",
		// ex. kettle of fish
		"a kettle of",
	),
	macros: panicMacros(
		// tbd: flags need more thought.
		grok.ManyToOne, "kind of", // for "a closed kind of container"
		grok.ManyToOne, "kinds of", // for "are closed containers"
		grok.ManyToOne, "a kind of", // for "a kind of container"
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
