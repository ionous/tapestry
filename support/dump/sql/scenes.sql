-- yields the domains with base ( tapestry ) first,
-- and the requested domain last.
-- params:
--   ?1: base domain name
--
select uses from domain_tree 
where base = ?1 
order by dist desc