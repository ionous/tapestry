-- select all kinds which are of type relation
-- that are in the specified scene
-- params:
--   ?1: base domain name
--   ?2: rel id
-- mdl_pair: domain (name), relKind (id), oneNoun (id), otherNoun (id)
  with domains as (
  select md.requires as domain 
  from mdl_domain md
  where md.domain = ?1
  union all  
  select ?1
)
select one.noun as oneName, other.noun as otherName
from mdl_pair mp
join domains md
  using (domain)
join mdl_noun one
  on (one.rowid = mp.oneNoun)
join mdl_noun other
  on (other.rowid = mp.otherNoun)
where mp.relKind = ?2
order by mp.rowid
