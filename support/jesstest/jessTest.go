// Package jesstest exercises implementations of jess.Query
// to ensure they produce good results.
// ( see for instance package jessdb_test, and package jess_test. )
package jesstest

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/files"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

var Phrases = []Phrase{

	// ------------------------------------------------------------------------
	// TimedRule
	// ( storing is one of the predefined patterns. )
	// ------------------------------------------------------------------------
	{
		// tbd: would it make more sense to have these end with colon here?
		// currently, its read by the tokenizer it assumes the next bit a a tell document
		// there's no tell document here; just the partial phrase for testing.
		test:   `Instead of storing`,
		assign: true,
		result: []string{
			"ExtendPattern:", "before storing",
			"Rule:", "<unnamed>, Stop:true, Jump:JumpNow, Updates:false",
		},
	},
	{
		// since actions are always actor initiated,
		// the LHS ( the "someone" between prefix and the verb ) can only really be an actor
		// when its nothing, it implies the player;
		// but it could be a specific subclass.
		test:   `Report someone storing (this is the noisy storage rule)`,
		assign: true,
		result: []string{
			"ExtendPattern:", "after storing",
			"Rule:", "noisy storage, Stop:true, Jump:JumpNow, Updates:false",
		},
	},
	{
		// matching works ... now need to generate some ... filters....
		// and log them.... serialize to json maybe? ugh.
		test:   `Report someone storing the message`,
		assign: true,
		result: []string{
			"ExtendPattern:", "after storing",
			"Rule:", "<unnamed>, Stop:true, Jump:JumpNow, Updates:false",
		},
	},

	// tbd: "when" is the default; should we force it to be specified?
	// if so, how does that interact with "when <domain> begins|ends"

	// after advancing time, then continue
	// before someone attacking, then stop

	// todo:
	// can test the other prefixes and suffixes directly somewhere:
	// matching one set here is enough.

	// ------------------------------------------------------------------------
	// MapDirections
	// ------------------------------------------------------------------------
	{
		// FIX: some punctuation should be allowed in the names i guess.
		// ( "woodcutter's")
		test: `North of the Meadow is the woodcutter hut.`,
		result: []string{
			// nouns; defaulting to rooms:
			"AddNounName:", "meadow", "Meadow",
			"AddNounName:", "woodcutter hut", "woodcutter hut",
			"AddNounKind:", "meadow", "rooms",
			"AddNounKind:", "woodcutter hut", "rooms",
			// movement via a private door:
			"AddFact:", "dir", "meadow", "north", "woodcutter hut",
			"AddNounValue:", "meadow", "compass.north", textKind("meadow-north-door", "doors"),
			"AddNounKind:", "meadow-north-door", "doors",
			"AddNounName:", "meadow-north-door", "meadow-north-door",
			"AddNounTrait:", "meadow-north-door", "scenery",
			"AddNounTrait:", "meadow-north-door", "privately named",
			"AddNounValue:", "meadow-north-door", "destination", textKind("woodcutter hut", "rooms"),
			"AddNounPair:", "whereabouts", "meadow", "meadow-north-door",
			// implies reverse direction via another private door:
			"AddFact:", "dir", "woodcutter hut", "south", "meadow",
			"AddNounValue:", "woodcutter hut", "compass.south", textKind("woodcutter-hut-south-door", "doors"),
			"AddNounKind:", "woodcutter-hut-south-door", "doors",
			"AddNounName:", "woodcutter-hut-south-door", "woodcutter-hut-south-door",
			"AddNounTrait:", "woodcutter-hut-south-door", "scenery",
			"AddNounTrait:", "woodcutter-hut-south-door", "privately named",
			"AddNounValue:", "woodcutter-hut-south-door", "destination", textKind("meadow", "rooms"),
			"AddNounPair:", "whereabouts", "woodcutter hut", "woodcutter-hut-south-door",
		},
	},
	{
		// redirection
		test: `West of the Garden is south of the Meadow.`,
		result: []string{
			// nouns; defaulting to rooms:
			"AddNounName:", "garden", "Garden",
			"AddNounName:", "meadow", "Meadow",
			"AddNounKind:", "garden", "rooms",
			"AddNounKind:", "meadow", "rooms",
			// movement via a private door:
			"AddFact:", "dir", "garden", "west", "meadow",
			"AddNounValue:", "garden", "compass.west", textKind("garden-west-door", "doors"),
			"AddNounKind:", "garden-west-door", "doors",
			"AddNounName:", "garden-west-door", "garden-west-door",
			"AddNounTrait:", "garden-west-door", "scenery",
			"AddNounTrait:", "garden-west-door", "privately named",
			"AddNounValue:", "garden-west-door", "destination", textKind("meadow", "rooms"),
			"AddNounPair:", "whereabouts", "garden", "garden-west-door",
			// redirect via a another private door:
			"AddFact:", "dir", "meadow", "south", "garden",
			"AddNounValue:", "meadow", "compass.south", textKind("meadow-south-door", "doors"),
			"AddNounKind:", "meadow-south-door", "doors",
			"AddNounName:", "meadow-south-door", "meadow-south-door",
			"AddNounTrait:", "meadow-south-door", "scenery",
			"AddNounTrait:", "meadow-south-door", "privately named",
			"AddNounValue:", "meadow-south-door", "destination", textKind("garden", "rooms"),
			"AddNounPair:", "whereabouts", "meadow", "meadow-south-door",
		},
	},
	// ------------------------------------------------------------------------
	// MapLocations
	// starts with a single room or door, and maps to rooms or doors.
	// note: in inform, doors can live in multiple places;
	// in tapestry, they cannot.
	// ------------------------------------------------------------------------
	{
		// unless otherwise specified, both lhs and rhs default to rooms.
		test: `The passageway is south of the kitchen.`,
		result: []string{
			// nouns; defaulting to rooms:
			"AddNounName:", "passageway", "passageway",
			"AddNounName:", "kitchen", "kitchen",
			"AddNounKind:", "passageway", "rooms",
			"AddNounKind:", "kitchen", "rooms",
			// room to room conflict detection.
			"AddFact:", "dir", "kitchen", "south", "passageway",
			// private door leading south to the passageway
			"AddNounValue:", "kitchen", "compass.south", textKind("kitchen-south-door", "doors"),
			"AddNounKind:", "kitchen-south-door", "doors",
			"AddNounName:", "kitchen-south-door", "kitchen-south-door",
			"AddNounTrait:", "kitchen-south-door", "scenery",
			"AddNounTrait:", "kitchen-south-door", "privately named",
			"AddNounValue:", "kitchen-south-door", "destination", textKind("passageway", "rooms"),
			"AddNounPair:", "whereabouts", "kitchen", "kitchen-south-door",
			// room to room conflict detection.
			"AddFact:", "dir", "passageway", "north", "kitchen",
			// private door leading north to the kitchen
			"AddNounValue:", "passageway", "compass.north", textKind("passageway-north-door", "doors"),
			"AddNounKind:", "passageway-north-door", "doors",
			"AddNounName:", "passageway-north-door", "passageway-north-door",
			"AddNounTrait:", "passageway-north-door", "scenery",
			"AddNounTrait:", "passageway-north-door", "privately named",
			"AddNounValue:", "passageway-north-door", "destination", textKind("kitchen", "rooms"),
			"AddNounPair:", "whereabouts", "passageway", "passageway-north-door",
		},
	},
	{
		// ensure that a door can be used on the lhs.
		// but note: doors cant be mapped to multiple places.
		// ( ie, ng: The door is south of one place and north of another place. )
		// fix: verify the use of nowhere
		test: `The passageway is south of the kitchen. The passageway is a door.`,
		result: []string{
			"AddNounName:", "passageway", "passageway",
			"AddNounName:", "kitchen", "kitchen",
			"AddNounKind:", "passageway", "doors",
			"AddNounKind:", "kitchen", "rooms",
			"AddNounPair:", "whereabouts", "kitchen", "passageway",
			// the direction is what the direction says.
			"AddNounValue:", "kitchen", "compass.south", textKind("passageway", "doors"),
		},
	},
	{
		// ensure that a door can be used on the rhs,
		test: `The kitchen is south of the passageway. The passageway is a door.`,
		result: []string{
			"AddNounName:", "kitchen", "kitchen",
			"AddNounName:", "passageway", "passageway",
			"AddNounKind:", "passageway", "doors",
			"AddNounKind:", "kitchen", "rooms",
			"AddNounPair:", "whereabouts", "kitchen", "passageway",
			// if the room is south from the door; then the door is north from room.
			"AddNounValue:", "kitchen", "compass.north", textKind("passageway", "doors"),
		},
	},
	{
		// the door can located in another room.
		test: `The kitchen is south of the passageway. The door called the passageway is in the room called the basement.`,
		result: []string{
			"AddNounName:", "kitchen", "kitchen",
			"AddNounName:", "passageway", "passageway",
			"AddNounKind:", "passageway", "doors",
			"AddNounKind:", "basement", "rooms", // it can immediately determine this, so kind is before name.
			"AddNounName:", "basement", "basement",
			"AddNounPair:", "whereabouts", "basement", "passageway",
			"AddNounKind:", "kitchen", "rooms", // kitchen defaults to room because its mentioned in a directional phrase
			// movement via the existing door
			// the kitchen is south of the passageway,
			// and the passageway is in the basement,
			// so moving south from the basement arrives in the kitchen.
			"AddFact:", "dir", "basement", "south", "kitchen",
			"AddNounValue:", "passageway", "destination", textKind("kitchen", "rooms"),
			"AddNounValue:", "basement", "compass.south", textKind("passageway", "doors"),
			// implies reverse direction via a private door:
			// ie. north from the kitchen is the basement.
			"AddFact:", "dir", "kitchen", "north", "basement",
			"AddNounValue:", "kitchen", "compass.north", textKind("kitchen-north-door", "doors"),
			"AddNounKind:", "kitchen-north-door", "doors",
			"AddNounName:", "kitchen-north-door", "kitchen-north-door",
			"AddNounTrait:", "kitchen-north-door", "scenery",
			"AddNounTrait:", "kitchen-north-door", "privately named",
			"AddNounValue:", "kitchen-north-door", "destination", textKind("basement", "rooms"),
			"AddNounPair:", "whereabouts", "kitchen", "kitchen-north-door",
			// ugh. so horrible.
			// inform has various checks about doors needing locations and exits,
			// so it'd be hard or impossible to declare a worn door due to other restrictions.
			"AddNounTrait:", "passageway", "not worn",
			// fix? inform seems to *always* set indefinite article for "called the" ( re: passageway )
			// but tapestry only does so when the noun is new.
			// it's possible that inform elevates "called the" sub-phrases ahead of everything else.
			"AddNounValue:", "basement", "indefinite article", text("the"),
		},
	},
	{
		// doors can't exist on both sides of this equation.
		test:   `The gate is south of the passageway. The passageway is a door. The gate is a door.`,
		result: errors.New("nope"),
	},
	{
		// the phrase also works with multiple rooms.
		// fix? the uniform dash join creates some small potential for conflict
		test: `The mystery spot is west of the waterfall and south of the sea.`,
		result: []string{
			// the places mentiond:
			"AddNounName:", "mystery spot", "mystery spot",
			"AddNounName:", "waterfall", "waterfall",
			"AddNounName:", "sea", "sea",
			// all default to rooms:
			"AddNounKind:", "mystery spot", "rooms",
			"AddNounKind:", "waterfall", "rooms",
			"AddNounKind:", "sea", "rooms",
			// movement via a private door:
			"AddFact:", "dir", "waterfall", "west", "mystery spot",
			// maybe leading the name with a dash would be good enough.
			"AddNounValue:", "waterfall", "compass.west", textKind("waterfall-west-door", "doors"),
			"AddNounKind:", "waterfall-west-door", "doors",
			"AddNounName:", "waterfall-west-door", "waterfall-west-door",
			"AddNounTrait:", "waterfall-west-door", "scenery",
			"AddNounTrait:", "waterfall-west-door", "privately named",
			"AddNounValue:", "waterfall-west-door", "destination", textKind("mystery spot", "rooms"),
			"AddNounPair:", "whereabouts", "waterfall", "waterfall-west-door",
			/// implies reverse direction via another private door:
			"AddFact:", "dir", "mystery spot", "east", "waterfall",
			"AddNounValue:", "mystery spot", "compass.east", textKind("mystery-spot-east-door", "doors"),
			"AddNounKind:", "mystery-spot-east-door", "doors",
			"AddNounName:", "mystery-spot-east-door", "mystery-spot-east-door",
			"AddNounTrait:", "mystery-spot-east-door", "scenery",
			"AddNounTrait:", "mystery-spot-east-door", "privately named",
			"AddNounValue:", "mystery-spot-east-door", "destination", textKind("waterfall", "rooms"),
			"AddNounPair:", "whereabouts", "mystery spot", "mystery-spot-east-door",
			// movement via a private door:
			"AddFact:", "dir", "sea", "south", "mystery spot",
			"AddNounValue:", "sea", "compass.south", textKind("sea-south-door", "doors"),
			"AddNounKind:", "sea-south-door", "doors",
			"AddNounName:", "sea-south-door", "sea-south-door",
			"AddNounTrait:", "sea-south-door", "scenery",
			"AddNounTrait:", "sea-south-door", "privately named",
			"AddNounValue:", "sea-south-door", "destination", textKind("mystery spot", "rooms"),
			"AddNounPair:", "whereabouts", "sea", "sea-south-door",
			// implies reverse direction via another private door:
			"AddFact:", "dir", "mystery spot", "north", "sea",
			"AddNounValue:", "mystery spot", "compass.north", textKind("mystery-spot-north-door", "doors"),
			"AddNounKind:", "mystery-spot-north-door", "doors",
			"AddNounName:", "mystery-spot-north-door", "mystery-spot-north-door",
			"AddNounTrait:", "mystery-spot-north-door", "scenery",
			"AddNounTrait:", "mystery-spot-north-door", "privately named",
			"AddNounValue:", "mystery-spot-north-door", "destination", textKind("sea", "rooms"),
			"AddNounPair:", "whereabouts", "mystery spot", "mystery-spot-north-door",
		},
	},
	// ------------------------------------------------------------------------
	// MapConnection
	// ------------------------------------------------------------------------
	{
		test: `Through the long slide is nowhere.`,
		result: []string{
			"AddNounKind:", "long slide", "doors", // explicitly a door
			"AddNounName:", "long slide", "long slide",
			"AddNounValue:", "long slide", "destination", textKind("", "rooms"),
		},
	},
	{
		test: `Through the blue door is the Flat Landing.`,
		result: []string{
			"AddNounKind:", "flat landing", "rooms",
			"AddNounName:", "flat landing", "Flat Landing",
			"AddNounKind:", "blue door", "doors",
			"AddNounName:", "blue door", "blue door",
			"AddNounValue:", "blue door", "destination", textKind("flat landing", "rooms"),
		},
	},
	{
		test: `Through the gate and the hatch is a dark room called An End.`,
		result: []string{
			"AddNounKind:", "an end", "rooms",
			"AddNounName:", "an end", "An End",
			"AddNounKind:", "gate", "doors",
			"AddNounName:", "gate", "gate",
			"AddNounKind:", "hatch", "doors",
			"AddNounName:", "hatch", "hatch",
			"AddNounTrait:", "an end", "proper named",
			"AddNounTrait:", "an end", "dark",
			"AddNounValue:", "gate", "destination", textKind("an end", "rooms"),
			"AddNounValue:", "hatch", "destination", textKind("an end", "rooms"),
		},
	},
	// ----------------
	// multiple sentences
	// ----------------
	{
		// the definition as a container should win out over the default of a thing.
		// names_verb_names, names_are_like_verbs.
		test: `The bottle is in the kitchen. The bottle is a container. The kitchen is a room.`,
		result: []string{
			// the registration of the placeholder kind is hidden
			"AddNounName:", "bottle", "bottle",
			"AddNounName:", "kitchen", "kitchen",
			// before applying any relations the kinds are finalized
			// ( sorted by name )
			"AddNounKind:", "bottle", "containers",
			"AddNounKind:", "kitchen", "rooms",
			// all done.
			"AddNounPair:", "whereabouts", "kitchen", "bottle",
			"AddNounTrait:", "bottle", "not worn",
		},
	},
	{
		// understandings need nouns;
		// define a new noun, and provide an understanding for it.
		// as per inform: the nouns can be defined after the understanding references them.
		test: `Understand "donut" as the doughnut. The doughnut is a thing.`,
		result: []string{
			"AddNounKind:", "doughnut", "things",
			"AddNounName:", "doughnut", "doughnut",
			"AddNounAlias:", "doughnut", "donut",
		},
	},

	// ------------------------------------------------------------------------
	// Understandings
	// ------------------------------------------------------------------------
	{
		// message and missive are built-in test nouns.
		test: `Understand "floor" or "sawdust" as the message.`,
		result: []string{
			"AddNounAlias:", "message", "floor",
			"AddNounAlias:", "message", "sawdust",
		},
	},
	{
		// message and missive are built-in test nouns.
		test: `Understand "missives" as the plural of missive and message.`,
		result: []string{
			"AddPlural:", "missives", "missive",
			"AddPlural:", "missives", "message",
		},
	},
	{
		// storing is a built-in for testing.
		test: `Understand "hang [objects] on/onto/-- [objects]" as storing.`,
		result: []string{
			"AddGrammar:",
			"hang [objects] on/onto/-- [objects]",
			`{"One word:":["hang"]}`,
			`{"One noun:":"objects"}`,
			`{"One word:":["on","onto",""]}`,
			`{"One noun:":"objects"}`,
			`{"Action:":"storing"}`,
		},
	},
	// ------------------------------------------------------------------------
	// PropertyPropertyPropertyNounValue
	// ------------------------------------------------------------------------
	{
		test: `The title of the story is "A Secret."`,
		result: []string{
			// no noun declaration because the story is a known noun ( in these tests )
			"AddNounValue:", "story", "title", text("A Secret."),
		},
	},
	{
		// note: we don't validate properties while matching
		// weave validates them when attempting to write them.
		test: `The age of the bottle is 42.`,
		result: []string{
			"AddNounName:", "bottle", "bottle",
			"AddNounKind:", "bottle", "things",
			"AddNounValue:", "bottle", "age", number(42),
		},
	},
	{
		// create a new noun "story teller";
		// shouldnt match the existing noun "story."
		test: `The age of the story teller is 42.`,
		result: []string{
			"AddNounName:", "story teller", "story teller",
			"AddNounKind:", "story teller", "things",
			"AddNounValue:", "story teller", "age", number(42),
		},
	},
	{
		// fix: currently succeeds with a noun called "thing called the cat"
		// inform gets confused, but we could handle this okay
		test: `The description of the thing called the cat is "meow."`,
		// result: errors.New("can't use property noun value this way."),
	},
	// ------------------------------------------------------------------------
	// NounPropertyValue
	// note: we don't validate the values of properties while matching;
	// weave validates values when attempting to write them to the db.
	// ------------------------------------------------------------------------
	{
		// fix? mixed feelings on the trailing full-stop.
		test: `The story has the title "{15|print_num!}".`,
		result: []string{
			// test that it can convert a template
			"AddNounValue:", "story", "title", `{"FromText:":{"Print num:":{"Num value:":15}}}`,
		},
	},
	{
		// inform specifically disallows this:
		// "the 'of' here appears superfluous"
		// that seems silly to me.
		test: `The bottle has an age of 42.`,
		result: []string{
			"AddNounName:", "bottle", "bottle",
			"AddNounKind:", "bottle", "things",
			"AddNounValue:", "bottle", "age", number(42),
		},
	},
	{
		// fix: inform allows this, and jess does not.
		test: `The thing called the cat has the description "meow."`,
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
			"AddAspectTraits:", "color", "red", "blue", "cobalt",
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
			"AddKindTrait:", "containers", "closed",
		},
	},
	{
		test: `Containers and supporters are usually fixed in place.`,
		result: []string{
			"AddKindTrait:", "containers", "fixed in place",
			"AddKindTrait:", "supporters", "fixed in place",
		},
	},
	// ------------------------------------------------------------------------
	// KindsOf
	// ------------------------------------------------------------------------
	{
		// when the requested kind being declared isn't yet known:
		test: `Devices are a kind of thing.`,
		result: []string{
			"AddKind:", "devices", "things",
		},
	},
	{
		// when the kind being declared is already known
		test: `A container is a kind of thing.`,
		result: []string{
			"AddKind:", "containers", "things",
		},
	},
	{
		// determine pluralness based on trailing is/are
		test: `A device is a kind of thing.`,
		result: []string{
			"AddKind:", "devices", "things",
		},
	},
	{
		// adding trailing properties
		test: `A casket is a kind of closed container.`,
		result: []string{
			"AddKind:", "caskets", "containers",
			"AddKindTrait:", "caskets", "closed",
		},
	},
	{
		// complex parsing
		test: `The closed containers called the safes are a kind of fixed in place thing.`,
		result: []string{
			"AddKind:", "safes", "things",
			"AddKind:", "safes", "containers",
			"AddKindTrait:", "safes", "closed",
			"AddKindTrait:", "safes", "fixed in place",
		},
	},
	{
		// correctly producing unexpected results.
		test: `The closed casket is a kind of container.`,
		result: []string{
			"AddKind:", "closed caskets", "containers",
		},
	},
	// ------------------------------------------------------------------------
	// KindsAreEither
	// A thing can be [either]
	// Things can be [either] tall or short.
	// ------------------------------------------------------------------------
	{
		test: "Containers are either opened or closed.",
		result: []string{
			// inform names these: "container status", "container status 2", etc.
			"AddKind:", "opened status", "aspects",
			"AddAspectTraits:", "opened status", "opened", "closed",
			"AddKindFields:", "containers", "opened status", "text", "opened status",
		},
	},
	{
		// can be [either] ...
		test: "A thing can be opened or closed or ajar.",
		result: []string{
			"AddKind:", "opened status", "aspects",
			"AddAspectTraits:", "opened status", "opened", "closed", "ajar",
			"AddKindFields:", "things", "opened status", "text", "opened status",
		},
	},
	{
		// inform doesnt allow [either] here; i'm fine with whatever works out.
		test: "A thing can be scenery.",
		result: []string{
			"AddKindFields:", "things", "scenery", "bool", "",
		},
	},
	// ------------------------------------------------------------------------
	// KindsHaveProperties
	// inform allows: things have some text called the one, and the two.
	// i feel like if that's the case. then better would be forcing the type as well
	// so you can mix numbers and text; maybe even "can be".
	// ------------------------------------------------------------------------
	{
		test:   "Things have some text called a description.",
		result: []string{"AddKindFields:", "things", "description", "text", ""},
	},
	{
		test:   "Things have some text.",
		result: errors.New("unnamed text fields are prohibited"),
	},
	{
		test:   "Things have a number.",
		result: errors.New("unnamed number fields are prohibited"),
	},
	{
		test:   "A supporter has a number called carrying capacity.",
		result: []string{"AddKindFields:", "supporters", "carrying capacity", "num", ""},
	},
	{
		// except for number and text, inform allows "bare" properties: a "list of text" creates a member called "list of text"
		test:   "Things have a list of text called frenemies.",
		result: []string{"AddKindFields:", "things", "frenemies", "text_list", ""},
	},
	{
		test:   "Things have a list of numbers called the lotto numbers.",
		result: []string{"AddKindFields:", "things", "lotto numbers", "num_list", ""},
	},
	{
		// groups are a pre-defined type of record; anonymous fields are allowed.
		// references to kinds become text; except for records which are embedded.
		test:   "Things have a color.",
		result: []string{"AddKindFields:", "things", "color", "text", "color"},
	},
	{
		test:   "Things have a group.",
		result: []string{"AddKindFields:", "things", "group", "record", "groups"},
	},
	{
		test:   "Things have a group called the set.",
		result: []string{"AddKindFields:", "things", "set", "record", "groups"},
	},
	{
		test:   "Things have a list of groups.",
		result: []string{"AddKindFields:", "things", "groups", "record_list", "groups"},
	},
	// ------------------------------------------------------------------------
	// NamesVerbNames
	// ------------------------------------------------------------------------
	{
		// "in" should use "containers" by default.
		test: `Two things are in the bottle.`,
		result: []string{
			"AddNounKind:", "thing-1", "things",
			"AddNounName:", "thing-1", "thing-1",
			"AddNounKind:", "thing-2", "things",
			"AddNounName:", "thing-2", "thing-2",
			"AddNounName:", "bottle", "bottle",
			//
			"AddNounPair:", "whereabouts", "bottle", "thing-1",
			"AddNounPair:", "whereabouts", "bottle", "thing-2",
			"AddNounKind:", "bottle", "containers", // fallback implied by "in"
			//
			"AddNounAlias:", "thing-1", "thing",
			"AddNounTrait:", "thing-1", "counted",
			"AddNounTrait:", "thing-1", "not worn",
			"AddNounValue:", "thing-1", "printed name", text("thing"),
			//
			"AddNounAlias:", "thing-2", "thing",
			"AddNounTrait:", "thing-2", "counted",
			"AddNounTrait:", "thing-2", "not worn",
			"AddNounValue:", "thing-2", "printed name", text("thing"),
		},
	},
	{
		test: `Hershel is carrying scissors and a pen.`,
		result: []string{
			"AddNounName:", "hershel", "Hershel",
			"AddNounName:", "scissors", "scissors",
			"AddNounName:", "pen", "pen",
			//
			"AddNounKind:", "hershel", "actors",
			"AddNounKind:", "scissors", "things",
			"AddNounKind:", "pen", "things",
			//
			"AddNounPair:", "whereabouts", "hershel", "scissors",
			"AddNounPair:", "whereabouts", "hershel", "pen",
			//
			"AddNounTrait:", "hershel", "proper named",
			"AddNounTrait:", "scissors", "proper named", // yes; this conforms with inform.
			"AddNounTrait:", "scissors", "not worn",
			"AddNounTrait:", "scissors", "portable",
			"AddNounTrait:", "pen", "not worn",
			"AddNounTrait:", "pen", "portable",
		},
	},
	{
		// reverse carrying relation.
		test: `The scissors and a pen are carried by Hershel.`,
		result: []string{
			"AddNounName:", "scissors", "scissors",
			"AddNounName:", "pen", "pen",
			"AddNounName:", "hershel", "Hershel",
			//
			"AddNounKind:", "scissors", "things",
			"AddNounKind:", "pen", "things",
			"AddNounKind:", "hershel", "actors",
			//
			"AddNounPair:", "whereabouts", "hershel", "scissors",
			"AddNounPair:", "whereabouts", "hershel", "pen",
			//
			"AddNounTrait:", "scissors", "not worn",
			"AddNounTrait:", "scissors", "portable",
			"AddNounTrait:", "pen", "not worn",
			"AddNounTrait:", "pen", "portable",
			"AddNounTrait:", "hershel", "proper named",
		},
	},
	{
		test: `The unhappy man is in the closed bottle.`,
		result: []string{
			"AddNounName:", "unhappy man", "unhappy man",
			"AddNounName:", "closed bottle", "closed bottle",
			"AddNounKind:", "unhappy man", "things",
			"AddNounPair:", "whereabouts", "closed bottle", "unhappy man",
			"AddNounKind:", "closed bottle", "containers", // the implications of the verb "is in"
			"AddNounTrait:", "unhappy man", "not worn",
		},
	},
	{
		// called both before and after the macro
		// note: The closed openable container called the trunk and the box is in the lobby.
		// would create a noun called "the trunk and the box"
		test: `The thing called the stake is on the supporter called the altar.`,
		result: []string{
			"AddNounKind:", "stake", "things",
			"AddNounName:", "stake", "stake",
			//
			"AddNounKind:", "altar", "supporters",
			"AddNounName:", "altar", "altar",
			//
			"AddNounPair:", "whereabouts", "altar", "stake",
			//
			"AddNounTrait:", "stake", "not worn",
			"AddNounValue:", "stake", "indefinite article", text("the"),
			"AddNounValue:", "altar", "indefinite article", text("the"),
		},
	},
	{
		// add leading properties using 'called'
		// "is" left of the macro "in".
		// slightly different parsing than "kind/s of":
		// those expect only expect one set of nouns; these have two.
		test: `A closed openable container called the trunk is in the lobby.`,
		result: []string{
			"AddNounKind:", "trunk", "containers",
			"AddNounName:", "trunk", "trunk",
			//
			"AddNounName:", "lobby", "lobby",
			"AddNounPair:", "whereabouts", "lobby", "trunk",
			"AddNounKind:", "lobby", "containers", // in defaults to containers
			//
			"AddNounTrait:", "trunk", "closed",
			"AddNounTrait:", "trunk", "openable",
			"AddNounTrait:", "trunk", "not worn",
			"AddNounValue:", "trunk", "indefinite article", text("the"),
		},
	},
	{
		// multiple primary: "is" left of the macro "in".
		test: `Some coins, a notebook, and the gripping hand are in the coffin.`,
		result: []string{
			"AddNounName:", "coins", "coins",
			"AddNounName:", "notebook", "notebook",
			"AddNounName:", "gripping hand", "gripping hand",
			"AddNounName:", "coffin", "coffin",
			//
			"AddNounKind:", "coins", "things",
			"AddNounKind:", "notebook", "things",
			"AddNounKind:", "gripping hand", "things",
			//
			"AddNounPair:", "whereabouts", "coffin", "coins",
			"AddNounPair:", "whereabouts", "coffin", "notebook",
			"AddNounPair:", "whereabouts", "coffin", "gripping hand",
			"AddNounKind:", "coffin", "containers", // fallback implied by "in"
			//
			"AddNounTrait:", "coins", "plural named",
			"AddNounTrait:", "coins", "not worn",
			"AddNounTrait:", "notebook", "not worn",
			"AddNounTrait:", "gripping hand", "not worn",
		},
	},
	{
		// the special nxn description: no properties are allowed.
		test: `Hector and Maria are suspicious of Santa and Santana.`,
		result: []string{
			"AddNounName:", "hector", "Hector",
			"AddNounName:", "maria", "Maria",
			"AddNounName:", "santa", "Santa",
			"AddNounName:", "santana", "Santana",
			//
			"AddNounKind:", "hector", "actors",
			"AddNounKind:", "maria", "actors",
			"AddNounKind:", "santa", "actors",
			"AddNounKind:", "santana", "actors",
			//
			"AddNounPair:", "suspicion", "hector", "santa",
			"AddNounPair:", "suspicion", "hector", "santana",
			"AddNounPair:", "suspicion", "maria", "santa",
			"AddNounPair:", "suspicion", "maria", "santana",
			//
			"AddNounTrait:", "hector", "proper named",
			"AddNounTrait:", "maria", "proper named",
			"AddNounTrait:", "santa", "proper named",
			"AddNounTrait:", "santana", "proper named",
		},
	},
	{
		test:   `A container is in the lobby.`,
		result: errors.New("this is specifically disallowed, and should generate an error"),
	},
	// ------------------------------------------------------------------------
	// NamesAreLikeVerbs
	// ------------------------------------------------------------------------
	{
		// simple trait:
		test: `The bottle is closed.`,
		result: []string{
			// FIX? inform would error on this saying "Properties depend on kind"
			// because it would auto define a bottle as a thing; and things cant be "closed."
			"AddNounName:", "bottle", "bottle",
			"AddNounKind:", "bottle", "things",
			"AddNounTrait:", "bottle", "closed",
		},
	},
	{
		// multi word trait:
		test: `The tree is fixed in place.`,
		result: []string{
			"AddNounName:", "tree", "tree",
			"AddNounKind:", "tree", "things",
			"AddNounTrait:", "tree", "fixed in place",
		},
	},
	{
		// multiple trailing properties, using the kind as a property.
		test: `The bottle is a transparent, open, container.`,
		result: []string{
			"AddNounKind:", "bottle", "containers",
			"AddNounName:", "bottle", "bottle",
			"AddNounTrait:", "bottle", "transparent",
			"AddNounTrait:", "bottle", "open",
		},
	},
	{
		// multiple trailing properties without commas.
		test: `The bottle is a transparent open container.`,
		result: []string{
			"AddNounKind:", "bottle", "containers",
			"AddNounName:", "bottle", "bottle",
			"AddNounTrait:", "bottle", "transparent",
			"AddNounTrait:", "bottle", "open",
		},
	},
	{
		// multiple nouns of different kinds
		test: `The box and the top are closed containers.`,
		result: []string{
			"AddNounKind:", "box", "containers",
			"AddNounName:", "box", "box",
			//
			"AddNounKind:", "top", "containers",
			"AddNounName:", "top", "top",
			//
			"AddNounTrait:", "box", "closed",
			"AddNounTrait:", "top", "closed",
		},
	},
	{
		// using 'called' without a macro
		// fix: inform specifically errors on commas "the kitchen contains a thing called the one, two."
		test: `The container called the sarcophagus is open.`,
		result: []string{
			"AddNounKind:", "sarcophagus", "containers",
			"AddNounName:", "sarcophagus", "sarcophagus",
			"AddNounTrait:", "sarcophagus", "open",
			"AddNounValue:", "sarcophagus", "indefinite article", text("the"),
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
			"AddKind:", "buckets", "containers",
			"AddKind:", "baskets", "containers",
		},
	},
	{
		// inform doesnt allow this;
		// it requires: The closed container called the coffin is in the antechamber.
		test: `The coffin is a closed container in the antechamber.`,
		result: []string{
			"AddNounKind:", "coffin", "containers",
			"AddNounName:", "coffin", "coffin",
			"AddNounName:", "antechamber", "antechamber",
			"AddNounPair:", "whereabouts", "antechamber", "coffin",
			"AddNounKind:", "antechamber", "containers", // fallback implied by "in"
			"AddNounTrait:", "coffin", "closed",
			"AddNounTrait:", "coffin", "not worn",
		},
	},
	{
		// allowed even though it implies something different than what is written:
		test: `The bottle is openable in the kitchen.`,
		result: []string{
			"AddNounName:", "bottle", "bottle",
			"AddNounName:", "kitchen", "kitchen",
			//
			"AddNounKind:", "bottle", "things",
			"AddNounPair:", "whereabouts", "kitchen", "bottle",
			"AddNounKind:", "kitchen", "containers", // fallback implied by "in"
			//
			"AddNounTrait:", "bottle", "openable",
			"AddNounTrait:", "bottle", "not worn",
		},
	},
	// ------------------------------------------------------------------------
	// VerbNamesAreNames
	// ------------------------------------------------------------------------
	{
		// multiple primary with a leading macro
		test: `In the coffin are some coins, a notebook, and the gripping hand.`,
		result: []string{
			"AddNounName:", "coins", "coins",
			"AddNounName:", "notebook", "notebook",
			"AddNounName:", "gripping hand", "gripping hand",
			"AddNounName:", "coffin", "coffin",
			//
			"AddNounKind:", "coins", "things",
			"AddNounKind:", "notebook", "things",
			"AddNounKind:", "gripping hand", "things",
			//
			"AddNounPair:", "whereabouts", "coffin", "coins",
			"AddNounPair:", "whereabouts", "coffin", "notebook",
			"AddNounPair:", "whereabouts", "coffin", "gripping hand",
			"AddNounKind:", "coffin", "containers", // fallback implied from "in"
			//
			"AddNounTrait:", "coins", "plural named",
			"AddNounTrait:", "coins", "not worn",
			"AddNounTrait:", "notebook", "not worn",
			"AddNounTrait:", "gripping hand", "not worn",
		},
	},
	{
		// multiple anonymous nouns.
		test: `In the lobby are a supporter and a container.`,
		result: []string{
			"AddNounKind:", "supporter-1", "supporters",
			"AddNounName:", "supporter-1", "supporter-1",
			"AddNounKind:", "container-1", "containers",
			"AddNounName:", "container-1", "container-1",
			// fix? at this point, lobby is an object
			// whereabouts happens to be be object, objects
			// but .... what if it weren't?
			// we'd have to apply the relationship types to make things work.
			// -- connections need to read pairs to determine parent
			// -- so pairs need to be written before / as part of connections
			// -- and fallbacks have to happen after connections
			"AddNounName:", "lobby", "lobby",
			"AddNounPair:", "whereabouts", "lobby", "supporter-1",
			"AddNounPair:", "whereabouts", "lobby", "container-1",
			"AddNounKind:", "lobby", "containers", // generated by fallback phase
			//
			"AddNounAlias:", "supporter-1", "supporter",
			"AddNounTrait:", "supporter-1", "counted",
			"AddNounTrait:", "supporter-1", "not worn",
			"AddNounValue:", "supporter-1", "printed name", text("supporter"),
			//
			"AddNounAlias:", "container-1", "container",
			"AddNounTrait:", "container-1", "counted",
			"AddNounTrait:", "container-1", "not worn",
			"AddNounValue:", "container-1", "printed name", text("container"),
		},
	},
}

