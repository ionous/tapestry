/* for saving, restoring a player's game session. */
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
active_domains as 
select md.rowid as domain, md.domain as name
from run_domain rd 
join mdl_domain md 
	on rd.active > 0 and rd.domain = md.rowid;

create view if not exists
active_kinds as 
select ds.name as domain, mk.rowid as kind, mk.kind as name, mk.path
from active_domains ds
join mdl_kind mk 
	using (domain);

/* domain name, noun id, noun name, and kind id
* the domain name is a nod towards needing the domain name to fully scope the noun */ 
create view if not exists
active_nouns as 
select ds.name as domain, mn.rowid as noun, mn.noun as name, mn.kind
from active_domains ds
join mdl_noun mn 
	using (domain);

create view if not exists
active_plurals as 
select ds.name as domain, mp.many, mp.one, mp.at
from active_domains ds
join mdl_plural mp 
	using (domain);	

create view if not exists
active_rev as 
select ds.name as domain, mp.oneWord, mp.otherWord, mp.at
from active_domains ds
join mdl_rev mp 
	using (domain);	

/* for finding relatives and reciprocals: returns relName, nounName, otherName */
create view if not exists
rp_names as
select rp.domain, mk.kind as relName, one.noun as oneName, other.noun as otherName
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