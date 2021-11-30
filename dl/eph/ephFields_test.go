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
	var out testOut
	if cat, e := buildFields(dt); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, testOut{
		"a:k:t:text:k:2",
		"a:k:n:number::3",
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
	var out testOut
	if cat, e := buildFields(dt); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, testOut{
		"a:k:n:number::2",
		"a:k:b:bool::5",
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
	var out testOut
	if cat, e := buildFields(dt); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings.shift()); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings.shift()); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, testOut{
		"a:k:n:number::1",
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
		t.Fatal("expected a different error", e)
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
	var out testOut
	if cat, e := buildFields(dt); e != nil {
		t.Fatal(e)
	} else if e := okDomainConflict("a", Duplicated, warnings.shift()); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, testOut{
		"a:k:t:text::4",
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
		t.Fatal("expected error")
	} else if e := okDomainConflict("a", Redefined, e); e != nil {
		t.Fatal("expected a different error", e)
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
	if _, e := buildFields(dt); e == nil || e.Error() != `unknown field class m for kind "k" in domain "a"` {
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
		AncestryPhase: PhaseAction{
			PhaseFlags{NoDuplicates: true},
			func(d *Domain) error {
				_, e := d.ResolveKinds()
				return e
			},
		},
	}); e != nil {
		err = e
	} else {
		ret = &cat
	}
	return
}
