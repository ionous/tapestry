package eph

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/grammar"
	"git.sr.ht/~ionous/iffy/tables/mdl"
	"github.com/kr/pretty"
)

// /there's not much to test for directives right now
// just verify some text comes out?
func TestGrammarDirectives(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("b"),
		&EphDirectives{
			Name: `jump/skip/hop`,
			Directive: grammar.Directive{
				Lede: []string{"jump", "skip", "hop"},
				Scans: []grammar.ScannerMaker{
					&grammar.Action{"jumping"},
				},
			},
		},
	)
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(nil); e != nil {
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
