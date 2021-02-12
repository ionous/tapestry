package list_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/list"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"git.sr.ht/~ionous/iffy/test/testutil"
	"github.com/ionous/sliceOf"
	"github.com/kr/pretty"
)

func TestSort(t *testing.T) {
	var kinds testutil.Kinds
	kinds.AddKinds((*Things)(nil), (*Values)(nil))
	objectNames := sliceOf.String("mildred", "apple", "pen", "eve", "Pan")
	if objs, e := testutil.Objects(kinds.Kind("Things"), objectNames...); e != nil {
		t.Fatal(e)
	} else {
		values := kinds.New("Values", "Objects", objectNames)
		lt := testutil.Runtime{
			Kinds:     &kinds,
			ObjectMap: objs,
			Stack: []rt.Scope{
				g.RecordOf(values),
			},
		}
		// create a new value of type "Values" containing "Objects:objectNames"
		for key, obj := range objs {
			obj.SetNamedField("Key", g.StringOf(key))
		}
		// sorts in place
		sorter := &list.SortText{
			Var:     core.Variable{Str: "Objects"},
			ByField: &list.SortByField{Name: "Key"},
		}
		if e := safe.Run(&lt, sorter); e != nil {
			t.Fatal(e)
		}
		if diff := pretty.Diff(objectNames, []string{"apple", "eve", "mildred", "Pan", "pen"}); len(diff) > 0 {
			t.Fatal(objectNames)
		}
		//
		sorter.Case = true
		if e := safe.Run(&lt, sorter); e != nil {
			t.Fatal(e)
		}
		if diff := pretty.Diff(objectNames, []string{"Pan", "apple", "eve", "mildred", "pen"}); len(diff) > 0 {
			t.Fatal(objectNames)
		}
		//
		sorter.Order = true
		if e := safe.Run(&lt, sorter); e != nil {
			t.Fatal(e)
		}
		if diff := pretty.Diff(objectNames, []string{"pen", "mildred", "eve", "apple", "Pan"}); len(diff) > 0 {
			t.Fatal(objectNames)
		}
	}
}

// a simple noun
type Things struct{ Key string }

// global variables for grouping tests
type Values struct{ Objects []string }
