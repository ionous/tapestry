package mdl

import (
	"fmt"
	"strconv"
	"strings"
)

// Writer - turns write request using the simple model table definitions
// into writes that can lookup and store ids.
// ex. so a caller can Write(mdl.Check, ...) and have the args changed into ids.
// coincidentally, the returned object happens to have the same interface needed for the eph catalog writer.
type Writer func(q string, args ...interface{}) error

func (m Writer) Write(q string, args ...interface{}) (err error) {
	var out string
	if sel, ok := idWriter[q]; ok {
		out = sel
	} else {
		out = q
	}
	return m(out, args...)
}

// create a virtual table consisting of the paths part names turned into comma separated ids:
// NOTE: this winds up flipping the order of the paths: root is towards the end.
func materialize(key string, arg int) string {
	return fmt.Sprintf(`
with recursive
-- str is a list of comma separated parts, each time dropping the left-most part.
parts(str, ids) AS (
select ?%[2]d || ',',  ''
union all
select substr(str, 1+instr(str, ',')), ids || ( 
	-- turn the left most part into a rowid
	select rowid from mdl_%[1]s 
	where %[1]s is substr(str, 0, instr(str, ','))
) || ','
from parts
-- the last str printed is empty, and it contains the full id path.
where length(str) > 1
-- stop any accidental infinite recursion
limit 23) `, key, arg)
}

// select the id where all of the parts have been consumed, or if there were no parts (the root) select the empty string.
const materialized = `(select ids from parts where length(str) == 0 union all select '' limit 1)`

// selects from mdl_<key> where <key>=?<arg>
func idOf(key string, arg int) string {
	return fmt.Sprintf("(select rowid from mdl_%[1]s where %[1]s = ?%[2]d)", key, arg)
}

func insert(name string, args ...string) string {
	var ins strings.Builder
	ins.WriteString("insert into ")
	ins.WriteString(name)
	ins.WriteRune('(')
	for i, cnt := 0, len(args); i < cnt; i += 2 {
		key := args[i]
		if i > 0 {
			ins.WriteRune(',')
		}
		ins.WriteString(key)
	}
	ins.WriteRune(')')
	ins.WriteString(" values (")
	for i, cnt := 1, len(args); i < cnt; i += 2 {
		arg := args[i]
		if i > 1 {
			ins.WriteRune(',')
		}
		if len(arg) > 0 {
			ins.WriteString(arg)
		} else {
			ins.WriteRune('?')
			ins.WriteString(strconv.Itoa((i + 1) / 2))
		}
	}
	ins.WriteRune(')')
	return ins.String()
}

// select the id of a key ( ex. mdl_kind.kind ) indexed by domain name d and key name n
func simpleScope(key string, d, n int) string {
	return fmt.Sprintf(
		`(select key.rowid 
		from mdl_%[1]s key
		join mdl_domain md
		on (key.domain=md.rowid)
		where (md.domain=?%[2]d) and (key.%[1]s=?%[3]d))`, key, d, n)
}

// same as simple scope, but the domain d can be a child of the key's domain
// so to find the correct key, we have to look through all domains
// where its materialized path contains the key's domain
//
//	ex. does path ( ",3,2,1," ) contain ( ",3," )
//	or, at the root ( ",1,," ) contain ( ",1," ).
//
// the prefixed and suffix commas are to avoid partial matches:
//
//	ex. "21," against "1," or ",12" against ",1"
func derivedScope(key string, d, n int) string {
	return fmt.Sprintf(
		`(select key.rowid
		from mdl_%[1]s key
		join mdl_domain md
		where (key.%[1]s = ?%[3]d)
		and instr(',' || md.rowid || ',' || md.path, ',' || key.domain || ','))`, key, d, n)
}

// if i have something defined in domain 2 it should be visible in 3
// we look through all domains and build paths
// ,1,,
// ,1,2,
// ,1,2,3,
// each time asking if it contains ,2,

const unchanged = ""

func arg(i int) string { return "?" + strconv.Itoa(i) }

