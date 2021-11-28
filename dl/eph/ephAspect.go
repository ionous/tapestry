package eph

// fix: actually adding the aspect to a target kind
// ( this is just the traits to the aspect )
func (el *EphAspect) Phase() Phase { return AncestryPhase }

func (el *EphAspect) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if aspect, ok := UniformString(el.Aspect); !ok {
		err = InvalidString(el.Aspect)
	} else if traits, e := UniformStrings(el.Traits); e != nil {
		err = e
	} else if e := addKind(c, d, at, aspect, AspectKinds); e != nil {
		err = e
	} else {
		err = c.AddEphemera(EphAt{at, PhaseFunction{AspectPhase,
			func(c *Catalog, d *Domain, at string) error {
				return c.AddFields(d, aspect, &traitDef{at, aspect, traits})
			}}})
	}
	return
}
