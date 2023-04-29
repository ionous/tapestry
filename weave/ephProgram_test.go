package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/kr/pretty"
)

// there's not much to test for directives right now
// just verify some text comes out?
func TestGrammarDirectives(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("b"),
		&eph.Directives{
			Name: `jump/skip/hop`,
			Directive: grammar.Directive{
				Lede: []string{"jump", "skip", "hop"},
				Scans: []grammar.ScannerMaker{
					&grammar.Action{Action: "jumping"},
				},
			},
		},
	)

	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl.Grammar}
		if e := cat.WriteDirectives(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			`b:jump/skip/hop:{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}:x`,
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}
