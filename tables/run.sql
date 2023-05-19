/* 
 * for saving, restoring a player's game session. 
 * a global application active count increases monotonically 
 * for every request to activate any domain.
 * a particular domain is in scope when its own active value is non-zero; 
 * its value gets reset to zero when the domain falls out of scope.
 */
create table if not exists 
	run_domain( domain text, active int, primary key( domain )); 

/* we dont need an "active" -- we can join against run_domain, or write 0 to domain to disable a pair. */
create table if not exists 
	run_pair( domain text, relKind int, oneNoun int, otherNoun int, unique( relKind, oneNoun, otherNoun ) ); 

/**
 * the set of active domain names
 */
create view if not exists
active_domains as 
select * 
from run_domain rd 
where rd.active > 0;

/**
 * the set of active domain names
 */
create view if not exists
active_domains as 
select * 
from run_domain rd 
where rd.active > 0;

create view if not exists
active_kinds as 
select ds.domain, mk.rowid as kind, mk.kind as name, mk.path
from active_domains ds
join mdl_kind mk 
	using (domain);

/* domain name, noun id, noun name, and kind id
* the domain name is a nod towards needing the domain name to fully scope the noun */ 
create view if not exists
active_nouns as 
select ds.domain, mn.rowid as noun, mn.noun as name, mn.kind
from active_domains ds
join mdl_noun mn 
	using (domain);

create view if not exists
active_plurals as 
select ds.domain, mp.many, mp.one, mp.at
from active_domains ds
join mdl_plural mp 
	using (domain);	

create view if not exists
active_rev as 
select ds.domain, mp.oneWord, mp.otherWord, mp.at
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