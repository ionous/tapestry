-- generate pairs domain, dependency, and the depth of the dependency
-- ( the larger the depth, the closer to root )
select base, uses from domain_tree
where dist != 0 -- exclude itself, for once.
order by base, dist, uses
