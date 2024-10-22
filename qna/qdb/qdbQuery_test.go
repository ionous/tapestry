package qdb_test

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/kr/pretty"
)

// write test data to the database, and ensure we can query it back.
// this exercises the asm.writer.Add( xforming from strings to ids )
// and the various runtime queries we need.
func TestQueries(t *testing.T) {
	db := tables.CreateTest(t.Name(), true)
	defer db.Close()

	const domain = "main"
	const subDomain = "sub"
	const pattern = "common ancestor"
	const relation = "whereabouts"
	const otherRelation = "other rel"
	const plural = "clutches"
	const singular = "purse"
	const kind = "k"
	const subKind = "j"
	const aspect = "a"
	if m, e := mdl.NewModeler(db); e != nil {
		t.Fatal(e)
	} else if e := mdlDomain(m,
		// name, path
		// -------------------------
		domain, "",
		subDomain, domain); e != nil {
		t.Fatal(e)
	} else if e := mdlPlural(m,
		// name, path
		// -------------------------
		domain, plural, singular); e != nil {
		t.Fatal(e)
	} else if e := mdlKind(m,
		append(defaultKinds(domain),
			// "domain", "kind", "path"
			// ---------------------------
			domain, kind, kindsOf.Kind.String(),
			domain, aspect, kindsOf.Aspect.String(),
			// super confusing, in the db: root is towards end of the path; the row id is at the start.
			// when read: root is hit first ( its in earlier *rows* ) so it becomes first.
			subDomain, subKind, kind,
			// patterns:
			domain, pattern, kindsOf.Pattern.String(),
			// relations:
			domain, relation, kindsOf.Relation.String(),
			subDomain, otherRelation, kindsOf.Relation.String(),
		)...,
	); e != nil {
		t.Fatal(e)
	} else if e := mdlField(m,
		// domain, kind, field, affinity, type
		// ---------------------------------------
		// traits of an aspect
		addMember, domain, aspect, "brief", affine.Bool, "",
		addMember, domain, aspect, "verbose", affine.Bool, "",
		addMember, domain, aspect, "superbrief", affine.Bool, "",
		// kind that uses that aspect
		addMember, domain, kind, aspect, affine.Text, aspect,
		// patterns
		addParameter, domain, pattern, "object", affine.Text, kind,
		addParameter, domain, pattern, "other object", affine.Text, kind,
		addResult, domain, pattern, "ancestor", affine.Text, kind,
	); e != nil {
		t.Fatal(e)
	} else if e := mdlNoun(m,
		// domain, noun, kind
		// ---------------------------------------
		domain, "apple", kind,
		domain, "empire apple", kind,
		subDomain, "table", subKind,
	); e != nil {
		t.Fatal(e)
		t.Fatal(e)
	} else if e := mdlValue(m,
		// "domain", "noun", "field", "value", "at"
		// ---------------------------------------
		domain, "apple", aspect, "brief",
		domain, "empire apple", aspect, "verbose",
		subDomain, "table", aspect, "superbrief",
	); e != nil {
		t.Fatal(e)
		t.Fatal(e)
	} else if e := mdlName(m,
		// domain, noun, name, rank
		// ---------------------------------------
		domain, "empire apple", "empire apple", 0,
		domain, "empire apple", "apple", 1,
		domain, "empire apple", "empire", 2,
		domain, "apple", "apple", 0, // a different noun with a similar name
	); e != nil {
		t.Fatal(e)
		t.Fatal(e)
	} else if e := mdlRule(m,
		// "domain", "kind", "target", "rank", "prog", "at"
		// ---------------------------------------
		domain, pattern, 1, "prog1",
		domain, pattern, 2, "prog2",
		domain, pattern, 3, "prog3",
	); e != nil {
		t.Fatal(e)
		t.Fatal(e)
	} else if e := mdlRel(m,
		// domain, rel, kind, otherKind, cardinality
		// ---------------------------------------------
		domain, relation, kind, kind, false, true,
		// ( something random )
		subDomain, otherRelation, kind, aspect, false, false,
	); e != nil {
		t.Fatal(e)
		t.Fatal(e)
	} else if e := mdlPair(m,
		// "domain", "relKind", "oneNoun", "otherNoun", "at"
		// ---------------------------------------------
		subDomain, relation, "table", "empire apple",
		subDomain, relation, "table", "apple",
	); e != nil {
		t.Fatal(e)
		t.Fatal(e)
	}

	// start querying
	if q, e := qdb.NewQueries(db, qdb.DecodeNone("testing")); e != nil {
		t.Fatal(e)
	} else if domainPoke, e := db.Prepare(
		// turn on / off a domain regardless of hierarchy
		`insert or replace into run_domain(domain, active)
			values( ?1, ?2 )`,
	); e != nil {
		t.Fatal(e)
	} else if _, e := domainPoke.Exec(domain, true); e != nil {
		t.Fatal(e)
	} else if one, e := q.PluralToSingular(plural); e != nil || one != singular {
		t.Fatal("singular", one, e)
	} else if many, e := q.PluralFromSingular(singular); e != nil || many != plural {
		t.Fatal("plural", many, e)
	} else if one, e := q.PluralToSingular("x" + plural); e != nil || one != "" {
		t.Fatal("singular", one, e)
	} else if many, e := q.PluralFromSingular("x" + singular); e != nil || many != "" {
		t.Fatal("plural", many, e)
	} else if fd, e := q.FieldsOf(kind); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(fd, []rt.Field{
		{Name: aspect, Affinity: affine.Text, Type: aspect},
	}); len(diff) > 0 {
		t.Fatal(fd, diff)
	} else if fd, e := q.FieldsOf(aspect); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(fd, []rt.Field{
		{Name: "brief", Affinity: affine.Bool},
		{Name: "verbose", Affinity: affine.Bool},
		{Name: "superbrief", Affinity: affine.Bool},
	}); len(diff) > 0 {
		t.Fatal(fd, diff)
	} else if path, e := q.KindOfAncestors("j"); e != nil || len(path) != 0 {
		t.Fatal("KindOfAncestors", path, e)
	} else if _, e := domainPoke.Exec(subDomain, true); e != nil {
		t.Fatal(e)
	} else if path, e := q.KindOfAncestors("j"); e != nil || strings.Join(path, ",") != "j,k,kinds" {
		got := strings.Join(path, ",")
		t.Fatal("KindOfAncestors", got, e)
	} else if _, e := domainPoke.Exec(subDomain, false); e != nil {
		t.Fatal(e)
	} else /*if ok, e := q.NounActive("apple"); e != nil || ok != true {
		t.Fatal(e, ok)
	} else if ok, e := q.NounActive("table"); e != nil || ok != false {
		t.Fatal(e, ok)
	} else if kindOfApple, e := q.NounKind("apple"); e != nil || kindOfApple != kind {
		t.Fatal(kindOfApple, e)
	} else if kindOfTable, e := q.NounKind("table"); e != nil || kindOfTable != "" {
		t.Fatal(kindOfTable, e) // should be blank because the table is out of scope
	} else if pairs, e := q.NounValue("apple", aspect); e != nil ||
		!reflect.DeepEqual(pairs, []query.ValueData{{Field: aspect, Value: []byte("brief")}}) {
		t.Fatal("the aspect's current value should be the value 'brief'", pairs, e)
		// } else if pairs, e := q.NounValue("table", aspect); e != nil || pairs != nil {
		// 	t.Fatal(pairs, e) // should be out of scope
		// } else if name, e := q.NounName("empire apple"); e != nil || name != "empire apple" {
		// 	t.Fatal(name, e)
	} else */if id, e := q.NounInfo("apple"); e != nil || id != (query.NounInfo{Domain: domain, Noun: "apple", Kind: kind}) {
		t.Fatal(e, id)
	} else if id, e := q.NounInfo("empire"); e != nil || id != (query.NounInfo{Domain: domain, Noun: "empire apple", Kind: kind}) {
		t.Fatal(e, id)
	} else if id, e := q.NounInfo("table"); e != nil || id != (query.NounInfo{}) {
		t.Fatal(id, e) // should be blank because the table is out of scope
	} else if got, e := q.PatternLabels(pattern); e != nil {
		t.Fatal("patternLabels:", e)
	} else if diff := pretty.Diff(got, []string{"object", "other object", "ancestor"}); len(diff) > 0 {
		t.Fatal(e, diff)
	} else /*else if got, e := q.RulesFor(pattern); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(got, []query.RuleData{
		{Name: "rule 1", Prog: []byte("prog1")},
		{Name: "rule 2", Prog: []byte("prog2")},
		{Name: "rule 3", Prog: []byte("prog3")},
	}); len(diff) > 0 {
		t.Fatal(got, diff)
	} else if got, e := q.Relation(relation); e != nil {
		t.Fatal("relation:", e)
	} else if diff := pretty.Diff(got, RelationData{
		kind, kind, tables.ONE_TO_MANY,
	}); len(diff) > 0 {
		t.Fatal(got, diff)
	} else*/if _, _, e := q.ActivateDomains(subDomain); e != nil {
		t.Fatal("ActivateDomain", e) // enable the sub domain again to get reasonable pairs
		// note: we never previously fully activated a domain, so prev is empty.
	} else if rel, e := q.RelativesOf(relation, "table"); e != nil ||
		len(rel) != 2 || rel[1] != "empire apple" || rel[0] != "apple" {
		t.Fatal("RelativesOf: table", e, rel)
	} else if rel, e := q.ReciprocalsOf(relation, "apple"); e != nil ||
		len(rel) != 1 || rel[0] != "table" {
		t.Fatal("ReciprocalsOf: apple", e, rel)
	} else if e := q.Relate(relation, "apple", "empire apple"); e != nil {
		t.Fatal("Relate", e, rel)
	} else if rel, e := q.RelativesOf(relation, "table"); e != nil ||
		len(rel) != 1 || rel[0] != "apple" {
		t.Fatal("RelativesOf: table", e, rel)
	} else if rel, e := q.RelativesOf(relation, "apple"); e != nil ||
		len(rel) != 1 || rel[0] != "empire apple" {
		t.Fatal("RelativesOf: apple", e, rel)
	}
}

