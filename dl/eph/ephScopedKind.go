package eph

import (
	"errors"

	"github.com/ionous/errutil"
)

type ScopedKind struct {
	name, at string
	domain   *Domain
	reqs     Requires // references to ancestors ( at most it can have one direct parent )
	traits   []traitDef
	fields   []fieldDef
}

// name of a kind to assembly info
// ready after phase Ancestry
type ScopedKinds map[string]*ScopedKind

// implement the Dependency interface
func (k *ScopedKind) Name() string                           { return k.name }
func (k *ScopedKind) AddRequirement(name string)             { k.reqs.AddRequirement(name) }
func (d *ScopedKind) GetDependencies() (Dependencies, error) { return d.reqs.GetDependencies() }

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
	if dep, e := k.GetDependencies(); e != nil {
		err = e
	} else if as := dep.Ancestors(); len(name) == 0 && len(as) == 0 {
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
	return
}

// the kind must have been resolved for this to work
func (k *ScopedKind) AddFields(field FieldDefinition) (err error) {
	if deps, e := k.GetDependencies(); e != nil {
		err = e
	} else {
		// the full tree includes the kind itself; its a bit weird, but it keeps this loop simple.
		for _, dep := range deps.FullTree() {
			kind := dep.(*ScopedKind)
			var conflict *Conflict
			if e := field.CheckConflict(kind); errors.As(e, &conflict) && conflict.Reason == Duplicated {
				LogWarning(e) // warn if it was a duplicated definition
			} else if e != nil {
				err = KindError{kind.name, e}
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
