package jess_test

import (
	"reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"github.com/kr/pretty"
)

func TestValidSplit(t *testing.T) {
	phrase := "reach underneath/under/beneath/-- [objects]"
	if res, e := jess.BuildPhrase(phrase); e != nil {
		t.Fatal(e)
	} else {
		expected := &grammar.Sequence{
			Series: []grammar.ScannerMaker{
				words("reach"),
				words("underneath", "under", "beneath", ""),
				&grammar.Noun{Kind: "objects"},
			},
		}
		if !reflect.DeepEqual(res, expected) {
			pretty.Println(res)
			t.Fatal("mismatch")
		}
	}
}

func TestInvalidChars(t *testing.T) {
	if _, e := jess.BuildPhrase("---"); e == nil {
		t.Fatal("expected failure")
	} else if _, e := jess.BuildPhrase("-"); e == nil {
		t.Fatal("expected failure")
	} else if _, e := jess.BuildPhrase("--"); e != nil {
		t.Fatal(e)
	} else if _, e := jess.BuildPhrase("hello [there"); e == nil {
		t.Fatal("expected failure")
	} else if _, e := jess.BuildPhrase("hello/[there]"); e == nil {
		t.Fatal("expected failure")
	}
}

// backslash would look odd, especially odd inside a json or tell file.
// so, arbitrarily, this uses equals as the escape.
// fix? brackets [\] would probably be better
func TestEscape(t *testing.T) {
	phrase := `=//==/=-/=[`
	if res, e := jess.BuildPhrase(phrase); e != nil {
		t.Fatal(e)
	} else {
		expected := &grammar.Sequence{
			Series: []grammar.ScannerMaker{
				words(`/`, `=`, `-`, `[`),
			},
		}
		if !reflect.DeepEqual(res, expected) {
			pretty.Println(res)
			t.Fatal("mismatch")
		}
	}
}

func words(strs ...string) *grammar.Words {
	return &grammar.Words{Words: strs}
}
