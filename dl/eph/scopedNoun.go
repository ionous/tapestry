package eph

import (
	"github.com/ionous/errutil"
)

type ScopedNoun struct {
	Requires // kinds ( when resolved, can have one direct parent )
	domain   *Domain
}

func (n *ScopedNoun) Resolve() (ret Dependencies, err error) {
	if len(n.at) == 0 {
		err = NounError{n.name, errutil.New("never defined")}
	} else if ks, e := n.resolve(n, (*kindFinder)(n.domain)); e != nil {
		err = NounError{n.name, e}
	} else {
		ret = ks
	}
	return
}

func (n *ScopedNoun) Kind() (ret *ScopedKind, err error) {
	if dep, e := n.GetDependencies(); e != nil {
		err = e
	} else if ks := dep.Parents(); len(ks) != 1 {
		err = errutil.Fmt("noun %q has unexpected %d parents", n.name, len(ks))
	} else {
		ret = ks[0].(*ScopedKind)
	}
	return
}
