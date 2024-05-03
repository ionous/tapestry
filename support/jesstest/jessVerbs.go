package jesstest

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
	"github.com/ionous/sliceOf"
)

// fix? maybe add "wearing" instead of carrying to test implications better?
var KnownVerbs = MockVerbs{
	"carrying": {
		Subject:  "actors",
		Object:   "things",
		Relation: "whereabouts",
		Implies:  sliceOf.String("not worn", "portable"),
		Reversed: false, // (parent) is carrying (child)
	},
	"carried by": {
		Subject:  "actors",
		Object:   "things",
		Relation: "whereabouts",
		Implies:  sliceOf.String("not worn", "portable"),
		Reversed: true, // (child) is carried by (parent)
	},
	"in": {
		Subject:   "containers",
		Alternate: "rooms", // alternate
		Object:    "things",
		Relation:  "whereabouts",
		Implies:   sliceOf.String("not worn"),
		Reversed:  true, // (child) is in (parent)
	},
	"on": {
		Subject:  "supporters",
		Object:   "things",
		Relation: "whereabouts",
		Implies:  sliceOf.String("not worn"),
		Reversed: true, // (child) is on (parent)
	},
	"suspicious of": {
		Subject:  "actors",
		Object:   "actors",
		Relation: "suspicion",
		Reversed: false, // (parent) is suspicious of (child)
	},
}

type MockVerbs map[string]jess.VerbDesc

func (vs MockVerbs) GetVerbValue(name, field string) (ret rt.Value, err error) {
	if v, ok := vs[name]; !ok {
		err = fmt.Errorf("%w %q %q", weaver.Missing, name, field)
	} else {
		switch field {
		case jess.VerbSubject:
			ret = rt.StringOf(v.Subject)
		case jess.VerbAlternate:
			ret = rt.StringOf(v.Alternate)
		case jess.VerbObject:
			ret = rt.StringOf(v.Object)
		case jess.VerbRelation:
			ret = rt.StringOf(v.Relation)
		case jess.VerbImplies:
			ret = rt.StringsOf(v.Implies)
		case jess.VerbReversed:
			var str string
			if v.Reversed {
				str = "reversed"
			} else {
				str = "not reversed"
			}
			ret = rt.StringOf(str)
		default:
			err = fmt.Errorf("%w %q %q", weaver.Missing, name, field)
		}
	}
	return
}
