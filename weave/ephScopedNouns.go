package weave

import (
	"math"

	"github.com/ionous/errutil"
)

// name of a noun to assembly info
// ready after phase Ancestry
//
// fix? this is an exact copy of scopeKinds -- but its difficult to share.
type ScopedNouns map[string]*ScopedNoun

// find the noun with the ( exact ) name in this scope
func (d *Domain) GetNoun(name string) (ret *ScopedNoun, okay bool) {
	if e := d.visit(func(scope *Domain) (err error) {
		if n, ok := scope.nouns[name]; ok {
			ret, okay, err = n, true, Visited
		}
		return
	}); e != nil && e != Visited {
		LogWarning(e)
	}
	return
}

// find the noun with the closest name in this scope
func (d *Domain) GetClosestNoun(name string) (ret *ScopedNoun, okay bool) {
	bestRank, bestNoun := math.MaxInt, ret
	if e := d.visit(func(scope *Domain) (err error) {
		// used the resolved nouns to generate a consistent ordering
		if nouns, e := scope.resolvedNouns.GetTable(); e != nil {
			err = e
		} else {
			for _, el := range nouns {
				noun := el.Leaf().(*ScopedNoun)
				names := noun.Names()
				for i, cnt := 0, len(names); i < cnt && i < bestRank; i++ {
					if name == names[i] {
						bestRank, bestNoun = i, noun
						// cant do better than the best
						if i == 0 {
							err = Visited
							break
						}
					}
				}
			}
		}
		return
	}); e != nil && e != Visited {
		LogWarning(e)
	} else if bestRank < math.MaxInt {
		ret, okay = bestNoun, true
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

// distill the nouns from this domain into a table sorted by kind.
func (d *Domain) ResolveDomainNouns() (ret DependencyTable, err error) {
	if _, e := d.resolveKinds(); e != nil {
		err = errutil.Append(err, e)
	} else {
		ret, err = d.resolvedNouns.resolve(func() (ret DependencyTable, err error) {
			m := TableMaker(len(d.nouns))
			for _, n := range d.nouns {
				if parentName, ok := m.ResolveParent(n); ok {
					if e := d.AddDefinition(MakeKey("nouns", n.name), n.at, parentName); e != nil {
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
