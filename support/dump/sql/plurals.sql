-- one, many pairs
-- in alphabetical order
-- params:
--   ?1: base domain name
with domains as (
  select md.requires as domain 
  from mdl_domain  md
  where md.domain = ?1
  union all  
  select ?1
), 
onemany as (
select one, many 
  from mdl_plural mp
  join domains 
    using (domain)

union all 
  select singular, kind 
  from mdl_kind mk
  join domains 
    using (domain)
  where singular is not null
)
select one, many from onemany
order by one