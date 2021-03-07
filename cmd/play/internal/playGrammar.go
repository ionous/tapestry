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
			act("look"),
			allOf(noun(), act("look")),
			allOf(words("at"), noun(), act("examine")),
			// before "look inside", since inside is also direction.
			// allOf(noun(&parser.HasClass{"directions"}), act("examine")),
			// allOf(words("to"), noun(&parser.HasClass{"directions"}), act("examine")),
			// allOf(words("inside/in/into/through/on"), noun(), act("search")),
			// allOf(words("under"), noun(), act("look_under")),
		)),
	allOf(words("pick"),
		anyOf(
			allOf(words("up"), things(), act("take")),
			allOf(things(), words("up"), act("take")),
		)),
	allOf(words("get"),
		&parser.Target{[]parser.Scanner{things(), words("from/off"), thing()}},
		act("remove"),
	),
	allOf(
		words("examine/x/watch/describe/check"), noun(), act("examine"),
	),
	allOf(words("read"),
		allOf(noun(), act("examine")),
	),
)
