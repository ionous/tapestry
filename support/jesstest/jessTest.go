// Package jesstest exercises implementations of jess.Query
// to ensure they produce good results.
// ( see for instance package jessdb_test, and package jess_test. )
package jesstest

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/files"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// verify a standard set of phrases using some function that takes each of those phrases
// and produces a result map, or a string list.
func RunPhraseTests(t *testing.T, interpret func(string) (jess.Generator, error)) {
	var phrases = []struct {
		test   string
		result any
		skip   any
	}{

		// ------------------------
		// PropertyNounValue
		{
			test: `The title of the story is "A Secret."`,
			result: []string{
				"AddNounValue", "story", "title", text("A Secret."),
			},
		},
		{
			// note: we don't validate properties while matching
			// weave validates them when attempting to write them.
			test: `The age of the bottle is 42.`,
			result: []string{
				"AddNounValue", "bottle", "age", number(42),
			},
		},
		// ------------------------
		// NounPropertyValue
		{
			test: `The story has the title "{15|print_num!}"`,
			result: []string{
				// test that it can convert a template
				"AddNounValue", "story", "title", `{"FromText:":{"Numeral:":{"Num value:":15}}}`,
			},
		},
		{
			// note: we don't validate properties while matching
			// weave validates them when attempting to write them.
			test: `The bottle has an age of 42.`,
			result: []string{
				"AddNounValue", "bottle", "age", number(42),
			},
		},
		// ------------------------
		// Aspects are traits
		// note: "The colors can be..." isn't allowed.
		// ( for testing, "colors" is an established aspect with zero traits )
		{
			test: `The colors are red, blue, and cobalt.`,
			result: []string{
				"AddTraits", "color", "red", "blue", "cobalt",
			},
		},
		// ------------------------
		// Kinds have properties
		{
			test:   "Containers are either opened or closed.",
			result: nil,
		},
		{
			// can be [either] ...
			test:   "A thing can be opened or closed.",
			result: nil,
		},
		{
			// inform doesnt allow [either] here; i'm fine however it works.
			test:   "A thing can be scenery.",
			result: nil,
		},
		{
			test:   "Things have some text called a description.",
			result: []string{"AddFields", "things", "description", "text", ""},
		},
		{
			test:   "Things have some text.",
			result: errors.New("unnamed text fields are prohibited."),
		},
		{
			test:   "Things have a number.",
			result: errors.New("unnamed number fields are prohibited."),
		},
		{
			test:   "A supporter has a number called carrying capacity.",
			result: []string{"AddFields", "supporters", "carrying capacity", "number", ""},
		},
		{
			// except for number and text, inform allows "bare" properties: a "list of text" creates a member called "list of text"
			test:   "Things have a list of text called frenemies.",
			result: []string{"AddFields", "things", "frenemies", "text_list", ""},
		},
		{
			test:   "Things have a list of numbers called the lotto numbers.",
			result: []string{"AddFields", "things", "lotto numbers", "num_list", ""},
		},
		{
			// groups are a pre-defined type of record; anonymous fields are allowed.
			// references to kinds become text; except for records which are embedded.
			test:   "Things have a color.",
			result: []string{"AddFields", "things", "color", "text", "color"},
		},
		{
			test:   "Things have a group.",
			result: []string{"AddFields", "things", "group", "record", "groups"},
		},
		{
			test:   "Things have a group called the set.",
			result: []string{"AddFields", "things", "set", "record", "groups"},
		},
		{
			test:   "Things have a list of groups.",
			result: []string{"AddFields", "things", "groups", "record_list", "groups"},
		},

		// -------------------------
		{
			// note: "Devices are fixed in place" will parse properly
			// but weave will assume that the name "devices" refers to a noun
			// and ( probably, hopefully ) some error will occur.
			// I like that usually specifically indicates and separates kinds from nouns --
			// im not sure the other certainties (never, always) are really needed:
			// if so: the "final" field for mdl_value, mdl_value_kind could be used.
			test: `Containers are usually closed.`,
			result: []string{
				"AddKindTrait", "containers", "closed",
			},
		},
		{
			test: `Containers and supporters are usually fixed in place.`,
			result: []string{
				"AddKindTrait", "containers", "fixed in place",
				"AddKindTrait", "supporters", "fixed in place",
			},
		},
		{
			test: `Two things are in the kitchen.`,
			result: []string{
				"AddNoun", "things-1", "things-1", "things",
				"AddNounAlias", "things-1", "thing",
				"AddNounTrait", "things-1", "counted",
				"AddNounValue", "things-1", "printed name", text("thing"),
				//
				"AddNoun", "things-2", "things-2", "things",
				"AddNounAlias", "things-2", "thing",
				"AddNounTrait", "things-2", "counted",
				"AddNounValue", "things-2", "printed name", text("thing"),
				//
				"ApplyMacro", "contain", "kitchen", "things-1", "things-2",
			},
		},
		{
			test: `Hershel is carrying scissors and a pen.`,
			result: []string{
				"AddNounTrait", "hershel", "proper named",
				"ApplyMacro", "carry", "hershel", "scissors", "pen",
			},
		},
		// reverse carrying relation.
		{
			test: `The scissors and a pen are carried by Hershel.`,
			result: []string{
				"AddNounTrait", "hershel", "proper named",
				"ApplyMacro", "carry", "hershel", "scissors", "pen",
			},
		},

		// simple trait:
		{
			test: `The bottle is closed.`,
			result: []string{
				"AddNounTrait", "bottle", "closed",
			},
		},
		// multi Word trait:
		{
			test: `The tree is fixed in place.`,
			result: []string{
				"AddNounTrait", "tree", "fixed in place",
			},
		},
		// multiple trailing properties, using the kind as a property.
		{
			test: `The bottle is a transparent, open, container.`,
			result: []string{
				"AddNoun", "bottle", "", "containers",
				"AddNounTrait", "bottle", "transparent",
				"AddNounTrait", "bottle", "open",
			},
		},
		// multiple nouns of different kinds
		{
			test: `The box and the top are closed containers.`,
			result: []string{
				"AddNoun", "box", "", "containers",
				"AddNounTrait", "box", "closed",
				"AddNoun", "top", "", "containers",
				"AddNounTrait", "top", "closed",
			},
		},
		// using 'called' without a macro
		{
			test: `The container called the sarcophagus is open.`,
			result: []string{
				"AddNoun", "sarcophagus", "", "containers",
				"AddNounTrait", "sarcophagus", "open",
			},
		},
		// ------------------------------------------------------------------------
		// KindsOf
		// ------------------------------------------------------------------------
		{
			// when the requested kind being declared isn't yet known:
			test: `Devices are a kind of thing.`,
			result: []string{
				"AddKind", "devices", "things",
			},
		},
		{
			// when the kind being declared is already known
			test: `A container is a kind of thing.`,
			result: []string{
				"AddKind", "containers", "things",
			},
		},
		{
			// determine pluralness based on trailing is/are
			test: `A device is a kind of thing.`,
			result: []string{
				"AddKind", "devices", "things",
			},
		},
		{
			// adding trailing properties
			test: `A casket is a kind of closed container.`,
			result: []string{
				"AddKind", "caskets", "containers",
				"AddKindTrait", "caskets", "closed",
			},
		},
		{
			// complex parsing
			test: `The closed containers called the safes are a kind of fixed in place thing.`,
			result: []string{
				"AddKind", "safes", "things",
				"AddKind", "safes", "containers", // the parent can be singular or plural
				"AddKindTrait", "safes", "closed",
				"AddKindTrait", "safes", "fixed in place",
			},
		},
		{
			// correctly producing unexpected results.
			test: `The closed casket is a kind of container.`,
			result: []string{
				"AddKind", "closed caskets", "containers",
			},
		},
		{
			// adding middle properties is not allowed ( should it be? )
			test:   `The casket is a closed kind of container.`,
			result: errutil.New("not allowed"),
		},
		{
			// in inform, these become the plural kind "Bucketss" and "basketss"
			test: `Buckets and baskets are kinds of container.`,
			result: []string{
				"AddKind", "buckets", "containers",
				"AddKind", "baskets", "containers",
			},
		},
		{
			test:   `A container is in the lobby.`,
			result: errors.New("this is specifically disallowed, and should generate an error"),
		},
		// rhs-contains; "in" is
		{
			test: `The unhappy man is in the closed bottle.`,
			result: []string{
				"ApplyMacro", "contain", "closed bottle", "unhappy man",
			},
		},
		// same pattern as the middle properties above; but not using kind of
		{
			test: `The coffin is a closed container in the antechamber.`,
			result: []string{
				"AddNoun", "coffin", "", "containers",
				"AddNounTrait", "coffin", "closed",
				"ApplyMacro", "contain", "antechamber", "coffin",
			},
		},
		// note, this is allowed even though it implies something different than what is written:
		{
			test: `The bottle is openable in the kitchen.`,
			result: []string{
				"AddNounTrait", "bottle", "openable",
				"ApplyMacro", "contain", "kitchen", "bottle",
			},
		},
		// called both before and after the macro
		// note: The closed openable container called the trunk and the box is in the lobby.
		// would create a noun called "the trunk and the box"
		{
			test: `The thing called the stake is on the supporter called the altar.`,
			result: []string{
				"AddNoun", "altar", "", "supporters",
				"AddNoun", "stake", "", "things",
				"ApplyMacro", "support", "altar", "stake",
			},
		},
		// add leading properties using 'called'
		// "is" left of the macro "in".
		// slightly different parsing than "kind/s of":
		// those expect only expect one set of nouns; these have two.
		{
			test: `A closed openable container called the trunk is in the lobby.`,
			result: []string{
				"AddNoun", "trunk", "", "containers",
				"AddNounTrait", "trunk", "closed",
				"AddNounTrait", "trunk", "openable",
				"ApplyMacro", "contain", "lobby", "trunk",
			},
		},
		// multiple primary:
		// "is" left of the macro "in".
		{
			test: `Some coins, a notebook, and the gripping hand are in the coffin.`,
			result: []string{
				"AddNounValue", "coins", "indefinite article", text("some"),
				"ApplyMacro", "contain", "coffin", "coins", "notebook", "gripping hand",
			},
		},
		// multiple primary with a leading macro
		{
			test: `In the coffin are some coins, a notebook, and the gripping hand.`,
			result: []string{
				"AddNounValue", "coins", "indefinite article", text("some"),
				"ApplyMacro", "contain", "coffin", "coins", "notebook", "gripping hand",
			},
		},

		// multiple anonymous nouns.
		{
			test: `In the lobby are a supporter and a container.`,
			result: []string{
				"AddNoun", "", "", "supporters",
				"AddNoun", "", "", "containers",
				"ApplyMacro", "contain", "lobby", "", "",
			},
		},
		// the special nxn description: no properties are allowed.
		{
			test: `Hector and Maria are suspicious of Santa and Santana.`,
			result: []string{
				"AddNounTrait", "hector", "proper named",
				"AddNounTrait", "maria", "proper named",
				"AddNounTrait", "santa", "proper named",
				"AddNounTrait", "santana", "proper named",
				"ApplyMacro", "suspect", "hector", "maria", "santa", "santana",
			},
		},
	}
	var skipped int
	for i, p := range phrases {
		if len(p.test) == 0 || p.result == nil {
			skipped++
		} else {
			var haveRes []string
			got, haveError := interpret(p.test)
			if haveError == nil {
				var m Mock
				if e := got.Generate(&m); e != nil {
					haveError = e
				} else {
					haveRes = m.out
				}
			}
			if expectError, ok := p.result.(error); ok {
				if haveError != nil {
					t.Log("ok, test", i, p.test, haveError)
				} else {
					//
					t.Log("NG! test", i, p.test, "expected an error", expectError, "but succeeded with", pretty.Sprint(haveRes))
					t.Fail()
				}
			} else if haveError != nil {
				t.Log("NG! test", i, p.test, haveError)
				t.Fail()
			} else {
				if d := pretty.Diff(p.result, haveRes); len(d) > 0 {
					t.Log("NG! test", i, p.test, "got:\n", pretty.Sprint(haveRes))
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

func Marshal(a rt.Assignment) (ret string, err error) {
	var enc encode.Encoder
	if v, ok := a.(typeinfo.Instance); !ok {
		err = errors.New("not a generated type?")
	} else if d, e := enc.Encode(v); e != nil {
		err = e
	} else {
		var str strings.Builder
		if e := files.JsonEncoder(&str, files.RawJson).Encode(d); e != nil {
			err = e
		} else {
			ret = str.String()
		}
	}
	return
}

func text(str string) string {
	return fmt.Sprintf(`{"FromText:":{"Text value:":%q}}`, str)
}

func number(num float64) string {
	return fmt.Sprintf(`{"FromNumber:":{"Num value:":%g}}`, num)
}
