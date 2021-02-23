package internal

import "git.sr.ht/~ionous/iffy/parser"

var lookGrammar = allOf(words("look/l"), anyOf(
	allOf(act("look")),
	allOf(words("at"), noun(), act("examine")),
	// before "look inside", since inside is also direction.
	allOf(noun(&parser.HasClass{"directions"}), act("examine")),
	allOf(words("to"), noun(&parser.HasClass{"directions"}), act("examine")),
	allOf(words("inside/in/into/through/on"), noun(), act("search")),
	allOf(words("under"), noun(), act("look_under")),
))

var pickGrammar = allOf(words("pick"), anyOf(
	allOf(words("up"), things(), act("take")),
	allOf(things(), words("up"), act("take")),
))

var Grammar = anyOf(lookGrammar, pickGrammar)
