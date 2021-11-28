package eph

type Dependency interface {
	Name() string
	// make this dependency depend on the passed name
	// ( clear any cached requirements )
	AddRequirement(name string)
	// change all requirements into a filtered sorted list of parents and ancestors
	Resolve() (Dependencies, error)
	// return any previously resolved dependencies
	GetDependencies() (Dependencies, error)
}

// contains all dependencies of dependencies and all dependencies not listed in another dependency.
type Dependencies struct {
	ancestors []Dependency // sorted root/s first, leaf last.
	parents   []Dependency
}

// list of direct parents, not in any particular order.
func (d *Dependencies) Parents() []Dependency {
	return d.parents
}

// complete list of all ancestors, sorted root/s first, direct parents last.
// the parents may or may not be contiguous ( depending on whether their ancestors overlap. )
func (d *Dependencies) Ancestors() (ret []Dependency) {
	if cnt := len(d.ancestors); cnt > 0 {
		ret = d.ancestors[:cnt-1]
	}
	return
}

// the ancestors + the leaf ( the leaf appears last )
// when created via the dependency interface should return an array of at least one value: the leaf.
func (d *Dependencies) FullTree() []Dependency {
	return d.ancestors
}

// the dependency for which this set of dependencies was generated.
// when created via the dependency interface should return a valid dependency;
// the zero value returns nil.
func (d *Dependencies) Leaf() (ret Dependency) {
	if cnt := len(d.ancestors); cnt > 0 {
		ret = d.ancestors[cnt-1]
	}
	return
}
