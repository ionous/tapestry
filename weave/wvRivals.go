package weave

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

// exposed for testing:
// tbd: maybe this could pull in the newly relevant domains;
// ( ex. use domain = active count )
// currently it happens after the domains have been activated
// and therefore compares everything to everything each time.
// note: fields don't have rivals because they all exist in the same domain as their owner kind.
func findRivals(db tables.Querier, onConflict func(group, domain, key, value, at string) error) (err error) {
	if rows, e := db.Query(`
	with active_grammar as (
		select mg.*
		from mdl_grammar mg 
		join active_domains
		using (domain)
	),

	active_facts as (
		select mx.*
		from mdl_fact mx
		join active_domains
		using (domain)
	),

	select 'fact', a.domain, a.at, a.fact, a.value
		from active_facts as a 
		join active_facts as b 
			using(fact)
		where a.domain != b.domain 
		and a.value != b.value
	union all

	select 'kind', a.domain, a.at, a.kind, ''
		from active_kinds as a 
		join active_kinds as b 
			using(name)
		where a.domain != b.domain
	union all

	select 'grammar', a.domain, a.at, a.name, ''
		from active_grammar as a 
		join active_grammar as b 
			using(name)
		where a.domain != b.domain 
		and a.prog != b.prog
	union all

	select 'phrase', a.domain, a.at, a.phrase, ''
		from active_phrases as a 
		join active_phrases as b 
			using(phrase)
		where a.domain != b.domain
	union all

	select 'opposite', a.domain, a.at, a.oneWord, a.otherWord 
		from active_rev as a 
		join active_rev as b 
			using(oneWord)
		where a.domain != b.domain 
		and a.otherWord != b.otherWord
	union all

	select 'plural', a.domain, a.at, a.many, a.one
		from active_plurals as a 
		join active_plurals as b 
			using(many)
		where a.domain != b.domain
		and a.one != b.one
	`); e != nil {
		err = errutil.New("database error", e)
	} else {
		var group, domain, key, value string
		var at sql.NullString
		if e := tables.ScanAll(rows, func() error {
			return onConflict(group, domain, key, value, at.String)
		}, &group, &domain, &at, &key, &value); e != nil && e != sql.ErrNoRows {
			err = e
		}
	}
	return
}