type Phrase struct {
	test   string
	assign bool
	result any
}

// returns empty if no result is expected (ie. a skip)
func (p *Phrase) Test() (string, rt.Assignment, bool) {
	var a rt.Assignment
	if p.assign {
		// we only need to test matching and generation; parsing sub docs is elsewhere
		a = &call.FromExe{}
	}
	return p.test, a, p.result != nil
}

// verify a standard set of phrases using some function that takes each of those phrases
func (p *Phrase) Verify(haveRes []string, haveError error) (okay bool) {
	if expectError, ok := p.result.(error); ok {
		if haveError != nil {
			log.Println("ok, test", p.test, haveError)
			okay = true
		} else {
			log.Println("NG! test", p.test, "expected an error", expectError, "but succeeded with", pretty.Sprint(haveRes))
		}
	} else if haveError != nil {
		log.Println("NG! test", p.test, haveError)
	} else {
		if d := pretty.Diff(p.result, haveRes); len(d) == 0 {
			okay = true
		} else {
			var out strings.Builder
			for _, str := range haveRes {
				if strings.HasSuffix(str, ":") {
					out.WriteString("\n\t")
				}
				out.WriteRune('"')
				out.WriteString(str)
				out.WriteString(`", `)
			}
			log.Printf("NG! test %q got:%s", p.test, out.String())

			log.Println(d)
		}
	}
	return
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
	return fmt.Sprintf(`{"Text value:":%q}`, str)
}

func textKind(str, kind string) string {
	return fmt.Sprintf(`{"Text kind:value:":[%q,%q]}`, kind, str)
}

func number(num float64) string {
	return fmt.Sprintf(`{"Num value:":%g}`, num)
}
