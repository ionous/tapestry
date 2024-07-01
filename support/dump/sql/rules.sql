-- returns the executable rules for a given pattern
-- params:
--   ?1: base domain name
--   ?2: pattern id
with domains as (
  select md.requires as domain 
  from mdl_domain  md
  where md.domain = ?1
  union all  
  select ?1
)
select coalesce(mu.name, 'rule ' || mu.rowid),
       mu.stop, mu.jump, mu.updates, mu.prog
from domains 
join mdl_rule mu
  using (domain)
join mdl_kind mk 
  on (mk.rowid = mu.kind) 
where mk.rowid = ?2
order by 
mu.rank,
-- tbd: positive rank items sort first specified to last specified (asc)
-- zero and negative ranked items sort last specified to first specified  (desc)
--  mu.rowid * (case when mu.rank > 0 then 1 else -1 end)
mu.rowid desc