package parser_test

import (
	"testing"

	"github.com/ionous/sliceOf"
)

func TestReversal(t *testing.T) {
	grammar := showGrammar
	t.Run("forward", func(t *testing.T) {
		if e := parse(t, ctx, grammar,
			Phrases("show torch to bob"),
			&ActionGoal{"Show", sliceOf.String(
				"torch",
				"bob")}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("reverse", func(t *testing.T) {
		if e := parse(t, ctx, grammar,
			Phrases("show bob torch"),
			&ActionGoal{"Show", sliceOf.String(
				"torch",
				"bob")}); e != nil {
			t.Fatal(e)
		}
	})
}
