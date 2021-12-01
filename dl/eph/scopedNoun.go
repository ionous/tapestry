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
