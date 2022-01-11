package eph

import "github.com/ionous/errutil"

// name of a kind to assembly info
// ready after phase Ancestry
type ScopedKinds map[string]*ScopedKind

func (d *Domain) GetPluralKind(name string) (ret *ScopedKind, okay bool) {
	if a, ok := d.GetKind(name); ok {
		ret, okay = a, true
	} else if p, e := d.Pluralize(name); e == nil {
		ret, okay = d.GetKind(p)
	}
	return
}

// return the uniformly named domain ( if it exists )
func (d *Domain) GetKind(name string) (ret *ScopedKind, okay bool) {
	if e := VisitTree(d, func(dep Dependency) (err error) {
		if n, ok := dep.(*Domain).kinds[name]; ok {
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
	if k, ok := d.GetKind(name); ok {
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

// distill a tree of kinds into a set of names and their hierarchy
func (d *Domain) ResolveKinds() (DependencyTable, error) {
	return d.resolvedKinds.resolve(func() (ret DependencyTable, err error) {
		m := TableMaker(len(d.kinds))
		for _, k := range d.kinds {
			if parentName, ok := m.ResolveParent(k); ok {
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