// rewrite some tables to use ids
// the key of the table is the original, simplified insert statement
// the value is a more complex statement usually involving selects
var idWriter = map[string]string{
	// turn domain name into an id
	Check: insert("mdl_check",
		"domain", idOf("domain", 1),
		"name", unchanged,
		"value", unchanged,
		"affinity", unchanged,
		"prog", unchanged,
		"at", unchanged,
	),
	// turn the materialized path of domain ancestor names into ancestor idOfs
	Domain: materialize("domain", 2) +
		insert("mdl_domain",
			"domain", unchanged,
			"path", materialized,
			"at", unchanged,
		),
	// domain name + kind name select a specific kind entry.
	// domain is not written b/c fields exist per kind, not per domain.
	Field: insert("mdl_field",
		"kind", simpleScope("kind", 1, 2),
		"field", arg(3),
		"affinity", arg(4),
		"type", derivedScope("kind", 1, 5),
		"at", arg(6),
	),
	// turn domain name into an id
	Grammar: insert("mdl_grammar",
		"domain", idOf("domain", 1),
		"name", unchanged,
		"prog", unchanged,
		"at", unchanged,
	),
	// turn domain name into an id, and materialize the ancestor path
	Kind: materialize("kind", 3) +
		insert("mdl_kind",
			"domain", idOf("domain", 1),
			"kind", unchanged,
			"path", materialized,
			"at", unchanged,
		),
	// turn domain, kind, field into ids, associated with the local var's initial assignment.
	// domain and kind become redundant b/c fields exist at the scope of the kind.
	Assign: string(`with parts(did, domain, kid, kind, fid, field) as (
		select md.rowid, md.domain, mk.rowid, mk.kind, mf.rowid, mf.field
		from mdl_field mf
		join mdl_kind mk
			on (mk.rowid = mf.kind)
		join mdl_domain md
			on (mk.domain = md.rowid))
		insert into mdl_assign(field, value)
		select fid, ?4
		from parts where domain=?1 and kind=?2 and field=?3`,
	),
	Name: insert("mdl_name",
		"domain", idOf("domain", 1), // currently redundant, names have the same scope as their noun.
		"noun", simpleScope("noun", 1, 2),
		"name", unchanged,
		"rank", unchanged,
		"at", unchanged,
	),
	Noun: insert("mdl_noun",
		"domain", idOf("domain", 1), // domain where the noun was declared
		"noun", unchanged,
		"kind", derivedScope("kind", 1, 3),
		"at", unchanged,
	),
	Pair: insert("mdl_pair",
		"domain", idOf("domain", 1), // domain where the pair was declared
		"relKind", derivedScope("kind", 1, 2), // we point to the kind table not the relation table.
		"oneNoun", derivedScope("noun", 1, 3),
		"otherNoun", derivedScope("noun", 1, 4),
		"at", unchanged,
	),
	// the labels are fields of kind
	// the domain is dropped: its the same as the kind's scope.
	Pat: insert("mdl_pat",
		"kind", simpleScope("kind", 1, 2),
		"labels", arg(3), // fix? this are comma-separated field names, should it be field ids?
		"result", arg(4), // fix? this is a field, should it be a field id?
	),
	Plural: insert("mdl_plural",
		"domain", idOf("domain", 1),
		"many", unchanged,
		"one", unchanged,
		"at", unchanged,
	),
	Rel: insert("mdl_rel",
		"relKind", simpleScope("kind", 1, 2),
		"oneKind", derivedScope("kind", 1, 3),
		"otherKind", derivedScope("kind", 1, 4),
		"cardinality", arg(5),
		"at", arg(6),
	),
	Rule: insert("mdl_rule",
		"domain", idOf("domain", 1), // domain where the rule was declared
		"kind", derivedScope("kind", 1, 2),
		"target", derivedScope("kind", 1, 3),
		"phase", unchanged,
		"filter", unchanged,
		"prog", unchanged,
		"at", unchanged,
	),
	// first: build a virtual [domain, noun, field] table
	// note: values are written per noun, not per domain; so the domain column is redundant once the noun id is known.
	// to get to the field id, we have to look at all possible fields for the noun.
	// given the kind of the noun, accept all fields who's kind is in its materialized path.
	// fix? some values are references to objects in the form "#domain::noun" -- should the be changed to ids?
	Value: string(`with parts(did, domain, nin, noun, fid, field) as (
			select md.rowid, md.domain, mn.rowid, mn.noun, mf.rowid, mf.field
			from mdl_noun mn
			join mdl_domain md
				on (mn.domain = md.rowid)
			join mdl_kind mk
				on (mn.kind = mk.rowid)
			left join mdl_field mf
				where instr(',' || mk.rowid || ',' || mk.path, ',' || mf.kind || ','))
			insert into mdl_value(noun, field, value, at)
			select nin, fid, ?4, ?5
			from parts where domain=?1 and noun=?2 and field=?3`,
	),
}
