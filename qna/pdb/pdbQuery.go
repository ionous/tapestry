// Package pdb asks specific questions of the play.db for the package qna runtime.
// It relies on the model.sql and run.sql sqlite tables that are written to by package asm.
package pdb

import (
	"database/sql"
	"strings"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

// Read various data from the play database.
type Query struct {
	db *sql.DB
	domainActivation,
	domainScope,
	domainDelete,
	fieldsOf,
	kindOfAncestors,
	nounActive,
	nounInfo,
	nounKind,
	nounName,
	nounValue,
	nounsByKind,
	patternOf,
	pluralToSingular,
	pluralFromSingular,
	reciprocalOf,
	relateChanges,
	relateNames,
	relativesOf,
	rulesFor *sql.Stmt

	//
	domain     string // name of the most recently activated domain ( for scoping new run_pair entries )
	activation int    // number of domain activation requests ( to find new domains in run_domain )
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
func (q *Query) ActivateDomain(name string) (ret string, err error) {
	if name == q.domain {
		ret = q.domain
	} else if tx, e := q.db.Begin(); e != nil {
		err = e
	} else {
		// register the new scope(s)
		act := q.activation + 1
		if _, e := q.domainActivation.Exec(name, act); e != nil {
			err = e
		} else {
			// optional, delete pairs of nouns that fell out of scope
			// alt: join run_pair requests with domain_activate
			// ( would mean that "relate changes" is probably doing more work clearing old relations )
			if _, e := q.domainDelete.Exec(); e != nil {
				err = e
			} else {
				// add pairs that were just activated.
				if _, e := q.relateChanges.Exec(name, act); e != nil {
					err = e
				}
			}
		}
		if err == nil {
			q.domain, q.activation, ret = name, act, q.domain
			err = tx.Commit()
		} else if e := tx.Rollback(); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

type FieldData struct {
	Name     string
	Affinity affine.Affinity
	Class    string
	Init     []byte
}

func (q *Query) FieldsOf(kind string) (ret []FieldData, err error) {
	if rows, e := q.fieldsOf.Query(kind); e != nil {
		err = e
	} else {
		var field FieldData
		err = tables.ScanAll(rows, func() (err error) {
			ret = append(ret, field)
			return
		}, &field.Name, &field.Affinity, &field.Class, &field.Init)
	}
	return
}

// returns the ancestor hierarchy, not including the kind itself
func (q *Query) KindOfAncestors(kind string) ([]string, error) {
	return scanStrings(q.kindOfAncestors, kind)
}

func (q *Query) NounActive(id string) (okay bool, err error) {
	if e := q.nounActive.QueryRow(id).Scan(&okay); e != nil && e != sql.ErrNoRows {
		err = e
	}
	return
}

// given a name, find a noun ( and some useful other context )
func (q *Query) NounInfo(name string) (ret NounInfo, err error) {
	if e := q.nounInfo.QueryRow(name).Scan(&ret.Domain, &ret.Name, &ret.Kind); e != nil && e != sql.ErrNoRows {
		err = e
	}
	return
}

type NounInfo struct {
	Domain, Name, Kind string
}

func (n *NounInfo) IsValid() bool {
	return len(n.Name) != 0
}

// return the noun's kind ( or blank if not known or not in scope )
func (q *Query) NounKind(id string) (ret string, err error) {
	return scanString(q.nounKind, id)
}

// return the best "short name" for a noun ( or blank if the noun isnt known or isnt in scope )
func (q *Query) NounName(id string) (ret string, err error) {
	return scanString(q.nounName, id)
}

//
func (q *Query) NounValue(id, field string) (retAff affine.Affinity, retVal interface{}, err error) {
	if e := q.nounValue.QueryRow(id, field).Scan(&retVal, &retAff); e != nil && e != sql.ErrNoRows {
		err = e
	}
	return
}

func (q *Query) NounsByKind(kind string) ([]string, error) {
	return scanStrings(q.nounsByKind, kind)
}

func (q *Query) PluralToSingular(plural string) (ret string, err error) {
	return scanString(q.pluralToSingular, plural)
}

func (q *Query) PluralFromSingular(singular string) (ret string, err error) {
	return scanString(q.pluralFromSingular, singular)
}

type PatternLabels struct {
	Result string
	Labels []string
}

func (q *Query) PatternLabels(pat string) (ret PatternLabels, err error) {
	var labels, result string
	if e := q.patternOf.QueryRow(pat).Scan(&labels, &result); e != nil && e != sql.ErrNoRows {
		err = e
	} else {
		ret = PatternLabels{result, strings.Split(labels, ",")}
	}
	return
}

type Rules struct {
	Id           string // really an id, but we'll let the driver convert
	Phase        int
	Filter, Prog []byte
}

func (q *Query) RulesFor(pat, target string) (ret []Rules, err error) {
	if rows, e := q.rulesFor.Query(pat, target); e != nil {
		err = e
	} else {
		var rule Rules
		if e := tables.ScanAll(rows, func() (err error) {
			ret = append(ret, rule)
			return
		}, &rule.Id, &rule.Phase, &rule.Filter, &rule.Prog); e != nil {
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
	_, err = q.relateNames.Exec(q.domain, rel, noun, otherNoun)
	return
}

func NewQueries(db *sql.DB) (ret *Query, err error) {
	var ps tables.Prep
	q := &Query{
		db:         db,
		activation: 1,
		domainActivation: ps.Prep(db,
			// build a table of nxn domains indicating "wants" active and "was" active.
			// if activating the first domain activates the second, "want" is set;
			// if the second is *currently* active ( in run_domain ), "was" is set.
			// filter to *changes* for our requested domain and insert the requested value.
			`with shouldExist as (
				select
						  pd.domain as parent, pd.rowid as parentId,
						  cd.domain as child, cd.rowid as childId,
						  -- is Y (the complete path of something else) within X (our own path)
							-- then it's active when we are.
						  instr(',' ||pd.rowid || ',' || pd.path,   ',' ||cd.rowid || ',' || cd.path) and 1 as want,
						  ifnull(rd.active, 0) as was
				from mdl_domain pd
				join mdl_domain cd
				left join run_domain rd
					on (rd.domain=cd.rowid)
			)
			insert or replace into run_domain(domain, active)
			select childId, iif(want, ?2, 0) from shouldExist
			where parent = ?1
			and (want>0) != (was>0)`,
		),
		domainDelete: ps.Prep(db,
			`delete from run_pair
			where domain not in ( 
				select domain from domain_scope 
			)`,
		),
		// determine if the named domain is currently in scope.
		domainScope: ps.Prep(db,
			`select 1 
			from domain_scope ds 
			where name = ?1`,
		),
		// every field of a given kind, and its initial assignment if any.
		// ( interestingly the order by seems to generate a shorter query than without )
		fieldsOf: ps.Prep(db,
			`select mf.field, mf.affinity, ifnull(mf.type, ''), ma.value
			from kind_scope ks
			join mdl_field mf
				using (kind)
			left join mdl_assign ma 
				on (ma.field = mf.rowid)
			where ks.name = ?1
			order by mf.rowid`,
		),
		// path is materialized ids so we return multiple values of resolved names
		kindOfAncestors: ps.Prep(db,
			`select mk.kind 
			from  kind_scope ks
			join mdl_kind mk
				-- is Y (is their name) a part of X (our path)
				on instr(',' || ks.path, 
								 ',' || mk.rowid || ',' )
			where ks.name = ?1`,
		),
		// FIX: what do i even need this for????
		// find all of the kinds of the named kind that are currently in scope.
		// ie. descendants
		// kindOfChildren: ps.Prep(db,
		// 	`select ks.name
		// 		from kind_scope ks
		// 		join kind_scope matching
		// 		on matching.name = ?1
		// 		where instr(',' ||ks.kind || ',' || ks.path, ','|| matching.kind)`,
		// ),
		// determine if the named noun is in scope
		// see also the nounInScope struct
		nounActive: ps.Prep(db,
			`select 1
			from noun_scope ns
			where name = ?1`,
		),
		// given a short name, find the noun's fullname.
		// we filter out parser understandings (which have ranks < 0)
		nounInfo: ps.Prep(db,
			`select ns.domain, ns.name, mk.kind
			from noun_scope ns
			join mdl_name my
				using (noun)
			join mdl_kind mk
				on (mk.rowid = ns.kind)
			where rank >= 0 and my.name = ?1
			order by rank
			limit 1`,
		),
		// select the kind of the given noun
		nounKind: ps.Prep(db,
			`select mk.kind
			from noun_scope ns
			join mdl_kind mk
				on mk.rowid = ns.kind
			where ns.name = ?1
			limit 1`,
		),
		// given the fullname of a noun, find the best short name
		nounName: ps.Prep(db,
			`select my.name as nameOf
			from noun_scope ns
			join mdl_name my
				using (noun)
			where my.rank >= 0 and ns.name = ?1
			order by rank
			limit 1`,
		),
		// given a named kind, find the nouns
		nounsByKind: ps.Prep(db,
			`select ns.name
			from kind_scope ks
			join noun_scope ns 
				using(kind)
			where ks.name=?1`,
		),
		// find the names of a given pattern ( kind's ) args and results
		patternOf: ps.Prep(db,
			`select mp.labels, mp.result
			from kind_scope ks 
			join mdl_pat mp
				using (kind)
			where ks.name = ?1
			limit 1`,
		),
		pluralToSingular: ps.Prep(db,
			`select one 
			from mdl_plural
			join domain_scope
				using(domain)
			where many=?
			limit 1`,
		),
		pluralFromSingular: ps.Prep(db,
			`select many 
			from mdl_plural
			join domain_scope
				using(domain)
			where one=?
			limit 1`,
		),
		// given the "right side" of some related nouns, return the left side noun(s).
		reciprocalOf: ps.Prep(db,
			`select oneName 
			from rp_names
			where relName=?1
			and otherName=?2`,
		),
		relateChanges: ps.Prep(db,
			newPairsFromChanges+relatePair,
		),
		relateNames: ps.Prep(db,
			newPairsFromNames+relatePair,
		),
		// type Relation struct {
		// 	Kind, OtherKind, Cardinality string
		// }
		// create view if not exists
		// rel_scope as
		// select ds.name as domain, mk.rowid as relKind, mk.kind as name, mr.oneKind, mr.otherKind, mr.cardinality
		// from domain_scope ds
		// join mdl_kind mk
		// 	using (domain)
		// join mdl_rel mr
		// 	on (mk.rowid = mr.relKind);
		//
		// find kinds and cardinality for the named relation.
		// relativeKinds: ps.Prep(db,
		// 	`select mk.kind as kind,
		// 		ok.kind as otherKind,
		// 		rs.cardinality
		// 	from rel_scope rs
		// 	join mdl_kind mk
		// 		on (mk.rowid = rs.oneKind)
		// 	join mdl_kind ok
		// 		on (ok.rowid = rs.otherKind)
		// 	where rs.name = ?1`,
		// ),
		// given the "left side" of some related nouns, return the right side noun(s).
		relativesOf: ps.Prep(db,
			`select otherName 
			from rp_names
			where relName=?1
			and oneName=?2`,
		),
		// returns the executable rules for a given kind and target
		rulesFor: ps.Prep(db,
			`select mu.rowid, mu.phase, mu.filter, mu.prog
			from domain_scope 
			join mdl_rule mu
				using (domain)
			join mdl_kind mk 
			  on (mk.rowid = mu.kind) 
			left join mdl_kind mt 
			  on (mt.rowid = mu.target)
			where mk.kind = ?1
			and ifnull(mt.kind,'') = ?2
			order by abs(mu.phase), mu.rowid`,
		),
		// query the db for the value of a given field for a given noun
		// fix: future, we will want to save values to a "run_value" table and union those in here.
		nounValue: ps.Prep(db,
			`select mv.value, mv.affinity
			from mdl_value mv
			join noun_scope ns
				using (noun)
			join mdl_field mf
				on (mf.rowid=mv.field)
			where ns.name=?1 and mf.field=?2
			limit 1`,
		),
	}
	if e := ps.Err(); e != nil {
		err = e
	} else {
		ret = q
	}
	return
}
