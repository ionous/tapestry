package jesstest

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
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
	GetVerbValue(name, field string) (rt.Value, error)
}

func (d *jessRt) ReciprocalsOf(b, relation string) (ret rt.Value, err error) {
	if relation != "whereabouts" {
		err = fmt.Errorf("jess ret, unexpected relation %q", relation)
	} else if a, ok := d.nounPairs[b]; !ok {
		ret = rt.StringsOf(nil)
	} else {
		ret = rt.StringsOf([]string{a})
	}
	return
}
func (d *jessRt) GetField(name, field string) (ret rt.Value, err error) {
	// ex. asking for north.opposite
	if field == "opposite" {
		var str string
		switch name {
		case "north":
			str = "south"
		case "south":
			str = "north"
		case "east":
			str = "west"
		case "west":
			str = "east"
		default:
			err = fmt.Errorf("jess rt, unexpected opposite for %q", name)
		}
		ret = rt.StringOf(str)
	} else if name == meta.KindAncestry {
		if field == "storing" {
			// root left, kind right.
			ret = rt.StringsOf([]string{"kinds", "actions", "storing"})
		}
	} else {
		ret, err = d.verbs.GetVerbValue(name, field)
	}
	return
}
func (d *jessRt) PluralOf(single string) string {
	return inflect.Pluralize(single)
}
func (d *jessRt) SingularOf(plural string) string {
	return inflect.Singularize(plural)
}
