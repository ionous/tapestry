package mdl

import (
	"database/sql"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

var mdl_field = tables.Insert("mdl_field", "domain", "kind", "field", "affinity", "type", "at")

func (pen *Pen) addField(kid, cls kindInfo, field string, aff affine.Affinity) (err error) {
	domain, at := pen.domain, pen.at
	if rows, e := pen.db.Query(`
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
)
select origin, name, affinity, typeName, aspect
from existingFields
join pendingFields
using(name)
`, kid.fullpath(), field, cls.id, pen.paths.aspectPath); e != nil {
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
		} else {
			// keep null instead of zero ids
			var clsid *int64
			if cls.id != 0 {
				clsid = &cls.id
			}
			if _, e := pen.db.Exec(mdl_field, domain, kid.id, field, aff, clsid, at); e != nil {
				err = errutil.New("database error", e)
			}
		}
	}
	return
}

type fieldInfo struct {
	id     int64
	name   string
	domain string
	cls    classInfo
	aff    affine.Affinity
}

func (f *fieldInfo) class() classInfo {
	return f.cls
}

// turn a trait boolean value into an aspect text value ( containing the name of the trait )
// or, for all other types of values, return the passed assignment back to the caller.
func (f *fieldInfo) rewriteTrait(name string, value assign.Assignment) (ret assign.Assignment, err error) {
	if name == f.name {
		ret = value
	} else {
		switch v := value.(type) {
		default:
			err = errutil.Fmt("incompatible assignment to trait, got %T", value)
		case *assign.FromBool:
			switch b := v.Value.(type) {
			default:
				err = errutil.Fmt("traits only support literal bools, got %T", value)
			case *literal.BoolValue:
				if !b.Value {
					err = errutil.New("opposite trait assignment not supported")
				} else {
					ret = &assign.FromText{Value: &literal.TextValue{Value: name}}
				}
			}
		}
	}
	return
}

// recursively descend through the specified fields
// returns the outer and inner most fields
func (pen *Pen) digField(noun nounInfo, path []string) (retout, retin fieldInfo, err error) {
	root, path := path[0], path[1:]
	if outer, e := pen.findField(noun.class(), root); e != nil {
		err = e
	} else {
		inner := outer
		for i := 0; i < len(path) && err == nil; i++ {
			subField := path[i]
			if inner.aff != affine.Record {
				err = errutil.Fmt("expected a field of type record for noun %q, kind %q, path %q(%d)",
					noun.name, noun.kind, strings.Join(path, "."), i)
			} else {
				inner, err = pen.findField(inner.class(), subField)
			}
		}
		if err == nil {
			retout, retin = outer, inner
		}
	}
	return
}

// check that the kind can store the requested value at the passed field
// returns the name of the field ( in case the originally specified field was a trait )
func (pen *Pen) findDefaultTrait(kind classInfo) (ret string, err error) {
	err = pen.db.QueryRow(`select field 
		from mdl_field where kind = ?1`,
		kind.id).Scan(&ret)
	return
}

// check that the kind can store the requested value at the passed field
// returns the name of the field ( in case the originally specified field was a trait )
func (pen *Pen) findField(kind classInfo, field string) (ret fieldInfo, err error) {
	if e := pen.db.QueryRow(` 
-- all possible traits:
with allTraits as (	
	select mk.rowid as kind,   -- id of the aspect,
				 field as name,      -- name of trait 
				 mf.domain,          -- domain of the field
	       mk.kind as aspect,  -- name of aspect
	       mk.domain           -- name of originating domain
	from mdl_field mf 
	join mdl_kind mk
		on(mf.kind = mk.rowid)
	-- where the field's kind (X) contains the aspect kind (Y)
	where instr(',' || mk.path, @aspects)
)
-- all fields of the targeted kind:
, fieldsInKind as (
	select mf.rowid as id,
				 field as name,        -- field name 
				 mf.domain,            -- domain name of the field 
				 affinity,             -- affinity 
				 mt.rowid as typeId,   -- type of the field 
				 mt.kind as typeName,  -- name of that type 
				 (',' || mt.rowid || ',' || mt.path) as fullpath
	from mdl_field mf 
	join mdl_kind mk 
		-- does our ancestry (X) contain any of these kinds (Y)
		on ((mf.kind = mk.rowid) and instr(@ancestry, ',' || mk.rowid || ',' ))
	left join mdl_kind mt 
		on (mt.rowid = mf.type)
)
-- fields in the target kind
-- if the field isnt a record; the type info (id,name,path) can be null
select id, name, kf.domain, affinity, coalesce(typeId,0), coalesce(typeName, ''), coalesce(fullpath, '')
from fieldsInKind kf
where name = @fieldName 
union all

-- traits in the target kind: return the aspect
select id, ma.aspect, ma.domain, 'text', 0, "", ""
from allTraits ma
join fieldsInKind fk
	on (ma.kind = fk.typeId)
where ma.name = @fieldName`,
		sql.Named("aspects", pen.paths.aspectPath),
		sql.Named("ancestry", kind.fullpath),
		sql.Named("fieldName", field)).
		Scan(&ret.id, &ret.name, &ret.domain, &ret.aff, &ret.cls.id, &ret.cls.name, &ret.cls.fullpath); e != nil {
		if e == sql.ErrNoRows {
			err = errutil.Fmt("%w field %q in kind %q domain %q", Missing, field, kind.name, pen.domain)
		} else {
			err = errutil.New("database error", e)
		}
	}
	return
}
