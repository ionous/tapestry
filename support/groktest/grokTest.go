// exercises implementations of grok.Grokker to ensure they produce good results.
package groktest

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/grok"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func Phrases(t *testing.T, g grok.Grokker) {
	var phrases = []struct {
		test   string
		result any
		skip   any
	}{
		{
			test: `Devices are a kind of prop.`,
			result: map[string]any{
				"macro": "inherit",
				"primary": []map[string]any{
					{
						"kinds": []string{"prop"},
						"name":  "Devices",
					},
				},
			},
		},
		{
			// note: "Devices are fixed in place" will parse properly
			// but storyGrok will assume that the name "devices" refers to a noun
			// and ( probably, hopefully ) some error will occur.
			// I like that usually specifically indicates and separates kinds from nouns --
			// im not sure the other certainties (never, always) are really needed:
			// if so: the "final" field for mdl_value, mdl_value_kind could be used.
			test: `Devices are usually closed.`,
			result: map[string]any{
				"macro": "implies",
				"primary": []map[string]any{
					{
						"name":   "Devices",
						"traits": []string{"closed"},
					},
				},
			},
		},
		{
			// note: in inform...	 §4.14. Duplicates
			// "Two circles are in the Lab."
			// it only works if "circles" is a known kind
			// otherwise, it assumes "two circles" is the complete name of a single noun.
			test: `Two things are in the kitchen.`,
			result: map[string]any{
				"macro": "contain",
				"secondary": []map[string]any{{
					"count": 2,
					"kinds": []string{"things"},
				}},
				"primary": []map[string]any{{
					"det":  "the",
					"name": "kitchen",
				}},
			},
		},

		{
			test: `Hershel is carrying scissors and a pen.`,
			result: map[string]any{
				"macro": "carry",
				"primary": []map[string]any{{
					"name": "Hershel",
				}},
				"secondary": []map[string]any{{
					"name": "scissors",
				}, {
					"det":  "a",
					"name": "pen",
				}},
			},
		},
		// reverse carrying relation.
		{
			test: `The scissors and a pen are carried by Hershel.`,
			result: map[string]any{
				"macro": "carry",
				"primary": []map[string]any{{
					"name": "Hershel",
				}},
				"secondary": []map[string]any{{
					"det":  "the",
					"name": "scissors",
				}, {
					"det":  "a",
					"name": "pen",
				}},
			},
		},

		// simple trait:
		{
			test: `The bottle is closed.`,
			result: map[string]any{
				"primary": []map[string]any{{
					"det":    "the",
					"name":   "bottle",
					"traits": []string{"closed"},
				}},
			},
		},
		// multi Word trait:
		{
			test: `The tree is fixed in place.`,
			result: map[string]any{
				"primary": []map[string]any{{
					"det":    "the",
					"name":   "tree",
					"traits": []string{"fixed in place"},
				}},
			},
		},
		// multiple trailing properties, using the kind as a property.
		{
			test: `The bottle is a transparent, open, container.`,
			result: map[string]any{
				"primary": []map[string]any{{
					"det":    "the",
					"name":   "bottle",
					"kinds":  []string{"container"},
					"traits": []string{"transparent", "open"},
				}},
			},
		},
		// multiple nouns of different kinds
		{
			test: `The box and the top are closed containers.`,
			result: map[string]any{
				"primary": []map[string]any{{
					"det":    "the",
					"name":   "box",
					"traits": []string{"closed"},
					"kinds":  []string{"containers"},
				}, {
					"det":    "the",
					"name":   "top",
					"traits": []string{"closed"},
					"kinds":  []string{"containers"},
				}},
			},
		},
		// using 'called' without a macro
		{
			test: `The container called the sarcophagus is open.`,
			result: map[string]any{
				"primary": []map[string]any{{
					"det":    "the", // note: this is the bit closes to the noun
					"name":   "sarcophagus",
					"exact":  true,
					"kinds":  []string{"container"},
					"traits": []string{"open"},
				}},
			},
		},
		// a kind of declaration ( uses a 'macro' verb )
		// "is" left of macro
		{
			test: `A casket is a kind of container.`,
			result: map[string]any{
				"macro": "inherit",
				"primary": []map[string]any{{
					"det":   "a",
					"name":  "casket",
					"kinds": []string{"container"},
				}},
			},
		},
		// kind of: adding trailing properties
		// "is" left of macro
		{
			test: `A casket is a kind of closed container.`,
			result: map[string]any{
				"macro": "inherit",
				"primary": []map[string]any{{
					"det":    "a",
					"name":   "casket",
					"traits": []string{"closed"},
					"kinds":  []string{"container"},
				}},
			},
		},
		// kind of, "correctly" producing incorrect results.
		// note: in inform, this also yields a noun named the "closed casket".
		// similarly, "The kind of the casket is a container", yields a name "kind of the casket".
		// "is" left of macro.
		{
			test: `The closed casket is a kind of container.`,
			result: map[string]any{
				"macro": "inherit",
				"primary": []map[string]any{{
					"det":   "the",
					"name":  "closed casket",
					"kinds": []string{"container"},
				}},
			},
		},
		// kind of: adding middle properties
		// tbd: not allowed, but maybe it should be....
		// "is" left of macro
		{
			test:   `The casket is a closed kind of container.`,
			result: errutil.New("not allowed"),
		},
		// in inform, these become the plural kind "Bucketss" and "basketss"
		{
			test: `Buckets and baskets are kinds of container.`,
			result: map[string]any{
				"macro": "inherit",
				"primary": []map[string]any{
					{
						"kinds": []string{"container"},
						"name":  "Buckets",
					},
					{
						"kinds": []string{"container"},
						"name":  "baskets",
					},
				}},
		},
		{
			test:   `A container is in the lobby.`,
			result: errutil.New("this is specifically disallowed, and should generate an error"),
		},
		// rhs-contains; "in" is
		{
			test: `The unhappy man is in the closed bottle.`,
			result: map[string]any{
				"macro": "contain",
				"secondary": []map[string]any{{
					"det": "the",
					// giving properties to the rhs and right secondary isnt permitted:
					// tbd: but it might be possible...
					"name": "unhappy man",
				}},
				"primary": []map[string]any{{
					"det":  "the",
					"name": "closed bottle",
				}},
			}},
		// same pattern as the middle properties above; but not using kind of
		{
			test: `The coffin is a closed container in the antechamber.`,
			result: map[string]any{
				"macro": "contain",
				"secondary": []map[string]any{{
					"det":    "the",
					"name":   "coffin",
					"traits": []string{"closed"},
					"kinds":  []string{"container"},
				}},
				"primary": []map[string]any{{
					"det":  "the",
					"name": "antechamber",
				}},
			},
		},
		// note, this is allowed even though it implies something different than what is written:
		{
			test: `The bottle is openable in the kitchen.`,
			result: map[string]any{
				"macro": "contain",
				"secondary": []map[string]any{{
					"det":    "the",
					"traits": []string{"openable"},
					"name":   "bottle",
				}},
				"primary": []map[string]any{{
					"det":  "the",
					"name": "kitchen",
				}},
			},
		},
		// called both before and after the macro
		// note: The closed openable container called the trunk and the box is in the lobby.
		// would create a noun called "the trunk and the box"
		{
			test: `The thing called the stake is on the supporter called the altar.`,
			result: map[string]any{
				"macro": "support",
				"secondary": []map[string]any{{
					"det":   "the",
					"name":  "stake",
					"exact": true,
					"kinds": []string{"thing"},
				}},
				"primary": []map[string]any{{
					"det":   "the",
					"name":  "altar",
					"exact": true,
					"kinds": []string{"supporter"},
				}},
			},
		},
		// add leading properties using 'called'
		// "is" left of the macro "in".
		// slightly different parsing than "kind/s of":
		// those expect only expect one set of nouns; these have two.
		{
			test: `A closed openable container called the trunk is in the lobby.`,
			result: map[string]any{
				"macro": "contain",
				"primary": []map[string]any{{
					"det":  "the",
					"name": "lobby",
				}},
				"secondary": []map[string]any{{
					"det":    "the", // closest to the trunk
					"name":   "trunk",
					"exact":  true,
					"traits": []string{"closed", "openable"},
					"kinds":  []string{"container"},
				}},
			},
		},
		// multiple primary:
		// "is" left of the macro "in".
		{
			test: `Some coins, a notebook, and the gripping hand are in the coffin.`,
			result: map[string]any{
				"macro": "contain",
				"primary": []map[string]any{{
					"det":  "the", // closest to the coffin
					"name": "coffin",
				}},
				"secondary": []map[string]any{{
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
		// multiple primary with a leading macro
		{
			test: `In the coffin are some coins, a notebook, and the gripping hand.`,
			result: map[string]any{
				"macro": "contain",
				"primary": []map[string]any{{
					"det":  "the", // lowercase, the closest to the trunk
					"name": "coffin",
				}},
				"secondary": []map[string]any{{
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
		// multiple anonymous nouns.
		{
			test: `In the lobby are a supporter and a container.`,
			result: map[string]any{
				"macro": "contain",
				"primary": []map[string]any{{
					"det":  "the",
					"name": "lobby",
				}},
				"secondary": []map[string]any{{
					"kinds": []string{"supporter"},
				}, {
					"kinds": []string{"container"},
				}},
			},
		},
		// the special nxn description: no properties are allowed.
		{
			test: `Hector and Maria are suspicious of Santa and Santana.`,
			result: map[string]any{
				"macro": "suspect",
				"primary": []map[string]any{{
					"name": "Hector",
				}, {
					"name": "Maria",
				}},
				"secondary": []map[string]any{{
					"name": "Santa",
				}, {
					"name": "Santana",
				}},
			},
		},
		// fix: trailing properties applying to the lhs
		{
			test: `The bottle in the kitchen is openable.`,
			skip: map[string]any{
				"macro": "contain",
				"primary": []map[string]any{{
					"det":    "the",
					"traits": []string{"openable"},
					"name":   "bottle",
				}},
				"secondary": []map[string]any{{
					"det":  "the",
					"name": "kitchen",
				}},
			},
		},
		// TODO: values.
		{
			//test:  `The bottle in the kitchen is openable and has age 42.`,
		},
		{
			//test: `The age of the bottle is 42.`,
		},
		// todo:  the device called the detonator is on the supporter called the shelf and is proper named"
		// todo: In the lobby are two supporters" ( and "Two supporters are in..." works fine. )
		// note: "In the lobby are two supporters the bat and the hat." generates a noun called "two supporters..."
	}
	var skipped int
	for i, p := range phrases {
		if len(p.test) > 0 && p.result == nil {
			skipped++
		} else {
			res, haveError := grok.Grok(g, p.test)
			if expectError, ok := p.result.(error); ok {
				if haveError != nil {
					t.Log("ok, test", i, p.test, haveError)
				} else {
					t.Log("ng, test", i, p.test, "expected an error", expectError, "but succeeded with", pretty.Sprint(resultMap(res)))
					t.Fail()
				}
			} else if haveError != nil {
				t.Log("ng, test", i, p.test, haveError)
				t.Fail()
			} else if expectMap, ok := p.result.(map[string]any); ok {
				m := resultMap(res)
				if d := pretty.Diff(expectMap, m); len(d) > 0 {
					t.Log("ng, test", i, p.test, "got:\n", pretty.Sprint(m))
					//t.Log("want:", pretty.Sprint(p.result))
					t.Log(d)
					t.Fail()
				}
			}
		}
	}
	if skipped > 0 {
		t.Logf("skipped %d tests", skipped)
	}
}

func Traits(t *testing.T, g grok.Grokker) {
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
		span, e := grok.MakeSpan(p.test)
		if e != nil {
			t.Fatal(e)
		}
		ts, haveError := grok.ParseTraitSet(g, span)
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
