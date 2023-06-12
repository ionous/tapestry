package mdl

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

var fieldSource = ` 
-- all possible traits:
with allTraits as (	
	select mk.rowid as kind,  -- id of the aspect,
				 field as name,     -- name of trait 
	       mk.kind as aspect, -- name of aspect
	       mk.domain          -- name of originating domain
	from mdl_field mf 
	join mdl_kind mk
		on(mf.kind = mk.rowid)
	-- where the field's kind (X) contains the aspect kind (Y)
	where instr(',' || mk.path, ?4 )
)
-- all fields of the targeted kind:
, fieldsInKind as (
	select mk.domain, field as name, affinity, mf.type as typeId, mt.kind as typeName
	from mdl_field mf 
	join mdl_kind mk 
		-- does our ancestry (X) contain any of these kinds (Y)
		on ((mf.kind = mk.rowid) and instr(?1, ',' || mk.rowid || ',' ))
	left join mdl_kind mt 
		on (mt.rowid = mf.type)
)
-- fields and traits in the target kind
-- ( all of them, because we dont know what might conflict with a pending aspect )
, existingFields( origin, name, affinity, typeName ) as (
	-- fields in the target kind
	select format('domain "%z"', domain), name, affinity, typeName
	from fieldsInKind

	union all

	-- traits in the target kind
	select format('aspect "%z"', ma.aspect), -- report the aspect as the origin 
					ma.name,   -- trait name 
					'bool',    -- fk.affinity is 'text', each trait is 'bool'
					null       -- traits have null type currently.
	from fieldsInKind fk
	join allTraits ma
		on (fk.typeId = ma.kind)
)
, pendingFields(name, aspect) as ( 
	-- the name of the field we're adding;
	select ?2, null
	union all 

	-- the names of traits when adding a field of type aspect; if any.
	select name, aspect
	from allTraits ma
	where (?3 = ma.kind)
)`

func (m *Modeler) addField(domain string, kid, cls kindInfo, field string, aff affine.Affinity, at string) (err error) {
	// println("=== adding field", domain, kid.name, field, cls.name)
	// if existing, e := tables.QueryStrings(m.db, fieldSource+`
	// 	select origin|| ', ' || name || ', '|| affinity|| ', ' || typeName
	// 	from existingFields`,
	// 	kid.fullpath(), field, cls.id, m.aspectPath); e != nil {
	// 	panic(e)
	// } else if pending, e := tables.QueryStrings(m.db, fieldSource+`
	// 	select '-' || name || ', ' || coalesce(aspect, 'nil')
	// 	from pendingFields`,
	// 	kid.fullpath(), field, cls.id, m.aspectPath); e != nil {
	// 	panic(e)
	// } else {
	// 	println("existing", strings.Join(existing, ";\n "))
	// 	println("pending", strings.Join(pending, ";\n "))
	// }

	if rows, e := m.db.Query(fieldSource+`
select origin, name, affinity, typeName, aspect
from existingFields
join pendingFields
using(name)
`, kid.fullpath(), field, cls.id, m.aspectPath); e != nil {
		err = errutil.New("database error", e)
	} else {
		var prev struct {
			name   string          // trait or field causing a conflict
			aspect sql.NullString  // aspect if any of the pending name
			origin string          // aspect or kind of the existing field
			aff    affine.Affinity // affinity of the existing field ( ex. 'bool' for aspects )
			cls    sql.NullString  // type name ( or null ) of existing field
		}
		if e := tables.ScanAll(rows, func() (err error) {
			// if the names differ, then the conflict is due to a trait ( being added or already existing )
			if prev.name != field {
				// adding an aspect: the conflict reports the pending aspect so this case can be detected
				if prev.aspect.String == cls.name {
					// is there a way to determine whether the origin is a domain or aspect
					err = errutil.Fmt("%w new field for kind %q of aspect %q conflicts with existing field %q from %s",
						Conflict, kid.name, field, prev.name, prev.origin)
				} else if prev.aspect.Valid {
					err = errutil.Fmt("%w new field for kind %q of aspect %q conflicts with trait %q from aspect %q",
						Conflict, kid.name, field,
						prev.name, prev.aspect.String)
				} else {
					// when does this show up?
					err = errutil.Fmt("%w field %q for kind %q was %s(%s) from %s, now %s(%s) in %q",
						Conflict, field, kid.name,
						prev.aff, prev.cls.String, prev.origin,
						aff, cls.name, domain)
				}
			} else if aff == prev.aff && cls.name == prev.cls.String {
				// if the affinity and typeName are the same, then its a duplicate
				err = errutil.Fmt("%w field %q for kind %q of %s(%s) from %s and now domain %q",
					Duplicate, field, kid.name,
					aff, cls.name,
					prev.origin, domain)
			} else {
				// otherwise, its a conflict
				err = errutil.Fmt("%w field %q for kind %q of %s(%s) from %s was redefined as %s(%s) in domain %q",
					Conflict, field, kid.name,
					prev.aff, prev.cls.String, prev.origin,
					aff, cls.name, domain)
			}
			return
		}, &prev.origin, &prev.name, &prev.aff, &prev.cls, &prev.aspect); e != nil {
			err = e
		} else if _, e := m.field.Exec(domain, kid.id, field, aff, cls.id, at); e != nil {
			err = errutil.New("database error", e)
		}
	}
	return
}

