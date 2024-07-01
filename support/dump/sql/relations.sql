-- select all kinds which are of type relation
-- that are in the specified scene
-- params:
--   ?1: base domain name
  with domains as (
  select md.requires as domain 
  from mdl_domain  md
  where md.domain = ?1
  union all  
  select ?1
)
-- id, name
select mr.relKind, mk.kind, one.kind, other.kind, mr.cardinality
from mdl_kind mk 
join domains md
  on (md.domain = mk.domain)
join mdl_rel mr
  on (mr.relKind = mk.rowid)
join mdl_kind one 
  on (one.rowid = mr.oneKind)
join mdl_kind other 
  on (other.rowid = mr.otherKind)

