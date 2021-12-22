package pdb

import (
	"database/sql"
	"strings"
	"testing"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/asm"
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/rt/kindsOf"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/tables/mdl"
	"git.sr.ht/~ionous/iffy/test/testdb"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// write test data to the database, and ensure we can query it back.
// this exercises the asm.writer ( xforming from strings to ids )
// and the various runtime queries we need.
func TestQueries(t *testing.T) {
	db := testdb.Open(t.Name(), testdb.Memory, "")
	// db := testdb.Open(t.Name(), "", "") // write to file
	defer db.Close()

	const at = ""
	const domain = "main"
	const subDomain = "sub"
	const pattern = "common_ancestor"
	const relation = "whereabouts"
	const otherRelation = "otherRel"
	const kind = "k"
	const subKind = "j"
	const aspect = "a"
	if e := createTable(db,
		func(w eph.Writer) (err error) {
			if e := write(w,
				// name, path, at
				// -------------------------
				mdl.Domain,
				domain, "", at,
				subDomain, domain, at); e != nil {
				err = e
			} else if e := write(w,
				mdl.Kind,
				append(defaultKinds(domain, at),
					// domain, kind, parent, at
					// ---------------------------
					domain, kind, kindsOf.Kind.String(), at,
					domain, aspect, kindsOf.Aspect.String(), at,
					subDomain, subKind, kind, at,
					// patterns:
					domain, pattern, kindsOf.Pattern.String(), at,
					// relations:
					domain, relation, kindsOf.Relation.String(), at,
					subDomain, otherRelation, kindsOf.Relation.String(), at,
				)...,
			); e != nil {
				err = e
			} else if e := write(w,
				mdl.Field,
				// domain, kind, field, affinity, type, at
				// ---------------------------------------
				// traits of an aspect
				domain, aspect, "brief", affine.Bool, nil, at,
				domain, aspect, "verbose", affine.Bool, nil, at,
				domain, aspect, "superbrief", affine.Bool, nil, at,
				// kind that uses that aspect
				domain, kind, aspect, affine.Text, aspect, at,
				// patterns
				domain, pattern, "object", affine.Text, kind, at,
				domain, pattern, "other_object", affine.Text, kind, at,
				domain, pattern, "ancestor", affine.Text, kind, at,
				// relations
				domain, relation, "kind", affine.Text, kind, at,
				domain, relation, "other_kinds", affine.Text, kind, at,
				// ( something random )
				subDomain, otherRelation, "kind", affine.Text, kind, at,
				subDomain, otherRelation, "other_kind", affine.Text, aspect, at,
			); e != nil {
				err = e
			} else if e := write(w,
				mdl.Noun,
				// domain, noun, kind, at
				// ---------------------------------------
				domain, "apple", kind, at,
				domain, "empire_apple", kind, at,
				subDomain, "table", subKind, at,
			); e != nil {
				err = e
				t.Fatal(e)
			} else if e := write(w,
				mdl.Value,
				// "domain", "noun", "field", "value", "affinity", "at"
				// ---------------------------------------
				domain, "apple", aspect, "brief", affine.Text, at,
				domain, "empire_apple", aspect, "verbose", affine.Text, at,
				subDomain, "table", aspect, "superbrief", affine.Text, at,
			); e != nil {
				err = e
				t.Fatal(e)
			} else if e := write(w,
				mdl.Name,
				// domain, noun, name, rank, at
				// ---------------------------------------
				domain, "empire_apple", "empire apple", 0, at,
				domain, "empire_apple", "empire_apple", 1, at,
				domain, "empire_apple", "apple", 2, at,
				domain, "empire_apple", "empire", 3, at,
				domain, "apple", "apple", 0, at, // a different noun with a similar name
			); e != nil {
				err = e
				t.Fatal(e)
			} else if e := write(w,
				mdl.Pat,
				// domain, kind, labels, result
				// ---------------------------------------
				domain, pattern, "object,other_object", "ancestor",
			); e != nil {
				err = e
				t.Fatal(e)
			} else if e := write(w,
				mdl.Rule,
				// "domain", "kind", "target", "phase", "filter", "prog", "at"
				// ---------------------------------------
				domain, pattern, "" /**/, 1, "filter1", "prog1", at,
				domain, pattern, kind, 2, "filter2", "prog2", at,
				domain, pattern, kind, 3, "filter3", "prog3", at,
			); e != nil {
				err = e
				t.Fatal(e)
			} else if e := write(w,
				mdl.Rel,
				// domain, rel, kind, cardinality, subKind, at
				// ---------------------------------------------
				domain, relation, kind, kind, tables.ONE_TO_MANY, at,
				subDomain, otherRelation, kind, aspect, tables.ONE_TO_ONE, at,
			); e != nil {
				err = e
				t.Fatal(e)
			} else if e := write(w,
				mdl.Pair,
				// "domain", "relKind", "oneNoun", "otherNoun", "at"
				// ---------------------------------------------
				subDomain, relation, "table", "empire_apple", at,
				subDomain, relation, "table", "apple", at,
			); e != nil {
				err = e
				t.Fatal(e)
			}
			return // done with writing
		}); e != nil {
		t.Fatal("writing failed", e)
	}

	// start querying
	if q, e := NewQueries(db); e != nil {
		t.Fatal(e)
	} else if domainPoke, e := db.Prepare(
		// turn on / off a domain regardless of hierarchy
		`insert or replace into run_domain(domain, active)
			select md.rowid as domain, ?2
			from mdl_domain md 
			where md.domain = ?1 
			limit 1`,
	); e != nil {
		t.Fatal(e)
	} else if _, e := domainPoke.Exec(domain, true); e != nil {
		t.Fatal(e)
	} else if fd, e := q.FieldsOf(aspect); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(fd, []FieldData{
		{Name: "brief", Affinity: affine.Bool},
		{Name: "verbose", Affinity: affine.Bool},
		{Name: "superbrief", Affinity: affine.Bool},
	}); len(diff) > 0 {
		t.Fatal(fd, diff)
	} else if ks, e := q.KindsOf(kind); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(ks, []string{
		kind,
	}); len(diff) > 0 {
		t.Fatal(ks, diff)
	} else if _, e := domainPoke.Exec(subDomain, true); e != nil {
		t.Fatal(e)
	} else if ks, e := q.KindsOf(kind); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(ks, []string{
		kind, "j",
	}); len(diff) > 0 {
		t.Fatal(ks, diff)
	} else if _, e := domainPoke.Exec(subDomain, false); e != nil {
		t.Fatal(e)
	} else if ok, e := q.IsNounInScope("apple"); e != nil || ok != true {
		t.Fatal(e, ok)
	} else if ok, e := q.IsNounInScope("table"); e != nil || ok != false {
		t.Fatal(e, ok)
	} else if kindOfApple, e := q.KindOfNoun("apple"); e != nil || kindOfApple != kind {
		t.Fatal(kindOfApple, e)
	} else if kindOfTable, e := q.KindOfNoun("table"); e != nil || kindOfTable != "" {
		t.Fatal(kindOfTable, e) // should be blank because the table is out of scope
	} else if aff, val, e := q.NounValue("apple", aspect); e != nil || aff != affine.Text || val.(string) != "brief" {
		t.Fatal(aff, val, e)
	} else if aff, val, e := q.NounValue("table", aspect); e != nil || aff != "" || val != nil {
		t.Fatal(aff, e) // should be out of scope
	} else if name, e := q.NameOfNoun("empire_apple"); e != nil || name != "empire apple" {
		t.Fatal(name, e)
	} else if id, e := q.NounByName("apple"); e != nil || id != "apple" {
		t.Fatal(e, id)
	} else if id, e := q.NounByName("empire"); e != nil || id != "empire_apple" {
		t.Fatal(e, id)
	} else if l, r, e := q.PatternLabels(pattern); e != nil ||
		l != "object,other_object" || r != "ancestor" {
		t.Fatal(e, l, r)
	} else if got, e := q.RulesFor(pattern, ""); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(got, []RuleData{
		{"1", 1, []byte("filter1"), []byte("prog1")},
	}); len(diff) > 0 {
		t.Fatal(got, diff)
	} else if got, e := q.RulesFor(pattern, kind); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(got, []RuleData{
		{"2", 2, []byte("filter2"), []byte("prog2")},
		{"3", 3, []byte("filter3"), []byte("prog3")},
	}); len(diff) > 0 {
		t.Fatal(got, diff)
	} else if got, e := q.Relation(relation); e != nil {
		t.Fatal("relation:", e)
	} else if diff := pretty.Diff(got, RelationData{
		kind, kind, tables.ONE_TO_MANY,
	}); len(diff) > 0 {
		t.Fatal(got, diff)
	} else if e := q.ActivateDomain(subDomain); e != nil {
		t.Fatal("ActivateDomain", e) // enable the sub domain again to get reasonable pairs
	} else if rel, e := q.RelativesOf(relation, "table"); e != nil ||
		len(rel) != 2 || rel[0] != "empire_apple" || rel[1] != "apple" {
		t.Fatal("RelativesOf: table", e, rel)
	} else if rel, e := q.ReciprocalsOf(relation, "apple"); e != nil ||
		len(rel) != 1 || rel[0] != "table" {
		t.Fatal("ReciprocalsOf: apple", e, rel)
	} else if rel := q.Relate(relation, "apple", "empire_apple"); e != nil {
		t.Fatal("Relate", e, rel)
	} else if rel, e := q.RelativesOf(relation, "table"); e != nil ||
		len(rel) != 1 || rel[0] != "apple" {
		t.Fatal("RelativesOf: table", e, rel)
	} else if rel, e := q.RelativesOf(relation, "apple"); e != nil ||
		len(rel) != 1 || rel[0] != "empire_apple" {
		t.Fatal("RelativesOf: apple", e, rel)
	}
}

var run_domain = tables.Insert("run_domain", "domain", "active")

func defaultKinds(domain, at string) (out []interface{}) {
	for _, k := range kindsOf.DefaultKinds {
		pk := k.Parent()
		out = append(out, domain, k.String(), pk.String(), at)
	}
	return
}

func write(w eph.Writer, q string, els ...interface{}) (err error) {
	width, cnt := strings.Count(q, "?"), len(els)
	if div := cnt / width; div*width != cnt {
		err = errutil.New("mismatched width", q)
	} else {
		for i, cnt := 0, len(els); i < cnt; i += width {
			row := els[i : i+width]
			if e := w.Write(q, row...); e != nil {
				onrow := pretty.Sprint("row:", i, row)
				err = errutil.New(q, onrow, e)
				break
			}
		}
	}
	return
}

func createTable(db *sql.DB, cb func(eph.Writer) error) (err error) {
	if e := tables.CreateAll(db); e != nil {
		err = errutil.New("couldnt create model", e)
	} else if tx, e := db.Begin(); e != nil {
		err = errutil.New("couldnt create transaction", e)
	} else if e := cb(asm.NewModelWriter(func(q string, args ...interface{}) (err error) {
		// nothing is confusing about these many layered functions... nothing at all...
		if _, e := tx.Exec(q, args...); e != nil {
			err = e
		}
		return
	})); e != nil {
		err = e
	} else {
		err = tx.Commit()
	}
	return
}
