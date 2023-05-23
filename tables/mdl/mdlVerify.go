package mdl

import (
	"database/sql"

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
	if e := m.db.QueryRow(`
	select domain, mk.rowid
	from mdl_kind mk
	join domain_tree
		on (uses = domain)
	where base = ?1
	and kind = ?2
	limit 1`, domain, kind).Scan(&retDomain, &retKind); e == sql.ErrNoRows {
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

// turn domain, kind, field into ids, associated with the local var's initial assignment.
// domain and kind become redundant b/c fields exist at the scope of the kind.

func (m *Writer) findField(domain, kind, field string) (ret int, err error) {
	if _, kid, e := m.findKind(domain, kind); e != nil {
		err = e
	} else if e := m.db.QueryRow(`
		select rowid
		from mdl_field mf
		where kind = ?1
		and field = ?2`, kid, field).Scan(&ret); e == sql.ErrNoRows {
		err = errutil.Fmt("no such field %q in kind %q in domain %q", field, kind, domain)
	} else {
		err = e
	}
	return
}
