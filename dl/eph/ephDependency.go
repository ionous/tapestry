package eph

import (
	"strings"

	"github.com/ionous/errutil"
)

type Dependency interface {
	// semi-unique name for the dependency ( uniqueness depends on type and scope of declaration )
	Name() string
	// location of the first found declaration which generated this dependency
	OriginAt() string
	// make this dependency depend on the passed name
	// ( and clear any previously resolved dependencies )
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

func (d *Dependencies) Strings(fullTree bool) string {
	var b strings.Builder
	var list []Dependency
	if fullTree {
		list = d.Ancestors()
	} else {
		list = d.Parents()
	}
	for i, cnt := 0, len(list); i < cnt; i++ {
		el := list[cnt-i-1]
		if i > 0 {
			b.WriteRune(',')
		}
		b.WriteString(el.Name())
	}
	return b.String()
}

const Visited = errutil.Error("git.sr.ht/~ionous/iffy/dl/eph/Visited")

func VisitTree(d Dependency, visit func(Dependency) error) (err error) {
	if deps, e := d.GetDependencies(); e != nil {
		err = e
	} else if ancestors := deps.ancestors; len(ancestors) == 0 {
		err = visit(d) // ugh - for testing where things may not actually have valid dependencies.
	} else {
		for i := len(ancestors) - 1; i >= 0; i-- {
			if e := visit(ancestors[i]); e != nil {
				err = e
				break
			}
		}
	}
	return
}
