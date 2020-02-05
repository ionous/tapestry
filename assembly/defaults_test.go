package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/kr/pretty"
)

// TestDefaultFieldAssigment to verify default values can be assigned to kinds.
func TestDefaultFieldAssigment(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		//
		if e := fakeHierarchy(t.modeler, []pair{
			{"K", ""},
			{"L", "K"},
			{"D", "K"},
			{"C", "L,K"},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeFields(t.modeler, []kfp{
			{"K", "d", ephemera.PRIM_DIGI},
			{"K", "t", ephemera.PRIM_TEXT},
			{"K", "t2", ephemera.PRIM_TEXT},
			{"L", "x", ephemera.PRIM_TEXT},
			{"D", "x", ephemera.PRIM_TEXT},
			{"C", "c", ephemera.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		} else if e := addDefaults(t.rec, []triplet{
			{"K", "t", "some text"},
			{"L", "t", "override text"},
			{"L", "t2", "other text"},
			{"L", "x", "x in p"},
			{"D", "x", "x in d"},
			{"C", "c", "c text"},
			{"C", "d", 123},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineDefaults(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchDefaults(t.db, []triplet{
			{"C", "c", "c text"},
			{"C", "d", int64(123)}, // re: int64 -- default scanner uses https://golang.org/pkg/database/sql/#Scanner
			{"D", "x", "x in d"},
			{"K", "t", "some text"},
			{"L", "t", "override text"},
			{"L", "t2", "other text"},
			{"L", "x", "x in p"},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestDefaultTraitAssignment to verify default traits can be assigned to kinds.
func TestDefaultTraitAssignment(t *testing.T) {
	if t, e := newDefaultsTest(t, memory, []triplet{
		{"K", "x", true},
		{"L", "y", true},
		{"L", "z", true},
		//
		{"N", "A", "w"},
		{"N", "B", "z"},
		{"N", "w", true},
	}); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		if e := DetermineDefaults(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchDefaults(t.db, []triplet{
			{"K", "A", "x"},
			{"L", "A", "y"},
			{"L", "B", "z"},
			{"N", "A", "w"},
			{"N", "B", "z"},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestDefaultDuplicates to verify that duplicate default specifications are okay
func TestDefaultDuplicates(t *testing.T) {
	if t, e := newDefaultsTest(t, memory, []triplet{
		{"K", "t", "text"},
		{"K", "t", "text"},
		{"L", "t", "text"},
		//
		{"K", "d", 123},
		{"K", "d", 123},
		{"L", "d", 123},
		//
		{"K", "A", "y"},
		{"K", "y", true},
		{"L", "x", true},
		{"L", "A", "x"},
	}); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		if e := DetermineDefaults(t.modeler, t.db); e != nil {
			t.Fatal(e)
		}
	}
}

// TestDefaultConflict to verify that conflicting values for the same default are not okay
func TestDefaultConflict(t *testing.T) {
	testConflict := func(t *testing.T, vals []triplet) (err error) {
		if t, e := newDefaultsTest(t, memory, vals); e != nil {
			t.Fatal(e)
		} else {
			defer t.Close()
			if e := DetermineDefaults(t.modeler, t.db); e == nil {
				err = errutil.New("expected error")
			} else {
				t.Log("okay:", e)
			}
		}
		return
	}
	if e := testConflict(t, []triplet{
		{"K", "t", "a"},
		{"K", "t", "b"},
	}); e != nil {
		t.Fatal(e)
	} else if e := testConflict(t, []triplet{
		{"K", "d", 1},
		{"K", "d", 2},
	}); e != nil {
		t.Fatal(e)
	}

	if e := testConflict(t, []triplet{
		{"K", "A", "x"},
		{"K", "A", "y"},
	}); e != nil {
		t.Fatal(e)
	} else if e := testConflict(t, []triplet{
		{"K", "x", true},
		{"K", "y", true},
	}); e != nil {
		t.Fatal(e)
	} else if e := testConflict(t, []triplet{
		{"K", "A", "x"},
		{"K", "y", true},
	}); e != nil {
		t.Fatal(e)
	}
}

// TestDefaultBadValue to verify that modeling requires appropriate values for defaults based on type
func TestDefaultBadValue(t *testing.T) {
	//- for now, we only allow text and number [ text and digi ]
	// - later we could add ambiguity for conversion [ 4 -> "4" ]
	testInvalid := func(t *testing.T, vals []triplet) (err error) {
		if t, e := newDefaultsTest(t, memory, vals); e != nil {
			err = e
		} else {
			defer t.Close()
			if e := DetermineDefaults(t.modeler, t.db); e == nil {
				err = errutil.New("expected error")
			} else {
				t.Log("okay:", e)
			}
		}
		return
	}

	if e := testInvalid(t, []triplet{
		{"K", "t", 1.2},
	}); e != nil {
		t.Fatal(e)
	} else if e := testInvalid(t, []triplet{
		{"K", "d", "1.2"},
	}); e != nil {
		t.Fatal(e)
	}
	// try to set trait like values

	if e := testInvalid(t, []triplet{
		{"K", "t", false},
	}); e != nil {
		t.Fatal(e)
	}

	/*
	   fix? somehow? bools in sqlite are stored as int64;
	   could switch to text ( "true", "false" ) perhaps and add some check/query
	   during determination
	   if e := testInvalid(t, []triplet{
	       {"K", "d", true},
	   }); e != nil {
	       t.Fatal(e)
	   }
	*/

	/* fix? aspects are set by matching traits
	1.2 is not a trait, so it's skipped.
	this might get handled by a "missing" check,
	or possibly by changing the determination query.

	if e := testInvalid(t, []triplet{
		{"K", "A", 1.2},
	}); e != nil {
		t.Fatal(e)
	}
	*/
}

// match generated model defaults
func matchDefaults(db *sql.DB, want []triplet) (err error) {
	var curr triplet
	var have []triplet
	if e := dbutil.QueryAll(db,
		`select kind, field, value 
			from mdl_default
			order by kind, field, value`,
		func() (err error) {
			have = append(have, curr)
			return
		},
		&curr.target, &curr.prop, &curr.value); e != nil {
		err = e
	} else if !reflect.DeepEqual(have, want) {
		err = errutil.New("mismatch",
			"have:", pretty.Sprint(have),
			"want:", pretty.Sprint(want))
	}
	return
}

// write ephemera describing some initial values
func addDefaults(rec *ephemera.Recorder, defaults []triplet) (err error) {
	for _, el := range defaults {
		namedKind := rec.Named(ephemera.NAMED_KIND, el.target, "test")
		namedField := rec.Named(ephemera.NAMED_PROPERTY, el.prop, "test")
		rec.NewDefault(namedKind, namedField, el.value)
	}
	return
}

func newDefaultsTest(t *testing.T, path string, defaults []triplet) (ret *assemblyTest, err error) {
	if t, e := newAssemblyTest(t, path); e != nil {
		err = e
	} else {
		if e := fakeHierarchy(t.modeler, []pair{
			{"K", ""},
			{"L", "K"},
			{"N", "K"},
		}); e != nil {
			err = e
		} else if e := fakeFields(t.modeler, []kfp{
			{"K", "d", ephemera.PRIM_DIGI},
			{"K", "t", ephemera.PRIM_TEXT},
			{"K", "A", ephemera.PRIM_ASPECT},
			{"L", "B", ephemera.PRIM_ASPECT},
			{"N", "B", ephemera.PRIM_ASPECT},
		}); e != nil {
			err = e
		} else if e := fakeTraits(t.modeler, []pair{
			{"A", "w"}, {"A", "x"}, {"A", "y"},
			{"B", "z"},
		}); e != nil {
			err = e
		} else if e := addDefaults(t.rec, defaults); e != nil {
			err = e
		}
		if err != nil {
			t.Close()
		} else {
			ret = t
		}
	}
	return
}
