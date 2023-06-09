package weave

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/weave/eph"
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
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if nouns, e := dt.readNouns(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(nouns, []string{
		"a:apple:k",
		"a:pear:k",
		"b:toy_boat:k",
	}); len(diff) > 0 {
		t.Log("nouns:", pretty.Sprint(nouns))
		t.Fatal(diff)
	} else if names, e := dt.readNames(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(names, []string{
		"a:apple:apple:0",
		"a:pear:pear:0",
		"b:toy_boat:toy boat:0",
		"b:toy_boat:toy_boat:1",
		"b:toy_boat:boat:2",
		"b:toy_boat:toy:3",
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
	_, e := dt.Assemble()
	if ok, e := okError(t, e, `Missing kind "t" in domain "a"`); !ok {
		t.Fatal("unexpected error:", e)
	} else {
		t.Log("ok:", e)
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
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if nouns, e := dt.readNouns(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(nouns, []string{
		// note: the new assembler takes the most specific type
		// the old one took the least specific type
		"a:apple:c",
		"a:bandanna:c",
		"a:pear:d",
	}); len(diff) > 0 {
		t.Log("nouns:", pretty.Sprint(nouns))
		t.Fatal(diff)
	} else if names, e := dt.readNames(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(names, []string{
		"a:apple:apple:0",
		"a:bandanna:bandanna:0",
		"a:pear:pear:0",
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
	_, e := dt.Assemble()
	if ok, e := okError(t, e, `Conflict can't redefine kind of "apple"`); !ok {
		t.Fatal("unexpected error:", e)
	} else {
		t.Log("ok:", e)
	}
}

func TestNounParts(t *testing.T) {
	dt := newTest(t.Name())
	defer dt.Close()
	dt.makeDomain(dd("a"),
		&eph.Kinds{Kind: "t"},
		&eph.Nouns{Noun: "collection of words", Kind: "t"},
	)
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if nouns, e := dt.readNouns(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(nouns, []string{
		"a:collection_of_words:t",
	}); len(diff) > 0 {
		t.Log("nouns:", pretty.Sprint(nouns))
		t.Fatal(diff)
	} else if names, e := dt.readNames(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(names, []string{
		"a:collection_of_words:collection of words:0",
		"a:collection_of_words:collection_of_words:1",
		"a:collection_of_words:words:2",
		"a:collection_of_words:of:3",
		"a:collection_of_words:collection:4",
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
	if _, e := dt.Assemble(); e != nil {
		t.Fatal(e)
	} else if names, e := dt.readNames(); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(names, []string{
		"b:apple:delicious:-1", // aliases first
		"b:apple:fruit:-1",
		"b:apple:apple:0",
		"b:toy_boat:model:-1", // aliases first
		"b:toy_boat:ship:-1",
		"b:toy_boat:toy boat:0", // spaces
		"b:toy_boat:toy_boat:1", // breaks
		"b:toy_boat:boat:2",     // left word
		"b:toy_boat:toy:3",      // right word
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
			if n, e := d.GetClosestNoun(name); e != nil {
				t.Error("couldnt get noun for name", name, e)
			} else if got := n.name; want != got {
				t.Errorf("wanted %q got %q", want, got)
			}
		}
	}
}
