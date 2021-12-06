package eph

import (
	"testing"

	"github.com/kr/pretty"
)

// add some fields to a kind
func TestFields(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
		&EphFields{Kinds: "k", Name: "t", Class: "k", Affinity: Affinity{Affinity_Text}},
		&EphFields{Kinds: "k", Name: "n", Affinity: Affinity{Affinity_Number}},
	)
	out := testOut{mdl_field}
	if cat, e := buildFields(dt); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:k:t:text:k:x",
		"a:k:n:number::x",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// we can define a kind in one domain, and its fields in another
func TestFieldsCrossDomain(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
		&EphFields{Kinds: "k", Name: "n", Affinity: Affinity{Affinity_Number}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphFields{Kinds: "k", Name: "b", Affinity: Affinity{Affinity_Bool}},
	)
	out := testOut{mdl_field}
	if cat, e := buildFields(dt); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:k:n:number::x",
		"a:k:b:bool::x",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// we can redefine fields in the same domain, and in another
func TestFieldsRedefine(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()
	//
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphFields{Kinds: "k", Name: "n", Affinity: Affinity{Affinity_Number}},
		&EphKinds{Kinds: "k"},
		&EphFields{Kinds: "k", Name: "n", Affinity: Affinity{Affinity_Number}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphFields{Kinds: "k", Name: "n", Affinity: Affinity{Affinity_Number}},
	)
	out := testOut{mdl_field}
	if cat, e := buildFields(dt); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings.shift()); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings.shift()); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:k:n:number::x",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// fields conflict in sub-domains
// we can redefine fields in the same domain, and in another
func TestFieldsConflict(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
		&EphFields{Kinds: "k", Name: "n", Affinity: Affinity{Affinity_Number}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphFields{Kinds: "k", Name: "n", Affinity: Affinity{Affinity_Text}},
	)
	if _, e := buildFields(dt); e == nil {
		t.Fatal("expected error")
	} else if e := okDomainConflict("a", Redefined, e); e != nil {
		t.Fatal(e)
	} else {
		t.Log("ok:", e)
	}
}

// rival fields are fine so long as they match
// ( really the fields exist all at the same time )
func TestFieldsMatchingRivals(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()

	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
	)
	dt.makeDomain(dd("c", "a"),
		&EphFields{Kinds: "k", Name: "t", Affinity: Affinity{Affinity_Text}},
	)
	dt.makeDomain(dd("d", "a"),
		&EphFields{Kinds: "k", Name: "t", Affinity: Affinity{Affinity_Text}},
	)
	dt.makeDomain(dd("z", "c", "d"))
	out := testOut{mdl_field}
	if cat, e := buildFields(dt); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings.shift()); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:k:t:text::x",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}

// fields in kinds exist all at once --
// there's not really "rival" fields
func TestFieldsMismatchingRivals(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()

	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
	)
	dt.makeDomain(dd("c", "a"),
		&EphFields{Kinds: "k", Name: "t", Affinity: Affinity{Affinity_Text}},
	)
	dt.makeDomain(dd("d", "a"),
		&EphFields{Kinds: "k", Name: "t", Affinity: Affinity{Affinity_Bool}},
	)
	// dt.makeDomain(dd("z", "c", "d")) <-- fails even without this.
	if _, e := buildFields(dt); e == nil {
		t.Fatal("expected error, got:")
	} else if e := okDomainConflict("a", Redefined, e); e != nil {
		t.Fatal(e)
	} else {
		t.Log("ok:", e)
	}
}

// classes cant refer to kinds that dont exist.
func TestFieldsUnknownClass(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
		&EphFields{Kinds: "k", Name: "t", Class: "m", Affinity: Affinity{Affinity_Text}},
	)
	dt.makeDomain(dd("c", "a"),
		&EphKinds{Kinds: "m"},
	)
	if _, e := buildFields(dt); e == nil || e.Error() != `unknown class "m" for field "t" for kind "k"` {
		t.Fatal("expected error", e)
	} else {
		t.Log("ok:", e)
	}
}

func buildFields(dt domainTest) (ret *Catalog, err error) {
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		err = e
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
	}); e != nil {
		err = e
	} else {
		ret = &cat
	}
	return
}
