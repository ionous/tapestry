package eph

import "github.com/ionous/errutil"

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
	} else if k.domain == nil {
		err = KindError{k.name, errutil.New("no domain")}
	} else {
		ret, err = k.reqs.Resolve(k, (*kindFinder)(k.domain))
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
			if e := field.CheckConflict(kind); e != nil {
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
