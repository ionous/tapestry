package eph

import (
	"sort"
	"strings"

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

func (t *tableMaker) onerror(e error) {
	t.err = errutil.Append(t.err, e)
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
func (ds DependencyTable) WriteTable(w Writer, target string, fullTree bool) (err error) {
	for _, d := range ds {
		var list []Dependency
		if fullTree {
			list = d.Ancestors()
		} else {
			list = d.Parents()
		}
		var b strings.Builder
		for i, cnt := 0, len(list); i < cnt; i++ {
			el := list[cnt-i-1]
			if i > 0 {
				b.WriteRune(',')
			}
			b.WriteString(el.Name())
		}
		name, row := d.Leaf().Name(), b.String()
		if e := w.Write(target, name, row); e != nil {
			err = errutil.Append(err, errutil.Fmt("couldn't write %q %w", name, e))
		}
	}
	return
}
