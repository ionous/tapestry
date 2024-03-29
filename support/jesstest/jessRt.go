package jesstest

import (
	"fmt"
	"io"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

// only expects a partial implementation;
// supporting a few bits of the runtime needed for jess
type jessRt struct {
	nounPairs map[string]string
	verbs     MockVerbs
}

func (d *jessRt) ActivateDomain(name string) error {
	panic("not implemented")
}
func (d *jessRt) GetKindByName(name string) (*g.Kind, error) {
	panic("not implemented")
}
func (d *jessRt) Call(name string, expectedReturn affine.Affinity, keys []string, vals []g.Value) (g.Value, error) {
	panic("not implemented")
}
func (d *jessRt) RelativesOf(a, relation string) (g.Value, error) {
	panic("not implemented")
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
func (d *jessRt) RelateTo(a, b, relation string) error {
	panic("not implemented")
}
func (d *jessRt) PushScope(rt.Scope) {
	panic("not implemented")
}
func (d *jessRt) PopScope() {
	panic("not implemented")
}
func (d *jessRt) GetField(name, field string) (ret g.Value, err error) {
	if str, e := d.verbs.GetVerbValue(name, field); e != nil {
		err = e
	} else {
		ret = g.StringOf(str)
	}
	return
}
func (d *jessRt) SetField(name, field string, value g.Value) error {
	panic("not implemented")
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
func (d *jessRt) Random(inclusiveMin, exclusiveMax int) int {
	panic("not implemented")
}
func (d *jessRt) Writer() io.Writer {
	panic("not implemented")
}
func (d *jessRt) SetWriter(io.Writer) (prev io.Writer) {
	panic("not implemented")
}
