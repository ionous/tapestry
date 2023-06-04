package weave

import (
	"database/sql"
	"strings"

	"github.com/ionous/errutil"
)

// name of a kind to assembly info
// ready after phase Ancestry
type ScopedKinds map[string]*ScopedKind

// return the uniformly named domain ( if it exists )
func (d *Domain) getKind(name string) (ret *ScopedKind, okay bool) {
	if e := d.visit(func(dep *Domain) (err error) {
		if n, ok := dep.kinds[name]; ok {
			ret, okay, err = n, true, Visited
		}
		return
	}); e != nil && e != Visited {
		LogWarning(e)
	}
	return
}

// return the uniformly named domain ( creating it in this domain if necessary )
func (d *Domain) EnsureKind(name, at string) (ret *ScopedKind) {
	if k, ok := d.getKind(name); ok {
		ret = k
	} else {
		k = &ScopedKind{Requires: Requires{name: name, at: at}, domain: d}
		if d.kinds == nil {
			d.kinds = map[string]*ScopedKind{name: k}
		} else {
			d.kinds[name] = k
		}
		ret = k
	}
	return
}

// distill a set of kinds into a set of names and their hierarchy
func (d *Domain) resolveKinds() (DependencyTable, error) {
	return d.resolvedKinds.resolve(func() (ret DependencyTable, err error) {
		m := TableMaker(len(d.kinds))
		for _, k := range d.kinds {
			if parentName, ok := m.ResolveParent(k); ok {
				// FIX: USE TABLE CONFLICTS INSTEAD.
				if e := d.AddDefinition(MakeKey("kinds", k.name), k.at, parentName); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
		if dt, e := m.GetSortedTable(); e != nil {
			err = errutil.Append(err, e)
		} else {
			ret = dt
		}
		return
	})
}

// -------------

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