// check that the kind can store the requested value at the passed field
// returns the name of the field ( in case the originally specified field was a trait )
// FIX: i think this would work better using the runtime kind cache.
func (m *Modeler) FindCompatibleField(domain, kind, field string, aff affine.Affinity) (retName, retClass string, err error) {
	var prev struct {
		name string
		aff  affine.Affinity
		cls  *string
	}
	if kid, e := m.findRequiredKind(domain, kind); e != nil {
		err = errutil.Fmt("%w trying to find field %q", e, field)
	} else if e := m.db.QueryRow(` 
-- all possible traits:
with allTraits as (	
	select mk.rowid as kind,    -- id of the aspect,
				 field as name,      -- name of trait 
	       mk.kind as aspect,  -- name of aspect
	       mk.domain          -- name of originating domain
	from mdl_field mf 
	join mdl_kind mk
		on(mf.kind = mk.rowid)
	-- where the field's kind (X) contains the aspect kind (Y)
	where instr(',' || mk.path, @aspects)
)
-- all fields of the targeted kind:
, fieldsInKind as (
	select mk.domain, field as name, affinity, mf.type as typeId, mt.kind as typeName
	from mdl_field mf 
	join mdl_kind mk 
		-- does our ancestry (X) contain any of these kinds (Y)
		on ((mf.kind = mk.rowid) and instr(@ancestry, ',' || mk.rowid || ',' ))
	left join mdl_kind mt 
		on (mt.rowid = mf.type)
)
-- fields in the target kind
select name, affinity, typeName
from fieldsInKind
where name = @fieldName 
union all

-- traits in the target kind: return the aspect
select ma.aspect, 'text', null
from allTraits ma
join fieldsInKind fk
where ma.name = @fieldName
and ma.kind = fk.typeId`,
		sql.Named("aspects", m.aspectPath),
		sql.Named("ancestry", kid.fullpath()),
		sql.Named("fieldName", field)).
		Scan(&prev.name, &prev.aff, &prev.cls); e != nil {
		if e == sql.ErrNoRows {
			err = errutil.Fmt("field %q not found in kind %q domain %q", field, kind, domain)
		} else {
			err = errutil.New("database error", e)
		}
	} else {
		// if the names don't match, than the search found a trait of an aspect:
		if prev.name != field {
			if aff != affine.Bool {
				err = errutil.Fmt("affinity %s is incompatible with trait %q of aspect %q in kind %q",
					aff, field, prev.name, kind)
			} else {
				retName = prev.name
			}
		} else {
			// otherwise the search returned a normal field:
			if prev.aff != aff {
				err = errutil.Fmt("affinity %s is incompatible with %s field %q in kind %q",
					aff, prev.aff, field, kind)
			} else {
				retName = prev.name
				if prev.cls != nil {
					retClass = *prev.cls
				}
			}
		}
	}
	return
}
