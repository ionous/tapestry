package mdl

import (
	"database/sql"
	"strconv"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

/**
 *
 */
func NewModeler(db *sql.DB, duplicateWarnings func(error)) (ret Modeler, err error) {
	var ps tables.Prep
	if duplicateWarnings == nil {
		duplicateWarnings = func(error) {}
	}
	m := &Writer{
		db:         db,
		warn:       duplicateWarnings,
		aspectPath: "XXX", // set to something that wont match until its set properly.
		assign: ps.Prep(db,
			tables.Insert("mdl_default", "field", "value"),
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
		opposite: ps.Prep(db,
			`insert into mdl_rev(domain, oneWord, otherWord, at) 
				values(?1, ?2, ?3, ?4), (?1, ?3, ?2, ?4)`,
		),
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
	db   *sql.DB
	warn func(error)
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
	aspectPath string // ex. ',4,'
}

// tbd: perhaps writing the aspect to its own table would be best
// join at runtime to synthesize fields.
// ( could potentially write both as a bridge )
func (m *Writer) Aspect(domain, aspect, at string, traits []string) (err error) {
	if _, ancestry, kid, e := m.pathOfKind(domain, aspect); e != nil {
		err = e
	} else if !strings.HasSuffix(ancestry, m.aspectPath) {
		err = errutil.Fmt("kind %q from %q is not an aspect", aspect, domain)
	} else if strings.Count(ancestry, ",") != 3 {
		// tbd: could loosen this; for now it simplifies writing the aspects;
		// no need to check for conflicting traits.
		err = errutil.Fmt("can't create aspect of %q; kinds of aspects can't be inherited", aspect, domain)
	} else {
		// has the aspect already been assigned traits?
		var existingTraits int
		if e := m.db.QueryRow(`
				select count(*) 
				from mdl_field mf 
				where mf.kind = ?1
				`, kid).Scan(&existingTraits); e != nil {
			err = errutil.New(e, "failed to count existing traits")
		} else if existingTraits > 0 {
			err = errutil.Fmt("aspect %q from %q already has traits", aspect, domain)
		} else {
			for _, t := range traits {
				if _, e := m.field.Exec(kid, t, affine.Bool, nil, at); e != nil {
					err = errutil.Fmt("%w for (%s.%s.%s)", e, domain, aspect, t)
					break
				}
			}
		}
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

func (m *Writer) Default(domain, kind, field string, value assign.Assignment) (err error) {
	if value, e := marshalout(value); e != nil {
		err = e
	} else if declaringDomain, f, e := m.findField(domain, kind, field); e != nil {
		err = e
	} else if rows, e := m.db.Query(
		`select affinity, type
			from mdl_field
			where kind = ?1
			and field = ?2`, f, field); e != nil {
		err = e
	} else {
		var prevValue string
		var dupe int
		if e := tables.ScanAll(rows, func() (err error) {
			if value == prevValue {
				m.warn(errutil.Fmt("duplicate default assignments \"%s.%s\" in %q and %q",
					kind, field,
					declaringDomain, domain))
				dupe++
			} else {
				err = errutil.Fmt("conflict: default assignment \"%s.%s\" in %q changed in %q",
					kind, field,
					declaringDomain, domain)
			}
			return
		}, &prevValue); e != nil {
			err = e
		} else if dupe == 0 {
			_, err = m.assign.Exec(f, value)
		}
	}
	return
}

// fix: are we forcing parent domains to exist before writing?
// that mgiht be cool .... but maybe this is the wrong level?
func (m *Writer) Domain(domain, requires, at string) (err error) {
	_, err = m.domain.Exec(domain, requires, at)
	return
}

func (m *Writer) Field(domain, kind, field string, aff affine.Affinity, cls, at string) (err error) {
	if _, ancestry, kid, e := m.pathOfKind(domain, kind); e != nil {
		err = errutil.Fmt("%w trying to add field %q", e, field)
	} else if _, typeId, e := m.findOptionalKind(domain, cls); e != nil {
		err = errutil.Fmt("%w trying to write field %q", e, field)
	} else if rows, e := m.db.Query(`
-- all possible traits:
with allTraits as (	
	select mk.rowid as kind,    -- id of the aspect,
				 field as name,      -- name of trait 
	       mk.kind as aspect,  -- name of aspect
	       mk.domain          -- name of originating domain
	from mdl_field mf 
	join mdl_kind mk
		on(mf.kind = mk.rowid)
	-- where the field's kind (X) contains the aspect kind (Y)
	where instr(',' || mk.path, ?4 )
),
-- all fields of the targeted kind:
fieldsInKind as (
	select mk.domain, field as name, affinity, mf.type as typeId, mt.kind as typeName
	from mdl_field mf 
	join mdl_kind mk 
		-- does our ancestry (X) contain any of these kinds (Y)
		on ((mf.kind = mk.rowid) and instr(?1, ',' || mk.rowid || ',' ))
	left join mdl_kind mt 
		on (mt.rowid = mf.type)

),
-- fields and traits in the target kind
-- ( all of them, because we dont know what might conflict with a pending aspect )
existingNames( origin, name, affinity, typeName ) as (
	-- fields in the target kind
	select format('domain "%z"', domain), name, affinity, typeName
	from fieldsInKind

	union all

	-- traits in the target kind
	select format('aspect "%z"', ma.aspect), -- report the aspect as the origin 
				 ma.name,   -- trait name 
				 'bool',    -- fk.affinity is 'text', each trait is 'bool'
				 null       -- traits have null type currently.
	from fieldsInKind fk
	join allTraits ma
		on (fk.typeId = ma.kind)
),
pendingNames(name, aspect) as ( 
	-- the name of the field we're adding;
	select ?2, null
	union all 

	-- the names of traits when adding a field of type aspect; if any.
	select name, aspect
	from allTraits ma
	where (?3 = ma.kind)
)
select origin, name, affinity, typeName, aspect
from existingNames
join pendingNames
using(name)
`, ancestry, field, typeId, m.aspectPath); e != nil {
		err = e
	} else {
		var con struct {
			name   string          // trait or field causing a conflict
			aspect sql.NullString  // aspect if any of the pending name
			origin string          // aspect or kind of the existing field
			aff    affine.Affinity // affinity of the existing field ( ex. 'bool' for aspects )
			cls    sql.NullString  // type name ( or null ) of existing field
		}
		var dupe error
		if e := tables.ScanAll(rows, func() (err error) {
			// if the names differ, then the conflict is due to a trait ( being added or already existing )
			if con.name != field {
				// adding an aspect: the conflict reports the pending aspect so this case can be detected
				if con.aspect.String == cls {
					// is there a way to determine whether the origin is a domain or aspect
					err = errutil.Fmt("conflict: new field for kind %q of aspect %q conflicts with existing field %q from %s",
						kind, field, con.name, con.origin)
				} else if con.aspect.Valid {
					err = errutil.Fmt("conflict: new field for kind %q of aspect %q conflicts with trait %q from aspect %q",
						kind, field,
						con.name, con.aspect.String)
				} else {
					// when does this show up?
					err = errutil.Fmt("conflict: field %q for kind %q was %s(%s) from %s, now %s(%s) in %q",
						field, kind,
						con.aff, con.cls.String, con.origin,
						aff, cls, domain)
				}
			} else if aff == con.aff && cls == con.cls.String {
				// if the affinity and typeName are the same, then its a duplicate
				dupe = errutil.Fmt("duplicate field %q for kind %q  of %s(%s) from %s and now domain %q",
					field, kind,
					aff, cls,
					con.origin, domain)
				err = dupe
			} else {
				// otherwise, its a conflict
				err = errutil.Fmt("conflict: field \"%s.%s\" was %s(%s) from %s, now %s(%s) in domain %q",
					kind, field,
					con.aff, con.cls.String, con.origin,
					aff, cls, domain)
			}
			return
		}, &con.origin, &con.name, &con.aff, &con.cls, &con.aspect); e != nil {
			if e == dupe {
				m.warn(e)
			} else {
				err = e // a conflict or sql error
			}
		} else {
			// err was nil, we can write the field:
			if _, e := m.field.Exec(kid, field, aff, typeId, at); e != nil {
				err = errutil.Fmt("%w for (%s.%s.%s)", e, domain, kind, field)
			}
		}
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
	} else if res, e := m.kind.Exec(d, kind, path, at); e != nil {
		err = e
	} else if kind == "aspects" {
		if i, e := res.LastInsertId(); e != nil {
			err = e
		} else {
			// fix? it would probably be better to have a separate table of: domain, aspect, trait
			// currently, the runtime expects that aspects are a kind, and its traits are fields.
			m.aspectPath = "," + strconv.FormatInt(i, 10) + ","
		}
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

func (m *Writer) Opposite(domain, a, b, at string) (err error) {
	if d, e := m.findDomain(domain); e != nil {
		err = e
	} else if rows, e := m.db.Query(
		`select oneWord, otherWord, domain
			from mdl_rev 
			join domain_tree
				on(uses=domain)
			where base = ?1`, d); e != nil {
		err = e
	} else {
		var x, y, from string
		var dupe int
		if e := tables.ScanAll(rows, func() (err error) {
			// the testing is a bit weird so we handle it all app side
			if (x == a && y == b) || (x == b && y == a) {
				m.warn(errutil.Fmt(
					"duplicate opposites: %q <=> %q in %q and %q",
					a, b, from, domain))
				dupe++
			} else if x == a || y == a || x == b || y == b {
				err = errutil.Fmt(
					"conflict: %q <=> %q defined as opposites in %q now %q <=> %q in %q",
					x, y, from, a, b, domain)
			}
			return
		}, &x, &y, &from); e != nil {
			err = e
		} else {
			if dupe == 0 {
				// writes the opposite paring as well
				_, err = m.opposite.Exec(d, a, b, at)
			}
		}
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
	} else if rows, e := m.db.Query(
		`select one, domain 
			from mdl_plural
			join domain_tree
				on(uses=domain)
			where base = ?1
			and many = ?2`, d, many); e != nil {
		err = e
	} else {
		var prev, from string
		var dupe int // log duplicates?
		if e := tables.ScanAll(rows, func() (err error) {
			if prev == one {
				m.warn(errutil.Fmt("duplicate plurals: %q was %q in %q and %q", many, one, from, domain))
				dupe++
			} else {
				err = errutil.Fmt("conflict: plural %q had singular %q (in %q) wants %q (in %q)", many, prev, from, one, domain)
			}
			return
		}, &prev, &from); e != nil {
			err = e
		} else if dupe == 0 {
			_, err = m.plural.Exec(d, many, one, at)
		}
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
