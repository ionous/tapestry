package parser_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/ident"
	. "git.sr.ht/~ionous/iffy/parser"
	"github.com/ionous/sliceOf"
)

var takeGrammar = allOf(
	Words("get"), anyOf(
		allOf(
			&Target{[]Scanner{things(), Words("from/off"), thing()}},
			&Action{"Remove"}),
		allOf(
			things(),
			&Action{"Take"}),
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
		Other: map[ident.Id]Bounds{
			ident.IdOf("apple-cart"): appleCart.SearchBounds,
			ident.IdOf("red-cart"):   redCart.SearchBounds},
	}

	t.Run("take from cart", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("get apple from/off red cart"),
			&ActionGoal{"Remove", sliceOf.String("red-cart", "yellow-apple")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("clarify from cart", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("get apple from/off cart"),
			&ClarifyGoal{"apple"},
			&ClarifyGoal{"red"},
			&ActionGoal{"Remove", sliceOf.String("apple-cart", "red-apple")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("exact take", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("get green"),
			&ActionGoal{"Take", sliceOf.String("green-apple")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("clarify take", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("get apple"),
			&ClarifyGoal{"cart"},
			&ActionGoal{"Take", sliceOf.String("apple-cart")})
		if e != nil {
			t.Fatal(e)
		}
	})
}
