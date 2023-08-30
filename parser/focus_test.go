package parser_test

import (
	"testing"

	. "git.sr.ht/~ionous/tapestry/parser"
	"github.com/ionous/errutil"
	"github.com/ionous/sliceOf"
)

var dropGrammar = allOf(words("drop"), anyOf(
	allOf(&Focus{Where: "held", What: things()}, &Action{"Drop"}),
))

type MyContext struct {
	MyBounds // world
	Player   map[string]Bounds
	Other    map[string]Bounds
	Log
}

func (m MyContext) GetBounds(who, where string) (ret Bounds, err error) {
	switch who {
	case "":
		if s, ok := m.Player[where]; ok {
			m.Log.Log("asking for bounds", who, where)
			ret = s
		} else {
			ret = m.SearchBounds
		}
	default:
		if s, ok := m.Other[who]; ok {
			m.Log.Log("asking for bounds", who, where)
			ret = s
		} else {
			err = errutil.New("unknown bounds", who, where)
		}
	}
	return
}

func TestFocus(t *testing.T) {
	grammar := dropGrammar
	bounds := MyBounds{
		makeObject("red apple", "apples"),
		makeObject("apple cart", "carts"),
		makeObject("red cart", "carts"),
	}
	invBounds := MyBounds{
		makeObject("torch", "devices"),
		makeObject("crab apple", "apples"),
	}
	ctx := MyContext{
		Log:      t,
		MyBounds: bounds,
		Player:   map[string]Bounds{"held": invBounds.SearchBounds},
	}

	t.Run("drop one", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("drop"),
			&ClarifyGoal{"apple"},
			&ActionGoal{"Drop", sliceOf.String("crab apple")})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("drop all", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("drop everything"),
			&ActionGoal{"Drop", sliceOf.String("torch", "crab apple")})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("drop error", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("drop cart"),
			&ErrorGoal{"You can't see any such thing."})
		if e != nil {
			t.Fatal("expected an error")
		}
	})
}
