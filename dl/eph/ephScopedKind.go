package eph

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"github.com/ionous/errutil"
)

type ScopedKind struct {
	Requires      // references to ancestors ( at most it can have one direct parent )
	domain        *Domain
	aspects       []traitDef  // used for kinds of aspects *and* for fields which use those aspects.
	fields        []*fieldDef // otherwise, everything is a field
	patternHeader patternHeader
	pendingFields []UniformField
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

// utility to check if this kind has the named trait
func (k *ScopedKind) FindTrait(name string) (ret traitDef, okay bool) {
	for _, a := range k.aspects {
		if a.HasTrait(name) {
			ret, okay = a, true
			break
		}
	}
	return
}

// check that the kind can store the requested value at the passed field
// returns the name of the field ( in case the originally specified field was a trait )
func (k *ScopedKind) findCompatibleField(field string, aff affine.Affinity) (ret *fieldDef, err error) {
	if aff != affine.Bool {
		if which, e := k.findScopedField(field, aff); e != nil {
			err = e
		} else if which == nil {
			err = errutil.Fmt("field not found '%s.%s'", k.name, field)
		} else {
			ret = which
		}
	} else {
		if aspect, e := k.findScopedTrait(field); e != nil {
			err = e
		} else if aspect == nil {
			err = errutil.Fmt("field not found '%s.%s'", k.name, field)
		} else {
			ret = aspect
		}
	}
	return
}

// returns true if the named field was found in this kind
// returns an error if it encountered a critical error along the way.
func (k *ScopedKind) findScopedField(field string, aff affine.Affinity) (ret *fieldDef, err error) {
	if e := VisitTree(k, func(dep Dependency) (err error) {
		// search through hierarchy's fields
		search := dep.(*ScopedKind)
		for _, def := range search.fields {
			if def.name == field {
				if def.affinity == aff.String() {
					ret, err = def, Visited
				} else {
					err = errutil.Fmt("affinity %s incompatible with '%s.%s:%s'",
						aff, search.name, field, def.affinity)
				}
				break
			}
		}
		return
	}); e != Visited {
		err = e
	}
	return
}

// returns the name of the aspect if the trait was found in this kind, or the empty string if not found.
// returns an error if it encountered a critical error along the way.
func (k *ScopedKind) findScopedTrait(field string) (ret *fieldDef, err error) {
	if e := VisitTree(k, func(dep Dependency) (err error) {
		// search through hierarchy's fields
		search := dep.(*ScopedKind)
		for _, def := range search.fields {
			// for fields that (probably) refer to a kind
			if cls := def.class; len(cls) > 0 {
				// see if that kind is an aspect
				if a, ok := k.domain.GetKind(cls); ok && a.HasParent(kindsOf.Aspect) {
					// when the name of the field is the same as the name of the aspect
					// that is our special "acts as trait" field.
					if a.name == def.name && def.affinity == affine.Text.String() {
						// and search through its traits
						if _, ok := a.FindTrait(field); ok {
							ret, err = def, Visited // to exit hierarchy search
							break                   // to exit the field search
						}
					}
				}
			}
		}
		return
	}); e != Visited {
		err = e
	}
	return
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
