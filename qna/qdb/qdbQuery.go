// Package qdb asks specific questions of the play.db for the package qna runtime.
// It relies on the model.sql and run.sql sqlite tables that are written to by package asm.
package qdb

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/tables"
)

// /verify that query none implements every method
var _ query.Query = (*Query)(nil)

// Read various data from the play database.
type Query struct {
	db   *sql.DB
	dec  CommandDecoder
	rand query.Randomizer
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
	nounKind,
	nounName,
	nounValues,
	nounAliases,
	nounsByKind,
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
	activation int         // number of domain activation requests ( to find new domains in run_domain )
	constVals  query.Cache // readonly info cached from the db
}

// implements closer
func (q *Query) Close() {
	q.db.Close()
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
func (q *Query) ActivateDomains(name string) (retEnds, retBegins []string, err error) {
	if name == q.domain {
		// do nothing.
	} else if tx, e := q.db.Begin(); e != nil {
		err = e
	} else {
		d, act := q.createActivator(tx), q.activation+1
		if ends, begins, e := d.run(name, act); e != nil {
			err = nil
		} else {
			q.domain, q.activation = name, act
			retEnds, retBegins = ends, begins
		}
		// rollback on any error.
		if err == nil {
			err = tx.Commit()
		} else if e := tx.Rollback(); e != nil {
			err = fmt.Errorf("%w caused rollback which caused %w", err, e)
		}
		q.constVals.Reset() // fix? focus cache clear to just the domains that became inactive?
	}
	return
}

type activator struct {
	domainChange, domainDelete, newPairsFromDomain *sql.Stmt
	domainDeactivate, domainActivate               *sql.Stmt
}

func (q *Query) createActivator(tx *sql.Tx) activator {
	return activator{
		domainChange:       tx.Stmt(q.domainChange),
		domainDelete:       tx.Stmt(q.domainDelete),
		newPairsFromDomain: tx.Stmt(q.newPairsFromDomain),
		domainDeactivate:   tx.Stmt(q.domainDeactivate),
		domainActivate:     tx.Stmt(q.domainActivate),
	}
}

// tbd: it might be better to deactivate outside when the end events are called.

func (q *activator) run(name string, act int) (retEnds, retBegins []string, err error) {
	if ends, e := scanStrings(q.domainDeactivate, name); e != nil {
		err = e
	} else if e := q.deactive(ends); e != nil {
		err = e
	} else if begins, e := scanStrings(q.domainActivate, name); e != nil {
		err = e
	} else if e := q.activate(act, begins); e != nil {
		err = e
	} else {
		retEnds, retBegins = ends, begins
	}
	return
}

func (q *activator) deactive(domains []string) (err error) {
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

func (q *activator) activate(act int, domains []string) (err error) {
	for i, cnt := 0, len(domains); i < cnt; i++ {
		d := domains[i]
		if _, e := q.domainChange.Exec(d, act); e != nil {
			err = e
			break
		} else if _, e := q.newPairsFromDomain.Exec(d); e != nil {
			// fix? this replaced newPairsFromChanges.
			// it tried to handle multiple domains at once but
			// didnt handle conflicts between multiple newly activating domains.
			// applying one domain at a time solves that problem.
			err = e
			break
		}
	}
	return
}

// returns the ancestor hierarchy, starting with the kind itself.
// empty if the kind doesnt exist, errors on a db error.
// accepts both the plural and singular kind.
func (q *Query) KindOfAncestors(kind string) (ret []string, err error) {
	if ks, e := scanStrings(q.kindOfAncestors, kind); e != nil && !errors.Is(e, sql.ErrNoRows) {
		err = e
	} else if e == nil {
		ret = ks
	}
	return
}

// given a short name, find a noun ( and some useful other context )
func (q *Query) NounInfo(name string) (ret query.NounInfo, err error) {
	if e := q.nounInfo.QueryRow(name).Scan(&ret.Domain, &ret.Noun, &ret.Kind); e != nil && e != sql.ErrNoRows {
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

func (q *Query) NounsByKind(kind string) ([]string, error) {
	return scanStrings(q.nounsByKind, kind)
}

func (q *Query) PluralToSingular(plural string) (string, error) {
	return scanString(q.pluralToSingular, plural)
}

func (q *Query) PluralFromSingular(singular string) (string, error) {
	return scanString(q.pluralFromSingular, singular)
}

// the last value is always the result, blank for execute statements
func (q *Query) PatternLabels(pat string) (ret []string, err error) {
	var labels sql.NullString
	var result sql.NullString
	switch e := q.patternOf.QueryRow(pat).Scan(&labels, &result); e {
	case sql.ErrNoRows:
		// returns blank with no error
	case nil:
		if labels.Valid {
			ret = strings.Split(labels.String, ",")
		}
		ret = append(ret, result.String)
	default:
		err = e
	}
	return
}

// get the rules from the cache, or build them and add them to the cache
func (q *Query) RulesFor(pat string) (ret query.RuleSet, err error) {
	key := query.MakeKey("rules", pat, "")
	if c, e := q.constVals.Ensure(key, func() (any, error) {
		return q.getRules(pat)
	}); e != nil {
		err = e
	} else {
		ret = c.(query.RuleSet)
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
		err = fmt.Errorf("nothing related for %q %q %q %q", q.domain, rel, noun, otherNoun)
	}
	return
}

func NewQueries(db *sql.DB, dec CommandDecoder) (*Query, error) {
	return NewQueryOptions(db, dec, query.RandomizedTime(), true)
}

func NewQueryOptions(db *sql.DB, dec CommandDecoder, rand query.Randomizer, cacheErrors bool) (ret *Query, err error) {
	var ps tables.Prep
	q := &Query{
		db:        db,
		rand:      rand,
		dec:       dec,
		constVals: query.MakeCache(cacheErrors),
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
		// every field exclusive to the passed kind
		// alt: could use a partition to select the final value max(final) over (PARTITION by kind,field)
		// instead this expects some filtering to ignore repeated values
		// alternatively, we could maybe use update / conflict resolution to only store the best final during weave.
		// fix:  store the id of the kind and pass it back in...
		fieldsOf: ps.Prep(db,
			`select mf.field, mf.affinity, ifnull(mt.kind, '') as type, mv.value
from active_kinds ks 
join mdl_kind ma
	-- is Y (is their name) a part of X (our path)
	on instr(',' || ks.path, 
					 ',' || ma.rowid || ',' )
	or (ks.kind = ma.rowid) 
join mdl_field mf
  on (ma.rowid = mf.kind)
-- pull in all values of any matching field
left join mdl_value_kind mv 
	-- this matches all the kinds from other trees
	-- ( ex. doors and supporters both have portability. )
  on (mv.field = mf.rowid)
	-- so filter initializers by the requested kind's fullpath
  and instr(
   ',' || ks.kind || ',' || ks.path, -- full path
   ',' || mv.kind || ',')
-- finally determine the name of the field's type
left join mdl_kind mt 
	on (mt.rowid = mf.type)
where (ks.name = ?1)
-- sort to get fields in definition order
-- ( that's implicitly also kind order: all fields in earlier kinds are defined first )
-- then by the initializer nearest to our requested kind 
-- and, finally, put final values first.
order by mf.rowid, mv.kind desc, mv.final desc`,
		),
		// path is materialized ids so we return multiple values of resolved names
		kindOfAncestors: ps.Prep(db,
			`select mk.kind 
			from active_kinds ks  -- the kinds in domain scope 
			join mdl_kind mk
				-- is Y (is their name) a part of X (our path)
				on instr(',' || ks.path, 
								 ',' || mk.rowid || ',' )
				or (ks.kind == mk.rowid) -- the kind itself ( to get its real name )
			where (ks.name = ?1) or (ks.singular = ?1)
			order by mk.rowid desc`,
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
		nounKind: ps.Prep(db,
			`select mk.kind
			from active_nouns ns
			join mdl_kind mk
				on (mk.rowid = ns.kind)
			where ns.name = ?1
			limit 1`,
		),
		// does a noun have some specific name?
		nounAliases: ps.Prep(db,
			`select distinct my.name
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
		// query the db for the value(s) of a given field for a given noun
		// fix: future, we will want to save values to a "run_value" table and union those in here.
		nounValues: ps.Prep(db,
			`select coalesce(mv.dot,'') as dot, mv.value
			from mdl_value mv
			join active_nouns ns
				using (noun)
			join mdl_field mf
				on (mf.rowid = mv.field)
			where (ns.name = ?1) and (mf.field = ?2) and (mv.noun = ns.noun)
			order by length(dot), mv.final desc`,
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
		// given the "right side" of some related nouns, return the left side noun(s).
		// for the sake of backwards compat with printing, names are returned in alphabetical order.
		// might be better to move that into an explicit sort in script
		reciprocalOf: ps.Prep(db,
			`select oneName 
			from active_names
			where relName=?1
			and otherName=?2
			order by oneName`,
		),
		// given the "left side" of some related nouns, return the right side noun(s).
		// for the sake of backwards compat with printing, names are returned in alphabetical order.
		// might be better to move that into an explicit sort in script
		relativesOf: ps.Prep(db,
			`select otherName 
			from active_names
			where relName=?1
			and oneName=?2
			order by otherName`,
		),
		// returns the executable rules for a given kind and target
		rulesFor: ps.Prep(db,
			`select coalesce(mu.name, 'rule ' || mu.rowid), mu.stop, mu.jump, mu.updates, mu.prog
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
	}
	if e := ps.Err(); e != nil {
		err = e
	} else {
		ret = q
	}
	return
}
