package weave

import (
	"sort"

	"github.com/ionous/errutil"
)

// given a name, return an object which describes the other names on which it depends.
type DependencyFinder interface {
	FindDependency(name string) (Dependency, bool)
}

type DependencyTable []Dependencies

type tableMaker struct {
	deps []Dependencies
	err  error
}

// a helper for creating a DependencyTable regardless of whether its Dependencies are stored in slices or maps.
func TableMaker(cnt int) tableMaker {
	return tableMaker{make([]Dependencies, 0, cnt), nil}
}

func (t *tableMaker) GetSortedTable() (ret DependencyTable, err error) {
	if e := t.err; e != nil {
		err = e
	} else {
		ret = t.deps
		ret.SortTable()
	}
	return
}

// accumulates errors for reporting via GetSortedTable, but returns true/false in the meantime
func (t *tableMaker) ResolveDep(dep Dependency) (ret Dependencies, okay bool) {
	if res, e := dep.Resolve(); e != nil {
		t.onerror(e)
	} else {
		t.deps = append(t.deps, res)
		ret, okay = res, true
	}
	return
}

// accumulates errors for reporting via GetSortedTable, but returns true/false in the meantime
func (t *tableMaker) ResolveReq(finder DependencyFinder, req string) (ret Dependencies, okay bool) {
	if finder == nil {
		t.onerror(errutil.New("unknown dependencies"))
	} else if dep, ok := finder.FindDependency(req); !ok {
		t.onerror(errutil.Fmt("unknown dependency %q", req))
	} else {
		ret, okay = t.ResolveDep(dep)
	}
	return
}

// same as resolve dep -- but adds an error if there's more than one parent
func (t *tableMaker) ResolveParent(dep Dependency) (ret string, okay bool) {
	if res, ok := t.ResolveDep(dep); ok {
		switch ps := res.Parents(); len(ps) {
		case 0:
			okay = true
		case 1:
			ret, okay = ps[0].Name(), true
		default:
			t.onerror(errutil.Fmt("%q has more than one parent", dep.Name()))
		}
	}
	return
}

func (t *tableMaker) onerror(e error) {
	t.err = errutil.Append(t.err, e)
}

type cachedTable struct {
	res    DependencyTable
	status error
}

// return previously resolved dependencies
func (t *cachedTable) GetTable() (ret DependencyTable, err error) {
	if e := t.status; e == nil {
		err = errutil.New("dependency table not resolved")
	} else if e != xResolved {
		err = e
	} else {
		ret = t.res // okay
	}
	return
}

func (t *cachedTable) resolve(do func() (DependencyTable, error)) (ret DependencyTable, err error) {
	switch t.status {
	case xResolved:
		ret = t.res
	case nil:
		t.status = xProcessing
		if res, e := do(); e != nil {
			err, t.status = e, e
		} else {
			ret, t.res, t.status = res, res, xResolved
		}
	default:
		err = t.status
	}
	return
}

// build a list of just the "column names" -- the resolved objects.
func (ds DependencyTable) Names() []string {
	out := make([]string, len(ds))
	for i, d := range ds {
		out[i] = d.Leaf().Name()
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
				case a.Name() < b.Name():
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
// func (ds DependencyTable) WriteTable(m mdl.Modeler, target string, fullTree bool) (err error) {
// 	for _, dep := range ds {
// 		name, row, at := dep.Leaf().Name(), dep.Strings(fullTree), dep.Leaf().OriginAt()
// 		if e := w.Write(target, name, row, at); e != nil {
// 			err = e
// 			break
// 		}
// 	}
// 	return
// }