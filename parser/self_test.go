package parser_test

import (
	"testing"

	"github.com/ionous/sliceOf"
)

func TestSelf(t *testing.T) {
	grammar := showGrammar
	if e := parse(t, ctx, grammar,
		Phrases("show torch to self"),
		&ActionGoal{"Examine", sliceOf.String(
			"torch")}); e != nil {
		t.Fatal(e)
	}
}
