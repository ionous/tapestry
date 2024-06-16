package debug_test

import (
	"testing"

	_ "embed"

	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

//go:embed fruity.json
var fruityJson string

//go:embed recordy.json
var recordyJson string

func TestStringify(t *testing.T) {
	var kinds testutil.Kinds
	type Rec struct {
		String  string
		Strings []string
		Float   float64
		Bool    bool
	}
	kinds.AddKinds((*Rec)(nil))
	// from TestRecordRecursion:
	type Fruit struct {
		Name string
		*Fruit
	}
	kinds.AddKinds((*Fruit)(nil))
	l1 := literal.RecordValue{
		KindName: "fruit",
		Fields: []literal.FieldValue{{
			FieldName: "name",
			Value:     &literal.TextValue{Value: "pomegranate"},
		}, {
			FieldName: "fruit",
			Value: &literal.RecordValue{
				KindName: "fruit",
				Fields: []literal.FieldValue{{
					FieldName: "name",
					Value:     &literal.TextValue{Value: "aril"},
				}},
			},
		}},
	}
	run := &testutil.Runtime{Kinds: &kinds}
	if rec, e := l1.GetRecord(run); e != nil {
		t.Fatal(e)
	} else if x := debug.Stringify(rec); x != fruityJson {
		t.Log("x", x)
		t.Fatal(files.Indent(x))
	}
	// other
	rec := kinds.NewRecord("rec",
		"string", "text",
		"strings", []string{"a", "b", "the end"},
		"float", 23.2,
		"bool", true,
	)
	if x := debug.Stringify(rt.RecordOf(rec)); x != recordyJson {
		t.Fatal(files.Indent(x))
	}
}
