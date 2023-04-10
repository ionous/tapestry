package eph

import (
	"errors"
	"git.sr.ht/~ionous/tapestry/imp/assert"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

// ensure fields which reference aspects use the necessary formatting
func AspectParam(aspectName string) EphParams {
	return EphParams{Name: aspectName, Affinity: Affinity{Affinity_Text}, Class: aspectName}
}

// uses the ancestry phase because it generates kinds ( one per aspect. )
// the assembly statement generates new ephemera for the aspect phase
// ( to fill the aspect's kind with bool fields representing the traits. )
func (op *EphAspects) Phase() assert.Phase { return assert.AncestryPhase }

// generates traits and adds them to a custom aspect kind.
func (op *EphAspects) Assemble(c *Catalog, d *Domain, at string) (err error) {
	// we dont singularize aspects even thought its a kind;
	// most are really singularizable anyway, and some common things like "darkness" dont singularize correctly.
	if aspect, ok := UniformString(op.Aspects); !ok {
		err = InvalidString(op.Aspects)
	} else if traits, e := UniformStrings(op.Traits); e != nil {
		err = e
	} else {
		kid := d.EnsureKind(aspect, at)
		kid.AddRequirement(kindsOf.Aspect.String())
		if len(traits) > 0 {
			err = d.AddEphemera(at, PhaseFunction{assert.AspectPhase,
				func(c *Catalog, d *Domain, at string) (err error) {
					var conflict *Conflict // checks for conflicts, allows duplicates.
					if e := kid.AddField(&traitDef{at, aspect, traits}); errors.As(e, &conflict) && conflict.Reason == Duplicated {
						LogWarning(e) // warn if it was a duplicated definition
					} else {
						err = e // some other error ( or nil )
					}
					return
				}})
		}
	}
	return
}
