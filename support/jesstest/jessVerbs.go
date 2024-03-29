package jesstest

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

type MockVerb struct {
	Subject, Object, Alternate, Relation, Implies string
	Reversed                                      bool
}

// fix? maybe add "wearing" instead of carrying, to test implication better?
var KnownVerbs = MockVerbs{
	"carrying": {
		Subject:  "actors",
		Object:   "things",
		Relation: "whereabouts",
		Implies:  "not worn",
		Reversed: false, // (parent) is carrying (child)
	},
	"carried by": {
		Subject:  "actors",
		Object:   "things",
		Relation: "whereabouts",
		Implies:  "not worn",
		Reversed: true, // (child) is carried by (parent)
	},
	"in": {
		Subject:   "containers",
		Alternate: "rooms", // alternate
		Object:    "things",
		Relation:  "whereabouts",
		Implies:   "not worn",
		Reversed:  true, // (child) is in (parent)
	},
	"on": {
		Subject:  "supporters",
		Object:   "things",
		Relation: "whereabouts",
		Implies:  "not worn",
		Reversed: true, // (child) is on (parent)
	},
	"suspicious of": {
		Subject:  "actors",
		Object:   "actors",
		Relation: "suspicion",
		Reversed: false, // (parent) is suspicious of (child)
	},
}

type MockVerbs map[string]MockVerb

func (vs MockVerbs) GetVerbValue(name, field string) (ret string, err error) {
	if v, ok := vs[name]; !ok {
		err = fmt.Errorf("%w %q %q", weaver.Missing, name, field)
	} else {
		str := "$bad"
		switch field {
		case jess.VerbSubject:
			str = v.Subject
		case jess.VerbAlternate:
			str = v.Alternate
		case jess.VerbObject:
			str = v.Object
		case jess.VerbRelation:
			str = v.Relation
		case jess.VerbImplication:
			str = v.Implies
		case jess.VerbReversed:
			if v.Reversed {
				str = "reversed"
			} else {
				str = "not reversed"
			}
		}
		if str == "bad" {
			err = fmt.Errorf("%w %q %q", weaver.Missing, name, field)
		} else {
			ret = str
		}
	}
	return
}
