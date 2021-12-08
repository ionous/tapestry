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

// check that the kind can store the requested value at the passed field
// returns the name of the field ( in case the originally specified field was a trait )
func (k *ScopedKind) findCompatibleValue(field string, value literal.LiteralValue) (ret string, err error) {
	if value.Affinity() == affine.Bool {
		if aspect, ok := k.findCompatibleTrait(field); ok {
			ret = aspect
		} else {
			err = errutil.New("trait not found '%s.%s'", k.name, field)
		}
	} else {
		if ok, e := k.findCompatibleField(field, value); e != nil {
			err = e
		} else if ok {
			ret = field
		} else {
			err = errutil.New("field not found '%s.%s'", k.name, field)
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
				err = errutil.New("value of affinity %s incompatible with '%s.%s:%s'",
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
