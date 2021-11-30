package eph

import "strings"

// uses the ancestry phase because it generates kinds ( one per aspect. )
// the assembly statement generates new ephemera for the aspect phase
// ( to fill the aspect's kind with bool fields representing the traits. )
func (el *EphAspects) Phase() Phase { return AncestryPhase }

// generates traits and adds them to a custom aspect kind.
func (el *EphAspects) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if singleAspect, e := c.Singularize(d.name, strings.TrimSpace(el.Aspects)); e != nil {
		err = e
	} else if aspect, ok := UniformString(singleAspect); !ok {
		err = InvalidString(el.Aspects)
	} else if traits, e := UniformStrings(el.Traits); e != nil {
		err = e
	} else {
		kid := d.EnsureKind(aspect, at)
		kid.AddRequirement(AspectKinds)
		err = c.AddEphemera(EphAt{at, PhaseFunction{AspectPhase,
			func(c *Catalog, d *Domain, at string) error {
				return kid.AddFields(&traitDef{at, aspect, traits})
			}}})
	}
	return
}
