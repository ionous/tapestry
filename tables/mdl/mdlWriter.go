package mdl

import (
	"database/sql"
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
func NewModeler(db *sql.DB) (ret *Modeler, err error) {
	ret = &Modeler{
		db:          tables.NewCache(db),
		aspectPath:  "XXX", // set to something that wont match until its set properly.
		patternPath: "XXX",
	}
	return
}

// Modeler wraps writing to the model table
// so the implementation can handle verifying dependent names when needed.
type Modeler struct {
	db *tables.Cache
	// some ugly caching:
	patternPath,
	aspectPath string // ex. ',4,'
}

// tbd: perhaps writing the aspect to its own table would be best
// join at runtime to synthesize fields.
// ( could potentially write both as a bridge )
func (m *Modeler) Aspect(domain, aspect, at string, traits []string) (err error) {
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
			if _, e := m.db.Exec(mdl_field, domain, kid.id, t, affine.Bool, nil, at); e != nil {
				err = errutil.New("database error", e)
				break
			}
		}
	}
	return
}

var mdl_check = tables.Insert("mdl_check", "domain", "name", "value", "affinity", "prog", "at")

// author tests of stories
func (m *Modeler) Check(domain, name string, value literal.LiteralValue, prog []rt.Execute, at string) (err error) {
	if d, e := m.findDomain(domain); e != nil {
		err = e
	} else {
		var prev struct {
			id          int64
			domain      string
			prog, value bool
		}
		if e := m.db.QueryRow(
			`select mc.rowid, 
					mc.domain,
					length(coalesce(mc.prog,'')),
					length(coalesce(mc.value,'')) 
			from mdl_check mc
			join domain_tree dt
				on (dt.uses = mc.domain)
			where base = ?1
			and name = ?2
		`, domain, name).Scan(&prev.id, &prev.domain, &prev.prog, &prev.value); e != nil && e != sql.ErrNoRows {
			err = errutil.New("database error", e)
		} else {
			// the user can write the check in parts if they so desire:
			if prev.id != 0 {
				if prog != nil && prev.prog {
					e := errutil.Fmt("%w new program for check %q in %q already specified in %q",
						Conflict, name, domain, prev.domain)
					err = errutil.Append(err, e)
				}
				if value != nil && prev.value {
					e := errutil.Fmt("%w new expectation for check %q in %q already specified in %q",
						Conflict, name, domain, prev.domain)
					err = e
				}
			}

			if err == nil {
				if prog, e := marshalprog(prog); e != nil {
					err = e
				} else if out, e := marshalout(value); e != nil {
					err = e
				} else {
					var aff affine.Affinity
					if value != nil {
						aff = value.Affinity()
					}
					if prev.id == 0 {
						_, err = m.db.Exec(mdl_check, d, name, out, aff, prog, at)
					} else {
						if res, e := m.db.Exec(`update mdl_check 
						set prog = (case when ?2 then ?2 else prog end),
						set value = (case when ?3 then ?3 else value end),
						where rowid = ?1`,
							prev.id, prev.prog, prev.value); e != nil {
							err = e
						} else if cnt, e := res.RowsAffected(); cnt != 1 {
							err = e
						} else {
							err = errutil.New("unexpected error updating check %q; %d rows affected.",
								name, cnt)
						}
					}
				}
			}
		}
	}
	return
}

var mdl_default = tables.Insert("mdl_default", "field", "value")

// the pattern half of Start; domain, kind, field are a pointer into Field
// value should be a marshaled compact value
func (m *Modeler) Default(domain, kind, field string, v assign.Assignment) (err error) {
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
				_, err = m.db.Exec(mdl_default, prev.id, out)
			}

		}
	}
	return
}

var mdl_domain = tables.Insert("mdl_domain", "domain", "requires", "at")

// pairs of domain name and (domain) dependencies
// fix: are we forcing/checking parent domains to exist before writing?
// that might be cool .... but maybe this is the wrong level?
func (m *Modeler) Domain(domain, requires, at string) (err error) {
	_, err = m.db.Exec(mdl_domain, domain, requires, at)
	return
}

var mdl_fact = tables.Insert("mdl_fact", "domain", "fact", "value", "at")

