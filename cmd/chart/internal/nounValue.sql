select mn.noun, mv.value
-- all of the requested domains 
from mdl_domain pd
join mdl_domain cd
-- is Y (the domain we want) within X (within some other domain)
	on instr(',' ||cd.rowid || ',' || cd.path,
		',' ||pd.rowid || ',' )
-- all of the nouns for those domains 
join mdl_noun mn
	on (mn.domain = cd.rowid)
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
-- and get the rank 0 friendlier name
-- ( ... printed name might even more friendly ... )
where pd.domain like ?1
and ks.kind = ?2 
and mf.field= ?3 
and (not length(?4) or mn.noun=?4)