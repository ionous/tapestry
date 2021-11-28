package eph

// implements FieldDefinition
type traitDef struct {
	at     string
	aspect string
	traits []string
}

func (td *traitDef) AddToKind(k *Kind) {
	k.traits = append(k.traits, *td)
}

func (td *traitDef) CheckConflict(k *Kind) (err error) {
	if e := td.checkProps(k); e != nil {
		err = e
	} else if td.checkTraits(k); e != nil {
		err = e
	}
	return
}

func (td *traitDef) checkProps(k *Kind) (err error) {
	for _, kf := range k.fields {
		for _, in := range td.traits {
			if in == kf.name {
				err = &Conflict{Redefined, Definition{kf.at, kf.name}}
				break
			}
		}
	}
	return
}

func (td *traitDef) checkTraits(k *Kind) (err error) {
	for _, ka := range k.traits {
		for _, t := range ka.traits {
			for _, in := range td.traits {
				if t == in {
					err = &Conflict{Redefined, Definition{ka.at, ka.aspect}}
					break
				}
			}
		}
	}
	return
}
