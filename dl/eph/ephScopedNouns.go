package eph

import "github.com/ionous/errutil"

// name of a noun to assembly info
// ready after phase Ancestry
//
// fix? this is an exact copy of scopeKinds -- but its difficult to share.
type ScopedNouns map[string]*ScopedNoun

// return the uniformly named domain ( if it exists )
// fix? move to table? ( ex. search for row "d", see if it contains "name" )
func (d *Domain) GetNoun(name string) (ret *ScopedNoun, okay bool) {
	if n, ok := d.nouns[name]; ok {
		ret, okay = n, true
	} else if deps, e := d.GetDependencies(); e != nil {
		// if not in this domain, then maybe in a parent domain....
		// ( dont force resolve here, if its not resolved... then stop trying )
		LogWarning(e)
	} else {
		list := deps.Ancestors()
		for i, cnt := 0, len(list); i < cnt; i++ {
			el := list[cnt-i-1].(*Domain)
			if n, ok := el.nouns[name]; ok {
				ret, okay = n, true
				break
			}
		}
	}
	return
}

// return the uniformly named domain ( creating it in this domain if necessary )
func (d *Domain) EnsureNoun(name, at string) (ret *ScopedNoun) {
	if n, ok := d.GetNoun(name); ok {
		ret = n
	} else {
		n = &ScopedNoun{Requires: Requires{name: name, at: at}, domain: d}
		if d.nouns == nil {
			d.nouns = map[string]*ScopedNoun{name: n}
		} else {
			d.nouns[name] = n
		}
		ret = n
	}
	return
}

// distill a tree of nouns into a set of names and their hierarchy
func (d *Domain) ResolveNouns() (ret DependencyTable, err error) {
	if _, e := d.ResolveKinds(); e != nil {
		err = errutil.Append(err, e)
	} else {
		ret, err = d.resolvedNouns.resolve(func() (ret DependencyTable, err error) {
			m := TableMaker(len(d.nouns))
			for _, n := range d.nouns {
				if parentName, ok := m.ResolveParent(n); ok {
					if e := d.AddDefinition(n.name, n.at, parentName); e != nil {
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
	return
}
