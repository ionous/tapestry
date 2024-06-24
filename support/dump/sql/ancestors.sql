-- ancestors of a given kind
-- lists the kind first, followed by all of its ancestors root ("kinds") last.
-- params:
--   ?1: id of kind
--
select mk.kind 
from mdl_kind ks
join mdl_kind mk
  -- is Y (is their name) a part of X (our path)
  on instr(',' || ks.path, 
           ',' || mk.rowid || ',' )
  or (ks.rowid == mk.rowid)  -- to include the kind itself
where (ks.rowid = ?1) 
order by mk.rowid desc