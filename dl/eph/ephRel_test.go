package eph

import (
	"testing"

	"git.sr.ht/~ionous/iffy/rt/kindsOf"
	"git.sr.ht/~ionous/iffy/tables/mdl"
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
		&EphRelations{"r", EphCardinality{
			EphCardinality_OneOne_Opt,
			&OneOne{"p", "q"},
		}},
		&EphRelations{"s", EphCardinality{
			EphCardinality_ManyOne_Opt,
			&ManyOne{"p", "q"},
		}},
	)
	dt.makeDomain(dd("c", "b"),
		&EphRelations{"t", EphCardinality{
			EphCardinality_OneMany_Opt,
			&OneMany{"p", "p"},
		}},
		&EphRelations{"u", EphCardinality{
			EphCardinality_OneMany_Opt,
			&ManyMany{"q", "p"},
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