func defaultKinds(domain string) (out []any) {
	for _, k := range kindsOf.DefaultKinds {
		pk := k.Parent()
		out = append(out, domain, k.String(), pk.String())
	}
	return
}

// adapt old style tests to new interface
func mdlDomain(m *mdl.Modeler, els ...any) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 2 {
		row := els[i:]
		domain, requires :=
			row[0].(string),
			row[1].(string)
		if e := m.Pin(domain).AddDependency(requires); e != nil {
			err = e
			break
		}
	}
	return
}
func mdlField(m *mdl.Modeler, els ...any) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 6 {
		row := els[i:]
		fn, domain, kind, field, affinity, typeName :=
			row[0].(func(pen *mdl.Pen, kind, field string, affinity affine.Affinity, typeName string) error),
			row[1].(string),
			row[2].(string),
			row[3].(string),
			row[4].(affine.Affinity),
			row[5].(string)
		pen := m.Pin(domain)
		if e := fn(pen, kind, field, affinity, typeName); e != nil {
			err = e
			break
		}
	}
	return
}
func addMember(pen *mdl.Pen, kind, field string, aff affine.Affinity, cls string) error {
	return pen.AddTestField(kind, field, aff, cls)
}
func addParameter(pen *mdl.Pen, kind, field string, aff affine.Affinity, cls string) error {
	return pen.AddTestParameter(kind, field, aff, cls)
}
func addResult(pen *mdl.Pen, kind, field string, aff affine.Affinity, cls string) error {
	return pen.AddTestResult(kind, field, aff, cls)
}
func mdlKind(m *mdl.Modeler, els ...any) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 3 {
		row := els[i:]
		domain, kind, path :=
			row[0].(string),
			row[1].(string),
			row[2].(string)
		if e := m.Pin(domain).AddKind(kind, path); e != nil {
			err = e
			break
		}
	}
	return
}

