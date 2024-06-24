-- initial values of fields for a given noun
-- params:
--   ?1: id of noun
--
select mf.field, coalesce(mv.dot, '') as dotstr, mv.value
  from mdl_noun mn
  join mdl_value mv
    on (mv.noun = mn.rowid)
  join mdl_field mf
    on (mf.rowid = mv.field)
  where (mn.rowid = ?1)
  order by mf.field, length(dotstr), mv.final desc

