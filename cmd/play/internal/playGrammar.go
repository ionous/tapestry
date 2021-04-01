package internal

import "git.sr.ht/~ionous/iffy/parser"

// this specifies a tree of parsing attempts
//
// you will need a "merge" tree at some point
// ex. so you could specify "look/l" multiple times as a starting node
// and then they get merged into a single tree.
var Grammar = anyOf(
	allOf(words("look/l"),
		anyOf(
			allOf(words("at"), noun(), act("examining")),
			allOf(words("in/inside/into/through"), noun(), act("searching")),
			allOf(words("under"), noun(), act("peeking")),
			allOf(noun(), act("examining")),
			act("looking"),

			// before "look inside", since inside is also direction.
			// allOf(noun(&parser.HasClass{"directions"}), act("examine")),
			// allOf(words("to"), noun(&parser.HasClass{"directions"}), act("examine")),
			// allOf(words("inside/in/into/through/on"), noun(), act("search")),
		)),
	// fix? inform defines these as synonyms, meaning this happens:
	// >carry off cat => "You aren't wearing the cat"
	allOf(words("take"), anyOf(
		allOf(words("inventory"), act("inventorying")),
		allOf(words("off"), thing(), act("removing")),
		anyOf(noun(),
			allOf(words("off"), act("removing")),
			act("taking"),
		))),
	allOf(words("carry/hold"), anyOf(
		allOf(words("inventory"), act("inventorying")),
		allOf(noun(), act("taking")),
	)),
	//  "take [things]" as taking.
	// "take inventory" as taking inventory.
	// "take off [something]" as taking off.
	// "take [something] off" as taking off.
	// "take [things inside] from [something]" as removing it from.
	// "take [things inside] off [something]" as removing it from.
	allOf(words("pick"), anyOf(
		allOf(words("up"), things(), act("taking")),
		allOf(things(), words("up"), act("taking")),
	)),
	allOf(words("put"), anyOf(
		allOf(words("on"), thing(), act("wearing")),
		allOf(thing(), words("on"), act("wearing")),
	)),
	allOf(words("get"),
		&parser.Target{[]parser.Scanner{things(), words("from/off"), thing()}},
		act("removing"),
	),
	allOf(
		words("examine/x/watch/describe/check"), noun(), act("examining"),
	),
	allOf(
		words("i/inv/inventory"), act("inventorying"),
	),
	allOf(
		words("read"), noun(), act("examining"),
	),
	allOf(
		words("search"), noun(), act("searching"),
	),
	allOf(words("remove/shed/doff/disrobe"),
		allOf(noun(), act("removing")),
	),
	allOf(words("open/unwrap/uncover"),
		allOf(noun(), act("opening")),
	),
	allOf(words("close/shut/cover"), anyOf(
		allOf(words("up"), act("closing")),
		allOf(noun(), act("closing"))),
	),
	allOf(words("wear/don"), anyOf(
		allOf(thing(), act("wearing"))),
	),
	allOf(words("switch/rotate/twist/unscrew/screw"), anyOf(
		allOf(words("on"), thing(), act("activating")),
		allOf(words("off"), thing(), act("deactivating")),
		allOf(thing(), anyOf(
			allOf(words("on"), act("activating")),
			allOf(words("off"), act("deactivating")),
		)),
	)),
	allOf(words("close/shut"), anyOf(
		allOf(words("off"), thing(), act("deactivating")),
		allOf(thing(), anyOf(
			allOf(words("off"), act("deactivating")),
		)),
	)),
)
