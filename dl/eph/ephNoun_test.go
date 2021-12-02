package eph

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// nouns
func TestNounNames(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()

	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "k"},
		&EphNouns{Noun: "apple", Kind: "k"},
		&EphNouns{Noun: "pear", Kind: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&EphNouns{Noun: "toy boat", Kind: "k"},
	)
	nouns, names := testOut{mdl_noun}, testOut{mdl_name}
	if e := writeNouns(dt, &nouns, &names); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(nouns[1:], testOut{
		"a:apple:k:x",
		"a:pear:k:x",
		"b:toy_boat:k:x",
	}); len(diff) > 0 {
		t.Log("nouns:", pretty.Sprint(nouns))
		t.Fatal(diff)
	} else if diff := pretty.Diff(names[1:], testOut{
		"a:apple:apple:0:x",
		"a:pear:pear:0:x",
		"b:toy_boat:toy boat:0:x",
		"b:toy_boat:toy_boat:1:x",
		"b:toy_boat:boat:2:x",
		"b:toy_boat:toy:3:x",
	}); len(diff) > 0 {
		t.Log("names:", pretty.Sprint(names))
		t.Fatal(diff)
	}
}

func TestNounFailure(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphNouns{Noun: "bad apple", Kind: "t"},
	)
	if e := writeNouns(dt, nil, nil); e == nil ||
		e.Error() != `error writing nouns: unknown dependency "t" for noun "bad_apple"` {
		t.Fatal("unexpected failure", e)
	} else {
		t.Log("ok", e)
	}
}

func TestNounHierarchy(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "t"},
		&EphKinds{Kinds: "p", From: "t"},
		&EphKinds{Kinds: "c", From: "p"},
		&EphKinds{Kinds: "d", From: "p"},
		&EphNouns{Noun: "apple", Kind: "c"},
		&EphNouns{Noun: "apple", Kind: "p"},
		&EphNouns{Noun: "pear", Kind: "d"},
		&EphNouns{Noun: "pear", Kind: "t"},
		&EphNouns{Noun: "bandanna", Kind: "c"},
	)
	nouns, names := testOut{mdl_noun}, testOut{mdl_name}
	if e := writeNouns(dt, &nouns, &names); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(nouns[1:], testOut{
		// note: the new assembler takes the most specific type
		// the old one took the least specific type
		"a:apple:c:x",
		"a:bandanna:c:x",
		"a:pear:d:x",
	}); len(diff) > 0 {
		t.Log("nouns:", pretty.Sprint(nouns))
		t.Fatal(diff)
	} else if diff := pretty.Diff(names[1:], testOut{
		"a:apple:apple:0:x",
		"a:bandanna:bandanna:0:x",
		"a:pear:pear:0:x",
	}); len(diff) > 0 {
		t.Log("names:", pretty.Sprint(names))
		t.Fatal(diff)
	}
}

func TestNounHierarchyFailure(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "t"},
		&EphKinds{Kinds: "c", From: "t"},
		&EphKinds{Kinds: "d", From: "t"},
		&EphNouns{Noun: "apple", Kind: "c"},
		&EphNouns{Noun: "apple", Kind: "d"},
	)
	if e := writeNouns(dt, nil, nil); e == nil ||
		e.Error() != `error writing nouns: "apple" has more than one parent` {
		t.Fatal("unexpected failure", e)
	} else {
		t.Log("ok", e)
	}
}

func TestNounParts(t *testing.T) {
	var dt domainTest
	dt.makeDomain(dd("a"),
		&EphKinds{Kinds: "t"},
		&EphNouns{Noun: "collection of words", Kind: "t"},
	)
	nouns, names := testOut{mdl_noun}, testOut{mdl_name}
	if e := writeNouns(dt, &nouns, &names); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(nouns[1:], testOut{
		"a:collection_of_words:t:x",
	}); len(diff) > 0 {
		t.Log("nouns:", pretty.Sprint(nouns))
		t.Fatal(diff)
	} else if diff := pretty.Diff(names[1:], testOut{
		"a:collection_of_words:collection of words:0:x",
		"a:collection_of_words:collection_of_words:1:x",
		"a:collection_of_words:words:2:x",
		"a:collection_of_words:of:3:x",
		"a:collection_of_words:collection:4:x",
	}); len(diff) > 0 {
		t.Log("names:", pretty.Sprint(names))
		t.Fatal(diff)
	}
}

func writeNouns(dt domainTest, nouns, names *testOut) (err error) {
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		err = e
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryPhaseActions,
	}); e != nil {
		err = e
	} else if e := cat.WriteNouns(nouns); e != nil {
		err = errutil.New("error writing nouns:", e)
	} else if e := cat.WriteNames(names); e != nil {
		err = errutil.New("error writing names:", e)
	}
	return
}
