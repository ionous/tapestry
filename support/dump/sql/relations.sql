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
select mk.rowid, mk.kind
from mdl_kind mk
join domains md
  using (domain)
join mdl_kind ks 
    -- is Y (is their name) a part of X (our path)
  on instr(',' || mk.path, 
         ',' || ks.rowid || ',' )
where ks.kind == "relations"

