package mdl

import (
	"database/sql"
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/tables"
)

func ReadAssignments(db *sql.DB) ([]string, error) {
	return tables.QueryStrings(db, `
    select domain ||':'|| many ||':'|| one
from mdl_plural`)
}

// an array of dependencies
func ReadDomain(db *sql.DB, n string) (ret []string, err error) {
	if ds, e := tables.QueryStrings(db,
		`select uses from domain_tree 
		where base = ?1 
		order by dist desc`,
		n); e != nil {
		err = e
	} else {
		// the domain tree includes the domain itself
		// the tests dont expect that.
		ret = ds[:len(ds)-1]
	}
	return
}

func ReadDomains(db *sql.DB) ([]string, error) {
	return tables.QueryStrings(db,
		`select domain from domain_order`)
}

// domain, kind, field name, affinity, subtype
// sorted by kind and within each kind, the field name
func ReadFields(db *sql.DB) ([]string, error) {
	return tables.QueryStrings(db, `
  select mk.domain ||':'|| mk.kind ||':'|| mf.field  ||':'|| 
         mf.affinity  ||':'|| coalesce(mt.kind, '') 
	from mdl_field mf 
	join mdl_kind mk 
		on(mf.kind = mk.rowid)
	join mdl_domain md 
		on(md.domain = mk.domain)
	left join mdl_kind mt 
		on(mf.type = mt.rowid)
	order by md.rowid, mk.kind, mf.field`)
}

// domain, input, serialized program
func ReadGrammar(db *sql.DB) ([]string, error) {
	return tables.QueryStrings(db, `
  select domain ||':'|| name ||':'|| prog
	from mdl_grammar`)
}

// domain, kind, expanded materialized path
// ordered by domain, length of path, and name
// ( that erases their natural, dependency order --
//   but independent siblings dont otherwise have a consistent order for testing )
func ReadKinds(db *sql.DB) (ret []string, err error) {
	type kind struct {
		id           int
		domain, kind string
	}
	var kinds []kind
	var k kind
	if e := tables.QueryAll(db,
		`select mk.rowid, domain, kind
		from mdl_kind mk
		join mdl_domain md
			using(domain)
		order by md.rowid, length(mk.path), mk.kind
		`, func() (_ error) {
			kinds = append(kinds, k)
			return
		}, &k.id, &k.domain, &k.kind); e != nil {
		err = e
	} else {
		// do this the manual way for now because its easier
		// fix? use a recursive query to expand the path.
		for _, k := range kinds {
			// just to be confusing, this is the opposite order of KindOfAncestors
			// root is on the right here.
			if path, e := tables.QueryStrings(db,
				`select mk.kind
				from mdl_kind ks
				join mdl_kind mk
					-- is Y (is their name) a part of X (our path)
					on instr(',' || ks.path,
									 ',' || mk.rowid || ',' )
				where ks.rowid = ?1
				order by mk.rowid desc`, k.id); e != nil {
				err = e
				break
			} else {
				row := fmt.Sprintf("%s:%s:%s", k.domain, k.kind, strings.Join(path, ","))
				ret = append(ret, row)
			}
		}
	}
	return
}

// domain, kind, name, serialized initialization
func ReadLocals(db *sql.DB) ([]string, error) {
	// read kind assignments for those that have initializers
	// the original order was per domain, and then in dependency order withing that domain
	//  wonder if itd be enough to simply go by kind id since the order of writing should follow
	return tables.QueryStrings(db, `
	select mk.domain ||':'|| mk.kind ||':'|| mf.field ||':'|| ma.value
	from mdl_default ma
	join mdl_field mf 
		on(ma.field = mf.rowid)
	join mdl_kind mk 
		on(mf.kind = mk.rowid)
	order by mk.rowid, mf.rowid, ma.rowid`)
}

// domain, noun, name, rank
func ReadNames(db *sql.DB) ([]string, error) {
	return tables.QueryStrings(db, `
	select my.domain ||':'|| mn.noun ||':'|| my.name ||':'|| my.rank
	from mdl_name my
	join mdl_noun mn
		on(my.noun = mn.rowid)
	join mdl_domain md
		on(mn.domain = md.domain)
	order by md.domain, mn.noun, my.rank, my.name`)
}

