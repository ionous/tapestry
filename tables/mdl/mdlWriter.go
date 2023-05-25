package mdl

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

func NewModeler(db *sql.DB) (ret Modeler, err error) {
	var ps tables.Prep
	m := &Writer{
		db: db,
		assign: ps.Prep(db,
			tables.Insert("mdl_assign", "field", "value"),
		),
		check:  ps.Prep(db, Check),
		domain: ps.Prep(db, Domain),
		field: ps.Prep(db,
			tables.Insert("mdl_field", "kind", "field", "affinity", "type", "at"),
		),
		grammar: ps.Prep(db, Grammar),

		// create a virtual table consisting of the paths part names turned into comma separated ids:
		// NOTE: this winds up flipping the order of the paths: root is towards the end.
		kind: ps.Prep(db,
			`with recursive
			-- str is a list of comma separated parts, 
			-- each time dropping the left-most part.
			parts(str, ids) as (
			select ?3 || ',',  ''
			union all
			select substr(str, 1+instr(str, ',')), ids || ( 
				-- turn the left most part into a rowid
				select rowid from mdl_kind 
				where kind is substr(str, 0, instr(str, ','))
			) || ','
			from parts
			-- the last str printed is empty, and it contains the full id path.
			where length(str) > 1
			-- stop any accidental infinite recursion
			limit 23)
			insert into mdl_kind( domain, kind, path, at ) 
			values ( 
				?1, 
				?2, 
				-- select the id where all of the parts have been consumed, 
				-- or if there were no parts (the root) select the empty string.
				(select ids from parts where length(str) == 0 union all select '' limit 1), 
				?4 
			)`),
		name: ps.Prep(db, Name),
		noun: ps.Prep(db,
			// kind is transformed, but the number of parameters remains the same.
			Noun,
		),
		opposite: ps.Prep(db, Opposite),
		pair: ps.Prep(db,
			tables.Insert("mdl_pair", "domain", "relKind", "oneNoun", "otherNoun", "at"),
		),
		pat: ps.Prep(db,
			tables.Insert("mdl_pat", "kind", "labels", "result"),
		),
		plural: ps.Prep(db, Plural),
		rel: ps.Prep(db,
			tables.Insert("mdl_rel", "relKind", "oneKind", "otherKind", "cardinality", "at"),
		),
		rule: ps.Prep(db,
			tables.Insert("mdl_rule", "domain", "kind", "target", "phase", "filter", "prog", "at"),
		),
		// first: build a virtual [domain, noun, field] table
		// note: values are written per noun, not per domain; so the domain column is redundant once the noun id is known.
		// to get to the field id, we have to look at all possible fields for the noun.
		// given the kind of the noun, accept all fields who's kind is in its materialized path.
		// fix? some values are references to objects in the form "#domain::noun" -- should the be changed to ids?
		value: ps.Prep(db,
			`with parts(domain, nin, noun, fid, field) as (
			select mk.domain, mn.rowid, mn.noun, mf.rowid, mf.field
			from mdl_noun mn
			join mdl_kind mk
				on (mn.kind = mk.rowid)
			left join mdl_field mf
				where instr(',' || mk.rowid || ',' || mk.path, ',' || mf.kind || ','))
			insert into mdl_value(noun, field, value, at)
			select nin, fid, ?4, ?5
			from parts where domain=?1 and noun=?2 and field=?3`,
		),
	}
	if e := ps.Err(); e != nil {
		err = e
	} else {
		ret = m
	}
	return
}

type Writer struct {
	db *sql.DB
	assign,
	check,
	domain,
	field,
	grammar,
	kind,
	name,
	noun,
	opposite,
	pair,
	pat,
	plural,
	rel,
	rev,
	rule,
	value *sql.Stmt
}

func (m *Writer) Assign(domain, kind, field, value string) (err error) {
	if f, e := m.findField(domain, kind, field); e != nil {
		err = e
	} else {
		_, err = m.assign.Exec(f, value)
	}
	return
}

func (m *Writer) Check(domain, name, value string, affinity affine.Affinity, prog, at string) (err error) {
	if d, e := m.findDomain(domain); e != nil {
		err = e
	} else {
		_, err = m.check.Exec(d, name, value, affinity, prog, at)
	}
	return
}

// fix: are we forcing parent domains to exist before writing?
// that mgiht be cool .... but maybe this is the wrong level?
func (m *Writer) Domain(domain, requires, at string) (err error) {
	_, err = m.domain.Exec(domain, requires, at)
	return
}

