package eph

import "github.com/ionous/errutil"

// name of a kind to assembly info
// ready after phase Ancestry
type ScopedKinds map[string]*ScopedKind

// return the uniformly named domain ( if it exists )
func (d *Domain) GetKind(name string) (ret *ScopedKind, okay bool) {
	if k, ok := d.kinds[name]; ok {
		ret, okay = k, true
	} else if deps, e := d.GetDependencies(); e != nil {
		// if not in this domain, then maybe in a parent domain....
		// ( dont force resolve here, if its not resolved... then stop trying )
		LogWarning(e)
	} else {
		list := deps.Ancestors()
		for i, cnt := 0, len(list); i < cnt; i++ {
			el := list[cnt-i-1].(*Domain)
			if k, ok := el.kinds[name]; ok {
				ret, okay = k, true
				break
			}
		}
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
				if e := d.AddDefinition(k.name, k.at, parentName); e != nil {
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
