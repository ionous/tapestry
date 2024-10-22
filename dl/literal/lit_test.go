package literal_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

func TestRecord(t *testing.T) {
	var kinds testutil.Kinds
	type Fruit struct {
		Name    string
		Variety string
	}
	kinds.AddKinds((*Fruit)(nil))
	l1 := literal.RecordValue{
		KindName: "fruit",
		Fields: textFields(
			"name", "apple",
		),
	}
	l2 := literal.RecordValue{
		KindName: "fruit",
		Fields: textFields(
			"variety", "gala",
			"name", "apple",
		),
	}
	run := &testutil.Runtime{Kinds: &kinds}
	if rec, e := l1.GetRecord(run); e != nil {
		t.Fatal(e)
	} else if v, e := rec.FieldByName("variety"); e != nil {
		t.Fatal(e)
	} else if n, e := rec.FieldByName("name"); e != nil {
		t.Fatal(e)
	} else if s := v.String(); s != "" {
		t.Fatal(s)
	} else if n := n.String(); n != "apple" {
		t.Fatal(s)
	} else if rec, e := l2.GetRecord(run); e != nil {
		t.Fatal(e)
	} else if v, e := rec.FieldByName("variety"); e != nil {
		t.Fatal(e)
	} else if n, e := rec.FieldByName("name"); e != nil {
		t.Fatal(e)
	} else if s := v.String(); s != "gala" {
		t.Fatal(s)
	} else if n := n.String(); n != "apple" {
		t.Fatal(s)
	}
	// fail
	l3 := literal.RecordValue{
		KindName: "fruit",
		Fields: textFields(
			"variety", "gala",
			"variety", "gala",
		),
	}
	if _, e := l3.GetRecord(run); e == nil {
		t.Fatal("expected failure with duplicate fields", e)
	}
}

func TestRecordRecursion(t *testing.T) {
	var kinds testutil.Kinds
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
	} else if v, e := rec.FieldByName("fruit"); e != nil {
		t.Fatal(e)
	} else if n, e := v.FieldByName("name"); e != nil {
		t.Fatal(e)
	} else if s := n.String(); s != "aril" {
		t.Fatal(s)
	}
}

func TestInnerRecord(t *testing.T) {
	var kinds testutil.Kinds
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
			Value: &literal.FieldList{
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
	} else if v, e := rec.FieldByName("fruit"); e != nil {
		t.Fatal(e)
	} else if n, e := v.FieldByName("name"); e != nil {
		t.Fatal(e)
	} else if s := n.String(); s != "aril" {
		t.Fatal(s)
	}
}

func textFields(els ...string) (ret []literal.FieldValue) {
	for i, cnt := 0, len(els); i < cnt; i += 2 {
		a, b := els[i], els[i+1]
		ret = append(ret, literal.FieldValue{
			FieldName: a,
			Value:     &literal.TextValue{Value: b},
		})
	}
	return
}