func (m *Writer) Field(domain, kind, field string, affinity affine.Affinity, typeName, at string) (err error) {
	if _, k, e := m.findKind(domain, kind); e != nil {
		err = errutil.Fmt("%w trying to add field %q", e, field)
	} else if _, t, e := m.findOptionalKind(domain, typeName); e != nil {
		err = errutil.Fmt("%w trying to write field %q", e, field)
	} else if _, e := m.field.Exec(k, field, affinity, t, at); e != nil {
		err = errutil.Fmt("%w for (%s.%s.%s)", e, domain, kind, field)
	}
	return
}

func (m *Writer) Grammar(domain, name string, prog *grammar.Directive, at string) (err error) {
	if str, e := marshalout(prog); e != nil {
		err = e
	} else if d, e := m.findDomain(domain); e != nil {
		err = e
	} else {
		_, err = m.grammar.Exec(d, name, str, at)
	}
	return
}

func (m *Writer) Kind(domain, kind, path, at string) (err error) {
	if d, e := m.findDomain(domain); e != nil {
		err = e
	} else {
		_, err = m.kind.Exec(d, kind, path, at)
	}
	return
}

func (m *Writer) Name(domain, noun, name string, rank int, at string) (err error) {
	if _, n, e := m.findNoun(domain, noun); e != nil {
		err = e // ^ for now, this can be a derived domain
	} else {
		// uses the domain of the declaration
		_, err = m.name.Exec(domain, n, name, rank, at)
	}
	return
}

func (m *Writer) Noun(domain, noun, kind, at string) (err error) {
	if _, k, e := m.findKind(domain, kind); e != nil {
		err = e
	} else {
		// uses the domain of the declaration
		_, err = m.noun.Exec(domain, noun, k, at)
	}
	return
}

func (m *Writer) Opposite(domain, oneWord, otherWord, at string) (err error) {
	if d, e := m.findDomain(domain); e != nil {
		err = e
	} else {
		_, err = m.opposite.Exec(d, oneWord, otherWord, at)
	}
	return
}

func (m *Writer) Pair(domain, relKind, oneNoun, otherNoun, at string) (err error) {
	if _, k, e := m.findKind(domain, relKind); e != nil {
		err = e
	} else if _, one, e := m.findNoun(domain, oneNoun); e != nil {
		err = e
	} else if _, other, e := m.findNoun(domain, otherNoun); e != nil {
		err = e
	} else {
		// uses the domain of the declaration
		_, err = m.pair.Exec(domain, k, one, other, at)
	}
	return
}

func (m *Writer) Pat(domain, kind, labels, result string) (err error) {
	// tbd: labels are are comma-separated field names, should it be field ids?
	// similarly, result is a field, should it be a field id?
	// and... either way... should they be validated
	if d, k, e := m.findKind(domain, kind); e != nil {
		err = e
	} else if d != domain {
		err = errutil.New("pattern signature expected in the same domain as the pattern declaration")
	} else {
		_, err = m.pat.Exec(k, labels, result)
	}
	return
}

func (m *Writer) Plural(domain, many, one, at string) (err error) {
	if d, e := m.findDomain(domain); e != nil {
		err = e
	} else {
		_, err = m.plural.Exec(d, many, one, at)
	}
	return
}

func (m *Writer) Rel(domain, relKind, oneKind, otherKind, cardinality, at string) (err error) {
	if d, rel, e := m.findKind(domain, relKind); e != nil {
		err = e
	} else if d != domain {
		err = errutil.New("relation signature expected in the same domain as relation declaration")
	} else if _, one, e := m.findKind(domain, oneKind); e != nil {
		err = e
	} else if _, other, e := m.findKind(domain, otherKind); e != nil {
		err = e
	} else {
		_, err = m.rel.Exec(rel, one, other, cardinality, at)
	}
	return
}

func (m *Writer) Rule(domain, pattern, target string, phase int, filter, prog, at string) (err error) {
	if _, k, e := m.findKind(domain, pattern); e != nil {
		err = e
	} else if _, t, e := m.findOptionalKind(domain, target); e != nil {
		err = e
	} else {
		_, err = m.rule.Exec(domain, k, t, phase, filter, prog, at)
	}
	return
}

func (m *Writer) Value(domain, noun, field, value, at string) (err error) {
	_, err = m.value.Exec(domain, noun, field, value, at)
	return
}

// shared generic marshal prog to text
func marshalout(cmd interface{}) (ret string, err error) {
	if cmd != nil {
		if m, ok := cmd.(jsn.Marshalee); !ok {
			err = errutil.Fmt("can only marshal autogenerated types (%T)", cmd)
		} else {
			ret, err = cout.Marshal(m, literal.CompactEncoder)
		}
	}
	return
}
