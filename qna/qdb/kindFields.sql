
select mf.field, mf.affinity, ifnull(mt.kind, '') as type, mv.value
-- starting with all kinds in our kind hierarchy
from mdl_kind mk  
join mdl_kind ma
	-- is Y (is their name) a part of X (our path)
	on instr(',' || mk.path, 
					 ',' || ma.rowid || ',' )
	or (mk.rowid = ma.rowid) 
-- select all fields that those kinds have 
-- ( ma.rowid becomes the origin of the field definition )
join mdl_field mf
  on (ma.rowid = mf.kind)
-- pull in all values of any matching field
left join mdl_value_kind mv 
	-- this matches all the kinds from other trees
	-- ( ex. doors and supporters both have portability. )
  on (mv.field = mf.rowid)
	-- so filter initializers by the requested kind's fullpath
  and instr(
   ',' || mk.rowid || ',' || mk.path, -- full path
   ',' || mv.kind || ',')
-- finally determine the name of the field's type
left join mdl_kind mt 
	on (mt.rowid = mf.type)
where (mk.kind = "supporters")
-- sort to get fields in definition order
-- ( that's implicitly also kind order: all fields in earlier kinds are defined first )
-- then by the initializer nearest to our requested kind 
-- and, finally, put final values first.
order by mf.rowid, mv.kind desc, mv.final desc