-- initial values of simple fields for a given noun
-- params:
--   ?1: id of noun
--
select mf.field, mv.value
  from mdl_noun mn
  join mdl_value mv
    on (mv.noun = mn.rowid)
  join mdl_field mf
    on (mf.rowid = mv.field)
  where (mn.rowid = ?1) and coalesce(mv.dot, '') == ''
  order by mf.field, mv.final desc

