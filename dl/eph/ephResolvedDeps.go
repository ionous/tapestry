package eph

import (
	"sort"
	"strings"

	"github.com/ionous/errutil"
)

type DependencyTable []Dependents

// contains all dependencies of dependencies and all dependencies not listed in another dependency.
type Dependents struct {
	ancestors []string // sorted in root/s first order, name last
	parents   []string
}

func (d *Dependents) Name() string      { return d.ancestors[len(d.ancestors)-1] }
func (d *Dependents) NumParents() int   { return len(d.parents) }
func (d *Dependents) NumAncestors() int { return len(d.ancestors) - 1 }

// if not the full tree returns just the parents.
// direct parents are at slice's start; ancestors ( and the root ) at the end.
func (d *Dependents) Ancestors(fullTree bool) (ret []string) {
	if fullTree {
		ret = d.ancestors[:len(d.ancestors)-1]
	} else {
		ret = d.parents
	}
	return
}

func MakeTable(reqs []string, names DependencyFinder) (ret DependencyTable, err error) {
	var ds DependencyTable
	for _, req := range reqs {
		if res, e := GetResolvedDependencies(req, names); e != nil {
			err = errutil.Append(err, e)
		} else {
			ds = append(ds, res)
		}
	}
	if err == nil {
		ret = ds
	}
	return
}

// build a list of just the "column names" -- the resolved objects.
func (ds DependencyTable) Names() []string {
	out := make([]string, len(ds))
	for i, d := range ds {
		out[i] = d.Name()
	}
	return out
}

// ensure that dependencies of a name appear before the name is used as a dependency
// and that groups of related dependencies are otherwise kept together
func (ds DependencyTable) SortTable() {
	sort.Slice(ds, func(i, j int) (less bool) {
		is, js := ds[i].ancestors, ds[j].ancestors
		in, jn := len(is), len(js)
		iv, jv := is[in-1], js[jn-1]
		// if i is required by j, then i is most certainly less.
		if req := findStr(js, iv, 0); req >= 0 {
			less = true
		} else if req := findStr(is, jv, 0); req < 0 {
			// so long as the opposite isn't true, then:
			// start from the root looking towards the leaves to find different groupings.
		Match:
			for n := 0; ; {
				switch a, b := is[n], js[n]; {
				case a == b:
					// they are the same, so dig deeper:
					if n = n + 1; n >= in {
						less = true
						break Match // in is shorter, and therefore less
					} else if n >= jn {
						break Match // jn is shorter, and therefore not less
					}
					// ( keep looping. we will finish eventually. )
				case a < b:
					less = true
					break Match // a is lexically lesser
				default:
					break Match // b is lexically greater
				}
			}
		}
		//println("sort:", iv, jv, less)
		return
	})
}

// for each domain in the passed list, output its full ancestry tree ( or just its parents )
func (ds DependencyTable) WriteTable(w Writer, target string, fullTree bool) (err error) {
	for _, d := range ds {
		name, list := d.Name(), d.Ancestors(fullTree)
		var b strings.Builder
		for i, cnt := 0, len(list); i < cnt; i++ {
			el := list[cnt-i-1]
			if i > 0 {
				b.WriteRune(',')
			}
			b.WriteString(el)
		}
		row := b.String()
		if e := w.Write(target, name, row); e != nil {
			err = errutil.Append(err, errutil.Fmt("couldn't write %q %e", name, e))
		}
	}
	return
}
