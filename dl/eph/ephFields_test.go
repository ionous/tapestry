package eph

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/kr/pretty"
)

// add some fields to a kind
func TestFields(t *testing.T) {
	dt := domainTest{noShuffle: true} // fields arent sorted
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: "k"},
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "t", Affinity: Affinity{Affinity_Text}, Class: "k"}}},
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "n", Affinity: Affinity{Affinity_Number}}}},
	)
	out := testOut{mdl.Field}
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e != nil {
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
	defer dt.Close() // fields arent sorted, but are probably added in domain order so...
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: "k"},
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "n", Affinity: Affinity{Affinity_Number}}}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "b", Affinity: Affinity{Affinity_Bool}}}},
	)
	out := testOut{mdl.Field}
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e != nil {
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
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "n", Affinity: Affinity{Affinity_Number}}}},
		&EphKinds{Kind: "k"},
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "n", Affinity: Affinity{Affinity_Number}}}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "n", Affinity: Affinity{Affinity_Number}}}},
	)
	out := testOut{mdl.Field}
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e != nil {
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
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: "k"},
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "n", Affinity: Affinity{Affinity_Number}}}},
	)
	dt.makeDomain(dd("b", "a"),
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "n", Affinity: Affinity{Affinity_Text}}}},
	)
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e == nil {
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
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: "k"},
	)
	dt.makeDomain(dd("c", "a"),
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "t", Affinity: Affinity{Affinity_Text}}}},
	)
	dt.makeDomain(dd("d", "a"),
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "t", Affinity: Affinity{Affinity_Text}}}},
	)
	dt.makeDomain(dd("z", "c", "d"))
	out := testOut{mdl.Field}
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e != nil {
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
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: "k"},
	)
	dt.makeDomain(dd("c", "a"),
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "t", Affinity: Affinity{Affinity_Text}}}},
	)
	dt.makeDomain(dd("d", "a"),
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "t", Affinity: Affinity{Affinity_Bool}}}},
	)
	// dt.makeDomain(dd("z", "c", "d")) <-- fails even without this.
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e == nil {
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
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: "k"},
		&EphKinds{Kind: "k", Contain: []EphParams{{Name: "t", Affinity: Affinity{Affinity_Text}, Class: "m"}}},
	)
	dt.makeDomain(dd("c", "a"),
		&EphKinds{Kind: "m"},
	)
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e == nil || e.Error() != `unknown class "m" for field "t" for kind "k"` {
		t.Fatal("expected error", e)
	} else {
		t.Log("ok:", e)
	}
}

// note: the original code would push shared fields upwards; the new code doesnt
func TestFieldLca(t *testing.T) {
	dt := domainTest{noShuffle: true} // fields arent sorted
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&EphKinds{Kind: "t"},
		&EphKinds{Kind: "p", Ancestor: "t"},
		&EphKinds{Kind: "q", Ancestor: "t"},
		//
		&EphKinds{Kind: "p", Contain: []EphParams{{Name: "t", Affinity: Affinity{Affinity_Text}}}},
		&EphKinds{Kind: "q", Contain: []EphParams{{Name: "t", Affinity: Affinity{Affinity_Text}}}},
	)
	out := testOut{mdl.Field}
	cat := NewCatalog(dt.Open(t.Name()))
	if e := dt.addToCat(cat); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e != nil {
		t.Fatal(e)
	} else if e := cat.WriteFields(&out); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out[1:], testOut{
		"a:p:t:text::x",
		"a:q:t:text::x",
	}); len(diff) > 0 {
		t.Log(pretty.Sprint(out))
		t.Fatal(diff)
	}
}
