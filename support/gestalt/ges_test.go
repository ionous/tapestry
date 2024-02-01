package gestalt_test

import (
	"fmt"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/gestalt"
	j "git.sr.ht/~ionous/tapestry/support/gestalt"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/support/groktest"
)

func TestPhrases(t *testing.T) {
	tree := &j.Branch{[]j.Interpreter{
		// [the] <name> a kind of <kind>
		&j.Sequence{[]j.Interpreter{
			&j.Name{},
			&j.Are{},
			&j.Verb{Str: "a kind of", Action: "inherit"},
			&j.Article{Optional: true, Ignore: true},
			&j.Kind{},
		}},

		// [the] <kind> are usually <traits>
		&j.Sequence{[]j.Interpreter{
			&j.Article{Optional: true, Ignore: true},
			&j.Kind{},
			&j.Are{},
			&j.Verb{Str: "usually", Action: "implies"},
			&j.Traits{},
		}},

		// <num> <kind> are <location> <noun>
		&j.Sequence{[]j.Interpreter{
			&j.Target{Primary: false, Group: []j.Interpreter{
				&j.Count{},
				&j.Kind{},
			}},
			&j.Are{},
			&j.Target{Primary: true, Group: []j.Interpreter{
				&j.Verb{}, // tbd: limit to placement?
				&j.Name{}, // tbd: would a "noun exists" be better here?
			}},
		}},

		// [the] <kind> called <noun> is traits
		&j.Sequence{[]j.Interpreter{
			&j.Article{Optional: true, Ignore: true},
			&j.Kind{},   // kinds allow optional determiners
			&j.Called{}, // optionally
			&j.Name{},   // names allow optional determiners
			&j.Are{},
			&j.Traits{},
		}},

		// &j.Branch{[]j.Interpreter{
		// 	// ... is open.
		// 	// FIX: handle comma-and.
		// 	// &j.Traits{},
		// 	// // ... is on the .....
		// 	// &j.Sequence{[]j.Interpreter{
		// 	// 	&j.Verb{},
		// 	// 	called,
		// 	// }},
		// }},

	}}
	//
	groktest.RunTests(t, func(str string) (res grok.Results, err error) {
		if ws, e := grok.MakeSpan(str); e != nil {
			err = e
		} else {
			q, in := j.MakeQuery(&known), j.MakeInputState(ws)
			if all := tree.Match(q, []j.InputState{in}); len(all) != 1 {
				err = fmt.Errorf("have %d matches", len(all))
			} else {
				res, err = gestalt.Reduce(all[0])
			}
			return
		}
		return
	})
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

// "kinds of", "inherit", grok.Macro_PrimaryOnly, false
// "a kind of", "inherit", grok.Macro_PrimaryOnly, false
// "usually", "implies", grok.Macro_PrimaryOnly, false
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
