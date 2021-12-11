package eph

import (
	"testing"

	"github.com/kr/pretty"
)

// /there's not much to test for directives right now
// just verify some text comes out?
func TestGrammarDirectives(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("b"),
		&EphDirectives{
			Name: `jump/skip/hop`,
			Type: `Directive`,
			Prog: `{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}`,
		},
	)
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(nil); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl_prog}
		if e := cat.WriteDirectives(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			`jump/skip/hop:Directive:{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}:x`,
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}
