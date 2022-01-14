package eph

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/kr/pretty"
)

// ensure we can create relations
// doesn't test contradictions, etc. since those are already tested by kinds
// ( and relation is just a special autogenerated kind )
func TestRelAssembly(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()

	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: kindsOf.Relation.String()}, // declare the relation table
		&EphKinds{Kinds: "p"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kinds: "q"},
		&EphRelations{Rel: "r", Cardinality: EphCardinality{
			Choice: EphCardinality_OneOne_Opt,
			Value:  &OneOne{Kind: "p", OtherKind: "q"},
		}},
		&EphRelations{Rel: "s", Cardinality: EphCardinality{
			Choice: EphCardinality_ManyOne_Opt,
			Value:  &ManyOne{Kinds: "p", OtherKind: "q"},
		}},
	)
	dt.makeDomain(dd("c", "b"),
		&EphRelations{Rel: "t", Cardinality: EphCardinality{
			Choice: EphCardinality_OneMany_Opt,
			Value:  &OneMany{Kind: "p", OtherKinds: "p"},
		}},
		&EphRelations{Rel: "u", Cardinality: EphCardinality{
			Choice: EphCardinality_OneMany_Opt,
			Value:  &ManyMany{Kinds: "q", OtherKinds: "p"},
		}},
	)
	if cat, e := buildAncestors(dt); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl.Rel}
		if e := cat.WriteRelations(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			"b:r:p:q:one_one:x",
			"b:s:p:q:any_one:x",
			"c:t:p:p:one_any:x",
			"c:u:q:p:any_any:x",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}
