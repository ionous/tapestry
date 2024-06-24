-- list of names and the nouns they represent
-- params:
--   ?1: id of scene
--
with domains as (
  select md.requires as domain 
  from mdl_domain  md
  where md.domain = ?1
  union all  
  select ?1
)
select  mn.rowid, lower(my.name) as lname
from mdl_name my  
join domains
  using (domain)
join mdl_noun mn
  on (my.noun = mn.rowid)
where my.rank >= 0
order by lname, mn.rowid
