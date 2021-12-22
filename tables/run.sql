/**
 * for saving, restoring a player's game session.
 */
create table if not exists 
	run_domain( domain int, active int, primary key( domain )); 


/* we dont need an "active" -- we can join against run_domain, or write 0 to domain to disable a pair. */
create table if not exists 
	run_pair( domain int, relKind int, oneNoun int, otherNoun int, unique( relKind, oneNoun, otherNoun ) ); 

/**
 * find the id and name of all active domains
 * path isnt really needed because any parts of an activated domain are themselves individually active.
 */
create view if not exists
domain_scope as 
select md.rowid as domain, md.domain as name
from run_domain rd 
join mdl_domain md 
	on rd.active > 0 and rd.domain = md.rowid;

create view if not exists
kind_scope as 
select mk.rowid as kind, mk.kind as name, mk.path
from domain_scope
join mdl_kind mk 
	using (domain);

create view if not exists
noun_scope as 
select mn.rowid as noun, mn.noun as name, mn.kind
from domain_scope
join mdl_noun mn 
	using (domain);

/* returns the kind's id as 'relKind' */ 
create view if not exists
rel_scope as 
select mk.rowid as relKind, mk.kind as name, mr.oneKind, mr.otherKind, mr.cardinality
from domain_scope
join mdl_kind mk 
	using (domain)
join mdl_rel mr
	on (mk.rowid = mr.relKind);

/* */
create view if not exists
rp_names as
select rp.domain, mk.kind as relName, rp.relKind, one.noun as oneName, rp.oneNoun, other.noun as otherName, rp.otherNoun, rel.cardinality
from run_pair rp
join mdl_kind mk
	on (mk.rowid = rp.relKind)
join mdl_noun one
	on (one.rowid = rp.oneNoun)
join mdl_noun other
	on (other.rowid = rp.otherNoun)
join mdl_rel rel 
	on (rel.relKind = rp.relKind)
where rp.domain > 0;