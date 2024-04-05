package mdl

import (
	"database/sql"
	"slices"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

// tbd: perhaps writing the aspect to its own table would be best
// join at runtime to synthesize fields; would fix the questions of adding bad traits ( see comments )
// ( could potentially write both as a bridge )

func (pen *Pen) AddAspectTraits(aspect string, traits []string) (err error) {
	if kid, e := pen.findRequiredKind(aspect); e != nil {
		err = e // ^ hrm.
	} else if isAspect := strings.HasSuffix(kid.fullpath(), pen.getPath(kindsOf.Aspect)); !isAspect {
		err = errutil.Fmt("kind %q in domain %q is not an aspect", aspect, pen.domain)
	} else if strings.Count(kid.fullpath(), ",") != 3 {
		// tbd: could loosen this; for now it simplifies writing the aspects;
		// no need to check for conflicting fields if there's no derivation
		// doesn't stop someone from adding derivation later though ...
		err = errutil.Fmt("can't create aspect of %q; kinds of aspects can't be inherited", aspect)
	} else {
		err = pen.addTraits(kid, traits)
	}
	return
}

func (pen *Pen) addAspect(aspect string, traits []string) (ret kindInfo, err error) {
	if cls, e := pen.addKind(aspect, kindsOf.Aspect.String()); e != nil {
		err = e
	} else if e := pen.addTraits(cls, traits); e != nil {
		err = e
	} else {
		ret = cls
	}
	return
}

func (pen *Pen) addTraits(kid kindInfo, traits []string) (err error) {
	domain, at := pen.domain, pen.at
	if existingTraits, e := tables.QueryStrings(pen.db, `
			select mf.field	
			from mdl_field mf 
			where mf.kind = ?1
			order by mf.rowid`, kid.id); e != nil {
		err = errutil.New("database error", e)
	} else if len(existingTraits) > 0 {
		// fix? doesn't stop someone from adding new traits later though....
		// field builder could check that it only builds kindsOf.Kind
		if slices.Compare(traits, existingTraits) != 0 {
			err = errutil.Fmt("aspect %q from %q already has traits", kid.name, domain)
		}
	} else if kid.domain != domain {
		err = errutil.Fmt("cant add traits to aspect %q; traits are expected to exist in the same domain as the aspect. was %q now %q",
			kid.name, kid.domain, domain)
	} else {
		for _, t := range traits {
			if _, e := pen.db.Exec(mdl_field, domain, kid.id, t, affine.Bool, nil, at); e != nil {
				err = errutil.New("database error", e)
				break
			}
		}
	}
	return
}

var mdl_check = tables.Insert("mdl_check", "domain", "name", "value", "affinity", "prog", "at")

// author tests of stories
func (pen *Pen) AddCheck(name string, value literal.LiteralValue, prog []rt.Execute) (err error) {
	domain, at := pen.domain, pen.at
	if d, e := pen.findDomain(); e != nil {
		err = e
	} else {
		var prev struct {
			id          int64
			domain      string
			prog, value bool
		}
		if e := pen.db.QueryRow(
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
				} else if out, e := marshalLiteral(value); e != nil {
					err = e
				} else {
					var aff affine.Affinity
					if value != nil {
						aff = literal.GetAffinity(value)
					}
					if prev.id == 0 {
						_, err = pen.db.Exec(mdl_check, d, name, out, aff, prog, at)
					} else {
						if res, e := pen.db.Exec(`update mdl_check 
						set prog = (case when ?2 then ?2 else prog end),
						set value = (case when ?3 then ?3 else value end),
						where rowid = ?1`,
							prev.id, prev.prog, prev.value); e != nil {
							err = e
						} else if cnt, e := res.RowsAffected(); cnt != 1 {
							err = e
						} else {
							err = errutil.Fmt("unexpected error updating check %q; %d rows affected.",
								name, cnt)
						}
					}
				}
			}
		}
	}
	return
}

// when elements are missing sometimes the same domain requirement pair gets inserted twice
// fix? for now this ignores the duplicate values.
var mdl_domain = tables.InsertWith("mdl_domain", "on conflict do nothing", "domain", "requires", "at")

