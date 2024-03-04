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
	"git.sr.ht/~ionous/tapestry/support/files"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// verify a standard set of phrases using some function that takes each of those phrases
// and produces a result map, or a string list.
func TestPhrases(t *testing.T, q jess.Query) {
	var phrases = []struct {
		test   string
		result any
		skip   any
	}{
		// ------------------------------------------------------------------------
		// MapConnection
		// ------------------------------------------------------------------------
		{
			test: `Through the long slide is nowhere.`,
			result: []string{
				"AddNounKind", "long slide", "doors",
				"AddNounName", "long slide", "long slide",
				"AddNounValue", "long slide", "destination", textKind("", "rooms"),
			},
		},
		{
			test: `Through the blue door is the Flat Landing.`,
			result: []string{
				"AddNounKind", "flat landing", "rooms",
				"AddNounName", "flat landing", "Flat Landing",
				"AddNounKind", "blue door", "doors",
				"AddNounName", "blue door", "blue door",
				"AddNounValue", "blue door", "destination", textKind("flat landing", "rooms"),
			},
		},
		{
			test: `Through the gate and the hatch is a dark room called An End.`,
			result: []string{
				"AddNounKind", "an end", "rooms",
				"AddNounName", "an end", "An End",
				"AddNounTrait", "an end", "dark",
				"AddNounTrait", "an end", "proper named",
				"AddNounKind", "gate", "doors",
				"AddNounName", "gate", "gate",
				"AddNounValue", "gate", "destination", textKind("an end", "rooms"),
				"AddNounKind", "hatch", "doors",
				"AddNounName", "hatch", "hatch",
				"AddNounValue", "hatch", "destination", textKind("an end", "rooms"),
			},
		},
		// ------------------------------------------------------------------------
		// Understandings
		// ------------------------------------------------------------------------
		{
			test: `Understand "floor" or "sawdust" as the message.`,
			result: []string{
				"AddNounAlias", "message", "floor",
				"AddNounAlias", "message", "sawdust",
			},
		},
		{
			test: `Understand "missives" as the plural of missive and message.`,
			result: []string{
				"AddPlural", "missives", "missive",
				"AddPlural", "missives", "message",
			},
		},
		{
			test: `Understand "hang [objects] on/onto/-- [objects]" as storing.`,
			result: []string{
				"AddGrammar",
				"hang [objects] on/onto/-- [objects]",
				`{"One word:":["hang"]}`,
				`{"One noun:":"objects"}`,
				`{"One word:":["on","onto",""]}`,
				`{"One noun:":"objects"}`,
				`{"Action:":"storing"}`,
			},
		},
		// ------------------------------------------------------------------------
		// PropertyNounValue
		// ------------------------------------------------------------------------
		{
			test: `The title of the story is "A Secret."`,
			result: []string{
				// no noun declaration because the story is a known noun ( in these tests )
				"AddNounValue", "story", "title", text("A Secret."),
			},
		},
		{
			// note: we don't validate properties while matching
			// weave validates them when attempting to write them.
			test: `The age of the bottle is 42.`,
			result: []string{
				"AddNounKind", "bottle", "things",
				"AddNounName", "bottle", "bottle",
				"AddNounValue", "bottle", "age", number(42),
			},
		},
		{
			// this should properly create a new noun "story teller"
			// and not try to match the existing noun "story."
			test: `The age of the story teller is 42.`,
			result: []string{
				"AddNounKind", "story teller", "things",
				"AddNounName", "story teller", "story teller",
				"AddNounValue", "story teller", "age", number(42),
			},
		},
		{
			// fix: currently succeeds with "thing called the cat"
			// inform gets confused, but we could handle this okay
			// test:   `The description of the thing called the cat is "meow."`,
			// result: errors.New("can't use property noun value this way."),
		},
		// ------------------------------------------------------------------------
		// NounPropertyValue
		// ------------------------------------------------------------------------
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
				"AddNounKind", "bottle", "things",
				"AddNounName", "bottle", "bottle",
				"AddNounValue", "bottle", "age", number(42),
			},
		},
		{
			// fix: inform allows this, and jess does not.
			// 	test: `The thing called the cat has the description "meow."`,
			// 	result: []string{
			// 	},
		},

		// ------------------------------------------------------------------------
		// AspectsAreTraits
		// ( for testing, "colors" is an established aspect with zero traits )
		// ------------------------------------------------------------------------
		{
			test: `The colors are red, blue, and cobalt.`,
			result: []string{
				"AddTraits", "color", "red", "blue", "cobalt",
			},
		},
		{
			// not allowed. matches KindsAreEither, but "aspect" is prohibited for that phrase.
			test:   `The colors can be red, blue, and cobalt.`,
			result: errors.New("not allowed?"),
		},
		// ------------------------------------------------------------------------
		// KindsAreTraits
		// ------------------------------------------------------------------------
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
				"AddKind", "safes", "containers",
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
		// ------------------------------------------------------------------------
		// KindsAreEither
		// A thing is either tall or short.
		// ng: Things are either tall or short. [ i dont see a reason not to allow this ]
		// A thing can be [either]
		// Things can be [either] tall or short.
		// ------------------------------------------------------------------------
		{
			test: "Containers are either opened or closed.",
			result: []string{
				// inform names these: "container status", "container status 2", etc.
				"AddKind", "opened status", "aspects",
				"AddTraits", "opened status", "opened", "closed",
				"AddFields", "containers", "opened status", "text", "opened status",
			},
		},
		{
			// can be [either] ...
			test: "A thing can be opened or closed or ajar.",
			result: []string{
				"AddKind", "opened status", "aspects",
				"AddTraits", "opened status", "opened", "closed", "ajar",
				"AddFields", "things", "opened status", "text", "opened status",
			},
		},
		{
			// inform doesnt allow [either] here; i'm fine with whatever.
			test:   "A thing can be scenery.",
			result: []string{"AddFields", "things", "scenery", "bool", ""},
		},
		// ------------------------------------------------------------------------
		// KindsHaveProperties
		// ------------------------------------------------------------------------
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
		// ------------------------------------------------------------------------
		// NamesVerbNames
		// ------------------------------------------------------------------------
		{
			test: `Two things are in the kitchen.`,
			result: []string{
				"AddNounKind", "kitchen", "things", // uses things unless some other constraint exists
				"AddNounName", "kitchen", "kitchen",
				//
				"AddNounKind", "thing-1", "things",
				"AddNounName", "thing-1", "thing-1",
				"AddNounAlias", "thing-1", "thing",
				"AddNounTrait", "thing-1", "counted",
				"AddNounValue", "thing-1", "printed name", text("thing"),
				//
				"AddNounKind", "thing-2", "things",
				"AddNounName", "thing-2", "thing-2",
				"AddNounAlias", "thing-2", "thing",
				"AddNounTrait", "thing-2", "counted",
				"AddNounValue", "thing-2", "printed name", text("thing"),
				//
				"ApplyMacro", "contain", "kitchen", "thing-1", "thing-2",
			},
		},
		{
			test: `Hershel is carrying scissors and a pen.`,
			result: []string{
				"AddNounKind", "hershel", "things", // FIX: carries should be actor
				"AddNounName", "hershel", "Hershel",
				"AddNounTrait", "hershel", "proper named",
				//
				"AddNounKind", "scissors", "things",
				"AddNounName", "scissors", "scissors",
				"AddNounTrait", "scissors", "proper named", // yes; this conforms with inform.
				//
				"AddNounKind", "pen", "things",
				"AddNounName", "pen", "pen",
				"ApplyMacro", "carry", "hershel", "scissors", "pen",
			},
		},
		{
			// reverse carrying relation.
			test: `The scissors and a pen are carried by Hershel.`,
			result: []string{
				"AddNounKind", "hershel", "things",
				"AddNounName", "hershel", "Hershel",
				"AddNounTrait", "hershel", "proper named",
				//
				"AddNounKind", "scissors", "things",
				"AddNounName", "scissors", "scissors",
				//
				"AddNounKind", "pen", "things",
				"AddNounName", "pen", "pen",
				"ApplyMacro", "carry", "hershel", "scissors", "pen",
			},
		},
		{
			test: `The unhappy man is in the closed bottle.`,
			result: []string{
				"AddNounKind", "closed bottle", "things",
				"AddNounName", "closed bottle", "closed bottle",
				"AddNounKind", "unhappy man", "things",
				"AddNounName", "unhappy man", "unhappy man",
				"ApplyMacro", "contain", "closed bottle", "unhappy man",
			},
		},
		{
			// called both before and after the macro
			// note: The closed openable container called the trunk and the box is in the lobby.
			// would create a noun called "the trunk and the box"
			test: `The thing called the stake is on the supporter called the altar.`,
			result: []string{
				"AddNounKind", "altar", "supporters",
				"AddNounName", "altar", "altar",
				"AddNounValue", "altar", "indefinite article", text("the"),
				//
				"AddNounKind", "stake", "things",
				"AddNounName", "stake", "stake",
				"AddNounValue", "stake", "indefinite article", text("the"),
				//
				"ApplyMacro", "support", "altar", "stake",
			},
		},
		{
			// add leading properties using 'called'
			// "is" left of the macro "in".
			// slightly different parsing than "kind/s of":
			// those expect only expect one set of nouns; these have two.
			test: `A closed openable container called the trunk is in the lobby.`,
			result: []string{
				"AddNounKind", "lobby", "things",
				"AddNounName", "lobby", "lobby",
				//
				"AddNounKind", "trunk", "containers",
				"AddNounName", "trunk", "trunk",
				"AddNounTrait", "trunk", "closed",
				"AddNounTrait", "trunk", "openable",
				"AddNounValue", "trunk", "indefinite article", text("the"),
				"ApplyMacro", "contain", "lobby", "trunk",
			},
		},
		{
			// multiple primary: "is" left of the macro "in".
			test: `Some coins, a notebook, and the gripping hand are in the coffin.`,
			result: []string{
				"AddNounKind", "coffin", "things",
				"AddNounName", "coffin", "coffin",
				//
				"AddNounKind", "coins", "things",
				"AddNounName", "coins", "coins",
				"AddNounTrait", "coins", "plural named",
				//
				"AddNounKind", "notebook", "things",
				"AddNounName", "notebook", "notebook",
				//
				"AddNounKind", "gripping hand", "things",
				"AddNounName", "gripping hand", "gripping hand",
				//
				"ApplyMacro", "contain", "coffin", "coins", "notebook", "gripping hand",
			},
		},
		{
			// the special nxn description: no properties are allowed.
			test: `Hector and Maria are suspicious of Santa and Santana.`,
			result: []string{
				"AddNounKind", "hector", "things",
				"AddNounName", "hector", "Hector",
				"AddNounTrait", "hector", "proper named",
				//
				"AddNounKind", "maria", "things",
				"AddNounName", "maria", "Maria",
				"AddNounTrait", "maria", "proper named",
				//
				"AddNounKind", "santa", "things",
				"AddNounName", "santa", "Santa",
				"AddNounTrait", "santa", "proper named",
				//
				"AddNounKind", "santana", "things",
				"AddNounName", "santana", "Santana",
				"AddNounTrait", "santana", "proper named",
				//
				"ApplyMacro", "suspect", "hector", "maria", "santa", "santana",
			},
		},
		// ------------------------------------------------------------------------
		// NamesAreLikeVerbs
		// ------------------------------------------------------------------------
		{
			test:   `A container is in the lobby.`,
			result: errors.New("this is specifically disallowed, and should generate an error"),
		},
		{
			// simple trait:
			test: `The bottle is closed.`,
			result: []string{
				// FIX? inform would error on this saying "Properties depend on kind"
				// because it would auto define a bottle as a thing; and things cant be "closed."
				"AddNounKind", "bottle", "things",
				"AddNounName", "bottle", "bottle",
				"AddNounTrait", "bottle", "closed",
			},
		},
		{
			// multi word trait:
			test: `The tree is fixed in place.`,
			result: []string{
				"AddNounKind", "tree", "things",
				"AddNounName", "tree", "tree",
				"AddNounTrait", "tree", "fixed in place",
			},
		},
		{
			// multiple trailing properties, using the kind as a property.
			test: `The bottle is a transparent, open, container.`,
			result: []string{
				"AddNounKind", "bottle", "containers",
				"AddNounName", "bottle", "bottle",
				"AddNounTrait", "bottle", "transparent",
				"AddNounTrait", "bottle", "open",
			},
		},
		{
			// multiple trailing properties without commas.
			test: `The bottle is a transparent open container.`,
			result: []string{
				"AddNounKind", "bottle", "containers",
				"AddNounName", "bottle", "bottle",
				"AddNounTrait", "bottle", "transparent",
				"AddNounTrait", "bottle", "open",
			},
		},
		{
			// multiple nouns of different kinds
			test: `The box and the top are closed containers.`,
			result: []string{
				"AddNounKind", "box", "containers",
				"AddNounName", "box", "box",
				"AddNounTrait", "box", "closed",
				"AddNounKind", "top", "containers",
				"AddNounName", "top", "top",
				"AddNounTrait", "top", "closed",
			},
		},
		{
			// using 'called' without a macro
			test: `The container called the sarcophagus is open.`,
			result: []string{
				"AddNounKind", "sarcophagus", "containers",
				"AddNounName", "sarcophagus", "sarcophagus",
				"AddNounTrait", "sarcophagus", "open",
				"AddNounValue", "sarcophagus", "indefinite article", text("the"),
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
			test: `The coffin is a closed container in the antechamber.`,
			result: []string{
				"AddNounKind", "antechamber", "things",
				"AddNounName", "antechamber", "antechamber",
				//
				"AddNounKind", "coffin", "containers",
				"AddNounName", "coffin", "coffin",
				"AddNounTrait", "coffin", "closed",
				//
				"ApplyMacro", "contain", "antechamber", "coffin",
			},
		},
		{
			// allowed even though it implies something different than what is written:
			test: `The bottle is openable in the kitchen.`,
			result: []string{
				"AddNounKind", "kitchen", "things",
				"AddNounName", "kitchen", "kitchen",
				//
				"AddNounKind", "bottle", "things",
				"AddNounName", "bottle", "bottle",
				"AddNounTrait", "bottle", "openable",
				"ApplyMacro", "contain", "kitchen", "bottle",
			},
		},
		// ------------------------------------------------------------------------
		// VerbNamesAreNames
		// ------------------------------------------------------------------------
		{
			// multiple primary with a leading macro
			test: `In the coffin are some coins, a notebook, and the gripping hand.`,
			result: []string{
				"AddNounKind", "coffin", "things",
				"AddNounName", "coffin", "coffin",
				//
				"AddNounKind", "coins", "things",
				"AddNounName", "coins", "coins",
				"AddNounTrait", "coins", "plural named",
				//
				"AddNounKind", "notebook", "things",
				"AddNounName", "notebook", "notebook",
				//
				"AddNounKind", "gripping hand", "things",
				"AddNounName", "gripping hand", "gripping hand",
				//
				"ApplyMacro", "contain", "coffin", "coins", "notebook", "gripping hand",
			},
		},
		{
			// multiple anonymous nouns.
			test: `In the lobby are a supporter and a container.`,
			result: []string{
				"AddNounKind", "lobby", "things",
				"AddNounName", "lobby", "lobby",
				//
				"AddNounKind", "supporter-1", "supporters",
				"AddNounName", "supporter-1", "supporter-1",
				"AddNounAlias", "supporter-1", "supporter",
				"AddNounTrait", "supporter-1", "counted",
				"AddNounValue", "supporter-1", "printed name", text("supporter"),
				//
				"AddNounKind", "container-1", "containers",
				"AddNounName", "container-1", "container-1",
				"AddNounAlias", "container-1", "container",
				"AddNounTrait", "container-1", "counted",
				"AddNounValue", "container-1", "printed name", text("container"),
				//
				"ApplyMacro", "contain", "lobby", "supporter-1", "container-1",
			},
		},
	}
	var skipped int
	for i, p := range phrases {
		if len(p.test) == 0 || p.result == nil {
			skipped++
		} else {
			var m Mock
			var haveRes []string
			var haveError error
			if e := jess.Generate(q, &m, p.test); e != nil {
				haveError = e
			} else if e := m.RunPost(q); e != nil {
				haveError = e
			} else {
				haveRes = m.out
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
					t.Logf("NG! test %d %q got: %#v", i, p.test, haveRes)
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

func Marshal(a any) (ret string, err error) {
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

func textKind(str, kind string) string {
	return fmt.Sprintf(`{"FromText:":{"Text value:kind:":[%q,%q]}}`, str, kind)
}

func number(num float64) string {
	return fmt.Sprintf(`{"FromNumber:":{"Num value:":%g}}`, num)
}
