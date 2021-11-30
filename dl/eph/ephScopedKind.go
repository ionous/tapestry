package eph

import (
	"github.com/ionous/errutil"
)

type ScopedKind struct {
	name, at string
	domain   *Domain
	reqs     Requires // references to ancestors ( at most it can have one direct parent )
	traits   []traitDef
	fields   []fieldDef
}

// implement the Dependency interface
func (k *ScopedKind) Name() string                           { return k.name }
func (k *ScopedKind) AddRequirement(name string)             { k.reqs.AddRequirement(name) }
func (k *ScopedKind) GetDependencies() (Dependencies, error) { return k.reqs.GetDependencies() }

func (k *ScopedKind) Resolve() (ret Dependencies, err error) {
	if len(k.at) == 0 {
		err = KindError{k.name, errutil.New("never defined")}
	} else if ks, e := k.reqs.Resolve(k, (*kindFinder)(k.domain)); e != nil {
		err = KindError{k.name, e}
	} else {
		ret = ks
	}
	return
}

func (k *ScopedKind) HasAncestor(name string) (okay bool, err error) {
	return HasAncestor(k, name)
}

// the kind must have been resolved for this to work
func (k *ScopedKind) AddField(field FieldDefinition) (err error) {
	if deps, e := k.GetDependencies(); e != nil {
		err = e
	} else {
		// the full tree includes the kind itself; its a bit weird, but it keeps this loop simple.
		for _, dep := range deps.FullTree() {
			kind := dep.(*ScopedKind)
			if e := field.CheckConflict(kind); e != nil {
				err = DomainError{kind.domain.name, KindError{kind.name, e}}
				break
			}
		}
		// if everything worked out store definition locally
		if err == nil {
			field.AddToKind(k)
		}
	}
	return
}

// private helper to make the catalog compatible with the DependencyFinder ( for domains )
type kindFinder Domain

// look upwards through the domains to find the named kind
func (kf *kindFinder) FindDependency(name string) (ret Dependency, okay bool) {
	domain := (*Domain)(kf)
	if k, ok := domain.GetKind(name); ok {
		ret, okay = k, true
	}
	return
}
