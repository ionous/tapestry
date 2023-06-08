package weave

import (
	"database/sql"
	"strings"

	"github.com/ionous/errutil"
)

func (d *Domain) findKind(name string) (okay bool) {
	_, _, _, e := d.pathOfKind(name)
	return e == nil
}

func (d *Domain) findPluralKind(name string) (ret string, okay bool) {
	if ok := d.findKind(name); ok {
		ret, okay = name, true
	} else if p := d.catalog.run.PluralOf(name); p != name && d.findKind(p) {
		ret, okay = p, true
	}
	return
}

func (d *Domain) hasAncestor(name, parent string) (okay bool) {
	if _, fulltree, _, e := d.pathOfKind(name); e == nil {
		if _, uppertree, _, e := d.pathOfKind(parent); e == nil {
			okay = strings.HasSuffix(fulltree, uppertree)
		}
	}
	return
}

// fix: duplicated in mdlVerify
func (d *Domain) pathOfKind(kind string) (retDomain, retPath string, retKind int, err error) {
	q := d.catalog.db
	if e := q.QueryRow(`
	select domain, mk.rowid, ',' || mk.rowid || ',' || mk.path
	from mdl_kind mk
	join domain_tree
		on (uses = domain)
	where base = ?1
	and kind = ?2
	limit 1`, d.name, kind).Scan(&retDomain, &retKind, &retPath); e == sql.ErrNoRows {
		err = errutil.Fmt("no such kind %q in domain %q", kind, d.name)
	} else {
		err = e
	}
	return
}
