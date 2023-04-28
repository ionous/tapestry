package weave

import "github.com/ionous/errutil"

// generator for a dependency graph
// designed to be embedded in a map by pointer or in embedded in some other object to store dependencies about that object.
// partially implements "Dependency" ( missing the Resolve() method )
type Requires struct {
	name, at string
	reqs     UniqueNames  // original list of dependencies
	resolved Dependencies // valid after Resolve()
	status   error        // nil status means "unresolved"
}

// dependency status markers
const (
	xProcessing = errutil.Error("processing") // helper to detect circular references during Resolve()
	xResolved   = errutil.Error("resolved")   // marks a successfully completed Resolve()
	// Resolved = nil -- except go doesnt allow nil const
)

// implements Dependency
func (d *Requires) Name() string {
	return d.name
}

func (d *Requires) OriginAt() string {
	return d.at
}

// make the name or object this set of dependencies represents require the passed dep
// clears any previous cached resolution data or internal errors
func (d *Requires) AddRequirement(name string) {
	if d.reqs.AddName(name) >= 0 && d.status != nil {
		// println("clearing dependencies for", d.name, "while adding", name)
		d.status = nil // not seen before, then clear any cache.
	}
}

// return previously resolved dependencies
func (d *Requires) GetDependencies() (ret Dependencies, err error) {
	if e := d.status; e == nil {
		err = errutil.New(d.name, "dependencies not resolved")
	} else if e != xResolved {
		err = e
	} else {
		ret = d.resolved // okay
	}
	return
}

// ( must be previously resolved to work properly )
func (d *Requires) HasParent(name string) (okay bool) {
	if dep, e := d.GetDependencies(); e == nil {
		if as := dep.Parents(); len(as) > 0 && as[0].Name() == name {
			okay = true
		}
	}
	return
}

// ( must be previously resolved to work properly )
func (d *Requires) HasAncestor(name string) (okay bool) {
	if dep, e := d.GetDependencies(); e == nil {
		if as := dep.Ancestors(); len(name) == 0 && len(as) == 0 {
			okay = true // if an empty parent is required and there are no parents
		} else {
			// otherwise... make sure whatever kind the child domain is specifying lines up
			for _, a := range as {
				if a.Name() == name {
					okay = true
					break
				}
			}
		}
	}
	return
}

// return the graph of all dependencies ( recursively creating that graph when needed. )
func (d *Requires) resolve(node Dependency, names DependencyFinder) (ret Dependencies, err error) {
	switch d.status {
	case xResolved: // already resolved? return the list.
		ret = d.resolved

	case xProcessing:
		e := errutil.New("circular reference detected:", node.Name())
		err, d.status = e, e

	case nil: // unresolved? resolved.
		d.status = xProcessing
		if res, e := resolve(node, d.reqs, names); e != nil {
			err, d.status = e, e
		} else {
			ret, d.resolved, d.status = res, res, xResolved
			// println("resolved", d.name)
		}

	default: // otherwise return the cached error.
		err = d.status
	}
	return
}

func resolve(node Dependency, reqs []string, names DependencyFinder) (ret Dependencies, err error) {
	// capital-R resolve each specified dependency
	// ( which winds us up back here -- one requirement deeper -- unless its results are cached already )
	m := TableMaker(len(reqs))
	for i, req := range reqs {
		// accumulates any errors
		if d, ok := m.ResolveReq(names, req); ok {
			reqs[i] = d.Leaf().Name() // we rewrite the names to handle plural lookups
		}
	}
	if ds, e := m.GetSortedTable(); e != nil {
		err = e
	} else {
		var parents, ancestors []Dependency
		if len(ds) > 0 {
			// merge dependencies in decreasing dependency order ( longest first )
			for i, cnt := 0, len(ds)-1; i <= cnt; i++ {
				// skip this dependency if we have its chain already ( ie. from a previous longer requirement )
				dep := ds[cnt-i] // dep is a set of Dependencies
				if req := dep.Leaf(); findStr(ancestors, req, 0) < 0 {
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
		ret = Dependencies{
			ancestors: append(ancestors, node),
			parents:   parents,
		}
	}
	return
}

// linear, unsorted search
func findStr(a []Dependency, dep Dependency, from int) (ret int) {
	ret = -1 // provisionally
	for i, cnt := from, len(a); i < cnt; i++ {
		if d := a[i]; d == dep {
			ret = i
			break
		}
	}
	return
}
