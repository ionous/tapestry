-- list of parser aliases for a noun
-- params:
--  ?1: id of noun
--
select distinct name
from mdl_name 
where noun=?1
and rank <= 0
-- parser aliases are always rank -1
-- the default friendly name is rank 0
order by rank desc, rowid