// arbitrary key-value storage
func (m *Modeler) Fact(domain, fact, value, at string) (err error) {
	if domain, e := m.findDomain(domain); e != nil {
		err = e
	} else {
		var prev struct {
			domain, value string
		}
		q := m.db.QueryRow(`
		select mx.domain, mx.value
		from mdl_fact mx
		join domain_tree
			on (uses = domain)
		where base = ?1
		and fact = ?2`, domain, fact)
		switch e := q.Scan(&prev.domain, &prev.value); e {
		case sql.ErrNoRows:
			if _, e := m.db.Exec(mdl_fact, domain, fact, value, at); e != nil {
				err = errutil.New("database error", e)
			}
		case nil:
			if prev.value != value {
				err = errutil.Fmt("%w fact %q was %q in domain %q and now %q in domain %q",
					Conflict, fact, value, prev.domain, value, domain)
			} else {
				// do we want to warn about duplicate facts? isnt that kind of the point of them?
				// maybe eat at the weave level?
				err = errutil.Fmt("%w fact %q already declared in domain %q and now domain %q",
					Duplicate, fact, prev.domain, domain)
			}
		default:
			err = errutil.New("database error", e)
		}
	}
	return
}

var mdl_grammar = tables.Insert("mdl_grammar", "domain", "name", "prog", "at")

// player input parsing
func (m *Modeler) Grammar(domain, name string, prog *grammar.Directive, at string) (err error) {
	if prog, e := marshalout(prog); e != nil {
		err = e
	} else if d, e := m.findDomain(domain); e != nil {
		err = e
	} else {
		var prev struct {
			domain, prog string
		}
		e := m.db.QueryRow(
			`select mg.domain, mg.prog
			from mdl_grammar mg
			join domain_tree dt
				on (dt.uses = mg.domain)
			where base = ?1
			and name = ?2
		`, domain, name).Scan(&prev.domain, &prev.prog)
		switch e {
		case sql.ErrNoRows:
			_, err = m.db.Exec(mdl_grammar, d, name, prog, at)

		case nil:
			if prev.prog != prog {
				err = errutil.Fmt("%w grammar %q was %q in domain %q and now %q in domain %q",
					Conflict, name, prev.prog, prev.domain, prog, domain)
			} else {
				err = errutil.Fmt("%w grammar %q already declared in domain %q and now domain %q",
					Duplicate, name, prev.domain, domain)
			}

		default:
			err = errutil.New("database error", e)
		}
	}
	return
}

var mdl_kind = tables.Insert("mdl_kind", "domain", "kind", "singular", "path", "at")

