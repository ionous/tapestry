package grok

import (
	"strings"
	"testing"

	"github.com/kr/pretty"
)

var phrases = []struct {
	test   string
	result map[string]any
	skip   bool
}{
	// simple trait:
	{test: `The bottle is closed.`},
	// multi word trait:
	{test: `The tree is fixed in place.`},
	// a kind of declaration ( uses a 'macro' verb )
	{
		test: `The box is a kind of container.`,
		result: map[string]any{
			"macro": "a kind of",
			"subjects": []map[string]any{{
				"det":    "The", // uppercase because its the real value from the original string.
				"name":   "box",
				"traits": []string{"container"},
			}},
		},
	},
	// kind of: adding trailing properties
	{
		test: `The box is a kind of closed container.`,
		result: map[string]any{
			"macro": "a kind of",
			"subjects": []map[string]any{{
				"det":    "The",
				"name":   "box",
				"traits": []string{"closed", "container"},
			}},
		},
	},
	// kind of: adding middle properties
	{
		test: `The box is a closed kind of container.`,
		result: map[string]any{
			"macro": "kind of",
			"subjects": []map[string]any{{
				"det":    "The",
				"name":   "box",
				"traits": []string{"closed", "container"},
			}},
		},
	},
	// kind of, "correctly" failing prefixed properties
	// note: in inform, this also yields a noun named the "closed box".
	{
		test: `The closed box is a kind of container.`,
		result: map[string]any{
			"macro": "a kind of",
			"subjects": []map[string]any{{
				"det":    "The",
				"name":   "closed box",
				"traits": []string{"container"},
			}},
		},
	},
	// multiple kinds of things
	// fix: add multiple subject macro'd tests too
	{
		test: `The box and the top are closed kinds of things.`,
		result: map[string]any{
			"macro": "kinds of",
			"subjects": []map[string]any{{
				"det":    "The",
				"name":   "box",
				"traits": []string{"closed", "things"},
			}, {
				"det":    "the",
				"name":   "top",
				"traits": []string{"closed", "things"},
			}},
		},
	},
	// multiple trailing properties, using the kind as a property.
	{test: `The bottle is a transparent, open, container.`},
	// add leading properties in the lede by using 'called' ( 'in' is also a macro verb )
	{test: `The closed openable container called the trunk is in the lobby.`},
	// using 'called' without a macro
	{test: `The container called the sarcophagus is open.`},
	// a leading macro
	{test: `In the coffin are some coins, a notebook, and a gripping hand.`},
	// same pattern as the middle properties above; but not using kind of
	{test: `The coffin is a closed container in the antechamber.`},
	// giving properties to the rhs target:
	// note: this doesn't make him human? and, no, an anonymous man doesn't work here.
	{test: `The unhappy man is in the closed bottle.`},
	// trailing properties applying to the lhs
	{test: `The bottle in the kitchen is openable and has age 42.`},
	// note, this is allowed even though it implies something different than what is written:
	{test: `The bottle is openable in the kitchen.`},
	// called both before and after the macro
	{
		test: `The thing called the stake is on the supporter called the bookshelf.`,
		result: map[string]any{
			"macro": "on",
			"subjects": []map[string]any{{
				"det":    "the",
				"name":   "stake",
				"traits": []string{"thing"},
			}},
			"targets": []map[string]any{{
				"det":    "the",
				"name":   "bookshelf",
				"traits": []string{"supporter"},
			}},
		},
	},
	// the special nxn description: no properties are allowed.
	{test: `Hector and Maria are suspicious of Santa and Santana.`},
	// fix: in this case i think inform eats the first "is" and allows a subsequent one ( and also allows values )
	// fix: the device called the detonator is on the supporter called the shelf and is proper named"
}

func TestPhrases(t *testing.T) {
	var skipped int
	for i, p := range phrases {
		if p.result == nil {
			continue
		} else if p.skip {
			skipped++
		} else if res, e := Grok(p.test); e != nil {
			t.Fatal(e)
		} else {
			m := resultMap(res)
			if d := pretty.Diff(p.result, m); len(d) > 0 {
				t.Log(i, "got:", pretty.Sprint(m))
				//t.Log("want:", pretty.Sprint(p.result))
				t.Fatal(d)
			}
		}
	}
	if skipped > 0 {
		t.Fatal("skipped %d", skipped)
	}
}

func TestTraits(t *testing.T) {
	// fix merge []word and the original string? along with context?
	traits, e := parseTraits(panicHash("the open container"))
	if e != nil {
		t.Fatal(e)
	} else if cnt := len(traits); cnt != 2 {
		t.Fatal(cnt)
	} else {
		if t0, t1 := wordsToString(traits[0]), wordsToString(traits[1]); t0 != "open" || t1 != "container" {
			t.Fatal(t0, t1)
		}
	}
}

func resultMap(in results) map[string]any {
	m := make(map[string]any)
	nounsIntoMap(m, "subjects", in.subjects)
	nounsIntoMap(m, "targets", in.targets)
	wordsIntoMap(m, "macro", in.macro)
	return m
}

func nounsIntoMap(m map[string]any, field string, ns []noun) {
	if len(ns) > 0 {
		out := make([]map[string]any, len(ns))
		for i, n := range ns {
			out[i] = nounToMap(n)
		}
		m[field] = out
	}
}

func nounToMap(n noun) map[string]any {
	m := make(map[string]any)
	wordsIntoMap(m, "name", n.name)
	wordsIntoMap(m, "det", n.det)
	wordListIntoMap(m, "traits", n.traits)
	return m
}

func wordListIntoMap(m map[string]any, field string, ws [][]word) {
	if cnt := len(ws); cnt > 0 {
		out := make([]string, cnt)
		for i, w := range ws {
			out[i] = wordsToString(w)
		}
		m[field] = out
	}
	return
}

func wordsIntoMap(m map[string]any, field string, w []word) {
	if len(w) > 0 {
		m[field] = wordsToString(w)
	}
}

func wordsToString(w []word) (ret string) {
	var b strings.Builder
	for i, w := range w {
		if i > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(w.String())
	}
	return b.String()
}
