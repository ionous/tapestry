package eph

type FieldDefinition interface {
	CheckConflict(*ScopedKind) error
	AddToKind(*ScopedKind)
	Write(w Writer) error
}

type fieldDef struct {
	name, affinity, class, at string
}

func (fd *fieldDef) Write(w Writer) error {
	return w.Write(mdl_field, fd.name, fd.affinity, fd.class, fd.at)
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
			// fix? it might be nice to treat class as a dependency
			// then resolve to determine compatibility
			if kf.affinity == fd.affinity && kf.class == fd.class {
				reason = Duplicated
			} else {
				reason = Redefined
			}
			err = &Conflict{
				Reason: reason,
				Was:    Definition{kf.at, kf.name},
				Value:  fd.name,
			}
			break
		}
	}
	return
}

func (fd *fieldDef) checkTraits(k *ScopedKind) (err error) {
	for _, ka := range k.traits {
		for _, t := range ka.traits {
			if t == fd.name {
				err = &Conflict{
					Reason: Redefined,
					Was:    Definition{ka.at, ka.aspect},
					Value:  fd.name,
				}
				break
			}
		}
	}
	return
}
