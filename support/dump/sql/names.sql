-- nouns and short names, sorted by the latter.
--
-- note: parser names include short names, but not the other way around.
-- ie. "understandings" don't influence the names in a story,
-- mirroring inform's behavior. ( tbd: is it a good thing? )
--
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
select lower(my.name) as lname, mn.noun
from mdl_name my  
join domains
  using (domain)
join mdl_noun mn
  on (my.noun = mn.rowid)
where my.rank >= 0
order by lname, mn.noun
