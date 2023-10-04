package list_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/scope"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/ionous/sliceOf"
	"github.com/kr/pretty"
)

func TestSort(t *testing.T) {
	var kinds testutil.Kinds
	var objs testutil.Objects

	kinds.AddKinds((*Things)(nil), (*Locals)(nil))
	objectNames := sliceOf.String("mildred", "apple", "pen", "eve", "Pan")
	objs.AddObjects(kinds.Kind("things"), objectNames...)

	locals := kinds.NewRecord("locals", "objects", objectNames)
	lt := testutil.Runtime{
		Kinds:     &kinds,
		ObjectMap: objs,
	}
	lt.Chain = scope.MakeChain(scope.FromRecord(&lt, locals))

	// create a new value of type "locals" containing "Objects:objectNames"
	for key, obj := range objs {
		if e := obj.SetNamedField("key", g.StringOf(key)); e != nil {
			t.Fatal(e)
		}
	}
	// sorts in place
	sorter := &list.ListSortText{
		Target:  core.Variable("objects"),
		ByField: "key",
	}
	if e := safe.Run(&lt, sorter); e != nil {
		t.Fatal(e)
	}
	if diff := pretty.Diff(objectNames, []string{"apple", "eve", "mildred", "Pan", "pen"}); len(diff) > 0 {
		t.Fatal(objectNames)
	}
	//
	sorter.UsingCase = B(true)
	if e := safe.Run(&lt, sorter); e != nil {
		t.Fatal(e)
	}
	if diff := pretty.Diff(objectNames, []string{"Pan", "apple", "eve", "mildred", "pen"}); len(diff) > 0 {
		t.Fatal(objectNames)
	}
	//
	sorter.Descending = B(true)
	if e := safe.Run(&lt, sorter); e != nil {
		t.Fatal(e)
	}
	if diff := pretty.Diff(objectNames, []string{"pen", "mildred", "eve", "apple", "Pan"}); len(diff) > 0 {
		t.Fatal(objectNames)
	}
}

// a simple noun
type Things struct{ Key string }

// global variables for grouping tests
type Locals struct{ Objects []string }
