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
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
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
		kind: ps.Prep(db,
			tables.Insert("mdl_kind", "domain", "kind", "path", "at"),
		),
		name: ps.Prep(db,
			tables.Insert("mdl_name", "domain", "noun", "name", "rank", "at"),
		),
		noun: ps.Prep(db,
			// the domain tells the scope in which the noun was defined
			// ( the same as - or a child of - the domain of the kind )
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
	if kid, e := m.findRequiredKind(domain, aspect); e != nil {
		err = e
	} else if !strings.HasSuffix(kid.fullpath(), m.aspectPath) {
		err = errutil.Fmt("kind %q in domain %q is not an aspect", aspect, domain)
	} else if strings.Count(kid.fullpath(), ",") != 3 {
		// tbd: could loosen this; for now it simplifies writing the aspects;
		// no need to check for conflicting traits.
		err = errutil.Fmt("can't create aspect of %q; kinds of aspects can't be inherited", aspect)
	} else if e := m.db.QueryRow(`
				select count(*) 
				from mdl_field mf 
				where mf.kind = ?1
				`, kid.id).Scan(&existingTraits); e != nil {
		err = errutil.New("database error", e)
	} else if existingTraits > 0 {
		err = errutil.Fmt("aspect %q from %q already has traits", aspect, domain)
	} else if kid.domain != domain {
		err = errutil.Fmt("cant add traits to aspect %q; traits are expected to exist in the same domain as the aspect. was %q now %q",
			aspect, kid.domain, domain)
	} else {
		for _, t := range traits {
			if _, e := m.field.Exec(domain, kid.id, t, affine.Bool, nil, at); e != nil {
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
	} else if kid, e := m.findRequiredKind(domain, kind); e != nil {
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
		and mf.field = ?2`, kid.id, field).Scan(&prev.id, &prev.domain, &prev.aff, &prev.out); e == sql.ErrNoRows {
			err = errutil.Fmt("%w field in assignment %q of kind %q in domain %q",
				Missing, field, kind, domain)
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

// fix: are we forcing/checking parent domains to exist before writing?
// that might be cool .... but maybe this is the wrong level?
func (m *Writer) Domain(domain, requires, at string) (err error) {
	_, err = m.domain.Exec(domain, requires, at)
	return
}

func (m *Writer) Member(domain, kind, field string, aff affine.Affinity, cls, at string) (err error) {
	if kid, e := m.findRequiredKind(domain, kind); e != nil {
		err = errutil.Fmt("%w trying to add field %q", e, field)
	} else if cls, e := m.findOptionalKind(domain, cls); e != nil {
		err = errutil.Fmt("%w trying to write field %q", e, field)
	} else {
		err = m.addField(domain, kid, cls, field, aff, at)
	}
	return
}

func (m *Writer) Parameter(domain, kind, field string, aff affine.Affinity, cls, at string) (err error) {
	if kid, e := m.findRequiredKind(domain, kind); e != nil {
		err = errutil.Fmt("%w trying to add parameter %q", e, field)
	} else if cls, e := m.findOptionalKind(domain, cls); e != nil {
		err = errutil.Fmt("%w trying to write parameter %q", e, field)
	} else if kid.domain != domain {
		err = errutil.Fmt("%w new parameter %q of %q expected in the same domain as the original declaration; was %q now %q",
			Conflict, field, kind, kid.domain, domain)
	} else if e := m.addField(domain, kid, cls, field, aff, at); e != nil {
		err = e
	} else if res, e := m.db.Exec(`
		insert into mdl_pat(kind, labels, result)
		values(?1, ?2, null)
		on conflict do update 
		set labels = labels ||','|| ?2
		where result is null
		`, kid.id, field); e != nil {
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
	if kid, e := m.findRequiredKind(domain, kind); e != nil {
		err = errutil.Fmt("%w trying to add result %q", e, field)
	} else if cls, e := m.findOptionalKind(domain, cls); e != nil {
		err = errutil.Fmt("%w trying to write result %q", e, field)
	} else if kid.domain != domain {
		err = errutil.Fmt("%w new result %q of %q expected in the same domain as the original declaration; was %q now %q",
			Conflict, field, kind, kid.domain, domain)
	} else if e := m.addField(domain, kid, cls, field, aff, at); e != nil {
		err = e
	} else if res, e := m.db.Exec(`
		insert into mdl_pat(kind, labels, result)
		values(?1, null, ?2)
		on conflict do update 
		set result=?2
		where result is null
		`, kid.id, field); e != nil {
		err = e
	} else if rows, e := res.RowsAffected(); e != nil {
		err = e
	} else if rows == 0 {
		err = errutil.Fmt("unexpected result %q for kind %q in domain %q",
			field, kind, domain)
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

// this duplicates the algorithm used by Noun()
func (m *Writer) Kind(domain, kind, parent, at string) (err error) {
	if parent, e := m.findOptionalKind(domain, parent); e != nil {
		err = e
	} else if len(parent.name) > 0 && parent.id == 0 {
		err = errutil.Fmt("%w ancestor %q",
			Missing, parent.name)
	} else if kind, e := m.findKind(domain, kind); e != nil {
		err = e
	} else if kind.id == 0 {
		// easiest is if the name has never been mentioned before;
		// we verified the other inputs already, so insert:
		if res, e := m.kind.Exec(domain, kind.name, trimPath(parent.fullpath()), at); e != nil {
			err = errutil.New("database error", e)
		} else if kind.name == kindsOf.Aspect.String() {
			// fix? it would probably be better to have a separate table of: domain, aspect, trait
			// currently, the runtime expects that aspects are a kind, and its traits are fields.
			if i, e := res.LastInsertId(); e != nil {
				err = e
			} else {
				m.aspectPath = "," + strconv.FormatInt(i, 10) + ","
			}
		}
	} else if parent.id != 0 { // note: if the kind exists, ignore nil parents.
		// does the newly specified ancestor contain the existing parent?
		// then we are ratcheting down. (ex. new: ,c,b,a,)  ( existing: ,a, )
		if strings.HasSuffix(kind.path, parent.fullpath()) {
			// does the existing parent fully contain the new ancestor?
			// then its a duplicate request (ex. existing: ,c,b,a,)  ( new: ,a, )
			err = errutil.Fmt("%w %q already declared as an ancestor of %q.",
				Duplicate, kind.name, parent.name)
		} else if strings.HasSuffix(parent.fullpath(), kind.path) {
			if kind.domain != domain {
				// if it was declared in a different domain: we can't change it now.
				err = errutil.Fmt("%w can't redefine the ancestor of %q as %q; the domains differ: was %q, now %q.",
					Conflict, kind.name, parent.name, kind.domain, domain)
			} else if res, e := m.db.Exec(`update mdl_kind set path = ?2 where rowid = ?1`,
				kind.id, trimPath(parent.fullpath())); e != nil {
				err = e
			} else if cnt, e := res.RowsAffected(); cnt != 1 {
				err = errutil.New("unexpected error updating hierarchy of %q; %d rows affected.",
					kind.name, cnt)
			} else if e != nil {
				err = e
			}
		} else if kind.domain != domain {
			// unrelated completely? then its an error
			err = errutil.Fmt("%w can't redefine the ancestor of %q as %q; the domains differ: was %q, now %q.",
				Conflict, kind.name, parent.name, kind.domain, domain)
		} else {
			// its possible some future definition might allow this to happen.
			err = errutil.Fmt("%w a definition in domain %q that would allow %q to have the ancestor %q; the hierarchies differ.",
				Missing, domain, kind.name, parent.name)
		}
	}
	return
}

func (m *Writer) Name(domain, noun, name string, rank int, at string) (err error) {
	if noun, e := m.findRequiredNoun(domain, noun, nounSansKind); e != nil {
		err = e // ^ for now, this can be a derived domain
	} else {
		var exists bool
		if e := m.db.QueryRow(`
	select 1
	from mdl_name mn
	join domain_tree
		on (uses = domain)
	where base = ?1
	and noun = ?2
	and name = ?3`, domain, noun.id, name).Scan(&exists); e != nil && e != sql.ErrNoRows {
			err = errutil.New("database error", e)
		} else if exists {
			err = errutil.Fmt("%w %q already an alias of %q", Duplicate, name, noun.name)
		} else if _, e := m.name.Exec(domain, noun.id, name, rank, at); e != nil {
			err = errutil.Fmt("database error writing name %q for noun %q in domain %q, %v", name, noun.name, domain, e)
		}
	}
	return
}

// this duplicates the algorithm used by Kind()
func (m *Writer) Noun(domain, name, ancestor, at string) (err error) {
	if parent, e := m.findRequiredKind(domain, ancestor); e != nil {
		err = e
	} else if prev, e := m.findNoun(domain, name, nounWithKind); e != nil {
		err = e
	} else if prev.id == 0 {
		// easiest is if the name has never been mentioned before;
		// we verified the other inputs already, so just go ahead and insert:
		_, err = m.noun.Exec(domain, name, parent.id, at)
	} else if prev.domain != domain {
		// if it was declared in a different domain: we can't change it now.
		err = errutil.Fmt("%w new ancestor %q of %q expected in the same domain as the original declaration; was %q now %q",
			Conflict, ancestor, name, prev.domain, domain)
	} else {
		// does the newly specified ancestor contain the existing parent?
		// then we are ratcheting down. (ex. new: ,c,b,a,)  ( existing: ,a, )
		if strings.HasSuffix(parent.fullpath(), prev.fullpath) {
			if res, e := m.db.Exec(`update mdl_noun set kind = ?2 where rowid = ?1`,
				prev.id, parent.id); e != nil {
				err = e
			} else if cnt, e := res.RowsAffected(); cnt != 1 {
				err = errutil.New("unexpected error updating hierarchy of %q; %d rows affected",
					name, cnt)
			} else if e != nil {
				err = e
			}
		} else if strings.HasSuffix(prev.fullpath, parent.fullpath()) {
			// does the existing kind fully contain the new kind?
			// then its a duplicate request (ex. existing: ,c,b,a,)  ( new: ,a, )
			err = errutil.Fmt("%w %q already declared as a kind of %q",
				Duplicate, name, ancestor)
		} else {
			// unrelated completely? then its an error
			err = errutil.Fmt("%w can't redefine kind of %q as %q",
				Conflict, name, ancestor)
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
	if kid, e := m.findRequiredKind(domain, relKind); e != nil {
		err = e
	} else if one, e := m.findRequiredNoun(domain, oneNoun, nounSansKind); e != nil {
		err = e
	} else if other, e := m.findRequiredNoun(domain, otherNoun, nounSansKind); e != nil {
		err = e
	} else {
		// uses the domain of the declaration
		_, err = m.pair.Exec(domain, kid.id, one.id, other.id, at)
	}
	return
}

func (m *Writer) Pat(domain, kind, labels, result string) (err error) {
	// tbd: labels are are comma-separated field names, should it be field ids?
	// similarly, result is a field, should it be a field id?
	// and... either way... should they be validated
	if kid, e := m.findRequiredKind(domain, kind); e != nil {
		err = e
	} else if kid.domain != domain {
		err = errutil.Fmt("%w pattern %q signature expected in the same domain as the pattern declaration; was %q now %q",
			Conflict, kind, kid.domain, domain)
	} else {
		_, err = m.pat.Exec(kid.id, labels, result)
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
	if rel, e := m.findRequiredKind(domain, relKind); e != nil {
		err = e
	} else if rel.domain != domain {
		err = errutil.New("relation signature expected in the same domain as relation declaration")
	} else if one, e := m.findRequiredKind(domain, oneKind); e != nil {
		err = e
	} else if other, e := m.findRequiredKind(domain, otherKind); e != nil {
		err = e
	} else {
		_, err = m.rel.Exec(rel.id, one.id, other.id, cardinality, at)
	}
	return
}

func (m *Writer) Rule(domain, pattern, target string, phase int, filter rt.BoolEval, exe []rt.Execute, at string) (err error) {
	if kid, e := m.findRequiredKind(domain, pattern); e != nil {
		err = e
	} else if tgt, e := m.findOptionalKind(domain, target); e != nil {
		err = e
	} else {
		slice := rt.Execute_Slice(exe)
		if filter, e := marshalout(filter); e != nil {
			err = e
		} else if prog, e := marshalout(&slice); e != nil {
			err = e
		} else {
			_, err = m.rule.Exec(domain, kid.id, tgt.id, phase, filter, prog, at)
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
func marshalout(cmd any) (ret string, err error) {
	if cmd != nil {
		if m, ok := cmd.(jsn.Marshalee); !ok {
			err = errutil.Fmt("can only marshal autogenerated types (%T)", cmd)
		} else {
			ret, err = cout.Marshal(m, literal.CompactEncoder)
		}
	}
	return
}
