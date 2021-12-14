package eph

import "github.com/ionous/errutil"

// implements FieldDefinition
// when the owner kind is the kind of aspect --
// the owner will only have one trait def
// and the name of the aspect will match the name of the kind
type traitDef struct {
	at     string
	aspect string
	traits []string
}

func (td *traitDef) Write(w Writer) error {
	return w.Write(mdl_field, td.aspect, Affinity_Text, td.aspect, td.at)
}

func (td *traitDef) AddToKind(k *ScopedKind) {
	k.aspects = append(k.aspects, *td)
}

func (td *traitDef) HasTrait(n string) (ret bool) {
	for _, trait := range td.traits {
		if n == trait {
			ret = true
			break
		}
	}
	return
}

func (td *traitDef) CheckConflict(k *ScopedKind) (err error) {
	if k.HasParent(KindsOfAspect) && (len(k.aspects) > 0 || k.name != td.aspect) {
		err = errutil.New("kinds of aspect can only have one set of traits")
	} else if e := td.checkProps(k); e != nil {
		err = e
	} else if e := td.checkTraits(k); e != nil {
		err = e
	}
	return
}

// does this set of traits conflict with any existing fields?
func (td *traitDef) checkProps(k *ScopedKind) (err error) {
	for _, kf := range k.fields {
		// see if this set of traits contains the a field from the kind
		if td.HasTrait(kf.name) {
			err = newConflict(
				k.name,
				Redefined,
				Definition{kf.at, kf.name},
				td.aspect,
			)
			break
		}
	}
	return
}

// does this set of traits conflict with any existing set of traits?
func (td *traitDef) checkTraits(k *ScopedKind) (err error) {
	for _, t := range td.traits {
		if a, ok := k.FindTrait(t); ok {
			err = newConflict(
				k.name,
				Redefined,
				Definition{a.at, a.aspect},
				td.aspect,
			)
			break
		}
	}
	return
}
