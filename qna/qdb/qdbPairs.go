package qdb

// type PairData struct {
// 	// relation kind id, noun id, other noun id
// 	rel, noun, other string
// }

// func scanPairs(q *sql.Stmt, args ...interface{}) (ret []PairData, err error) {
// 	if rows, e := q.Query(args...); e != nil {
// 		err = e
// 	} else {
// 		var one PairData
// 		err = tables.ScanAll(rows, func() (err error) {
// 			ret = append(ret, one)
// 			return
// 		}, &one.rel, &one.noun, &one.other)
// 	}
// 	return
// }

// var newPairsFromChanges = `
// with newPairs as (
// 	select dn.rowid as domain, relKind, oneNoun, otherNoun, cardinality
// 	from mdl_pair mp
// 	join run_domain rd   -- run_domain instead of active_domains is a little faster.
// 	  on (mp.domain=rd.domain)
// 	join mdl_rel
// 	 using (relKind)
// 	join mdl_domain dn
// 	where rd.active = ?2 -- select all newly activated domains with rd
// 	and dn.domain = ?1   -- select just the current named domain with dn
// )`+relatePair

// used for changes due to domain changes
// only the domain name is passed.
var newPairsFromDomain = `
with newPairs as (
select domain, relKind, oneNoun, otherNoun, cardinality
	from mdl_pair mp
	join mdl_rel 
	 using (relKind)
	where mp.domain = ?1
)` + relatePair

// used for dynamic changes to relationships between nouns
// fix? this doesnt check to see whether the nouns are compatible with the relation
// ex. if oneNoun is compatible with mdl_rel.oneKind; the caller does that instead...
// ( see also: Runner.RelateTo )
var newPairsFromNames = `
with newPairs as (
	select ?1 as domain, rel.relKind, one.noun as oneNoun, other.noun as otherNoun, rel.cardinality
	from mdl_rel rel 
	join active_kinds ks
		on (ks.kind = rel.relKind)
	join active_nouns one
	join active_nouns other
	where ks.name = ?2
	and one.name = ?3
	and other.name = ?4
)` + relatePair

// zero out and mismatched pairs, and then write the new pairs
var relatePair = string(`
insert or replace into run_pair
	-- zero out the domain to hide the pair from future queries.
	select null, prev.relKind, prev.oneNoun, prev.otherNoun
	from newPairs
	join run_pair prev
		using (relKind)
	where ((prev.oneNoun = newPairs.oneNoun and newPairs.cardinality glob '*_one') or
			   (prev.otherNoun = newPairs.otherNoun and newPairs.cardinality glob 'one_*'))
union all
	select newPairs.domain, newPairs.relKind, newPairs.oneNoun, newPairs.otherNoun
	from newPairs
`)
