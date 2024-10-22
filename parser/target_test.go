package parser_test

import (
	"testing"

	. "git.sr.ht/~ionous/tapestry/parser"
	"github.com/ionous/sliceOf"
)

var takeGrammar = allOf(
	words("get"), anyOf(
		allOf(
			&Refine{[]Scanner{things(), words("from", "off"), thing()}},
			&Action{Name: "Remove"}),
		allOf(
			things(),
			&Action{Name: "Take"}),
	),
)

func TestTarget(t *testing.T) {
	grammar := takeGrammar
	bounds := MyBounds{
		makeObject("green apple", "apples"),
		makeObject("apple cart", "carts"),
		makeObject("red cart", "carts"),
	}
	appleCart := MyBounds{
		makeObject("crab apple", "apples"),
		makeObject("red apple", "apples"),
	}
	redCart := MyBounds{
		makeObject("yellow apple", "apples"),
	}
	//
	ctx := MyContext{
		Log:      t,
		MyBounds: bounds,
		Other: map[string]Bounds{
			"apple cart": appleCart.SearchBounds,
			"red cart":   redCart.SearchBounds},
	}

	t.Run("take from cart", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("get apple from/off red cart"),
			&ActionGoal{"Remove", sliceOf.String("red cart", "yellow apple")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("clarify from cart", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("get apple from/off cart"),
			&ClarifyGoal{"apple"},
			&ClarifyGoal{"red"},
			&ActionGoal{"Remove", sliceOf.String("apple cart", "red apple")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("exact take", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("get green"),
			&ActionGoal{"Take", sliceOf.String("green apple")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("clarify take", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("get apple"),
			&ClarifyGoal{"cart"},
			&ActionGoal{"Take", sliceOf.String("apple cart")})
		if e != nil {
			t.Fatal(e)
		}
	})
}
