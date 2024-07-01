select mn.noun, mv.value
-- all of the nouns for those domains 
from mdl_noun mn
-- the hierarchy of kinds used by those nouns 	
join mdl_kind mk 
	on (mk.rowid = mn.kind)
join mdl_kind ks 
    -- is Y (the kind we want) within X (the path of our noun's kind)
	on instr( ',' || mk.rowid || ',' || mk.path, 
		',' ||ks.rowid || ',' )
-- all fields for those kinds 
join mdl_field mf 
	on (mf.kind=ks.rowid)
-- all of the values for those fields used by that are used by our nouns. 
join mdl_value mv 
	on (mv.noun=mn.rowid) and (mv.field=mf.rowid)
-- all of the domains in which that noun is active
join mdl_domain md
	on (mn.domain=md.domain)
-- and get the rank 0 friendlier name
-- ( ... printed name might even more friendly ... )
where (md.domain=?1 or md.requires=?1)
and ks.kind = ?2 
and mf.field= ?3 
and (not length(?4) or mn.noun=?4)