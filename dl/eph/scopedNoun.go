package eph

import (
	"github.com/ionous/errutil"
)

type ScopedNoun struct {
	name, at string
	domain   *Domain
	reqs     Requires // kinds ( when resolved, can have one direct parent )
}

// implement the Dependency interface
func (n *ScopedNoun) Name() string                           { return n.name }
func (n *ScopedNoun) AddRequirement(name string)             { n.reqs.AddRequirement(name) }
func (n *ScopedNoun) GetDependencies() (Dependencies, error) { return n.reqs.GetDependencies() }

func (n *ScopedNoun) Resolve() (ret Dependencies, err error) {
	if len(n.at) == 0 {
		err = NounError{n.name, errutil.New("never defined")}
	} else if ks, e := n.reqs.Resolve(n, (*kindFinder)(n.domain)); e != nil {
		err = NounError{n.name, e}
	} else {
		ret = ks
	}
	return
}

func (n *ScopedNoun) HasAncestor(name string) (okay bool, err error) {
	return HasAncestor(n, name)
}
