package list_test

import (
	"reflect"
	"strconv"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/scope"
	"git.sr.ht/~ionous/tapestry/test"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

// test sorting by record fields
// doesn't test all sorting options
// ( relies on the list sorter tests for that )
func TestSortRecords(t *testing.T) {
	// records
	var kinds testutil.Kinds
	els := buildTestRecords(t, kinds.AddKind((*Item)(nil)))
	// sort by string
	var txt []string
	must(t, list.Sorter{Kind: "item", Field: "text"}.SortRecords(&kinds, els))
	for _, el := range els {
		v, e := el.GetNamedField("text")
		must(t, e)
		txt = append(txt, v.String())
	}
	if !reflect.DeepEqual(txt, []string{"xerox", "xylophone", "zebra"}) {
		t.Fatal("failed txt sorting", txt)
	}
	// sort by number
	var num []float64
	must(t, list.Sorter{Kind: "item", Field: "number"}.SortRecords(&kinds, els))
	for _, el := range els {
		v, e := el.GetNamedField("number")
		must(t, e)
		num = append(num, v.Float())
	}
	if !reflect.DeepEqual(num, []float64{2, 7, 23}) {
		t.Fatal("failed num sorting", num)
	}
}

// sort by a field within an object.
func TestSortObjects(t *testing.T) {
	var kinds testutil.Kinds
	// the objects ( test objects are records )
	els := buildTestRecords(t, kinds.AddKind((*Item)(nil)))
	objs := testutil.Objects{
		"desmond": els[2], // {Text: "xylophone", Number: 7}
		"mildred": els[0], // {Text: "zebra", Number: 2}
		"yuri":    els[1], // {Text: "xerox", Number: 23}
	}
	// local variable containing a list of object names
	names := []string{"yuri", "desmond", "mildred"}
	kinds.AddKinds((*Locals)(nil))
	locals := kinds.NewRecord("locals", "names", names)
	// a runtime ( for accessing the locals )
	run := testutil.Runtime{
		Kinds:   &kinds,
		Objects: objs,
		Chain:   scope.MakeChain(scope.FromRecord(&kinds, locals)),
	}
	// command to sort the list by string field
	must(t, safe.Run(&run, &list.ListSort{
		Target:    object.Variable("names"),
		KindName:  literal.T("item"),
		FieldName: literal.T("text"),
	}))
	// verify the results
	if !reflect.DeepEqual(names, []string{"yuri", "desmond", "mildred"}) {
		t.Fatal("failed sorting by string field", names)
	}
	// command to sort the list by number field
	must(t, safe.Run(&run, &list.ListSort{
		Target:    object.Variable("names"),
		KindName:  literal.T("item"),
		FieldName: literal.T("number"),
	}))
	if !reflect.DeepEqual(names, []string{"mildred", "desmond", "yuri"}) {
		t.Fatal("failed sorting by field of number", names)
	}
	// command to sort the list by object name
	must(t, safe.Run(&run, &list.ListSort{
		Target: object.Variable("names"),
	}))
	if !reflect.DeepEqual(names, []string{"desmond", "mildred", "yuri"}) {
		t.Fatal("failed sorting by field of string", names)
	}
}

// sort by the value of an aspect.
func TestSortByAspect(t *testing.T) {
	var kinds testutil.Kinds
	var objs testutil.Objects

	// register the classes needed for this test
	kinds.AddKinds((*test.Messages)(nil), (*Locals)(nil))
	msgKind := kinds.Kind("messages")

	// create some objects and set their neatness values
	// ( assign them backwards so we can sort them forward )
	// obj0:trampled(2), obj1:scuffed(1), obj2:neat(0)
	var names []string
	for i := range test.NumNeatness {
		name := "obj" + strconv.Itoa(i)
		obj := objs.AddObject(msgKind, name)
		neat := test.Neatness(test.NumNeatness - i - 1)
		obj.SetNamedField("neatness", rt.StringOf(neat.String()))
		names = append(names, name)
	}
	// runtime
	run := testutil.Runtime{
		Kinds:   &kinds,
		Objects: objs,
		Chain: scope.MakeChain(scope.FromRecord(&kinds,
			kinds.NewRecord("locals", "names", names))),
	}
	// sorts in place
	sorter := &list.ListSort{
		Target:    object.Variable("names"),
		KindName:  literal.T("messages"),
		FieldName: literal.T("neatness"),
	}
	must(t, safe.Run(&run, sorter))
	if !reflect.DeepEqual(names, []string{
		"obj2", "obj1", "obj0",
	}) {
		t.Fatal("ascending", names)
	}
	// swap to descending
	sorter.Descending = literal.B(true)
	must(t, safe.Run(&run, sorter))
	if !reflect.DeepEqual(names, []string{
		"obj0", "obj1", "obj2",
	}) {
		t.Fatal("descending", names)
	}
}

// a simple noun
type Item struct {
	Text   string
	Number float64
}

// global variables for grouping tests
type Locals struct {
	Names []string
}

func buildTestRecords(t *testing.T, k *rt.Kind) (ret []*rt.Record) {
	for _, item := range []Item{
		{Text: "zebra", Number: 2},
		{Text: "xerox", Number: 23},
		{Text: "xylophone", Number: 7},
	} {
		el := rt.NewRecord(k)
		must(t, el.SetNamedField("text", rt.StringOf(item.Text)))
		must(t, el.SetNamedField("number", rt.FloatOf(item.Number)))
		ret = append(ret, el)
	}
	return
}

func must(t *testing.T, e error) {
	if e != nil {
		t.Fatal(e)
	}
}
