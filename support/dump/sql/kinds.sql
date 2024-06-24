-- id, name for every kind
-- in alphabetical order
-- params:
--   ?1: base domain name
  with domains as (
  select md.requires as domain 
  from mdl_domain  md
  where md.domain = ?1
  union all  
  select ?1
)
select mk.rowid, mk.domain, mk.kind
from mdl_kind mk
join domains md
  using (domain)
order by mk.kind