// domain, noun, kind
func ReadNouns(db *sql.DB) ([]string, error) {
	return tables.QueryStrings(db, `
  select mn.domain ||':'|| mn.noun ||':'|| mk.kind
	from mdl_noun mn
	join mdl_kind mk
		on(mn.kind = mk.rowid)
	join mdl_domain md
		on(mn.domain = md.domain)
	order by md.domain, mk.kind, mn.noun`)
}

// domain, oneWord, otherWord
func ReadOpposites(db *sql.DB) ([]string, error) {
	return tables.QueryStrings(db, `
  select domain ||':'|| oneWord ||':'|| otherWord
	from mdl_rev mp 
	order by domain, oneWord`)
}

// domain, relation, noun, other noun
// original order was domain, and alpha relation name
func ReadPairs(db *sql.DB) ([]string, error) {
	return tables.QueryStrings(db, `
	with domain_rank(domain, rank) as (
		select *, row_number() over (order by 1) 
		from domain_order
	)
	select mp.domain ||':'|| mk.kind ||':'|| one.noun ||':'|| other.noun
	from domain_rank dr
	join mdl_pair mp 
	on(mp.domain = dr.domain)
	join mdl_kind mk
		on(mp.relKind = mk.rowid)
	join mdl_noun one
		on(oneNoun = one.rowid)
	join mdl_noun other
		on(otherNoun = other.rowid)
	order by dr.rank, mk.kind, one.noun, other.noun `)
}

// domain, pattern, labels, result field
func ReadPatterns(db *sql.DB) ([]string, error) {
	return tables.QueryStrings(db, `
  select mk.domain ||':'|| mk.kind ||':'|| coalesce(labels,'') ||':'|| coalesce(result,'')
	from mdl_pat mp
	join mdl_kind mk 
		on(mp.kind = mk.rowid)`)
}

// domain, many, one
func ReadPlurals(db *sql.DB) ([]string, error) {
	return tables.QueryStrings(db, `
  select domain ||':'|| many ||':'|| one
	from mdl_plural`)
}

// domain, relation, one kind, other kind, cardinality
// ordered by name of the relation for test consistency.
func ReadRelations(db *sql.DB) ([]string, error) {
	return tables.QueryStrings(db, `
  select mk.domain ||':'|| mk.kind ||':'|| one.kind ||':'|| other.kind ||':'|| mr.cardinality
  from mdl_rel mr 
  join mdl_kind mk
  	on(relKind = mk.rowid)
	join mdl_kind one
		on(oneKind = one.rowid)
	join mdl_kind other
		on(otherKind = other.rowid)
	order by mk.kind`)
}

// domain, pattern, target, rank, filter, prog
func ReadRules(db *sql.DB) ([]string, error) {
	return tables.QueryStrings(db, `
  select mr.domain ||':'|| mk.kind ||':'|| coalesce(mt.kind, '') ||':'|| rank ||':'|| filter ||':'|| prog
  from mdl_rule mr 
  join mdl_kind mk
  	on(mr.kind = mk.rowid)
  join mdl_domain md 
  	on(md.domain = mr.domain)
  left join mdl_kind mt 
  	on(mr.target = mt.rowid)
  order by md.rowid desc, mk.kind, abs(mr.rank), mr.rowid desc
  `)
}

// domain, noun, field, dot, value
func ReadValues(db *sql.DB) ([]string, error) {
	return tables.QueryStrings(db, `
	select mn.domain ||':'|| mn.noun ||':'|| mf.field ||':'|| coalesce(mv.dot, '') ||':'|| mv.value
	from mdl_value mv 
	join mdl_noun mn 
		on(mv.noun = mn.rowid)
	join mdl_field mf	
		on(mv.field = mf.rowid)
	join mdl_domain md 
		on(md.domain = mn.domain)
	order by md.rowid, mn.noun, mf.field
		`)
}
