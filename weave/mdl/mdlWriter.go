package mdl

import (
	"database/sql"
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

type Pen struct {
	db         *tables.Cache
	paths      *paths
	domain, at string
	warn       Log
}

type Log func(fmt string, parts ...any)

func eatDuplicates(l Log, e error) (err error) {
	if e == nil || !errors.Is(e, Duplicate) {
		err = e
	} else if l != nil {
		l(e.Error())
	}
	return
}

// tbd: perhaps writing the aspect to its own table would be best
// join at runtime to synthesize fields; would fix the questions of adding bad traits ( see comments )
// ( could potentially write both as a bridge )
func (pen *Pen) AddAspect(aspect string, traits []string) (err error) {
	_, err = pen.addAspect(aspect, traits)
	return
}

func (pen *Pen) addAspect(aspect string, traits []string) (ret kindInfo, err error) {
	domain, at := pen.domain, pen.at
	var existingTraits int
	if kid, e := pen.addKind(aspect, kindsOf.Aspect.String()); e != nil {
		err = e // ^ hrm.
	} else if strings.Count(kid.fullpath(), ",") != 3 {
		// tbd: could loosen this; for now it simplifies writing the aspects;
		// no need to check for conflicting fields if there's no derivation
		// doesn't stop someone from adding derivation later though ...
		err = errutil.Fmt("can't create aspect of %q; kinds of aspects can't be inherited", aspect)
	} else if e := pen.db.QueryRow(`
			select count(*) 
			from mdl_field mf 
			where mf.kind = ?1
			`, kid.id).Scan(&existingTraits); e != nil {
		err = errutil.New("database error", e)
	} else if existingTraits > 0 {
		// fix? doesn't stop someone from adding new traits later though....
		// field builder could check that it only builds kindsOf.Kind
		err = errutil.Fmt("aspect %q from %q already has traits", aspect, domain)
	} else if kid.domain != domain {
		err = errutil.Fmt("cant add traits to aspect %q; traits are expected to exist in the same domain as the aspect. was %q now %q",
			aspect, kid.domain, domain)
	} else {
		for _, t := range traits {
			if _, e := pen.db.Exec(mdl_field, domain, kid.id, t, affine.Bool, nil, at); e != nil {
				err = errutil.New("database error", e)
				break
			}
		}
		if err == nil {
			ret = kid
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
				} else if out, e := marshalout(value); e != nil {
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

var mdl_domain = tables.Insert("mdl_domain", "domain", "requires", "at")

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
					Conflict, fact, value, prev.domain, value, domain)
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

func (pen *Pen) AddFields(fields Fields) error {
	return fields.writeFields(pen)
}

var mdl_grammar = tables.Insert("mdl_grammar", "domain", "name", "prog", "at")

// player input parsing
func (pen *Pen) AddGrammar(name string, prog *grammar.Directive) (err error) {
	domain, at := pen.domain, pen.at
	if prog, e := marshalout(prog); e != nil {
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

// singular name of kind and materialized hierarchy of ancestors separated by commas
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
					id:     newid,
					name:   name,
					domain: domain,
					path:   path,
					exact:  true,
				}
				// cache result...
				switch name {
				case kindsOf.Aspect.String():
					err = updatePath(res, parent.fullpath(), &pen.paths.aspectPath)
				case kindsOf.Pattern.String():
					err = updatePath(res, parent.fullpath(), &pen.paths.patternPath)
				case kindsOf.Macro.String():
					err = updatePath(res, parent.fullpath(), &pen.paths.macroPath)
				default:
					// super hacky..... hmmm...
					// if we've declared a new kind of a pattern:
					// write blanks into the mdl_pat; parameters and results use update only.
					if strings.HasSuffix(parent.fullpath(), pen.paths.patternPath) {
						_, err = pen.db.Exec(`insert into mdl_pat(kind) values(?1)`, newid)
					}
				}
			}
		}
	}
	return
}

