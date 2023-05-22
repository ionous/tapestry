package weave

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

type FieldDefinition interface {
	CheckConflict(*ScopedKind) error
	AddToKind(*ScopedKind)
	Write(m mdl.Modeler, domain, kind string) error
}

type fieldDef struct {
	name      string
	affinity  affine.Affinity
	class     string
	at        string
	initially assign.Assignment
}

func (fd *fieldDef) Write(m mdl.Modeler, domain, kind string) error {
	return m.Field(domain, kind, fd.name, fd.affinity, fd.class, fd.at)
}

func (fd *fieldDef) AddToKind(k *ScopedKind) {
	k.fields = append(k.fields, fd)
}

func (fd *fieldDef) CheckConflict(k *ScopedKind) (err error) {
	if k.HasParent(kindsOf.Aspect) {
		err = errutil.New("can't add fields to kinds of aspect")
	} else if e := fd.checkProps(k); e != nil {
		err = e
	} else if e := fd.checkTraits(k); e != nil {
		err = e
	}
	return
}

// does this field conflict with any existing fields?
func (fd *fieldDef) checkProps(k *ScopedKind) (err error) {
	for _, kf := range k.fields {
		if kf.name == fd.name {
			var reason ReasonForConflict
			// fix? it might be nice to treat class as a dependency
			// then resolve to determine compatibility
			if kf.affinity == fd.affinity && kf.class == fd.class {
				reason = Duplicated
			} else {
				reason = Redefined
			}
			key := MakeKey(k.name, kf.name)
			err = newConflict(
				key,
				reason,
				Definition{key, kf.at, kf.name},
				fd.name,
			)
			break
		}
	}
	return
}

// does this field conflict with any existing traits?
func (fd *fieldDef) checkTraits(k *ScopedKind) (err error) {
	if a, ok := k.FindTrait(fd.name); ok {
		key := MakeKey(k.name, fd.name)
		err = newConflict(
			key,
			Redefined,
			Definition{key, a.at, a.aspect},
			fd.name,
		)
	}
	return
}