// pairs of domain name and (domain) dependencies
// fix: are we forcing/checking parent domains to exist before writing?
// that might be cool .... but maybe this is the wrong level?
func (pen *Pen) AddDependency(requires string) (err error) {
	domain, at := pen.domain, pen.at
	_, err = pen.db.Exec(mdl_domain, domain, requires, at)
	return
}

var mdl_fact = tables.Insert("mdl_fact", "domain", "fact", "value", "at")

func makeKeyValue(key string, partsAndValue []string) (k, v string, err error) {
	if end := len(partsAndValue) - 1; end < 0 {
		err = errutil.New("invalid fact", key)
	} else if end == 0 {
		k, v = key, partsAndValue[0]
	} else {
		var b strings.Builder
		b.WriteString(key)
		for i := 0; i < end; i++ {
			part := partsAndValue[i]
			b.WriteRune('/')
			b.WriteString(part)
		}
		k, v = b.String(), partsAndValue[end]
	}
	return
}

// arbitrary key-value storage
// returns true if its a new fact, false otherwise or on error.
func (pen *Pen) AddFact(key string, partsAndValue ...string) (okay bool, err error) {
	if fact, value, e := makeKeyValue(key, partsAndValue); e != nil {
		err = e
	} else if domain, e := pen.findDomain(); e != nil {
		err = e
	} else {
		var prev struct {
			domain, value string
		}
		q := pen.db.QueryRow(`
		select mx.domain, mx.value
		from mdl_fact mx
		join domain_tree
			on (uses = domain)
		where base = ?1
		and fact = ?2`, domain, fact)
		switch e := q.Scan(&prev.domain, &prev.value); e {
		case sql.ErrNoRows:
			if _, e := pen.db.Exec(mdl_fact, domain, fact, value, pen.at); e != nil {
				err = errutil.New("database error", e)
			} else {
				okay = true
			}
		case nil:
			if prev.value != value {
				err = errutil.Fmt("%w fact %q was %q in domain %q and now %q in domain %q",
					Conflict, fact, prev.value, prev.domain, value, domain)
			} else {
				// do we want to warn about duplicate facts? isnt that kind of the point of them?
				// maybe eat at the weave level?
				pen.warn("%w fact %q already declared in domain %q and now domain %q",
					Duplicate, fact, prev.domain, domain)
			}
		default:
			err = errutil.New("database error", e)
		}
	}
	return
}

func (pen *Pen) AddKindFields(kind string, fields []FieldInfo) error {
	return pen.writeFields(kind, fields)
}

var mdl_grammar = tables.Insert("mdl_grammar", "domain", "name", "prog", "at")

// player input parsing
func (pen *Pen) AddGrammar(name string, prog *grammar.Directive) (err error) {
	domain, at := pen.domain, pen.at
	if prog, e := marshal(prog); e != nil {
		err = e
	} else if d, e := pen.findDomain(); e != nil {
		err = e
	} else {
		var prev struct {
			domain, prog string
		}
		e := pen.db.QueryRow(
			`select mg.domain, mg.prog
			from mdl_grammar mg
			join domain_tree dt
				on (dt.uses = mg.domain)
			where base = ?1
			and name = ?2
		`, domain, name).Scan(&prev.domain, &prev.prog)
		switch e {
		case sql.ErrNoRows:
			_, err = pen.db.Exec(mdl_grammar, d, name, prog, at)

		case nil:
			if prev.prog != prog {
				err = errutil.Fmt("%w grammar %q was %q in domain %q and now %q in domain %q",
					Conflict, name, prev.prog, prev.domain, prog, domain)
			} else {
				pen.warn("%w grammar %q already declared in domain %q and now domain %q",
					Duplicate, name, prev.domain, domain)
			}

		default:
			err = errutil.New("database error", e)
		}
	}
	return
}

var mdl_kind = tables.Insert("mdl_kind", "domain", "kind", "singular", "path", "at")

// plural name of kind and materialized hierarchy of ancestors separated by commas
// this (somewhat) duplicates the algorithm used by Noun()
func (pen *Pen) AddKind(name, parent string) (err error) {
	_, err = pen.addKind(name, parent)
	return
}

