package eph

import "github.com/ionous/errutil"

// given a name, return an object which describes the other names on which it depends.
type DependencyFinder interface {
	GetRequirements(name string) (*Requires, bool)
}

// solves the dependencies of the passed name using the lookup
func GetResolvedDependencies(name string, names DependencyFinder) (ret Dependents, err error) {
	if dep, ok := names.GetRequirements(name); !ok {
		err = errutil.New("Unknown dependency", name)
	} else {
		ret, err = dep.Resolve(name, names)
	}
	return
}

// generator for a dependency graph
// designed to be embedded in a map by pointer or in embedded in some other object to store dependencies about that object.
type Requires struct {
	reqs     UniqueNames // original list of dependencies
	resolved Dependents  // valid after Resolve()
	status   error       // nil status means "unresolved"
}

// dependency status markers
const (
	xProcessing = errutil.Error("processing") // helper to detect circular references during Resolve()
	xResolved   = errutil.Error("resolved")   // marks a successfully completed Resolve()
	// Resolved = nil -- except go doesnt allow nil const
)

// make the name or object this set of dependencies represents require the passed dep
// clears any previous cached resolution data or internal errors
func (d *Requires) AddRequirement(dep string) {
	if d.reqs.AddName(dep) {
		d.status = nil
	}
}

// return previously resolved dependencies
func (d *Requires) GetDependencies() (ret Dependents, err error) {
	if e := d.status; e == nil {
		err = errutil.New("dependencies not resolved")
	} else if e != xResolved {
		err = e
	} else {
		ret = d.resolved // okay
	}
	return
}

// return the graph of all dependencies ( recursively creating that graph when needed. )
func (d *Requires) Resolve(name string, names DependencyFinder) (ret Dependents, err error) {
	switch d.status {
	case xResolved: // already resolved? return the list.
		ret = d.resolved

	case xProcessing:
		e := errutil.New("circular reference detected:", name)
		err, d.status = e, e

	case nil: // unresolved? resolved.
		d.status = xProcessing
		if res, e := resolve(name, d.reqs, names); e != nil {
			err, d.status = e, e
		} else {
			ret, d.resolved, d.status = res, res, xResolved
		}

	default: // otherwise return the cached error.
		err = d.status
	}
	return
}

func resolve(name string, reqs []string, names DependencyFinder) (ret Dependents, err error) {
	// capital-R resolve each specified dependency
	if ds, e := MakeTable(reqs, names); e != nil {
		err = e
	} else {
		var parents, ancestors []string
		if len(ds) > 0 {
			ds.SortTable()
			// merge dependencies in decreasing dependency order ( longest first )
			for i, cnt := 0, len(ds)-1; i <= cnt; i++ {
				// skip this dependency if we have its chain already ( ie. from a previous longer requirement )
				dep := ds[cnt-i]
				if req := dep.Name(); findStr(ancestors, req, 0) < 0 {
					parents = append(parents, req)
					prev, searchFrom := ancestors, 0 // re: prev, only search the elements we already had.
					for _, a := range dep.ancestors {
						if at := findStr(prev, a, searchFrom); at >= 0 {
							searchFrom = at
						} else {
							ancestors = append(ancestors, a)
						}
					}
				}
			}
		}
		// return...
		ret = Dependents{
			ancestors: append(ancestors, name),
			parents:   parents,
		}
	}
	return
}

// linear, unsorted search
func findStr(a []string, name string, from int) (ret int) {
	ret = -1 // provisionally
	for i, cnt := from, len(a); i < cnt; i++ {
		if n := a[i]; n == name {
			ret = i
			break
		}
	}
	return
}
