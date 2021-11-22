package eph

import "github.com/ionous/errutil"

// given a name, return an object which describes the other names on which it depends.
type DependencyFinder interface {
	GetDependencies(name string) (*Dependencies, bool)
}

// solves the dependencies of the passed name using the lookup
func GetResolvedDependencies(name string, names DependencyFinder) (ret ResolvedDependencies, err error) {
	if dep, ok := names.GetDependencies(name); !ok {
		err = errutil.New("Unknown dependency", name)
	} else {
		ret, err = dep.Resolve(name, names)
	}
	return
}

// contains all dependencies of dependencies and all dependencies not listed in another dependency.
type ResolvedDependencies struct {
	fullTree []string // direct parents are at slice's start; ancestors ( and the root ) at the end.
	parents  int
}

func (d *ResolvedDependencies) NumParents() int   { return d.parents }
func (d *ResolvedDependencies) NumAncestors() int { return len(d.fullTree) }

// if not the full tree returns just the parents.
// direct parents are at slice's start; ancestors ( and the root ) at the end.
func (d *ResolvedDependencies) Ancestors(fullTree bool) (ret []string) {
	if fullTree {
		ret = d.fullTree
	} else {
		ret = d.fullTree[:d.parents]
	}
	return
}

// generator for a dependency graph
// designed to be embedded in a map by pointer or in embedded in some other object to store dependencies about that object.
type Dependencies struct {
	deps     UniqueNames          // original list of dependencies
	resolved ResolvedDependencies // valid after Resolve()
	status   error                // nil status means "unresolved"
}

// dependency status markers
const (
	xProcessing = errutil.Error("processing") // helper to detect circular references during Resolve()
	xResolved   = errutil.Error("resolved")   // marks a successfully completed Resolve()
	// Resolved = nil -- except go doesnt allow nil const
)

// make the name or object this set of dependencies represents require the passed dep
// clears any previous cached resolution data or internal errors
func (d *Dependencies) AddDependency(dep string) {
	if d.deps.AddName(dep) {
		d.status = nil
	}
}

// return previously resolved dependencies
func (d *Dependencies) GetDependencies() (ret ResolvedDependencies, err error) {
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
func (d *Dependencies) Resolve(name string, names DependencyFinder) (ret ResolvedDependencies, err error) {
	switch d.status {
	default:
		err = d.status

	case xResolved:
		// already resolved?
		ret = d.resolved

	case xProcessing:
		err = errutil.New("Circular reference detected:", name)
		d.status = err

	case nil: // Unresolved
		d.status = xProcessing
		// resolved dependencies will go in "all" and either "subs" or "parents" depending
		var all UniqueNames
		var subs, parents []string
		for _, depName := range d.deps {
			if subdeps, e := GetResolvedDependencies(depName, names); e != nil {
				err = errutil.New(e, "->", name)
				break
			} else {
				// add the resolved dependencies of our dependency
				for _, sub := range subdeps.Ancestors(true) {
					if all.AddName(sub) {
						subs = append(subs, sub)
					}
				}
				// add the dependency itself
				// and if it didn't appear as a dependency of a dependency
				// add it to our exclusive "parent" list
				if all.AddName(depName) {
					parents = append(parents, depName)
				}
			}
		}
		if err != nil {
			d.status = err
		} else {
			d.resolved = ResolvedDependencies{
				parents:  len(parents),
				fullTree: append(parents, subs...),
			}
			d.status = xResolved
			ret = d.resolved
		}
	}
	return
}