func (pen *Pen) addKind(name, parent string) (ret kindInfo, err error) {
	domain, at := pen.domain, pen.at
	if parent, e := pen.findOptionalKind(parent); e != nil {
		err = e
	} else if len(parent.name) > 0 && parent.id == 0 {
		// a specified ancestor doesn't exist.
		err = errutil.Fmt("%w ancestor %q", Missing, parent.name)
	} else if kind, e := pen.findKind(name); e != nil {
		err = e
	} else if kind.id != 0 {
		if parent.id == 0 {
			ret = kind
		} else if e := pen.addAncestor(kind, parent); e != nil {
			err = e
		} else {
			ret = kind
		}
	} else {
		// manage singular and plural kinds
		// i don't like this much; especially be cause it depends so much on the first declaration
		// maybe better would be a name/names table that any named concept can use.
		// or just force everyone to use the right names.
		if singular, e := pen.singularize(name); e != nil {
			err = e
		} else {
			var optionalOne *string
			if singular != name {
				optionalOne = &singular
			}
			// easiest is if the name has never been mentioned before;
			// we verified the other inputs already, so insert:
			path := parent.fullpath()
			if res, e := pen.db.Exec(mdl_kind, domain, name, optionalOne, trimPath(path), at); e != nil {
				err = errutil.New("database error", e)
			} else if newid, e := res.LastInsertId(); e != nil {
				err = e
			} else {
				ret = kindInfo{
					id:           newid,
					name:         name,
					domain:       domain,
					path:         path,
					exactName:    true,
					newlyCreated: true,
				}
				// hacky: cache result...
				if k := kindsOf.FindDefaultKind(name); k != 0 {
					if path, e := updatePath(res, parent.fullpath()); e != nil {
						err = e
					} else {
						pen.paths[k] = pathEntry{path: path}
					}
				}
				// super hacky....
				// if we've declared a new kind of a pattern:
				// write blanks into the mdl_pat; parameters and results use update only.
				// would be better in "createPattern" but some tests ( TestQueries )
				// create fake patterns via AddKind :/
				if strings.HasSuffix(parent.fullpath(), pen.getPath(kindsOf.Pattern)) {
					_, err = pen.db.Exec(`insert into mdl_pat(kind) values(?1)`, newid)
				}
			}
		}
	}
	return
}

// hacky: if we've declared a new kind of a pattern:
// write blanks into the mdl_pat; parameters and results use update only.
func (pen *Pen) createPattern(name, parent string) (ret kindInfo, err error) {
	return pen.addKind(name, parent)
}

