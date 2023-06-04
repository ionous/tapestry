package mdl

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/affine"
	"github.com/ionous/errutil"
)

// findDomain validates that the named domain exists
// the returned name is the same as the passed name.
func (m *Writer) findDomain(domain string) (ret string, err error) {
	if e := m.db.QueryRow(`
	select domain 
	from mdl_domain 
	where domain = ?1`, domain).Scan(&ret); e == sql.ErrNoRows {
		err = errutil.Fmt("domain not found %q", domain)
	} else {
		err = e
	}
	return
}

func (m *Writer) findOptionalKind(domain, kind string) (retDomain string, retKind *int, err error) {
	if len(kind) > 0 {
		if d, k, e := m.findKind(domain, kind); e != nil {
			err = e
		} else {
			retDomain = d
			retKind = &k
		}
	}
	return
}

func (m *Writer) findKind(domain, kind string) (retDomain string, retKind int, err error) {
	retDomain, _, retKind, err = m.pathOfKind(domain, kind)
	return
}

// where the returned comma separated path includes the id of the kind ",kid,...,"
func (m *Writer) pathOfKind(domain, kind string) (retDomain, retPath string, retKind int, err error) {
	if e := m.db.QueryRow(`
	select domain, mk.rowid, ',' || mk.rowid || ',' || mk.path
	from mdl_kind mk
	join domain_tree
		on (uses = domain)
	where base = ?1
	and kind = ?2
	limit 1`, domain, kind).Scan(&retDomain, &retKind, &retPath); e == sql.ErrNoRows {
		err = errutil.Fmt("no such kind %q in domain %q", kind, domain)
	} else {
		err = e
	}
	return
}

func (m *Writer) findNoun(domain, noun string) (retDomain string, retKind int, err error) {
	if e := m.db.QueryRow(`
	select domain, mn.rowid
	from mdl_noun mn
	join domain_tree
		on (uses = domain)
	where base = ?1
	and noun = ?2
	limit 1`, domain, noun).Scan(&retDomain, &retKind); e == sql.ErrNoRows {
		err = errutil.Fmt("no such noun %q in domain %q", noun, domain)
	} else {
		err = e
	}
	return
}

func (m *Writer) pathOfOptionalNoun(domain, noun string) (retDomain string, retNoun int, retPath string, err error) {
	if e := m.db.QueryRow(`
	select mn.domain, mn.rowid, ',' || mk.rowid || ',' || mk.path
	from mdl_noun mn
	join mdl_kind mk 
		on (mn.kind = mk.rowid)
	join domain_tree dt
		on (dt.uses = mn.domain)
	where base = ?1
	and noun = ?2
	limit 1`, domain, noun).Scan(&retDomain, &retNoun, &retPath); e != sql.ErrNoRows {
		err = e
	}
	return
}

// turn domain, kind, field into ids, associated with the local var's initial assignment.
// domain and kind become redundant b/c fields exist at the scope of the kind.
func (m *Writer) findField(domain, kind, field string) (retDomain string, retField int, err error) {
	if declaringDomain, kid, e := m.findKind(domain, kind); e != nil {
		err = e
	} else if e := m.db.QueryRow(`
		select rowid
		from mdl_field mf
		where kind = ?1
		and field = ?2`, kid, field).Scan(&retField); e == sql.ErrNoRows {
		err = errutil.Fmt("no such field %q in kind %q in domain %q", field, kind, domain)
	} else if e != nil {
		err = e
	} else {
		retDomain = declaringDomain
	}
	return
}

// check that the kind can store the requested value at the passed field
// returns the name of the field ( in case the originally specified field was a trait )
// FIX: i think this would work better using the runtime kind cache.
func (m *Writer) FindCompatibleField(domain, kind, field string, aff affine.Affinity) (retName, retClass string, err error) {
	var prev struct {
		name string
		aff  affine.Affinity
		cls  *string
	}
	if _, ancestry, _, e := m.pathOfKind(domain, kind); e != nil {
		err = errutil.Fmt("%w trying to add field %q", e, field)
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
select ma.name, 'bool', null
from allTraits ma
join fieldsInKind fk
where ma.name = @fieldName
and ma.kind = fk.typeId`,
		sql.Named("aspects", m.aspectPath),
		sql.Named("ancestry", ancestry),
		sql.Named("fieldName", field)).
		Scan(&prev.name, &prev.aff, &prev.cls); e != nil {
		if e == sql.ErrNoRows {
			err = errutil.Fmt("field %q not found in kind %q domain %q", field, kind, domain)
		} else {
			err = errutil.New("database error", e)
		}
	} else if prev.aff != aff {
		err = errutil.Fmt("affinity %s is incompatible with %s field %q in kind %q",
			aff, prev.aff, field, kind)
	} else if prev.name != field {
		// if they weren't asking for a trait, error:
		// the return name is the aspect; no subclass to speak of.
		retName = prev.name
	} else {
		retName = prev.name
		if prev.cls != nil {
			retClass = *prev.cls
		}
	}
	return
}
