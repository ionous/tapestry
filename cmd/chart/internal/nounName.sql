select my.name 
-- all of the requested domains 
from mdl_domain pd
join mdl_domain cd
-- is Y (the domain we want) within X (within some other domain)
	on instr(',' ||cd.rowid || ',' || cd.path,
		',' ||pd.rowid || ',' )
-- all of the nouns for those domains 
join mdl_noun mn
	on (mn.domain = cd.rowid)
-- the rank 0 name of those nouns
join mdl_name my 
	on (my.noun = mn.rowid) and (not my.rank)
where mn.noun=?1