package eph

import (
	"errors"
	"strings"
)

// uses the ancestry phase because it generates kinds ( one per aspect. )
// the assembly statement generates new ephemera for the aspect phase
// ( to fill the aspect's kind with bool fields representing the traits. )
func (el *EphAspects) Phase() Phase { return AncestryPhase }

// generates traits and adds them to a custom aspect kind.
func (el *EphAspects) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if singleAspect, e := d.Singularize(strings.TrimSpace(el.Aspects)); e != nil {
		err = e
	} else if aspect, ok := UniformString(singleAspect); !ok {
		err = InvalidString(el.Aspects)
	} else if traits, e := UniformStrings(el.Traits); e != nil {
		err = e
	} else {
		kid := d.EnsureKind(aspect, at)
		kid.AddRequirement(KindsOfAspect)
		err = d.AddEphemera(EphAt{at, PhaseFunction{AspectPhase,
			func(c *Catalog, d *Domain, at string) (err error) {
				var conflict *Conflict // checks for conflicts, allows duplicates.
				if e := kid.AddField(&traitDef{at, aspect, traits}); errors.As(e, &conflict) && conflict.Reason == Duplicated {
					LogWarning(e) // warn if it was a duplicated definition
				} else {
					err = e // some other error ( or nil )
				}
				return
			}}})
	}
	return
}
