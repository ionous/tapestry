-- initial values of record fields for a given noun.
-- params:
--   ?1:  noun id
--
select mf.field, coalesce(mv.dot, '') as dotstr, mv.value
  from mdl_noun mn
  join mdl_value mv
    on (mv.noun = mn.rowid)
  join mdl_field mf
    on (mf.rowid = mv.field)
  where (mn.rowid = ?1) and dotstr != ''
  order by mf.field, length(dotstr), mv.final desc