// singular name of kind and materialized hierarchy of ancestors separated by commas
// this (somewhat) duplicates the algorithm used by Noun()
func (m *Modeler) Kind(domain, name, parent, at string) (err error) {
	if parent, e := m.findOptionalKind(domain, parent); e != nil {
		err = e
	} else if len(parent.name) > 0 && parent.id == 0 {
		// a specified ancestor doesn't exist.
		err = errutil.Fmt("%w ancestor %q", Missing, parent.name)
	} else if kind, e := m.findKind(domain, name); e != nil {
		err = e
	} else if kind.id == 0 {
		// manage singular and plural kinds
		// i don't like this much; especially be cause it depends so much on the first declaration
		// maybe better would be a name/names table that any named concept can use.
		// or just force everyone to use the right names.
		if singular, e := m.singularize(domain, name); e != nil {
			err = e
		} else {
			var optionalOne *string
			if singular != name {
				optionalOne = &singular
			}
			// easiest is if the name has never been mentioned before;
			// we verified the other inputs already, so insert:
			if res, e := m.db.Exec(mdl_kind, domain, name, optionalOne, trimPath(parent.fullpath()), at); e != nil {
				err = errutil.New("database error", e)
			} else {
				// cache result... sometimes
				switch name {
				case kindsOf.Aspect.String():
					err = updatePath(res, &m.aspectPath)
				case kindsOf.Pattern.String():
					err = updatePath(res, &m.patternPath)
				default:
					// super hacky..... hmmm...
					// if we've declared a new kind of a pattern:
					// write blanks into the mdl_pat; parameters and results use update only.
					if strings.HasSuffix(parent.fullpath(), m.patternPath) {
						if newid, e := res.LastInsertId(); e != nil {
							err = e
						} else {
							_, err = m.db.Exec(`insert into mdl_pat(kind) values(?1)`, newid)
						}
					}
				}
			}
		}
	} else if parent.id != 0 { // this ignore empty ancestors if the kind already existed.

		if !kind.exact && parent.numAncestors() < 2 {
			// we allow plural named kinds for nouns, etc. not for patterns and built in kinds.
			err = errutil.Fmt("%w ambiguously named kinds: %q (in domain %q) and %q (in %q)",
				Conflict, name, domain, kind.name, kind.domain)
		} else if strings.HasSuffix(parent.fullpath(), kind.fullpath()) {
			err = errutil.Fmt("%w circular reference detected %q already declared as an ancestor of %q.",
				Conflict, name, parent.name)
		} else if strings.HasSuffix(kind.path, parent.fullpath()) {
			// did the existing path fully contain the new ancestor?
			// then its a duplicate request (ex. `,c,b,a,` `,b,a,` )
			err = errutil.Fmt("%w %q already declared as an ancestor of %q.",
				Duplicate, name, parent.name)
		} else if strings.HasSuffix(parent.fullpath(), kind.path) {
			// is the newly specified ancestor more specific than the existing path?
			// then we are ratcheting down. (ex. `,c,b,a,` `,b,a,` )
			if kind.domain != domain {
				// if it was declared in a different domain: we can't change it now.
				err = errutil.Fmt("%w can't redefine the ancestor of %q as %q; the domains differ: was %q, now %q.",
					Conflict, name, parent.name, kind.domain, domain)
			} else if res, e := m.db.Exec(`update mdl_kind set path = ?2 where rowid = ?1`,
				kind.id, trimPath(parent.fullpath())); e != nil {
				err = e
			} else if cnt, e := res.RowsAffected(); cnt != 1 {
				err = errutil.New("unexpected error updating hierarchy of %q; %d rows affected.",
					name, cnt)
			} else if e != nil {
				err = e
			}
		} else if kind.domain != domain {
			// unrelated completely? then its an error
			err = errutil.Fmt("%w can't redefine the ancestor of %q as %q; the domains differ: was %q, now %q.",
				Conflict, name, parent.name, kind.domain, domain)
		} else {
			// its possible some future definition might allow this to happen.
			err = errutil.Fmt("%w a definition in domain %q that would allow %q to have the ancestor %q; the hierarchies differ.",
				Missing, domain, name, parent.name)
		}
	}
	return
}

// a generic field of the kind
func (m *Modeler) Member(domain, kind, field string, aff affine.Affinity, cls, at string) (err error) {
	if kid, e := m.findRequiredKind(domain, kind); e != nil {
		err = errutil.Fmt("%w trying to add field %q", e, field)
	} else if cls, e := m.findOptionalKind(domain, cls); e != nil {
		err = errutil.Fmt("%w trying to write field %q", e, field)
	} else {
		err = m.addField(domain, kid, cls, field, aff, at)
	}
	return
}

var mdl_name = tables.Insert("mdl_name", "domain", "noun", "name", "rank", "at")

// words for authors and game players refer to nouns
// follows the domain rules of Noun.
func (m *Modeler) Name(domain, noun, name string, rank int, at string) (err error) {
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
		} else if _, e := m.db.Exec(mdl_name, domain, noun.id, name, rank, at); e != nil {
			err = errutil.New("database error", e)
		}
	}
	return
}

// the domain tells the scope in which the noun was defined
// ( the same as - or a child of - the domain of the kind )
var mdl_noun = tables.Insert("mdl_noun", "domain", "noun", "kind", "at")

// the domain tells the scope in which the noun was defined
// ( the same as - or a child of - the domain of the kind ) error
// this duplicates the algorithm used by Kind()
func (m *Modeler) Noun(domain, name, ancestor, at string) (err error) {
	if parent, e := m.findRequiredKind(domain, ancestor); e != nil {
		err = e
	} else if prev, e := m.findNoun(domain, name, nounWithKind); e != nil {
		err = e
	} else if prev.id == 0 {
		// easiest is if the name has never been mentioned before;
		// we verified the other inputs already, so just go ahead and insert:
		_, err = m.db.Exec(mdl_noun, domain, name, parent.id, at)
	} else {
		// does the newly specified ancestor contain the existing parent?
		// then we are ratcheting down. (ex. new: ,c,b,a,)  ( existing: ,a, )
		if strings.HasSuffix(prev.fullpath, parent.fullpath()) {
			// does the existing kind fully contain the new kind?
			// then its a duplicate request (ex. existing: ,c,b,a,)  ( new: ,a, )
			err = errutil.Fmt("%w %q already declared as a kind of %q",
				Duplicate, name, ancestor)
		} else if !strings.HasSuffix(parent.fullpath(), prev.fullpath) {
			// unrelated completely? then its an error
			err = errutil.Fmt("%w can't redefine kind of %q as %q",
				Conflict, name, ancestor)
		} else if prev.domain != domain {
			// if it was declared in a different domain: we can't change it now.
			err = errutil.Fmt("%w new ancestor %q of %q expected in the same domain as the original declaration; was %q now %q",
				Conflict, ancestor, name, prev.domain, domain)
		} else {
			if res, e := m.db.Exec(`update mdl_noun set kind = ?2 where rowid = ?1`,
				prev.id, parent.id); e != nil {
				err = e
			} else if cnt, e := res.RowsAffected(); cnt != 1 {
				err = errutil.New("unexpected error updating hierarchy of %q; %d rows affected",
					name, cnt)
			} else if e != nil {
				err = e
			}
		}
	}
	return
}

