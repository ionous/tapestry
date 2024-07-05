package weave

import (
	"database/sql"
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/tables"
)

type rival struct {
	Category, Domain, Key, Value, At string
}

func (r rival) Error() string {
	return fmt.Sprintf("domain %q at %q for %s %q", r.Domain, r.At, r.Category, r.Value)
}

type rivalErrorList []rival

func (rs rivalErrorList) Error() string {
	var b strings.Builder
	for i, r := range rs {
		if i > 0 {
			b.WriteRune('\n')
		}
		b.WriteString(r.Error())
	}
	return b.String()
}

// exposed for testing:
// tbd: maybe this could pull in the newly relevant domains;
// ( ex. use domain = active count )
// currently it happens after the domains have been activated
// and therefore compares everything to everything each time.
// note: fields don't have rivals because they all exist in the same domain as their owner kind.
func findRivals(db tables.Querier) (ret []rival, err error) {
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
	)

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

	select 'plural', a.domain, a.at, a.many, a.one
		from active_plurals as a 
		join active_plurals as b 
			using(many)
		where a.domain != b.domain
		and a.one != b.one
	`); e != nil {
		err = fmt.Errorf("database error %s", e)
	} else {
		var group, domain, key, value string
		var at sql.NullString
		err = tables.ScanAll(rows, func() (_ error) {
			ret = append(ret, rival{group, domain, key, value, at.String})
			return
		}, &group, &domain, &at, &key, &value)
	}
	return
}
