package mdl

import (
	"database/sql"
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/test/testdb"
	"github.com/ionous/errutil"
)

// todo: test that traits write to their aspect's field
func TestValueWriting(t *testing.T) {
	const (
		oneField   = "oneField"
		otherField = "otherField"
		recName    = "record"
		kindName   = "kind"
		nounA      = "a"
		nounB      = "b"
		oneValue   = "one value"
		otherValue = "other value"
	)
	//
	var warnings Warnings
	unwarn := warnings.Catch(t.Fatal)
	defer unwarn()
	db := testdb.Create(t.Name())
	if e := tables.CreateModel(db); e != nil {
		t.Fatal(e)
	} else {
		// two kinds, for a noun and a record:
		var kind struct {
			id, one, other, rec int64
		}
		var rec struct {
			id, one, other, rec int64
		}
		var nouns struct {
			a, b int64
		}

		// write a single helper domain
		var mdl_domain = insert(t, "mdl_domain")
		mdl_domain.insert(db)

		// and some kinds:
		var mdl_kind = insert(t, "mdl_kind", "kind", "path")
		kind.id = mdl_kind.insert(db, kindName, "")
		rec.id = mdl_kind.insert(db, recName, "")

		// fields for the kind:
		var mdl_field = insert(t, "mdl_field", "kind", "field", "affinity", "type")
		kind.one = mdl_field.insert(db, kind.id, oneField, affine.Text, nil)
		kind.other = mdl_field.insert(db, kind.id, otherField, affine.Text, nil)
		kind.rec = mdl_field.insert(db, kind.id, recName, affine.Record, rec.id)
		// fields for the record
		rec.one = mdl_field.insert(db, rec.id, oneField, affine.Text, nil)
		rec.other = mdl_field.insert(db, rec.id, otherField, affine.Text, nil)
		rec.rec = mdl_field.insert(db, rec.id, recName, affine.Record, rec.id) // recursive!

		// noun of kind
		var mdl_noun = insert(t, "mdl_noun", "noun", "kind")
		nouns.a = mdl_noun.insert(db, nounA, kind.id)
		nouns.b = mdl_noun.insert(db, nounB, kind.id)

		//
		if m, e := NewModelerWithWarnings(db, func(fmt string, parts ...any) {
			LogWarning(errutil.Fmt(fmt, parts...))
		}); e != nil {
			t.Fatal(e)
		} else {
			pen := m.Pin("domain", "at")
			// some independent fields:
			if e := pen.AddTestValue(nounA, oneField, oneValue); e != nil {
				t.Fatal(e)
			} else if e := pen.AddTestValue(nounA, otherField, otherValue); e != nil {
				t.Fatal(e)
			}
			// the deepening
			if e := pen.AddTestValue(nounA, MakePath(recName, recName, oneField), oneValue); e != nil {
				t.Fatal(e)
			}
			// writing again should be okay:
			if e := pen.AddTestValue(nounA, MakePath(recName, recName, oneField), oneValue); e != nil {
				t.Fatal(e)
			} else if e := warnings.Expect("Duplicate value for 'a.record.record.oneField'."); e != nil {
				t.Fatal(e)
			}
			// make sure we can't now write at the record itself
			if e := pen.AddTestValue(nounA, MakePath(recName, ""), oneValue); e == nil {
				t.Fatal("shouldn't have written a whole record after writing one of its fields")
			} else {
				t.Log("ok", e)
			}
			// make sure we can however write deeper
			if e := pen.AddTestValue(nounA, MakePath(recName, recName, recName, oneField), oneValue); e != nil {
				t.Fatal(e)
			}
			if e := pen.AddTestValue(nounA, MakePath(recName, recName, recName, otherField), otherValue); e != nil {
				t.Fatal(e)
			}
			// make sure to detect the conflict deeper
			if e := pen.AddTestValue(nounA, MakePath(recName, recName), otherValue); e == nil {
				t.Fatal("shouldn't have written a whole record after writing one of its fields")
			}
			//
			// reverse the order of writing ( via the second noun ) and make sure that fails too.
			//
			if e := pen.AddTestValue(nounB, MakePath(recName, ""), oneValue); e != nil {
				t.Fatal(e)
			}
			if e := pen.AddTestValue(nounB, MakePath(recName, recName, oneField), oneValue); e == nil {
				t.Fatal("shouldn't have written a field after writing the whole record")
			} else {
				t.Log("ok", e)
			}
		}
	}
}

type inserter struct {
	t   *testing.T
	ins string
	id  int64
}

func insert(t *testing.T, table string, keys ...string) inserter {
	ins := tables.Insert(table, append([]string{"domain", "rowid"}, keys...)...)
	return inserter{t, ins, 0}
}

func (ins *inserter) insert(db *sql.DB, values ...any) int64 {
	id := 1 + ins.id
	ins.id = id
	if _, e := db.Exec(ins.ins, append([]any{"domain", id}, values...)...); e != nil {
		ins.t.Fatal(e)
	}
	return id
}