func (pen *Pen) addAncestor(kind, parent kindInfo) (err error) {
	name := kind.name
	domain := pen.domain
	if !kind.exactName && parent.numAncestors() < 2 {
		// we only allow plural named kinds for nouns ( kinds of kind )
		// see notes in jessAspects.go
		err = errutil.Fmt("%w plural singular conflict for %q (in %q)",
			Conflict, name, domain)
	} else if strings.HasSuffix(parent.fullpath(), kind.fullpath()) {
		err = errutil.Fmt("%w circular reference detected %q already declared as an ancestor of %q.",
			Conflict, parent.name, name)
	} else if strings.HasSuffix(kind.path, parent.fullpath()) {
		// did the existing path fully contain the new ancestor?
		// then its a duplicate request (ex. `,c,b,a,` `,b,a,` )
		pen.warn("%w %q already declared as an ancestor of %q.",
			Duplicate, parent.name, name)
	} else if strings.HasSuffix(parent.fullpath(), kind.path) {
		// is the newly specified ancestor more specific than the existing path?
		// then we are ratcheting down. (ex. `,c,b,a,` `,b,a,` )
		if kind.domain != domain {
			// if it was declared in a different domain: we can't change it now.
			err = errutil.Fmt("%w can't redefine the ancestor of %q as %q; the domains differ: was %q, now %q.",
				Conflict, name, parent.name, kind.domain, domain)
		} else if res, e := pen.db.Exec(`update mdl_kind set path = ?2 where rowid = ?1`,
			kind.id, trimPath(parent.fullpath())); e != nil {
			err = e
		} else if cnt, e := res.RowsAffected(); cnt != 1 {
			err = errutil.Fmt("unexpected error updating hierarchy of %q; %d rows affected.",
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
	return
}

var mdl_name = tables.Insert("mdl_name", "domain", "noun", "name", "rank", "at")

// the domain tells the scope in which the noun was defined
// ( the same as - or a child of - the domain of the kind )
var mdl_noun = tables.Insert("mdl_noun", "domain", "noun", "kind", "at")

// add a noun with the passed identifier and kind ( both normalized ) the kind must exist.
// note: returns mdl.Duplicate if the noun is already know.
// see also: the utility function AddNamedNoun()
func (pen *Pen) AddNounKind(noun, kind string) (err error) {
	_, err = pen.addNoun(noun, kind)
	return
}

// the domain tells the scope in which the noun was defined
// ( the same as - or a child of - the domain of the kind ) error
// this duplicates the algorithm used by Kind()
func (pen *Pen) addNoun(name, ancestor string) (ret nounInfo, err error) {
	domain, at := pen.domain, pen.at
	if parent, e := pen.findRequiredKind(ancestor); e != nil {
		err = e
	} else if prev, e := pen.findNoun(name, nounWithKind); e != nil {
		err = e
	} else if prev.id == 0 {
		// easiest is if the name has never been mentioned before;
		// we verified the other inputs already, so just go ahead and insert:
		if res, e := pen.db.Exec(mdl_noun, domain, name, parent.id, at); e != nil {
			err = errutil.New("database error", e)
		} else if newid, e := res.LastInsertId(); e != nil {
			err = e
		} else {
			ret = nounInfo{
				id:       newid,
				name:     name,
				domain:   domain,
				kid:      parent.id,
				kind:     parent.name,
				fullpath: parent.fullpath(),
			}
		}
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
			if res, e := pen.db.Exec(`update mdl_noun set kind = ?2 where rowid = ?1`,
				prev.id, parent.id); e != nil {
				err = e
			} else if cnt, e := res.RowsAffected(); cnt != 1 {
				err = errutil.Fmt("unexpected error updating hierarchy of %q; %d rows affected",
					name, cnt)
			} else if e != nil {
				err = e
			} else {
				ret = prev
			}
		}
	}
	return
}

// currently negative ranks are runtime aliases,
// and positive values are weave time.
func (pen *Pen) AddNounName(noun, name string, rank int) (err error) {
	if n, e := pen.findRequiredNoun(noun, nounSansKind); e != nil {
		err = e
	} else {
		err = pen.addName(n, name, rank)
	}
	return
}

// words for authors and game players refer to nouns
// follows the domain rules of Noun.
func (pen *Pen) addName(noun nounInfo, name string, rank int) (err error) {
	var exists bool
	if e := pen.db.QueryRow(`
	select 1
	from mdl_name mn
	join domain_tree
		on (uses = domain)
	where base = ?1
	and noun = ?2
	and name = ?3`, pen.domain, noun.id, name).Scan(&exists); e != nil && e != sql.ErrNoRows {
		err = errutil.New("database error", e)
	} else if exists {
		// tbd: silence duplicates?
		// since these are generated, there's probably very little the user could do about them.
		pen.warn("%w %q already an alias of %q", Duplicate, name, noun.name)
	} else if _, e := pen.db.Exec(mdl_name, pen.domain, noun.id, name, rank, pen.at); e != nil {
		err = errutil.New("database error", e)
	}
	return
}

// domain captures the scope in which the pairing was defined.
// within that scope: the noun, relation, and otherNoun are all unique names --
// even if they are not unique globally, and even if they a broader/different scope than the pair's domain.
var mdl_pair = tables.Insert("mdl_pair", "domain", "relKind", "oneNoun", "otherNoun", "at")

// currently assumes exact noun names
func (pen *Pen) AddNounPair(rel, oneNoun, otherNoun string) (err error) {
	if rel, e := pen.findRequiredKind(rel); e != nil {
		err = e
	} else if one, e := pen.findRequiredNoun(oneNoun, nounSansKind); e != nil {
		err = e
	} else if other, e := pen.findRequiredNoun(otherNoun, nounSansKind); e != nil {
		err = e
	} else if card, e := pen.findCardinality(rel); e != nil {
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
			if e := pen.checkPair(rel, one, other, reverse, multi); e != nil {
				err = eatDuplicates(pen.warn, e)
			} else {
				_, err = pen.db.Exec(mdl_pair, pen.domain, rel.id, one.id, other.id, pen.at)
			}
		}
	}
	return
}

func (pen *Pen) findCardinality(kind kindInfo) (ret string, err error) {
	if e := pen.db.QueryRow(`
	select cardinality
	from mdl_rel
	where relKind = ?1 
	limit 1
	`, kind.id).Scan(&ret); e == sql.ErrNoRows {
		err = errutil.Fmt("unknown or invalid cardinality for %q in %q", kind.name, kind.domain)
	} else {
		err = e
	}
	return
}

func (pen *Pen) AddPattern(pat Pattern) (err error) {
	if e := pat.writePattern(pen, true); e != nil {
		err = errutil.Fmt("%w adding pattern %q domain %q", e, pat.name, pen.domain)
	}
	return
}

func (pen *Pen) ExtendPattern(pat Pattern) (err error) {
	if pat.parent != kindsOf.Pattern.String() {
		err = errutil.Fmt("extend pattern %q didn't expect a newly defined parent %q", pat.name, pat.parent)
	} else if e := pat.writePattern(pen, false); e != nil {
		err = errutil.Fmt("%w extending pattern %q domain %q", e, pat.name, pen.domain)
	}
	return
}

// a field used for patterns as a calling parameter
func (pen *Pen) addParameter(kid, cls kindInfo, field string, aff affine.Affinity) (err error) {
	domain := pen.domain
	if kid.domain != domain {
		err = errutil.Fmt("%w new parameter %q of %q expected in the same domain as the original declaration; was %q now %q",
			Conflict, field, kid.name, kid.domain, domain)
	} else if e := pen.addField(kid, cls, field, aff); e != nil {
		err = eatDuplicates(pen.warn, e)
	} else if res, e := pen.db.Exec(`
		update mdl_pat
		set labels = case when labels is null then (?2) else (labels ||','|| ?2) end
		where kind = ?1 and result is null
		`, kid.id, field); e != nil {
		err = e
	} else if rows, e := res.RowsAffected(); e != nil {
		err = e
	} else if rows == 0 {
		// can happen if the result was already written.
		err = errutil.Fmt("pattern parameters should be written before results")
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
func (pen *Pen) AddPlural(many, one string) (err error) {
	if domain, e := pen.findDomain(); e != nil {
		err = e
	} else if rows, e := pen.db.Query(
		`select one, domain 
			from mdl_plural
			join domain_tree
				on(uses=domain)
			where base = ?1
			and many = ?2`, domain, many); e != nil {
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
			err = eatDuplicates(pen.warn, e)
		} else {
			_, err = pen.db.Exec(mdl_plural, domain, many, one, pen.at)
		}
	}
	return
}

//	fix? the data is duplicated in kinds and fields... should this be removed?
//
// might also consider adding a "cardinality" field to the relation kind, and then use init for individual relations
var mdl_rel = tables.Insert("mdl_rel", "relKind", "oneKind", "otherKind", "cardinality", "at")

// relation and constraint between two kinds
func (pen *Pen) AddRelation(name, oneKind, otherKind string, amany bool, bmany bool) (err error) {
	if one, e := pen.findRequiredKind(oneKind); e != nil {
		err = e
	} else if other, e := pen.findRequiredKind(otherKind); e != nil {
		err = e
	} else {
		info := relInfo{oneKind, otherKind, makeCard(amany, bmany)}
		var prev struct {
			id     int64
			domain string
			relInfo
		}
		if e := pen.db.QueryRow(
			`select rel.rowid, rk.domain, ak.kind, bk.kind, rel.cardinality
			from mdl_rel rel
			join mdl_kind rk
				on (rel.relKind = rk.rowid)
			join domain_tree dt
				on (dt.uses = rk.domain)
			left join mdl_kind ak 
				on (rel.oneKind = ak.rowid)
			left join mdl_kind bk
				on (rel.otherKind = bk.rowid)
			where base = ?1
			and rk.kind = ?2
		`, pen.domain, name).Scan(&prev.id, &prev.domain, &prev.one, &prev.other, &prev.cardinality); e != nil && e != sql.ErrNoRows {
			err = errutil.New("database error", e)
		} else {
			if prev.id != 0 {
				if prev.relInfo != info {
					err = errutil.Fmt("%w relation %q in %q defined as %s, now %s",
						Conflict, name, prev.domain, prev.relInfo, info)
				} else {
					pen.warn("%w relation %q in domain %q", Duplicate, name, pen.domain)
				}
			} else {
				if rel, e := pen.addKind(name, kindsOf.Relation.String()); e != nil {
					err = e
				} else {
					a, b := info.makeRel()
					if e := pen.addField(rel, one, a.lhs(), a.affinity()); e != nil {
						err = e
					} else if e := pen.addField(rel, other, b.rhs(), b.affinity()); e != nil {
						err = e
					} else if _, e := pen.db.Exec(mdl_rel, rel.id, one.id, other.id, info.cardinality, pen.at); e != nil {
						err = e // improve the error result if the relation existed vefore?
					}
				}
			}
		}
	}
	return
}

// a field used for patterns as a returned value
func (pen *Pen) addResult(kid, cls kindInfo, field string, aff affine.Affinity) (err error) {
	if kid.domain != pen.domain {
		err = errutil.Fmt("%w new result %q of %q expected in the same domain as the original declaration; was %q now %q",
			Conflict, field, kid.name, kid.domain, pen.domain)
	} else {
		// sneaky: if a result duplicates a parameter; we mark that parameter as the return.
		e := pen.addField(kid, cls, field, aff)
		if e := eatDuplicates(pen.warn, e); e != nil {
			err = e
		} else if res, e := pen.db.Exec(`
		update mdl_pat
		set result=?2
		where kind = ?1 and result is null
		`, kid.id, field); e != nil {
			err = e
		} else if rows, e := res.RowsAffected(); e != nil {
			err = e
		} else if rows == 0 {
			// was there a previous result stored in the db?
			// because of the pattern precache loop, the shouldnt be any duplicate results
			err = errutil.Fmt("result already exists for kind %q in domain %q", kid.name, pen.domain)
		}
	}
	return
}

var mdl_rule = tables.Insert("mdl_rule", "domain", "kind", "name", "rank", "stop", "jump", "updates", "prog", "at")

func (pen *Pen) addRule(pattern kindInfo, name string, rank int, stop bool, jump int, updates bool, prog string) (err error) {
	// fix name needs to check for conflicts;
	// unique withing domain?
	_, err = pen.db.Exec(mdl_rule,
		pen.domain,
		pattern.id,
		sql.NullString{
			String: name,
			Valid:  len(name) > 0,
		},
		rank,
		stop,
		jump,
		updates,
		prog,
		pen.at)
	return
}

func (pen *Pen) AddKindTrait(kind, trait string) (err error) {
	return pen.AddDefaultValue(kind, trait, &assign.FromBool{
		Value: &literal.BoolValue{Value: true},
	})
}

// the top level fields of kinds can hold runtime evaluated assignments.
func (pen *Pen) AddDefaultValue(kind, field string, value rt.Assignment) (err error) {
	if kind, e := pen.findRequiredKind(kind); e != nil {
		err = e
	} else {
		err = pen.addDefaultValue(kind, field, value)
	}
	return
}

// the top level fields of nouns can hold runtime evaluated assignments.
// wrap with "ProvisionalAssignment" to assign implied values
// that can be overridden by explicit statements.
// note: assumes noun is an exact name
func (pen *Pen) AddNounValue(noun, field string, value rt.Assignment) (err error) {
	if strings.IndexRune(field, '.') >= 0 {
		err = errutil.Fmt("unexpected dot in assigned value for noun %q field %q", noun, field)
	} else {
		err = pen.addFieldValue(noun, field, value)
	}
	return
}

// store a literal value somewhere within a record held by a noun.
// note: assumes noun is an exact name
// fix: merge with AddNounValue; use the bits inside marshalAssignment
// ... and strip assignment down to a literal value
func (pen *Pen) AddNounPath(noun string, path []string, value literal.LiteralValue) (err error) {
	if len(path) == 0 {
		err = errutil.New("can't set an empty path")
	} else if len(path) == 1 {
		err = pen.addFieldValue(noun, path[0], assign.Literal(value))
	} else {
		err = pen.addPathValue(noun, path, value)
	}
	return
}

type ProvisionalAssignment struct {
	rt.Assignment
}
type ProvisionalLiteral struct {
	literal.LiteralValue
}
