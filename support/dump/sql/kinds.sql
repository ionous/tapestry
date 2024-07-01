-- name for every kind in alphabetical order
-- fix: its not quite right to use names 
-- because the current names are only unique with a domain set.
-- params:
--   ?1: base domain name
  with domains as (
  select md.requires as domain 
  from mdl_domain  md
  where md.domain = ?1
  union all  
  select ?1
)
select mk.kind
from mdl_kind mk
join domains md
  using (domain)
order by mk.kind
