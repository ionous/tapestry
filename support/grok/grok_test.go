package grok

import (
	"strings"
	"testing"

	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func TestPhrases(t *testing.T) {
	var phrases = []struct {
		test   string
		result any
		skip   any
	}{
		// simple trait:
		{
			test: `The bottle is closed.`,
			result: map[string]any{
				"sources": []map[string]any{{
					"det":    "The", // uppercase because its the real value from the original string.
					"name":   "bottle",
					"traits": []string{"closed"},
				}},
			},
		},
		// multi Word trait:
		{
			test: `The tree is fixed in place.`,
			result: map[string]any{
				"sources": []map[string]any{{
					"det":    "The",
					"name":   "tree",
					"traits": []string{"fixed in place"},
				}},
			},
		},
		// multiple trailing properties, using the kind as a property.
		{
			test: `The bottle is a transparent, open, container.`,
			result: map[string]any{
				"sources": []map[string]any{{
					"det":    "The",
					"name":   "bottle",
					"kinds":  []string{"container"},
					"traits": []string{"transparent", "open"},
				}},
			},
		},
		// using 'called' without a macro
		{
			test: `The container called the sarcophagus is open.`,
			result: map[string]any{
				"sources": []map[string]any{{
					"det":    "the", // lowercase, its the bit closet to the Noun
					"name":   "sarcophagus",
					"kinds":  []string{"container"},
					"traits": []string{"open"},
				}},
			},
		},
		// a kind of declaration ( uses a 'macro' verb )
		// "is" left of macro
		{
			test: `The box is a kind of container.`,
			result: map[string]any{
				"macro": "a kind of",
				"sources": []map[string]any{{
					"det":   "The",
					"name":  "box",
					"kinds": []string{"container"},
				}},
			},
		},
		// kind of: adding trailing properties
		// "is" left of macro
		{
			test: `The box is a kind of closed container.`,
			result: map[string]any{
				"macro": "a kind of",
				"sources": []map[string]any{{
					"det":    "The",
					"name":   "box",
					"traits": []string{"closed"},
					"kinds":  []string{"container"},
				}},
			},
		},

		// kind of, "correctly" failing prefixed properties.
		// note: in inform, this also yields a Noun named the "closed box".
		// similarly, The kind of the box is a container, yields a name "kind of the box".
		// "is" left of macro.
		{
			test: `The closed box is a kind of container.`,
			result: map[string]any{
				"macro": "a kind of",
				"sources": []map[string]any{{
					"det":   "The",
					"name":  "closed box",
					"kinds": []string{"container"},
				}},
			},
		},
		// multiple kinds of things
		{
			// 	test: `The box and the top are closed kinds of containers.`,
			// 	result:   map[string]any{
			// 		"macro": "kinds of",
			// 		"sources": []map[string]any{{
			// 			"det":    "The",
			// 			"name":   "box",
			// 			"traits": []string{"closed"},
			// 		}, {
			// 			"det":    "the",
			// 			"name":   "top",
			// 			"traits": []string{"closed"},
			// 			"kinds":  []string{"things"},
			// 		}},
			// 	},
		},
		// kind of: adding middle properties
		// "is" left of macro
		// not allowed in inform.
		{
			// 	test: `The box is a closed kind of container.`,
			// 	result:   map[string]any{
			// 		"macro": "kind of",
			// 		"sources": []map[string]any{{
			// 			"det":    "The",
			// 			"name":   "box",
			// 			"traits": []string{"closed", "container"},
			// 		}},
			// 	},
		},
		{
			//test: `A device is in the lobby.`,
			//result:   errutil.New("this is specifically disallowed, and should generate an error"),
		},
		// giving properties to the rhs and right targets isnt permitted:
		// tbd: but it might be possible...
		{
			test: `The unhappy man is in the closed bottle.`,
			result: map[string]any{
				"macro": "in",
				"sources": []map[string]any{{
					"det":  "The",
					"name": "unhappy man",
				}},
				"targets": []map[string]any{{
					"det":  "the",
					"name": "closed bottle",
				}},
			}},
		// same pattern as the middle properties above; but not using kind of
		{
			test: `The coffin is a closed container in the antechamber.`,
			result: map[string]any{
				"macro": "in",
				"sources": []map[string]any{{
					"det":    "The",
					"name":   "coffin",
					"traits": []string{"closed"},
					"kinds":  []string{"container"},
				}},
				"targets": []map[string]any{{
					"det":  "the",
					"name": "antechamber",
				}},
			},
		},

		// trailing properties applying to the lhs
		{
			//test: `The bottle in the kitchen is openable and has age 42.`,
			// FIX: values.
		},
		// note, this is allowed even though it implies something different than what is written:
		{
			test: `The bottle is openable in the kitchen.`,
			result: map[string]any{
				"macro": "in",
				"sources": []map[string]any{{
					"det":    "The",
					"traits": []string{"openable"},
					"name":   "bottle",
				}},
				"targets": []map[string]any{{
					"det":  "the",
					"name": "kitchen",
				}},
			},
		},
		// called both before and after the macro
		// FIX: ITS GIE THE STAKE BOTH KINDS AND THE ALTAR DOESNT EXIST
		{
			test: `The thing called the stake is on the supporter called the altar.`,
			result: map[string]any{
				"macro": "on",
				"sources": []map[string]any{{
					"det":   "the",
					"name":  "stake",
					"kinds": []string{"thing"},
				}},
				"targets": []map[string]any{{
					"det":   "the",
					"name":  "altar",
					"kinds": []string{"supporter"},
				}},
			},
		},
		{
			// FIX: add multiple subject macro'd tests
		},
		// add leading properties in the lede by using 'called'
		// "is" left of the macro "in".
		// slightly different parsing than "kind/s of":
		// those expect only expect one set of nouns; these have two.
		{
			test: `The closed openable container called the trunk is in the lobby.`,
			result: map[string]any{
				"macro": "in",
				"sources": []map[string]any{{
					"det":    "the", // lowercase, the closest to the trunk
					"name":   "trunk",
					"traits": []string{"closed", "openable"},
					"kinds":  []string{"container"},
				}},
				"targets": []map[string]any{{
					"det":  "the",
					"name": "lobby",
				}},
			},
		},
		// a leading macro
		{
			test: `In the coffin are some coins, a notebook, and the gripping hand.`,
			result: map[string]any{
				"macro": "in",
				"targets": []map[string]any{{
					"det":  "the", // lowercase, the closest to the trunk
					"name": "coffin",
				}},
				"sources": []map[string]any{{
					"det":  "some",
					"name": "coins",
				}, {
					"det":  "a",
					"name": "notebook",
				}, {
					"det":  "the",
					"name": "gripping hand",
				}},
			},
		},
		{
			test: `In the lobby are a supporter and a container.`,
			result: map[string]any{
				"macro": "in",
				"targets": []map[string]any{{
					"det":  "the",
					"name": "lobby",
				}},
				"sources": []map[string]any{{
					"kinds": []string{"supporter"},
				}, {
					"kinds": []string{"container"},
				}},
			},
		},

		// the special nxn description: no properties are allowed.
		{
			//	test: `Hector and Maria are suspicious of Santa and Santana.`,
		},
		// fix: in this case i think inform eats the first "is" and allows a subsequent one ( and also allows values )
		// fix: the device called the detonator is on the supporter called the shelf and is proper named"

		// humorous: "In the lobby are two supporters the bat and the hat."
		// generates a Noun called "two supporters..."
		// "In the lobby are two supporters" ( and "Two supporters are in..." works fine. )
	}
	var skipped int
	for i, p := range phrases {
		if len(p.test) > 0 && p.result == nil {
			skipped++
		} else {
			res, haveError := Grok(p.test)
			if expectError, ok := p.result.(error); ok {
				if haveError == nil {
					t.Log(i, p.test, "expected an error", expectError, "but succeeded")
					t.Fail()
				}
			} else if haveError != nil {
				t.Log(i, p.test, haveError)
				t.Fail()
			} else if expectMap, ok := p.result.(map[string]any); ok {
				m := resultMap(res)
				if d := pretty.Diff(expectMap, m); len(d) > 0 {
					t.Log("test", i, p.test, "got:\n", pretty.Sprint(m))
					//t.Log("want:", pretty.Sprint(p.result))
					t.Log(d)
					t.Fail()
				}
			}
		}
	}
	if skipped > 0 {
		t.Fatalf("skipped %d tests", skipped)
	}
}

func TestTraits(t *testing.T) {
	var phrases = []struct {
		test   string
		result any
		skip   any
	}{{
		test: "open container",
		result: map[string]any{
			"kind":   "container",
			"traits": []string{"open"},
		},
	}, {
		test: "the open and an openable container",
		result: map[string]any{
			"kind":   "container",
			"traits": []string{"open", "openable"},
		},
	}, {
		test: "open, and openable",
		result: map[string]any{
			"traits": []string{"open", "openable"},
		},
	}, {
		test: "open, openable",
		result: map[string]any{
			"traits": []string{"open", "openable"},
		},
	}, {
		test:   "open and and openable",
		result: errutil.New("two ands should fail"),
	}, {
		test:   "open and, openable",
		result: errutil.New("backwards commas should fail"),
	}}
	var skipped int
	for i, p := range phrases {
		ts, haveError := parseTraitSet(panicHash(p.test))
		if len(p.test) > 0 && p.result == nil {
			skipped++
		} else if expectError, ok := p.result.(error); ok {
			if haveError == nil {
				t.Fatal(i, p.test, "expected an error", expectError, "but succeeded")
			}
		} else if haveError != nil {
			t.Fatal(i, p.test, haveError)
		} else if expectMap, ok := p.result.(map[string]any); ok {
			m := traitSetMap(ts)
			if d := pretty.Diff(expectMap, m); len(d) > 0 {
				t.Log("test", i, p.test, "got:\n", pretty.Sprint(m))
				//t.Log("want:", pretty.Sprint(p.result))
				t.Fatal(d)
			}
		}
	}
	if skipped > 0 {
		t.Fatalf("skipped %d tests", skipped)
	}
}

func resultMap(in Results) map[string]any {
	m := make(map[string]any)
	nounsIntoMap(m, "sources", in.sources)
	nounsIntoMap(m, "targets", in.targets)
	wordsIntoMap(m, "macro", in.macro)
	return m
}

func traitSetMap(ts traitSet) map[string]any {
	m := make(map[string]any)
	wordListIntoMap(m, "traits", ts.traits)
	wordsIntoMap(m, "kind", ts.kind)
	return m
}

func nounsIntoMap(m map[string]any, field string, ns []Noun) {
	if len(ns) > 0 {
		out := make([]map[string]any, len(ns))
		for i, n := range ns {
			out[i] = nounToMap(n)
		}
		m[field] = out
	}
}

func nounToMap(n Noun) map[string]any {
	m := make(map[string]any)
	wordsIntoMap(m, "name", n.name)
	wordsIntoMap(m, "det", n.det)
	wordListIntoMap(m, "traits", n.traits)
	wordListIntoMap(m, "kinds", n.kinds)
	return m
}

func wordListIntoMap(m map[string]any, field string, ws [][]Word) {
	if cnt := len(ws); cnt > 0 {
		out := make([]string, cnt)
		for i, w := range ws {
			out[i] = wordsToString(w)
		}
		m[field] = out
	}
	return
}

func wordsIntoMap(m map[string]any, field string, w []Word) {
	if len(w) > 0 {
		m[field] = wordsToString(w)
	}
}

func wordsToString(w []Word) (ret string) {
	var b strings.Builder
	for i, w := range w {
		if i > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(w.String())
	}
	return b.String()
}
