package eph

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"github.com/ionous/errutil"
)

type ScopedKind struct {
	Requires // references to ancestors ( at most it can have one direct parent )
	domain   *Domain
	traits   []traitDef
	fields   []fieldDef
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

// the kind must have been resolved for this to work
func (k *ScopedKind) AddField(field FieldDefinition) (err error) {
	if e := VisitTree(k, func(dep Dependency) (err error) {
		kind := dep.(*ScopedKind)
		if e := field.CheckConflict(kind); e != nil {
			err = DomainError{kind.domain.name, KindError{kind.name, e}}
		}
		return
	}); e != nil {
		err = e
	} else {
		field.AddToKind(k) // if everything worked out store definition locally
	}
	return
}

// check that the kind can store the requested value at the passed field
// returns the name of the field ( in case the originally specified field was a trait )
func (k *ScopedKind) findCompatibleValue(field string, value literal.LiteralValue) (ret string, err error) {
	if value.Affinity() == affine.Bool {
		if e := VisitTree(k, func(dep Dependency) (err error) {
			kind := dep.(*ScopedKind)
			if aspect, ok := kind.findCompatibleTrait(field); ok {
				ret, err = aspect, Visited
			}
			return
		}); e == nil {
			err = errutil.Fmt("field not found '%s.%s'", k.name, field)
		} else if e != Visited {
			err = e
		}
	} else {
		if e := VisitTree(k, func(dep Dependency) (err error) {
			kind := dep.(*ScopedKind)
			if ok, e := kind.findCompatibleField(field, value); e != nil {
				err = e
			} else if ok {
				ret, err = field, Visited
			}
			return
		}); e == nil {
			err = errutil.Fmt("trait not found '%s.%s'", k.name, field)
		} else if e != Visited {
			err = e
		}
	}
	return
}

func (k *ScopedKind) findCompatibleField(field string, value literal.LiteralValue) (okay bool, err error) {
	for _, def := range k.fields {
		if def.name == field {
			if aff := value.Affinity(); def.affinity == aff.String() {
				okay = true
			} else {
				err = errutil.Fmt("value of affinity %s incompatible with '%s.%s:%s'",
					aff, k.name, field, def.affinity)
			}
			break
		}
	}
	return
}

func (k *ScopedKind) findCompatibleTrait(field string) (ret string, okay bool) {
	for _, def := range k.traits {
		// the names of traits of that aspect
		for _, trait := range def.traits {
			if trait == field {
				ret, okay = def.aspect, true
				break
			}
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
