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
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

/**
 *
 */
func NewModeler(db *sql.DB) (ret Modeler, err error) {
	var ps tables.Prep
	m := &Writer{
		db:         db,
		aspectPath: "XXX", // set to something that wont match until its set properly.
		assign: ps.Prep(db,
			tables.Insert("mdl_default", "field", "value"),
		),
		check: ps.Prep(db,
			tables.Insert("mdl_check", "domain", "name", "value", "affinity", "prog", "at"),
		),
		domain: ps.Prep(db,
			tables.Insert("mdl_domain", "domain", "requires", "at"),
		),
		field: ps.Prep(db,
			tables.Insert("mdl_field", "domain", "kind", "field", "affinity", "type", "at"),
		),
		grammar: ps.Prep(db,
			tables.Insert("mdl_grammar", "domain", "name", "prog", "at"),
		),

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
		name: ps.Prep(db,
			tables.Insert("mdl_name", "domain", "noun", "name", "rank", "at"),
		),
		noun: ps.Prep(db,
			// the domain tells the scope in which the noun was defined
			// ( the same as - or a child of - the domain of the kind )

			// kind is transformed, but the number of parameters remains the same.
			tables.Insert("mdl_noun", "domain", "noun", "kind", "at"),
		),
		opposite: ps.Prep(db,
			`insert into mdl_rev(domain, oneWord, otherWord, at) 
				values(?1, ?2, ?3, ?4), (?1, ?3, ?2, ?4)`,
		),
		pair: ps.Prep(db,
			// domain captures the scope in which the pairing was defined.
			// within that scope: the noun, relation, and otherNoun are all unique names --
			// even if they are not unique globally, and even if they a broader/different scope than the pair's domain.
			tables.Insert("mdl_pair", "domain", "relKind", "oneNoun", "otherNoun", "at"),
		),
		plural: ps.Prep(db,
			// a plural word ("many") can have at most one singular definition per domain
			// ie. "people" and "persons" are valid plurals of "person",
			// but "people" as a singular can only be defined as "person" ( not also "human" )
			tables.Insert("mdl_plural", "domain", "many", "one", "at"),
		),
		rel: ps.Prep(db,
			// relation and constraint between two kinds of nouns
			//  fix? the data is duplicated in kinds and fields... should this be removed?
			// might also consider adding a "cardinality" field to the relation kind, and then use init for individual relations
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
	// some ugly caching:
	aspectPath string // ex. ',4,'
}

// tbd: perhaps writing the aspect to its own table would be best
// join at runtime to synthesize fields.
// ( could potentially write both as a bridge )
func (m *Writer) Aspect(domain, aspect, at string, traits []string) (err error) {
	var existingTraits int
	if declaringDomain, ancestry, kid, e := m.pathOfKind(domain, aspect); e != nil {
		err = e
	} else if !strings.HasSuffix(ancestry, m.aspectPath) {
		err = errutil.Fmt("kind %q from %q is not an aspect", aspect, domain)
	} else if strings.Count(ancestry, ",") != 3 {
		// tbd: could loosen this; for now it simplifies writing the aspects;
		// no need to check for conflicting traits.
		err = errutil.Fmt("can't create aspect of %q; kinds of aspects can't be inherited", aspect)
	} else if e := m.db.QueryRow(`
				select count(*) 
				from mdl_field mf 
				where mf.kind = ?1
				`, kid).Scan(&existingTraits); e != nil {
		err = errutil.New("database error", e)
	} else if existingTraits > 0 {
		err = errutil.Fmt("aspect %q from %q already has traits", aspect, domain)
	} else if declaringDomain != domain {
		err = errutil.Fmt("cant add traits to aspect %q; traits are expected to exist in the same domain as the aspect. was %q now %q", aspect, declaringDomain, domain)
	} else {
		for _, t := range traits {
			if _, e := m.field.Exec(domain, kid, t, affine.Bool, nil, at); e != nil {
				err = errutil.New("database error", e)
				break
			}
		}
	}
	return
}

func (m *Writer) Check(domain, name string, v literal.LiteralValue, exe []rt.Execute, at string) (err error) {
	if d, e := m.findDomain(domain); e != nil {
		err = e
	} else if out, e := marshalout(v); e != nil {
		err = e
	} else {
		slice := rt.Execute_Slice(exe)
		if prog, e := marshalout(&slice); e != nil {
			err = e
		} else {
			aff := v.Affinity()
			_, err = m.check.Exec(d, name, out, aff, prog, at)
		}
	}
	return
}

func (m *Writer) Default(domain, kind, field string, v assign.Assignment) (err error) {
	if out, e := marshalout(v); e != nil {
		err = e
	} else if _, kid, e := m.findKind(domain, kind); e != nil {
		err = e
	} else {
		var prev struct {
			id     int
			domain string
			aff    affine.Affinity
			out    *string
		}
		if e := m.db.QueryRow(`
		select mf.rowid, domain, affinity, value
		from mdl_field mf
		left join mdl_default md
			on(md.field = mf.rowid)
		where mf.kind = ?1
		and mf.field = ?2`, kid, field).Scan(&prev.id, &prev.domain, &prev.aff, &prev.out); e == sql.ErrNoRows {
			err = errutil.Fmt("assignment requested for unknown field %q of kind %q in domain %q", field, kind, domain)
		} else if e != nil {
			err = e
		} else {
			if domain != prev.domain {
				// currently assuming that fields are initialized in the same domain as they are declared
				// that wont always be true... ex. derived classes or constraints
				err = errutil.Fmt("%w new assignment for field %q of kind %q differs in domain; was %q now %q.",
					Conflict, field, kind, prev.domain, domain)
			} else if prev.out != nil {
				if out == *prev.out {
					err = errutil.Fmt("%w assignment for field %q of kind %q in domain %q",
						Duplicate, field, kind, domain)
				} else {
					err = errutil.Fmt("%w new assignment for field %q of kind %q differs",
						Conflict, field, kind)
				}
			} else if aff := assign.GetAffinity(v); aff != prev.aff {
				err = errutil.Fmt("%w mismatched assignment for field %q of kind %q; field is %s, assignment was %s",
					Conflict, field, kind,
					prev.aff, aff)
			} else {
				_, err = m.assign.Exec(prev.id, out)
			}

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

func (m *Writer) Member(domain, kind, field string, aff affine.Affinity, cls, at string) (err error) {
	_, err = m.Field(domain, kind, field, aff, cls, at)
	return
}

func (m *Writer) Parameter(domain, kind, field string, aff affine.Affinity, cls, at string) (err error) {
	if k, e := m.Field(domain, kind, field, aff, cls, at); e != nil {
		err = e
	} else if res, e := m.db.Exec(`
		insert into mdl_pat(kind, labels, result)
		values(?1, ?2, null)
		on conflict do update 
		set labels = labels ||','|| ?2
		where result is null
		`, k, field); e != nil {
		err = e
	} else if rows, e := res.RowsAffected(); e != nil {
		err = e
	} else if rows == 0 {
		err = errutil.Fmt("unexpected parameter %q for kind %q in domain %q",
			field, kind, domain)
	}
	return
}

func (m *Writer) Result(domain, kind, field string, aff affine.Affinity, cls, at string) (err error) {
	if k, e := m.Field(domain, kind, field, aff, cls, at); e != nil {
		err = e
	} else if res, e := m.db.Exec(`
		insert into mdl_pat(kind, labels, result)
		values(?1, null, ?2)
		on conflict do update 
		set result=?2
		where result is null
		`, k, field); e != nil {
		err = e
	} else if rows, e := res.RowsAffected(); e != nil {
		err = e
	} else if rows == 0 {
		err = errutil.Fmt("unexpected result %q for kind %q in domain %q",
			field, kind, domain)
	}
	return
}

func (m *Writer) Field(domain, kind, field string, aff affine.Affinity, cls, at string) (retKind int, err error) {
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
)
-- all fields of the targeted kind:
, fieldsInKind as (
	select mk.domain, field as name, affinity, mf.type as typeId, mt.kind as typeName
	from mdl_field mf 
	join mdl_kind mk 
		-- does our ancestry (X) contain any of these kinds (Y)
		on ((mf.kind = mk.rowid) and instr(?1, ',' || mk.rowid || ',' ))
	left join mdl_kind mt 
		on (mt.rowid = mf.type)
)
-- fields and traits in the target kind
-- ( all of them, because we dont know what might conflict with a pending aspect )
, existingFields( origin, name, affinity, typeName ) as (
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
)
, pendingFields(name, aspect) as ( 
	-- the name of the field we're adding;
	select ?2, null
	union all 

	-- the names of traits when adding a field of type aspect; if any.
	select name, aspect
	from allTraits ma
	where (?3 = ma.kind)
)
select origin, name, affinity, typeName, aspect
from existingFields
join pendingFields
using(name)
`, ancestry, field, typeId, m.aspectPath); e != nil {
		err = errutil.New("database error", e)
	} else {
		var prev struct {
			name   string          // trait or field causing a conflict
			aspect sql.NullString  // aspect if any of the pending name
			origin string          // aspect or kind of the existing field
			aff    affine.Affinity // affinity of the existing field ( ex. 'bool' for aspects )
			cls    sql.NullString  // type name ( or null ) of existing field
		}
		if e := tables.ScanAll(rows, func() (err error) {
			// if the names differ, then the conflict is due to a trait ( being added or already existing )
			if prev.name != field {
				// adding an aspect: the conflict reports the pending aspect so this case can be detected
				if prev.aspect.String == cls {
					// is there a way to determine whether the origin is a domain or aspect
					err = errutil.Fmt("%w new field for kind %q of aspect %q conflicts with existing field %q from %s",
						Conflict, kind, field, prev.name, prev.origin)
				} else if prev.aspect.Valid {
					err = errutil.Fmt("%w new field for kind %q of aspect %q conflicts with trait %q from aspect %q",
						Conflict, kind, field,
						prev.name, prev.aspect.String)
				} else {
					// when does this show up?
					err = errutil.Fmt("%w field %q for kind %q was %s(%s) from %s, now %s(%s) in %q",
						Conflict, field, kind,
						prev.aff, prev.cls.String, prev.origin,
						aff, cls, domain)
				}
			} else if aff == prev.aff && cls == prev.cls.String {
				// if the affinity and typeName are the same, then its a duplicate
				err = errutil.Fmt("%w field %q for kind %q of %s(%s) from %s and now domain %q",
					Duplicate, field, kind,
					aff, cls,
					prev.origin, domain)
			} else {
				// otherwise, its a conflict
				err = errutil.Fmt("%w field %q for kind %q of %s(%s) from %s was redefined as %s(%s) in domain %q",
					Conflict, field, kind,
					prev.aff, prev.cls.String, prev.origin,
					aff, cls, domain)
			}
			return
		}, &prev.origin, &prev.name, &prev.aff, &prev.cls, &prev.aspect); e != nil {
			err = e
		} else {
			// err was nil, we can write the field:
			if _, e := m.field.Exec(domain, kid, field, aff, typeId, at); e != nil {
				err = errutil.Fmt("%w for (%s.%s.%s)", e, domain, kind, field)
			} else {
				retKind = kid
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
		err = errutil.New("database error", e)
	} else if kind == "aspects" {
		if i, e := res.LastInsertId(); e != nil {
			err = errutil.New("database error", e)
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

// create table mdl_noun( domain text not null, noun text, kind int not null, at text )
func (m *Writer) Noun(domain, noun, kind, at string) (err error) {
	if _, newAncestry, kid, e := m.pathOfKind(domain, kind); e != nil {
		err = e
	} else if d, id, existingAncestry, e := m.pathOfOptionalNoun(domain, noun); e != nil {
		err = e
	} else if id == 0 {
		// easiest is if the noun has never been mentioned before;
		// we verified the kind first thing, so just go ahead and insert:
		_, err = m.noun.Exec(domain, noun, kid, at)
	} else if d != domain {
		// if it has been declared, and in a different domain: that's an error.
		err = errutil.Fmt("%w kind %q of noun %q expected in the same domain as the noun declaration; was %q now %q",
			Conflict, kind, noun, d, domain)
	} else {
		// does the newly specified kind contain the existing kind?
		// then we are ratcheting down. (ex. new: ,c,b,a,)  ( existing: ,a, )
		if strings.HasSuffix(newAncestry, existingAncestry) {
			if res, e := m.db.Exec(`update mdl_noun set kind = ?2 where rowid = ?1`, id, kid); e != nil {
				err = e
			} else if cnt, e := res.RowsAffected(); cnt != 1 {
				err = errutil.New("unexpected error updating noun hierarchy %d rows affected", cnt)
			} else if e != nil {
				err = e
			}
		} else if strings.HasSuffix(existingAncestry, newAncestry) {
			// does the existing kind fully contain the new kind?
			// then its a duplicate request (ex. existing: ,c,b,a,)  ( new: ,a, )
			err = errutil.Fmt("%w noun %q already declared as %q",
				Duplicate, noun, kind)
		} else {
			// unrelated completely? then its an error
			err = errutil.Fmt("%w can't redefine kind of noun %q as %q",
				Conflict, noun, kind)
		}
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
		err = errutil.New("database error", e)
	} else {
		var x, y, from string
		if e := tables.ScanAll(rows, func() (err error) {
			// the testing is a bit weird so we handle it all app side
			if (x == a && y == b) || (x == b && y == a) {
				err = errutil.Fmt("%w opposite %q <=> %q in %q and %q", Duplicate, a, b, from, domain)
			} else if x == a || y == a || x == b || y == b {
				err = errutil.Fmt(
					"%w %q <=> %q defined as opposites in %q now %q <=> %q in %q",
					Conflict, x, y, from, a, b, domain)
			}
			return
		}, &x, &y, &from); e != nil {
			err = e
		} else {
			// writes the opposite paring as well
			_, err = m.opposite.Exec(d, a, b, at)
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
		err = errutil.Fmt("%w pattern %q signature expected in the same domain as the pattern declaration; was %q now %q",
			Conflict, kind, d, domain)
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
		err = errutil.New("database error", e)
	} else {
		var prev, from string
		if e := tables.ScanAll(rows, func() (err error) {
			if prev == one {
				err = errutil.Fmt("%w plural %q was %q in %q and %q",
					Duplicate, many, one, from, domain)
			} else {
				err = errutil.Fmt("%w plural %q had singular %q (in %q) wants %q (in %q)",
					Conflict, many, prev, from, one, domain)
			}
			return
		}, &prev, &from); e != nil {
			err = e
		} else {
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

func (m *Writer) Rule(domain, pattern, target string, phase int, filter rt.BoolEval, exe []rt.Execute, at string) (err error) {
	if _, k, e := m.findKind(domain, pattern); e != nil {
		err = e
	} else if _, t, e := m.findOptionalKind(domain, target); e != nil {
		err = e
	} else {
		slice := rt.Execute_Slice(exe)
		if filter, e := marshalout(filter); e != nil {
			err = e
		} else if prog, e := marshalout(&slice); e != nil {
			err = e
		} else {
			_, err = m.rule.Exec(domain, k, t, phase, filter, prog, at)
		}
	}
	return
}

func (m *Writer) Value(domain, noun, field string, value literal.LiteralValue, at string) (err error) {
	if out, e := marshalout(value); e != nil {
		err = errutil.Append(err, e)
	} else {
		_, err = m.value.Exec(domain, noun, field, out, at)
	}
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
