package weave

import (
	"strconv"
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

func TestPatternRules(t *testing.T) {
	dt := newTestShuffle(t.Name(), false)
	defer dt.Close()
	dt.makeDomain(dd("a"),
		append(append([]eph.Ephemera{
			&eph.Kinds{Kind: kindsOf.Pattern.String()}, // declare the patterns table
			&eph.Patterns{
				PatternName: "p",
				Result: &eph.Params{
					Name:     "success",
					Affinity: affine.Bool,
				}},
		}, makeRules("p", "d", eph.Timing_During, 3)...),
			makeRules("p", "b", eph.Timing_Before, 2)...)...,
	)
	dt.makeDomain(dd("b", "a"),
		makeRules("p", "a", eph.Timing_After, 3)...,
	)
	//

	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl.Rule}
		if e := cat.WriteRules(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			// domain, pattern, target, phase, filter, prog, at
			`b:p::3:"filter-a2":["prog-a2"]:x`, // for now at least, we list domains backwards for rules.
			`b:p::3:"filter-a1":["prog-a1"]:x`,
			`b:p::3:"filter-a0":["prog-a0"]:x`,
			//
			`a:p::1:"filter-b1":["prog-b1"]:x`, // even though the `b` items were specified second, `before` rules get listed first.
			`a:p::1:"filter-b0":["prog-b0"]:x`,
			//
			`a:p::2:"filter-d2":["prog-d2"]:x`,
			`a:p::2:"filter-d1":["prog-d1"]:x`,
			`a:p::2:"filter-d0":["prog-d0"]:x`,
		}); len(diff) > 0 {
			t.Log("got:", pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}

// make some pattern like rules for testing
func makeRules(pattern, group, timing string, cnt int) (ret []eph.Ephemera) {
	for i := 0; i < cnt; i++ {
		ret = append(ret, &eph.Rules{
			PatternName: pattern,
			Filter:      filter{T("filter-" + group + strconv.Itoa(i))},
			When:        eph.Timing{timing},
			Exe:         []rt.Execute{prog{T("prog-" + group + strconv.Itoa(i))}},
		})
	}
	return
}

// a fake bool eval for testing that marshals itself as text
type filter struct{ *literal.TextValue }

func (filter) GetBool(rt.Runtime) (g.Value, error) { panic("not implemented") }

// a fake executable statement for testing that marshals itself as text
type prog struct{ *literal.TextValue }

func (prog) Execute(rt.Runtime) error { panic("not implemented") }