func (pen *Pen) addAncestor(kind, parent kindInfo) (err error) {
	name := kind.name
	domain := pen.domain
	if !kind.exact && parent.numAncestors() < 2 {
		// we allow plural named kinds for nouns, etc. not for patterns and built in kinds.
		err = errutil.Fmt("%w ambiguously named kinds: %q (in domain %q) and %q (in %q)",
			Conflict, name, domain, kind.name, kind.domain)
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

func (pen *Pen) AddNoun(short, long, kind string) (err error) {
	if n, e := pen.addNoun(short, kind); e != nil {
		err = eatDuplicates(pen.warn, e)
	} else if len(long) != 0 {
		parts := genNames(n.name, long)
		for i, name := range parts {
			if e := pen.addName(n, name, i); e != nil {
				err = e
				break
			}
		}
	}
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

// tbd: anyway to remove or improve?
// i especially don't like this is the one dependency on lang.
// and really... it should have explicit tests
// a Noun builder maybe? and let story be the one to genNames
func genNames(short, long string) []string {
	// if the original got transformed into underscores
	// write the original name (ex. "toy boat" vs "toy_boat" )
	var out []string
	if clip := strings.TrimSpace(long); clip != short {
		out = append(out, clip)
	}
	out = append(out, short)

	// generate additional names by splitting the name into parts
	split := lang.Fields(short)
	if cnt := len(split); cnt > 1 {
		// in case the name was reduced due to multiple separators
		if breaks := strings.Join(split, " "); breaks != short {
			out = append(out, breaks)
		}
		// write individual words in increasing rank ( ex. "boat", then "toy" )
		// note: trailing words are considered "stronger"
		// because adjectives in noun names tend to be first ( ie. "toy boat" )
		for i := len(split) - 1; i >= 0; i-- {
			word := split[i]
			out = append(out, word)
		}
	}
	return out
}

func (pen *Pen) AddName(noun, name string, rank int) (err error) {
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

var mdl_opposite = `insert into mdl_rev(domain, oneWord, otherWord, at) 
				values(?1, ?2, ?3, ?4), (?1, ?3, ?2, ?4)`

// domain captures the scope in which the pairing was defined.
// within that scope: the noun, relation, and otherNoun are all unique names --
// even if they are not unique globally, and even if they a broader/different scope than the pair's domain.
func (pen *Pen) AddOpposite(a, b string) (err error) {
	domain, at := pen.domain, pen.at
	if d, e := pen.findDomain(); e != nil {
		err = e
	} else if rows, e := pen.db.Query(
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
			err = eatDuplicates(pen.warn, e)
		} else {
			// writes the opposite paring as well
			_, err = pen.db.Exec(mdl_opposite, d, a, b, at)
		}
	}
	return
}

// domain captures the scope in which the pairing was defined.
// within that scope: the noun, relation, and otherNoun are all unique names --
// even if they are not unique globally, and even if they a broader/different scope than the pair's domain.
var mdl_pair = tables.Insert("mdl_pair", "domain", "relKind", "oneNoun", "otherNoun", "at")

// currently assumes exact noun names
func (pen *Pen) AddPair(rel, oneNoun, otherNoun string) (err error) {
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
		err = errutil.Fmt("%w in pattern %q domain %q", e, pat.name, pen.domain)
	}
	return
}

func (pen *Pen) ExtendPattern(pat Pattern) (err error) {
	if pat.parent != kindsOf.Pattern.String() {
		err = errutil.Fmt("extend pattern %q didn't expect a newly defined parent %q", pat.name, pat.parent)
	} else if e := pat.writePattern(pen, false); e != nil {
		err = errutil.Fmt("%w in pattern %q domain %q", e, pat.name, pen.domain)
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

var mdl_phrase = tables.Insert("mdl_phrase", "domain", "macro", "phrase", "reversed", "at")

// author text parsing
func (pen *Pen) AddPhrase(macro, phrase string, reversed bool) (err error) {
	domain, at := pen.domain, pen.at
	if kind, e := pen.findRequiredKind(macro); e != nil {
		err = e
	} else if isMacro := strings.HasSuffix(kind.fullpath(), pen.paths.macroPath); !isMacro {
		err = errutil.Fmt("kind %q in domain %q is not a macro", macro, domain)
	} else {
		// search for conflicting phrases within this domain.
		var prev struct {
			domain   string
			kind     string
			reversed bool
		}
		e := pen.db.QueryRow(
			`select mg.domain
			from mdl_phrase mg
			join domain_tree dt
				on (dt.uses = mg.domain)
			where base = ?1
			and phrase = ?2
		`, domain, phrase).Scan(&prev.domain, &prev.kind, &prev.reversed)
		switch e {
		case sql.ErrNoRows:
			_, err = pen.db.Exec(mdl_phrase, domain, kind.id, phrase, reversed, at)

		case nil:
			if prev.kind != kind.name || prev.reversed != reversed {
				err = errutil.Fmt("%w phrase %q meant %s in domain %q and now %s in domain %q",
					Conflict, phrase, fmtMacro(prev.kind, prev.reversed), prev.domain,
					fmtMacro(macro, reversed), domain)
			} else {
				pen.warn("%w phrase %q already declared in domain %q and now domain %q",
					Duplicate, phrase, prev.domain, domain)
			}
		default:
			err = errutil.New("database error", e)
		}
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
		// sneaky: if a result duplicates an existing field ( ie. a parameter )
		// no problem: we return that parameter.
		// tbd: possibly itd be better to flag this as a conflict;
		// noting that "collate groups" currently relies on sharing
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
			err = errutil.Fmt("unexpected result %q for kind %q in domain %q",
				field, kid.name, pen.domain)
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

// the top level fields of nouns can hold runtime evaluated assignments.
// note: assumes noun is an exact name
func (pen *Pen) AddFieldValue(noun, field string, value assign.Assignment) (err error) {
	if strings.IndexRune(field, '.') >= 0 {
		err = errutil.Fmt("unexpected dot in assigned value for noun %q field %q", noun, field)
	} else {
		err = pen.addFieldValue(noun, field, value)
	}
	return
}

// store a literal value somewhere within a record held by a noun.
// note: assumes noun is an exact name
func (pen *Pen) AddPathValue(noun, path string, value literal.LiteralValue) (err error) {
	if parts := strings.Split(path, "."); len(parts) == 1 {
		err = pen.addFieldValue(noun, path, assign.Literal(value))
	} else {
		err = pen.addPathValue(noun, parts, value)
	}
	return
}

type ProvisionalAssignment struct {
	assign.Assignment
}
type ProvisionalLiteral struct {
	literal.LiteralValue
}

func isProvisional(a any) (okay bool) {
	if _, ok := a.(ProvisionalAssignment); ok {
		okay = true
	} else if _, ok := a.(ProvisionalLiteral); ok {
		okay = true
	}
	return
}
