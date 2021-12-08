package eph

// implements FieldDefinition
type traitDef struct {
	at     string
	aspect string
	traits []string
}

func (fd *traitDef) Write(w Writer) error {
	// fix: this is "old" format... i really want the specific aspect as the class. (fd.aspect)
	// table: [domain, kind], name, affinity, class, pos
	return w.Write(mdl_field, fd.aspect, Affinity_Text, KindsOfAspect, fd.at)
}

func (td *traitDef) AddToKind(k *ScopedKind) {
	k.traits = append(k.traits, *td)
}

func (td *traitDef) CheckConflict(k *ScopedKind) (err error) {
	if e := td.checkProps(k); e != nil {
		err = e
	} else if td.checkTraits(k); e != nil {
		err = e
	}
	return
}

func (td *traitDef) checkProps(k *ScopedKind) (err error) {
	for _, kf := range k.fields {
		for _, in := range td.traits {
			if in == kf.name {
				err = &Conflict{
					Reason: Redefined,
					Was:    Definition{kf.at, kf.name},
					Value:  td.aspect,
				}
				break
			}
		}
	}
	return
}

func (td *traitDef) checkTraits(k *ScopedKind) (err error) {
	for _, ka := range k.traits {
		for _, t := range ka.traits {
			for _, in := range td.traits {
				if t == in {
					err = &Conflict{
						Reason: Redefined,
						Was:    Definition{ka.at, ka.aspect},
						Value:  td.aspect,
					}
					break
				}
			}
		}
	}
	return
}
