-- list of all possible patterns in alphabetical order
-- with their id, name, labels + result
-- params:
--   ?1: base domain name
with domains as (
  select md.requires as domain 
  from mdl_domain  md
  where md.domain = ?1
  union all  
  select ?1
)
select mk.rowid,  -- rules are stored associated with the kind, not the pattern
       mk.kind, 
       coalesce(labels || ',', '') || -- null concat results in null
       coalesce(result, '') as lres
from mdl_kind mk
join domains md
  using (domain)
join mdl_pat mp
  on (mp.kind == mk.rowid)
order by mk.kind