package eph

type FieldDefinition interface {
	CheckConflict(*ScopedKind) error
	AddToKind(*ScopedKind)
}

type fieldDef struct {
	name, affinity, class, at string
}

func (fd *fieldDef) AddToKind(k *ScopedKind) {
	k.fields = append(k.fields, *fd)
}

func (fd *fieldDef) CheckConflict(k *ScopedKind) (err error) {
	if e := fd.checkProps(k); e != nil {
		err = e
	} else if fd.checkTraits(k); e != nil {
		err = e
	}
	return
}

func (fd *fieldDef) checkProps(k *ScopedKind) (err error) {
	for _, kf := range k.fields {
		if kf.name == fd.name {
			var reason ReasonForConflict
			if kf.affinity == fd.affinity && kf.class == fd.class {
				reason = Duplicated
			} else {
				reason = Redefined
			}
			err = &Conflict{reason, Definition{kf.at, kf.name}}
			break
		}
	}
	return
}

func (fd *fieldDef) checkTraits(k *ScopedKind) (err error) {
	for _, ka := range k.traits {
		for _, t := range ka.traits {
			if t == fd.name {
				err = &Conflict{Redefined, Definition{fd.at, fd.name}}
				break
			}
		}
	}
	return
}