var mdl_opposite = `insert into mdl_rev(domain, oneWord, otherWord, at) 
				values(?1, ?2, ?3, ?4), (?1, ?3, ?2, ?4)`

// domain captures the scope in which the pairing was defined.
// within that scope: the noun, relation, and otherNoun are all unique names --
// even if they are not unique globally, and even if they a broader/different scope than the pair's domain.
func (m *Modeler) Opposite(domain, a, b, at string) (err error) {
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
			_, err = m.db.Exec(mdl_opposite, d, a, b, at)
		}
	}
	return
}

// domain captures the scope in which the pairing was defined.
// within that scope: the noun, relation, and otherNoun are all unique names --
// even if they are not unique globally, and even if they a broader/different scope than the pair's domain.
var mdl_pair = tables.Insert("mdl_pair", "domain", "relKind", "oneNoun", "otherNoun", "at")

func (m *Modeler) Pair(domain, rel, oneNoun, otherNoun, at string) (err error) {
	if rel, e := m.findRequiredKind(domain, rel); e != nil {
		err = e
	} else if one, e := m.findRequiredNoun(domain, oneNoun, nounSansKind); e != nil {
		err = e
	} else if other, e := m.findRequiredNoun(domain, otherNoun, nounSansKind); e != nil {
		err = e
	} else if card, e := m.findCardinality(rel); e != nil {
		err = e
	} else {
		var reverse, multi bool
		switch card {
		case tables.ONE_TO_ONE:
			// sort the names so that the left column is always less than the second
			// simplifies testing of the conflicts for one-to-one
			reverse = true
			if one.name > other.name {
				one, other = other, one
			}
		case tables.ONE_TO_MANY:
			// for a given rhs, there can be only one lhs
			reverse = false

		case tables.MANY_TO_ONE:
			// for a given lhs, there can be only one rhs
			reverse = true

		case tables.MANY_TO_MANY:
			multi = true

		default:
			// well, it should have been one of those.
			err = errutil.Fmt("invalid cardinality %q for %q", card, rel.name)
		}
		if err == nil {
			if e := m.checkPair(domain, rel, one, other, reverse, multi); e != nil {
				err = e
			} else {
				err = m.addPair(domain, rel, one, other, at)
			}
		}

	}
	return
}

// a field used for patterns as a calling parameter
func (m *Modeler) Parameter(domain, kind, field string, aff affine.Affinity, cls, at string) (err error) {
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
		update mdl_pat
		set labels = case when labels is null then (?2) else (labels ||','|| ?2) end
		where kind = ?1 and result is null
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

// a plural word ("many") can have at most one singular definition per domain
// ie. "people" and "persons" are valid plurals of "person",
// but "people" as a singular can only be defined as "person" ( not also "human" )
var mdl_plural = tables.Insert("mdl_plural", "domain", "many", "one", "at")

// a plural word (many) can have at most one singular definition per domain
// ie. people and persons are valid plurals of person,
// but people as a singular can only be defined as person ( not also human ) error
func (m *Modeler) Plural(domain, many, one, at string) (err error) {
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
			_, err = m.db.Exec(mdl_plural, d, many, one, at)
		}
	}
	return
}

//	fix? the data is duplicated in kinds and fields... should this be removed?
//
// might also consider adding a "cardinality" field to the relation kind, and then use init for individual relations
var mdl_rel = tables.Insert("mdl_rel", "relKind", "oneKind", "otherKind", "cardinality", "at")

// relation and constraint between two kinds of nouns

