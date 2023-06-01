package weave

// implements FieldDefinition
// when the owner kind is the kind of aspect --
// the owner will only have one trait def
// and the name of the aspect will match the name of the kind
type traitDef struct {
	at     string
	aspect string
	traits []string
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
