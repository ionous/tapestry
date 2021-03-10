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
			allOf(words("under"), noun(), act("peeking")),
			allOf(noun(), act("examining")),
			act("looking"),

			// before "look inside", since inside is also direction.
			// allOf(noun(&parser.HasClass{"directions"}), act("examine")),
			// allOf(words("to"), noun(&parser.HasClass{"directions"}), act("examine")),
			// allOf(words("inside/in/into/through/on"), noun(), act("search")),
		)),
	allOf(words("pick"),
		anyOf(
			allOf(words("up"), things(), act("taking")),
			allOf(things(), words("up"), act("taking")),
		)),
	allOf(words("get"),
		&parser.Target{[]parser.Scanner{things(), words("from/off"), thing()}},
		act("removing"),
	),
	allOf(
		words("examine/x/watch/describe/check"), noun(), act("examining"),
	),
	allOf(words("read"),
		allOf(noun(), act("examining")),
	),
	allOf(words("open/unwrap/uncover"),
		allOf(noun(), act("opening")),
	),
	allOf(words("close/shut/cover"), anyOf(
		allOf(words("up"), act("closing")),
		allOf(noun(), act("closing")),
	),
	),
)
