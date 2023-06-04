package weave

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"github.com/ionous/errutil"
)

type ScopedKind struct {
	Requires // references to ancestors ( at most it can have one direct parent )
	domain   *Domain
}

func (k *ScopedKind) HasParent(other kindsOf.Kinds) bool {
	return k.Requires.HasParent(other.String())
}

func (k *ScopedKind) HasAncestor(other kindsOf.Kinds) bool {
	return k.Requires.HasAncestor(other.String())
}

func (k *ScopedKind) Resolve() (ret Dependencies, err error) {
	if len(k.at) == 0 {
		err = KindError{k.name, errutil.New("never defined")}
	} else if ks, e := k.resolve(k, (*kindFinder)(k.domain)); e != nil {
		err = KindError{k.name, e}
	} else {
		ret = ks
	}
	return
}

func (k *ScopedKind) findCompatibleField(field string, affinity affine.Affinity) (retName, retCls string, err error) {
	w := k.domain.catalog.writer
	return w.FindCompatibleField(k.domain.name, k.name, field, affinity)
}

// private helper to make the catalog compatible with the DependencyFinder ( for domains )
type kindFinder Domain

// look upwards through the domains to find the named kind
func (kf *kindFinder) FindDependency(name string) (ret Dependency, okay bool) {
	domain := (*Domain)(kf)
	if k, ok := domain.GetPluralKind(name); ok {
		ret, okay = k, true
	}
	return
}
