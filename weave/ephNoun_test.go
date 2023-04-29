package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// test nouns and their names
func TestNounFormation(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()

	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Nouns{Noun: "apple", Kind: "k"},
		&eph.Nouns{Noun: "pear", Kind: "k"},
	)
	dt.makeDomain(dd("b", "a"),
		&eph.Nouns{Noun: "toy boat", Kind: "k"},
	)
	nouns, names := testOut{mdl.Noun}, testOut{mdl.Name}
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if e := writeNouns(cat, &nouns, &names); e != nil {
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
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Nouns{Noun: "bad apple", Kind: "t"},
	)
	if _, e := dt.Assemble(); e == nil ||
		e.Error() != `unknown kind t` {
		t.Fatal("unexpected failure", e)
	} else {
		t.Log("ok", e)
	}
}

func TestNounHierarchy(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "t"},
		&eph.Kinds{Kind: "p", Ancestor: "t"},
		&eph.Kinds{Kind: "c", Ancestor: "p"},
		&eph.Kinds{Kind: "d", Ancestor: "p"},
		&eph.Nouns{Noun: "apple", Kind: "c"},
		&eph.Nouns{Noun: "apple", Kind: "p"},
		&eph.Nouns{Noun: "pear", Kind: "d"},
		&eph.Nouns{Noun: "pear", Kind: "t"},
		&eph.Nouns{Noun: "bandanna", Kind: "c"},
	)
	nouns, names := testOut{mdl.Noun}, testOut{mdl.Name}
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if e := writeNouns(cat, &nouns, &names); e != nil {
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
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "t"},
		&eph.Kinds{Kind: "c", Ancestor: "t"},
		&eph.Kinds{Kind: "d", Ancestor: "t"},
		&eph.Nouns{Noun: "apple", Kind: "c"},
		&eph.Nouns{Noun: "apple", Kind: "d"},
	)
	if _, e := dt.Assemble(); e == nil ||
		e.Error() != `"apple" has more than one parent` {
		t.Fatal("unexpected failure", e)
	} else {
		t.Log("ok", e)
	}
}

func TestNounParts(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "t"},
		&eph.Nouns{Noun: "collection of words", Kind: "t"},
	)
	nouns, names := testOut{mdl.Noun}, testOut{mdl.Name}
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if e := writeNouns(cat, &nouns, &names); e != nil {
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

func TestNounAliases(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("b"),
		&eph.Kinds{Kind: "k"},
		&eph.Nouns{Noun: "toy boat", Kind: "k"},
		&eph.Nouns{Noun: "apple", Kind: "k"},
		&eph.Aliases{ShortName: "toy", Aliases: dd("model")},
		&eph.Aliases{ShortName: "boat", Aliases: dd("ship")},
		&eph.Aliases{ShortName: "apple", Aliases: dd("delicious", "fruit")},
	)
	names := testOut{mdl.Name}
	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if e := writeNouns(cat, nil, &names); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(names[1:], testOut{
		"b:apple:delicious:-1:x", // aliases first
		"b:apple:fruit:-1:x",
		"b:apple:apple:0:x",
		"b:toy_boat:model:-1:x", // aliases first
		"b:toy_boat:ship:-1:x",
		"b:toy_boat:toy boat:0:x", // spaces
		"b:toy_boat:toy_boat:1:x", // breaks
		"b:toy_boat:boat:2:x",     // left word
		"b:toy_boat:toy:3:x",      // right word
	}); len(diff) > 0 {
		t.Log("names:", pretty.Sprint(names))
		t.Fatal(diff)
	}
}

// simple words should pick out reasonable nouns
func TestNounDistance(t *testing.T) {
	var warnings Warnings
	unwarn := warnings.catch(t)
	defer unwarn()

	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "k"},
		&eph.Nouns{Noun: "toy boat", Kind: "k"},
		&eph.Nouns{Noun: "boat", Kind: "k"},
	)

	if cat, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if _, e := cat.ResolveNouns(); e != nil {
		t.Fatal(e)
	} else if d, ok := cat.GetDomain("a"); !ok {
		t.Fatal("unknown domain")
	} else {
		tests := []string{
			// word(s), noun(s)
			"toy", "toy_boat",
			"boat", "boat",
			"toy boat", "toy_boat",
		}
		for i, cnt := 0, len(tests); i < cnt; i += 2 {
			name, want := tests[i], tests[i+1]
			if n, ok := d.GetClosestNoun(name); !ok {
				t.Error("couldnt get noun for name", name)
			} else if got := n.Name(); want != got {
				t.Errorf("wanted %q got %q", want, got)
			}
		}
	}
}

func writeNouns(cat *Catalog, nouns, names *testOut) (err error) {
	if nouns != nil && err == nil {
		if e := cat.WriteNouns(nouns); e != nil {
			err = errutil.New("error writing nouns:", e)
		}
	}
	if names != nil && err == nil {
		if e := cat.WriteNames(names); e != nil {
			err = errutil.New("error writing names:", e)
		}
	}
	return
}
