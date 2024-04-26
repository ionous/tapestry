
/**
 * the set of active domain names
 */
create view
temp.active_domains as 
select * 
from run_domain rd 
where rd.active > 0;

create view
temp.active_kinds as 
select ds.domain, mk.rowid as kind, mk.kind as name, mk.singular, mk.path, mk.at
from active_domains ds
join mdl_kind mk 
	using (domain);

create view
temp.active_plurals as 
select ds.domain, mp.many, mp.one, mp.at
from active_domains ds
join mdl_plural mp 
	using (domain);	

/* domain name, noun id, noun name, and kind id
* the domain name is a nod towards needing the domain name to fully scope the noun */ 
create view
temp.active_nouns as 
select ds.domain, mn.rowid as noun, mn.noun as name, mn.kind, mn.at
from active_domains ds
join mdl_noun mn 
	using (domain);

/* for finding relatives and reciprocals: returns relName, nounName, otherName */
create view
temp.active_names as
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