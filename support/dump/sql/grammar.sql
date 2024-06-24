-- name, prog for every grammar
-- params:
--   ?1: base domain name
  with domains as (
  select md.requires as domain 
  from mdl_domain  md
  where md.domain = ?1
  union all  
  select ?1
)
select mg.name, mg.prog
from mdl_grammar mg
join domains md
  using (domain)
order by mg.name