// relation and constraint between two kinds of nouns
//
//	fix? the data is duplicated in kinds and fields... should this be removed?
//
// might also consider adding a cardinality field to the relation kind, and then use init for individual relations
func (m *Modeler) Rel(domain, relKind, oneKind, otherKind, cardinality, at string) (err error) {
	if rel, e := m.findRequiredKind(domain, relKind); e != nil {
		err = e
	} else if rel.domain != domain {
		err = errutil.New("relation signature expected in the same domain as relation declaration")
	} else if one, e := m.findRequiredKind(domain, oneKind); e != nil {
		err = e
	} else if other, e := m.findRequiredKind(domain, otherKind); e != nil {
		err = e
	} else if _, e := m.db.Exec(mdl_rel, rel.id, one.id, other.id, cardinality, at); e != nil {
		err = e // improve the error result if the relation existed vefore?
	} else {
		a, b := makeRel(oneKind, otherKind, cardinality)
		if e := m.addField(domain, rel, one, a.lhs(), a.affinity(), at); e != nil {
			err = e
		} else if e := m.addField(domain, rel, other, b.rhs(), b.affinity(), at); e != nil {
			err = e
		}
	}
	return
}

// a field used for patterns as a returned value
func (m *Modeler) Result(domain, kind, field string, aff affine.Affinity, cls, at string) (err error) {
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
		update mdl_pat
		set result=?2
		where kind = ?1 and result is null
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

var mdl_rule = tables.Insert("mdl_rule", "domain", "kind", "target", "phase", "filter", "prog", "at")

func (m *Modeler) Rule(domain, pattern, target string, phase int, filter rt.BoolEval, prog []rt.Execute, at string) (err error) {
	if filter, e := marshalout(filter); e != nil {
		err = e
	} else if prog, e := marshalprog(prog); e != nil {
		err = e
	} else {
		err = m.UnmarshaledRule(domain, pattern, target, phase, filter, prog, at)
	}
	return
}

// public for tests:
func (m *Modeler) UnmarshaledRule(domain, pattern, target string, phase int, filter, prog, at string) (err error) {
	if kid, e := m.findRequiredKind(domain, pattern); e != nil {
		err = e
	} else if !strings.HasSuffix(kid.fullpath(), m.patternPath) {
		err = errutil.Fmt("kind %q in domain %q is not a pattern", pattern, domain)
	} else if tgt, e := m.findOptionalKind(domain, target); e != nil {
		err = e
	} else {
		_, err = m.db.Exec(mdl_rule, domain, kid.id, tgt.id, phase, filter, prog, at)
	}
	return
}

// note: values are written per noun, not per domain
// fix? some values are references to objects in the form "#domain::noun" -- should the be changed to ids?
var mdl_value = tables.Insert("mdl_value", "noun", "field", "value", "at")

// public for tests:
func (m *Modeler) UnmarshaledValue(domain, noun, field, value, at string) (err error) {
	if noun, e := m.findRequiredNoun(domain, noun, nounWithKind); e != nil {
		err = e
	} else {
		var fieldId int
		if e := m.db.QueryRow(`
		select mf.rowid
		from mdl_field mf
		join mdl_kind mk 
			on(mf.kind = mk.rowid)
		where field = ?2
		and instr(?1, ','||mk.rowid||',')`,
			noun.fullpath, field).Scan(&fieldId); e == sql.ErrNoRows {
			err = errutil.Fmt("%w field %q in noun %q in domain %q", Missing, field, noun.name, noun.domain)
		} else if e != nil {
			err = errutil.New("database error", e)
		} else {
			_, err = m.db.Exec(mdl_value, noun.id, fieldId, value, at)
		}
	}
	return
}

// the noun half of what was Start.
// domain, noun, field reference a join of Noun and Kind to get a filtered Field.
// FIX: nouns should be able to store EVALS too
// example: an object with a counter in its description.
func (m *Modeler) Value(domain, noun, field string, value literal.LiteralValue, at string) (err error) {
	if value, e := marshalout(value); e != nil {
		err = e
	} else {
		err = m.UnmarshaledValue(domain, noun, field, value, at)
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

func marshalprog(prog []rt.Execute) (ret string, err error) {
	if len(prog) > 0 {
		slice := rt.Execute_Slice(prog)
		if out, e := marshalout(&slice); e != nil {
			err = e
		} else {
			ret = out
		}
	}
	return
}
