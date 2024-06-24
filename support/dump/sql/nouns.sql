-- list of names and the nouns they represent
-- params:
--   ?1: id of scene
--
with domains as (
  select md.requires as domain 
  from mdl_domain md
  where md.domain = ?1
  union all  
  select ?1
)
select mn.rowid, mn.domain, mn.noun, mk.kind 
from mdl_noun mn
join domains
  using (domain)
join mdl_kind mk 
  on (mk.rowid = mn.kind)