func mdlName(m *mdl.Modeler, els ...any) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 4 {
		row := els[i:]
		domain, noun, name, rank :=
			row[0].(string),
			row[1].(string),
			row[2].(string),
			row[3].(int)
		if e := m.Pin(domain).AddNounName(noun, name, rank); e != nil {
			err = e
			break
		}
	}
	return
}

func mdlNoun(m *mdl.Modeler, els ...any) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 3 {
		row := els[i:]
		domain, noun, kind :=
			row[0].(string),
			row[1].(string),
			row[2].(string)
		if e := m.Pin(domain).AddNounKind(noun, kind); e != nil {
			err = e
			break
		}
	}
	return
}

func mdlPair(m *mdl.Modeler, els ...any) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 4 {
		row := els[i:]
		domain, relKind, oneNoun, otherNoun :=
			row[0].(string),
			row[1].(string),
			row[2].(string),
			row[3].(string)
		if e := m.Pin(domain).AddNounPair(relKind, oneNoun, otherNoun); e != nil {
			err = e
			break
		}
	}
	return
}
func mdlPlural(m *mdl.Modeler, els ...any) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 3 {
		row := els[i:]
		domain, many, one :=
			row[0].(string),
			row[1].(string),
			row[2].(string)
		if e := m.Pin(domain).AddPlural(many, one); e != nil {
			err = e
			break
		}
	}
	return
}

func mdlRel(m *mdl.Modeler, els ...any) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 6 {
		row := els[i:]
		domain, relKind, oneKind, otherKind, oneMany, otherMany :=
			row[0].(string),
			row[1].(string),
			row[2].(string),
			row[3].(string),
			row[4].(bool),
			row[5].(bool)
		if e := m.Pin(domain).AddRelation(relKind, oneKind, otherKind, oneMany, otherMany); e != nil {
			err = e
			break
		}
	}
	return
}
func mdlRule(m *mdl.Modeler, els ...any) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 4 {
		row := els[i:]
		domain := row[0].(string)
		pattern := row[1].(string)
		rank := row[2].(int)
		prog := row[3].(string)
		if e := m.Pin(domain).AddTestRule(pattern, rank, prog); e != nil {
			err = e
			break
		}
	}
	return
}
func mdlValue(m *mdl.Modeler, els ...any) (err error) {
	for i, cnt := 0, len(els); i < cnt; i += 4 {
		row := els[i:]
		domain, noun, field, value :=
			row[0].(string),
			row[1].(string),
			row[2].(string),
			row[3].(string)
		if e := m.Pin(domain).AddTestValue(noun, false, field, value); e != nil {
			err = e
			break
		}
	}
	return
}
