package parser_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/parser"
	"github.com/ionous/sliceOf"
)

func TestParser(t *testing.T) {
	grammar := lookGrammar

	t.Run("look", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("look/l"),
			&ActionGoal{"Look", nil})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("examine", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("look/l at something"),
			&ActionGoal{
				"Examine", sliceOf.String("something"),
			})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("search", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("look/l inside/in/into/through/on something"),
			&ActionGoal{
				"Search", sliceOf.String("something"),
			})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("look under", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("look/l under something"),
			&ActionGoal{
				"LookUnder", sliceOf.String("something"),
			})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("look dir", func(t *testing.T) {
		look := Phrases("look/l")
		for _, d := range directions {
			d := sliceOf.String(d)
			if e := parse(t, ctx, grammar,
				permute(look, d),
				&ActionGoal{"Examine", d}); e != nil {
				t.Fatal(e)
				break
			}
		}
	})
	// note: changed the errors to return the last one of the deepest depth
	// rather than the first one. most of the time that seems to make more sense.
	t.Run("look something", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("look something"),
			// &ErrorGoal{"too many words"})
			&ErrorGoal{parser.MismatchedWord{"under", "something", 1}.Error()})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("look to dir", func(t *testing.T) {
		lookTo := Phrases("look/l to")
		for _, d := range directions {
			d := sliceOf.String(d)
			if e := parse(t, ctx, grammar,
				permute(lookTo, d),
				&ActionGoal{"Examine", d}); e != nil {
				t.Fatal(e)
				break
			}
		}
	})
}
