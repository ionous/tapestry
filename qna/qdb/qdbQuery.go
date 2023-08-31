// Package qdb asks specific questions of the play.db for the package qna runtime.
// It relies on the model.sql and run.sql sqlite tables that are written to by package asm.
package qdb

import (
	"database/sql"
	"strings"

	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

// Read various data from the play database.
type Query struct {
	db *sql.DB
	activeDomains,
	domainActivate,
	domainChange,
	domainDeactivate,
	domainDelete,
	domainScope,
	checks,
	fieldsOf,
	kindOfAncestors,
	nounInfo,
	nounName,
	nounValues,
	nounAliases,
	nounsByKind,
	oppositeOf,
	patternOf,
	newPairsFromDomain,
	newPairsFromNames,
	pluralMatches,
	pluralToSingular,
	pluralFromSingular,
	reciprocalOf,
	relativesOf,
	rulesFor *sql.Stmt
	//
	domain     string
	activation int // number of domain activation requests ( to find new domains in run_domain )
}

func (q *Query) IsDomainActive(name string) (okay bool, err error) {
	// no rows just nothing was found -- that name isnt active
	// fix? would it be worth it to check the name is a valid domain?
	if e := q.domainScope.QueryRow(name).Scan(&okay); e != nil && e != sql.ErrNoRows {
		err = e
	}
	return
}

// changing domains can establish new relations ( abandoning now conflicting ones )
// and cause nouns to fall out of scope
// returns the previous domain name
func (q *Query) ActivateDomain(name string) (ret string, err error) {
	if name == q.domain {
		ret = q.domain
	} else if tx, e := q.db.Begin(); e != nil {
		err = e
	} else if len(name) == 0 {
		was := q.domain
		q.domain = ""
		q.activation += 1
		if e := resetDomain(q.db, true); e != nil {
			err = e
		} else {
			ret = was
		}
	} else {
		act := q.activation + 1
		if de, e := scanStrings(q.domainDeactivate, name); e != nil {
			err = e
		} else if e := q.deactive(de); e != nil {
			err = e
		} else if re, e := scanStrings(q.domainActivate, name); e != nil {
			err = e
		} else if e := q.activate(act, re); e != nil {
			err = e
		}

		// so we can rollback on any error.
		if err == nil {
			q.domain, q.activation, ret = name, act, q.domain
			err = tx.Commit()
		} else if e := tx.Rollback(); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// tbd: what if certain kinds of changes could happen automatically
// while others would need to be in an "on enter/exit" style event handlers
func (q *Query) deactive(domains []string) (err error) {
	for i, cnt := 0, len(domains); i < cnt; i++ {
		d := domains[cnt-i-1] // work backwards from leaf to root
		if _, e := q.domainChange.Exec(d, 0); e != nil {
			err = e
			break
		} else if _, e := q.domainDelete.Exec(d); e != nil {
			err = e
			break
		}
	}
	return
}

// tbd: what if certain kinds of changes could happen automatically
// while others would need to be in an "on enter/exit" style event handlers
func (q *Query) activate(act int, domains []string) (err error) {
	for i, cnt := 0, len(domains); i < cnt; i++ {
		d := domains[i]
		if _, e := q.domainChange.Exec(d, act); e != nil {
			err = e
			break
		} else if _, e := q.newPairsFromDomain.Exec(d); e != nil {
			// fix: this replaced newPairsFromChanges b/c it didnt handle the case where,
			// if multiple domains are being activated at once,
			// that the more derived domain should clear conflicting pairs
			// instead: every listed pair in the all the new domains get set.
			err = e
			break
		}
	}
	return
}

// read all the matching tests from the db.
func (q *Query) ReadChecks(actuallyJustThisOne string) (ret []query.CheckData, err error) {
	if rows, e := q.checks.Query(actuallyJustThisOne); e != nil {
		err = e
	} else {
		var check query.CheckData
		err = tables.ScanAll(rows, func() (err error) {
			ret = append(ret, check)
			return
		}, &check.Name, &check.Domain, &check.Value, &check.Aff, &check.Prog)
	}
	return
}

func (q *Query) FieldsOf(kind string) (ret []query.FieldData, err error) {
	if rows, e := q.fieldsOf.Query(kind); e != nil {
		err = e
	} else {
		var field query.FieldData
		err = tables.ScanAll(rows, func() (err error) {
			ret = append(ret, field)
			return
		}, &field.Name, &field.Affinity, &field.Class, &field.Init)
	}
	return
}

// returns the ancestor hierarchy, not including the kind itself
// ex. for doors that might be: kinds, objects, things, props, openers.
func (q *Query) KindOfAncestors(kind string) ([]string, error) {
	return scanStrings(q.kindOfAncestors, kind)
}

// given a name, find a noun ( and some useful other context )
func (q *Query) NounInfo(name string) (ret query.NounInfo, err error) {
	if e := q.nounInfo.QueryRow(name).Scan(&ret.Domain, &ret.Id, &ret.Kind); e != nil && e != sql.ErrNoRows {
		err = e
	}
	return
}

// return the best "short name" for a noun ( or blank if the noun isnt known or isnt in scope )
func (q *Query) NounName(id string) (ret string, err error) {
	return scanString(q.nounName, id)
}

func (q *Query) NounNames(id string) (ret []string, err error) {
	return scanStrings(q.nounAliases, id)
}

// interpreting the value is left to the caller ( re: field affinity )
// returns pairs of path, (marshaled) value
func (q *Query) NounValues(id, field string) (ret []string, err error) {
	if rows, e := q.nounValues.Query(id, field); e != nil {
		err = e
	} else {
		var path sql.NullString
		var value string
		err = tables.ScanAll(rows, func() (_ error) {
			ret = append(ret, path.String, value)
			return
		}, &path, &value)
	}
	return
}

func (q *Query) NounsByKind(kind string) ([]string, error) {
	return scanStrings(q.nounsByKind, kind)
}

func (q *Query) PluralToSingular(plural string) (string, error) {
	return scanString(q.pluralToSingular, plural)
}

func (q *Query) PluralFromSingular(singular string) (string, error) {
	return scanString(q.pluralFromSingular, singular)
}

func (q *Query) OppositeOf(word string) (string, error) {
	return scanString(q.oppositeOf, word)
}

// the last value is always the result, blank for execute statements
func (q *Query) PatternLabels(pat string) (ret []string, err error) {
	var labels sql.NullString
	var result sql.NullString
	switch e := q.patternOf.QueryRow(pat).Scan(&labels, &result); e {
	case sql.ErrNoRows:
		// returns blank with no error
	case nil:
		parts := strings.Split(labels.String, ",")
		ret = append(parts, result.String)
	default:
		err = e
	}
	return
}

func (q *Query) RulesFor(pat string) (ret []query.RuleData, err error) {
	if rows, e := q.rulesFor.Query(pat); e != nil {
		err = e
	} else {
		var rule query.RuleData
		if e := tables.ScanAll(rows, func() (err error) {
			ret = append(ret, rule)
			return
		}, &rule.Name, &rule.Stop, &rule.Jump, &rule.Updates, &rule.Prog); e != nil {
			err = e // scan error...
		}
	}
	return
}

func (q *Query) ReciprocalsOf(rel, id string) ([]string, error) {
	return scanStrings(q.reciprocalOf, rel, id)
}

func (q *Query) RelativesOf(rel, id string) ([]string, error) {
	return scanStrings(q.relativesOf, rel, id)
}

func (q *Query) Relate(rel, noun, otherNoun string) (err error) {
	if res, e := q.newPairsFromNames.Exec(q.domain, rel, noun, otherNoun); e != nil {
		err = e
	} else if rows, e := res.RowsAffected(); e != nil {
		err = e
	} else if rows == 0 {
		err = errutil.New("nothing related for", q.domain, rel, noun, otherNoun)
	}
	return
}

func resetDomain(db *sql.DB, reset bool) (err error) {
	if reset {
		_, err = db.Exec(`delete from run_domain; delete from run_pair`)
	}
	return
}

func NewQueries(db *sql.DB, reset bool) (ret *Query, err error) {
	if e := tables.CreateRun(db); e != nil {
		err = e
	} else if e := resetDomain(db, reset); e != nil {
		err = e
	} else {
		ret, err = newQueries(db)
	}
	return
}

func newQueries(db *sql.DB) (ret *Query, err error) {
	var ps tables.Prep
	q := &Query{
		db: db,
		activeDomains: ps.Prep(db,
			`select domain 
			from active_domains
			order by domain`,
		),
		checks: ps.Prep(db,
			`select mc.name, mc.domain, mc.value, mc.affinity, mc.prog
			from mdl_check mc
			where mc.name = ?1 or length(?1) == 0
			order by mc.domain, mc.name`,
		),
		// return domains which should be active and are not
		// ( in order of increasing depth. )
		domainActivate: ps.Prep(db,
			`select uses
				from domain_tree
				left join run_domain
					on (domain=uses)
				where base = ?1
				and not coalesce(active, 0)
				order by dist desc`,
		),
		domainChange: ps.Prep(db,
			`insert or replace 
			into run_domain(domain, active) 
			values( ?1, ?2 )`,
		),
		// return domains which are active and should not be
		// ( in order of increasing depth. )
		domainDeactivate: ps.Prep(db,
			`select domain from 
			domain_order
			join run_domain
				using (domain)
			where domain not in (
				select uses
				from domain_tree 
				where base = ?1
			)`,
		),
		domainDelete: ps.Prep(db,
			`delete from run_pair 
			where domain = ?`,
		),
		// determine if the named domain is currently in scope.
		domainScope: ps.Prep(db,
			`select 1 
			from active_domains
			where domain = ?1`,
		),
		// every field of a given kind, and its initial assignment if any.
		// ( interestingly the order by seems to generate a shorter query than without )
		fieldsOf: ps.Prep(db,
			`select mf.field, mf.affinity, ifnull(mk.kind, '') as type, mv.value
			from active_kinds ks  -- search for the kind in question
			join mdl_field mf
				using (kind)
			left join mdl_kind mk  -- search for the kind of type 
				on (mk.rowid = mf.type)
			left join mdl_value mv
				on (mv.field = mf.rowid and mv.noun = 0)
			where ks.name = ?1
			order by mf.rowid`,
		),
		// path is materialized ids so we return multiple values of resolved names
		kindOfAncestors: ps.Prep(db,
			`select mk.kind 
			from active_kinds ks  -- the kinds in domain scope 
			join mdl_kind mk
				-- is Y (is their name) a part of X (our path)
				on instr(',' || ks.path, 
								 ',' || mk.rowid || ',' )
			where ks.name = ?1
			order by mk.rowid`,
		),
		// given a short name, find the noun's fullname.
		// we filter out parser understandings (which have ranks < 0)
		nounInfo: ps.Prep(db,
			`select ns.domain, ns.name, mk.kind
			from active_nouns ns
			join mdl_name my
				using (noun)
			join mdl_kind mk
				on (mk.rowid = ns.kind)
			where rank >= 0 and my.name = ?1
			order by rank
			limit 1`,
		),
		// does a noun have some specific name?
		nounAliases: ps.Prep(db,
			`select my.name
			from mdl_name my
			join active_nouns ns
				using (noun)
			where ns.name=?1
			order by my.name`,
		),
		// given the fullname of a noun, find the best short name
		nounName: ps.Prep(db,
			`select my.name as nameOf
			from active_nouns ns
			join mdl_name my
				using (noun)
			where my.rank >= 0 and ns.name = ?1
			order by rank
			limit 1`,
		),
		// given a named kind, find the nouns
		nounsByKind: ps.Prep(db,
			`select ns.name
			from active_kinds ks
			join active_nouns ns 
				using (kind)
			where ks.name=?1`, // order?
		),
		oppositeOf: ps.Prep(db,
			`select otherWord
			from active_rev
			where oneWord = ?1
			limit 1`,
		),
		// find the names of a given pattern ( kind's ) args and results
		patternOf: ps.Prep(db,
			`select mp.labels, mp.result
			from active_kinds ks 
			join mdl_pat mp
				using (kind)
			where ks.name = ?1
			limit 1`,
		),
		pluralMatches: ps.Prep(db,
			`select domain, one, at 
			from active_plurals
			where many = ?1`,
		),
		pluralToSingular: ps.Prep(db,
			`select one 
			from mdl_plural
			join active_domains
				using (domain)
			where many=?
			limit 1`,
		),
		pluralFromSingular: ps.Prep(db,
			`select many 
			from mdl_plural
			join active_domains
				using (domain)
			where one=?
			limit 1`,
		),
		// maybe unneeded now that domains are activated one by one?
		// newPairsFromChanges: ps.Prep(db,
		// 	newPairsFromChanges,
		// ),
		newPairsFromNames: ps.Prep(db,
			newPairsFromNames,
		),
		newPairsFromDomain: ps.Prep(db,
			newPairsFromDomain,
		),
		// given the "right side" of some related nouns, return the left side noun(s).
		// for the sake of backwards compat with printing, names are returned in alphabetical order.
		// might be better to move that into an explicit sort in script
		reciprocalOf: ps.Prep(db,
			`select oneName 
			from rp_names
			where relName=?1
			and otherName=?2
			order by oneName`,
		),
		// given the "left side" of some related nouns, return the right side noun(s).
		// for the sake of backwards compat with printing, names are returned in alphabetical order.
		// might be better to move that into an explicit sort in script
		relativesOf: ps.Prep(db,
			`select otherName 
			from rp_names
			where relName=?1
			and oneName=?2
			order by otherName`,
		),
		// returns the executable rules for a given kind and target
		rulesFor: ps.Prep(db,
			`select coalesce(mu.name, mu.rowid), mu.stop, mu.jump, mu.updates, mu.prog
			from active_domains 
			join mdl_rule mu
				using (domain)
			join mdl_kind mk 
			  on (mk.rowid = mu.kind) 
			where mk.kind = ?1
			order by 
				mu.rank,
				-- tbd: positive rank items sort first specified to last specified (asc)
				-- zero and negative ranked items sort last specified to first specified  (desc)
				--  mu.rowid * (case when mu.rank > 0 then 1 else -1 end)
				mu.rowid desc
			`,
		),
		// query the db for the value(s) of a given field for a given noun
		// fix: future, we will want to save values to a "run_value" table and union those in here.
		nounValues: ps.Prep(db,
			`select mv.dot, mv.value
			from mdl_value mv
			join active_nouns ns
				using (noun)
			join mdl_field mf
				on (mf.rowid = mv.field)
			where ns.name = ?1 and mf.field = ?2
			order by length(mv.dot)`,
		),
	}
	if e := ps.Err(); e != nil {
		err = e
	} else {
		ret = q
	}
	return
}
