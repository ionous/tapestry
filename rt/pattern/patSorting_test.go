package pattern

import (
	"testing"

	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

func TestRuleSorting(t *testing.T) {
	ps := []rt.Rule{
		{RawFlags: float64(rt.Infix), Execute: Text("1")},
		{RawFlags: float64(rt.Postfix), Execute: Text("2")},
		{RawFlags: float64(rt.Prefix), Execute: Text("3")},
		{RawFlags: float64(-1) /*Filter: Skip,*/, Execute: Text("0")},
		{RawFlags: float64(rt.Postfix), Execute: Text("4")},
	}
	inds, flags := SortRules(ps)
	if flags != (rt.Infix | rt.Prefix | rt.Postfix) {
		t.Fatal("expected all flags set", flags)
	} else {
		var got string
		for _, i := range inds {
			got += string(ps[i].Execute.(Text))
		}
		if got != "3142" {
			t.Fatal("got", got)
		}
	}
}

type Text string

func (Text) Execute(rt.Runtime) error { return nil }

type Bool bool

func (b Bool) GetBool(rt.Runtime) (g.Value, error) {
	return g.BoolOf(bool(b)), nil
}

var Skip = Bool(false)
