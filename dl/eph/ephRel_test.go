package eph

import (
	"testing"

	"github.com/kr/pretty"
)

func TestRelAssembly(t *testing.T) {
	// var warnings Warnings
	// unwarn := warnings.catch(t)
	// defer unwarn()
	//
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: KindsOfRelation}, // declare the relation table
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
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
	}); e != nil {
		t.Fatal(e)
	} else {
		out := testOut{mdl_rel}
		if cat.WriteRelations(&out); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(out[1:], testOut{
			"r:p:one_one:q:x",
			"s:p:any_one:q:x",
			"t:p:one_any:p:x",
			"u:q:any_any:p:x",
		}); len(diff) > 0 {
			t.Log(pretty.Sprint(out))
			t.Fatal(diff)
		}
	}
}
