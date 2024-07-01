select my.name 
from mdl_noun mn
join mdl_name my 
	on (my.noun = mn.rowid) and (not my.rank)
where mn.noun=?1