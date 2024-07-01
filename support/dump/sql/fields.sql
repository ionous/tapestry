-- lists complete field set of a given kind
-- params:
--   ?1: name kind
--
-- fix? merge with query fieldsOf ( ex. using a cte )
-- 
select mf.field, mf.affinity, coalesce(mt.kind, '') as type, mv.value
from mdl_kind ks 
join mdl_kind ma
  -- if their id (Y) is in our path (X)
  -- then they are an ancestor
  on instr(',' || ks.rowid || ',' || ks.path,  -- our full path
           ',' || ma.rowid || ',' )
join mdl_field mf
  on (ma.rowid = mf.kind)
-- pull in all values of any matching field
left join mdl_value_kind mv 
  -- this matches all the kinds from other trees
  -- ( ex. doors and supporters both have portability. )
  on (mv.field = mf.rowid)
  -- so filter initializers by the requested kind's fullpath
  and instr(
   ',' || ks.rowid || ',' || ks.path, -- full path
   ',' || mv.kind || ',')
-- finally determine the name of the field's type
left join mdl_kind mt 
  on (mt.rowid = mf.type)
where (ks.kind = ?1)
-- sort to get fields in definition order
-- ( that's implicitly also kind order: all fields in earlier kinds are defined first )
-- then by the initializer nearest to our requested kind 
-- and, finally, put final values first.
order by mf.rowid, mv.kind desc, mv.final desc