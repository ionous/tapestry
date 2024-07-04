/**
 * generate triplets of domains and all of their requirements.
 * the smaller dist(ances) are closer to the named base domain; 
 * larger distances are most root-like.
 */
create view domain_tree as 
-- note: insertion doesnt allow exactly redundant pairs
-- so we shouldnt need distinct here.
select base, uses, dist from
( 
  with recursive paths(base, child, uses, dist) as (
    -- seed the recursion with matching parents 
    select domain, domain, requires, 1
    from mdl_domain 
    -- tbd: if we're filtering this so frequently, maybe it never should have been added
    -- but then how to declare a domain?
    -- (maybe, depending on itself; and avoid the ugly join below )
    where requires  != ""
    union all
    -- search upwards until there are no more parents
    select base, domain, requires, dist+1
      from paths 
      join mdl_domain d
      on (uses = domain)
    where requires  != ""
  )
  -- if there are two different routes to reach a dependency
  -- we want the route with the most number of steps
  -- this matters for 
  select *, max(dist) over (partition by base, uses) as maxdist
  from paths 
)
where dist = maxdist
-- ugly, join with itself.
union all 
select distinct domain, domain, 0 
from mdl_domain;

/**
 * list domains so they appear before they are needed as a requirement of another domain. 
 */
create view domain_order as 
select distinct uses as domain 
from (
  select uses, dist 
  from domain_tree
  order by dist desc
);
