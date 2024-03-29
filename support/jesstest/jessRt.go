package jesstest

import (
	"fmt"

	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

// only expects a partial implementation;
// supporting a few bits of the runtime needed for jess
type jessRt struct {
	nounPairs map[string]string
	verbs     VerbLookup
	testutil.PanicRuntime
}

type VerbLookup interface {
	GetVerbValue(name, field string) (ret string, err error)
}

func (d *jessRt) ReciprocalsOf(b, relation string) (ret g.Value, err error) {
	if relation != "whereabouts" {
		err = fmt.Errorf("unexpected relation %q", relation)
	} else if a, ok := d.nounPairs[b]; !ok {
		ret = g.StringsOf(nil)
	} else {
		ret = g.StringsOf([]string{a})
	}
	return
}
func (d *jessRt) GetField(name, field string) (ret g.Value, err error) {
	if str, e := d.verbs.GetVerbValue(name, field); e != nil {
		err = e
	} else {
		ret = g.StringOf(str)
	}
	return
}
func (d *jessRt) PluralOf(single string) string {
	return inflect.Pluralize(single)
}
func (d *jessRt) SingularOf(plural string) string {
	return inflect.Singularize(plural)
}
func (d *jessRt) OppositeOf(word string) (ret string) {
	switch word {
	case "north":
		ret = "south"
	case "south":
		ret = "north"
	case "east":
		ret = "west"
	case "west":
		ret = "east"
	default:
		ret = word
	}
	return
}
