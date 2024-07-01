-- ancestors of a given kind
-- lists the kind first, followed by all of its ancestors root ("kinds") last.
-- params:
--   ?1: exact name of kind
--
select mk.kind 
from mdl_kind ks
join mdl_kind mk
   -- if their id (Y) is in our path (X)
  -- then they are an ancestor
  on instr(',' || ks.rowid || ',' || ks.path, -- our full path
           ',' || mk.rowid || ',' )
where (ks.kind = ?1) 
order by mk.rowid desc