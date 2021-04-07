package internal

// this specifies a tree of parsing attempts
//
// you will need a "merge" tree at some point
// ex. so you could specify "look/l" multiple times as a starting node
// and then they get merged into a single tree.
var Grammar = anyOf(
	// fix? inform defines these as synonyms, meaning this happens:
	// >carry off cat => "You aren't wearing the cat"
	allOf(words("carry/hold"), anyOf(
		allOf(words("inventory"), act("inventorying")),
		allOf(noun(), act("taking")),
	)),
	allOf(words("close"), anyOf(
		allOf(words("off"), thing(), act("deactivating")),
		allOf(words("up"), thing(), act("closing")),
		// note: this might split here based on "preferable traits", etc.
		allOf(noun(), anyOf(
			allOf(words("off"), act("deactivating")),
			allOf(words("up"), act("closing")),
			act("closing"),
		)),
	)),
	allOf(words("cover"), anyOf(
		allOf(words("up"), thing(), act("closing")),
		allOf(noun(), act("closing")),
	)),
	allOf(words("doff/disrobe/shed"),
		noun(), act("removing"),
	),
	allOf(words("drop"),
		thing(), anyOf(
			allOf(words("in/into/down/on/onto"), thing(), act("storing")),
			act("dropping"),
		)),
	allOf(words("examine/x/watch/describe/check"),
		noun(), act("examining"),
	),
	allOf(words("get"), anyOf(
		allOf(retarget(things(), words("from/off"), thing()), act("retrieving")),
		allOf(things(), act("taking")),
	)),
	allOf(words("i/inv/inventory"),
		act("inventorying"),
	),
	allOf(words("insert"),
		thing(), words("in/into"), thing(), act("storing"),
	),
	allOf(words("look/l"),
		anyOf(
			allOf(words("at"), noun(), act("examining")),
			allOf(words("in/inside/into/through/on"), noun(), act("searching")),
			allOf(words("under"), noun(), act("peeking")),
			allOf(noun(), act("examining")),
			act("looking"),

			// before "look inside", since inside is also direction.
			// allOf(noun(&parser.HasClass{"directions"}), act("examining")),
			// allOf(words("to"), noun(&parser.HasClass{"directions"}), act("examining")),
		)),
	allOf(words("open/unwrap/uncover"),
		allOf(noun(), act("opening")),
	),
	allOf(words("pick"), anyOf(
		allOf(words("up"), things(), act("taking")),
		allOf(things(), words("up"), act("taking")),
	)),
	allOf(words("put"), anyOf(
		allOf(words("on"), thing(), act("wearing")),
		// FIX: plural things. in inform this is "[other things; things inside, etc.]"
		allOf(thing(), anyOf(
			allOf(words("in/inside/into/onto"), thing(), act("storing")),
			allOf(words("on"), anyOf(
				allOf(thing(), act("storing")),
				act("wearing"),
			)),
		)),
	)),
	allOf(words("read"),
		noun(), act("examining"),
	),
	allOf(words("remove"), anyOf(
		allOf(retarget(noun(), words("from"), noun()), act("retrieving")),
		allOf(noun(), act("removing")),
	)),
	allOf(words("search"),
		noun(), act("searching"),
	),
	allOf(words("switch/rotate/twist/unscrew/screw"), anyOf(
		allOf(words("on"), thing(), act("activating")),
		allOf(words("off"), thing(), act("deactivating")),
		allOf(thing(), anyOf(
			allOf(words("on"), act("activating")),
			allOf(words("off"), act("deactivating")),
		)),
	)),
	allOf(words("shut"), anyOf(
		allOf(words("off"), thing(), act("deactivating")),
		allOf(words("up"), thing(), act("closing")),
		allOf(thing(), anyOf(
			allOf(words("off"), act("deactivating")),
			allOf(words("up"), act("closing")),
			act("closing"),
		)),
	)),
	//  "take [things]" as taking.
	// "take inventory" as taking inventory.
	// "take off [something]" as taking off.
	// "take [something] off" as taking off.
	// "take [things inside] from [something]" as removing it from.
	// "take [things inside] off [something]" as removing it from.
	allOf(words("take"), anyOf(
		allOf(words("inventory"), act("inventorying")),
		allOf(words("off"), thing(), act("removing")),
		allOf(retarget(noun(), words("off/from"), thing(), act("removing"))),
		allOf(noun(), act("taking")),
	)),
	allOf(words("wait/z"),
		act("waiting"),
	),
	allOf(words("wear/don"),
		thing(), act("wearing"),
	),
)
