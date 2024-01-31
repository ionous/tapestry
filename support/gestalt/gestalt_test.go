package gestalt_test

import (
	"reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/gestalt"
	j "git.sr.ht/~ionous/tapestry/support/gestalt"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/support/groktest"
	"github.com/kr/pretty"
)

func TestPhrases(t *testing.T) {
	// called can probably appear anywhere a noun can
	// so you will want a name and a reference ( flyweight i guess )
	called := &j.Sequence{[]j.Interpreter{
		&j.Kind{},
		&j.Words{Str: "called"},
		// TODO: optional determiner
		// and handle allow many (similar to traits... using comma-and rules
		&j.NounExactly{},
	}}
	test := &j.Sequence{[]j.Interpreter{
		called,
		&j.Is{},
		&j.Branch{[]j.Interpreter{
			// ... is open.
			// FIX: handle comma-and.
			&j.Traits{},
			// ... is on the .....
			// &j.Sequence{[]j.Interpreter{
			// 	&j.Verb{},
			// 	called,
			// }},
		}},
	}}

	str := `The container called the sarcophagus is open.`
	if ws, e := grok.MakeSpan(str); e != nil {
		t.Fatal(e)
	} else {
		q, in := j.MakeQuery(&known), j.MakeInputState(ws)
		if all := test.Match(q, []j.InputState{in}); len(all) != 1 {
			t.Fatal("expected a match")
		} else {
			n := all[0]
			res := gestalt.Reduce(n)
			m := groktest.ResultMap(res)
			if !reflect.DeepEqual(m, map[string]any{
				"primary": []map[string]any{{
					"det":    "the", // note: this is the bit closes to the noun
					"name":   "sarcophagus",
					"exact":  true,
					"kinds":  []string{"container"},
					"traits": []string{"open"},
				}},
			}) {
				t.Log(pretty.Sprint(m))
				t.Fatal("mismatched")
			} else {
				t.Log("ok")
			}

		}
	}
	// str2 := `The thing called the stake is on the supporter called the altar.`
	// if ws, e := grok.MakeSpan(str2); e != nil {
	// 	t.Fatal(e)
	// } else {
	// 	in := j.InputState{Words: ws, Query: &known}
	// 	if e := test.Match([]j.InputState{in}); e != nil {
	// 		t.Fatal(e)
	// 	}
	// }
}

type Word = grok.Word
type Span = grok.Span
type MacroType = grok.MacroType
type Macro = grok.Macro
type Match = grok.Match

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
		"carried by", "carry", grok.Macro_ManySecondary, true, //
		"carrying", "carry", grok.Macro_ManySecondary, false, //
		// source contains the targets
		// the targets are in the source ( rhs macro )
		// in the source are the targets ( lhs macro; re-reversed )
		"in", "contain", grok.Macro_ManySecondary, true,
		// kinds:
		"kinds of", "inherit", grok.Macro_PrimaryOnly, false, // for "are kinds of containers"
		"a kind of", "inherit", grok.Macro_PrimaryOnly, false, // for "a kind of container"
		// kind values
		"usually", "implies", grok.Macro_PrimaryOnly, false, // for "are usually closed"
		// source supports/is supporting the targets
		// so, "targets are on source" is reversed ( rhs macro )
		// and, "on source are targets" ( lhs macro; re-reversed )
		"on", "support", grok.Macro_ManySecondary, true,
		//
		"suspicious of", "suspect", grok.Macro_ManyMany, false,
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
