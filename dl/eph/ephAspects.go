package eph

import (
	"errors"
)

// ensure fields which reference aspects use the necessary formatting
func AspectParam(aspectName string) EphParams {
	return EphParams{Name: aspectName, Affinity: Affinity{Affinity_Text}, Class: aspectName}
}

func (c *Catalog) WriteAspects(w Writer) (err error) {
	if deps, e := c.ResolveKinds(); e != nil {
		err = e
	} else {
		for _, dep := range deps {
			// this is a kind of aspect:
			if k := dep.Leaf().(*ScopedKind); k.HasParent(KindsOfAspect) {
				a := k.aspects[0] // we only expect to see 1 -- probably not worth error checking it.
				for i, t := range a.traits {
					if e := w.Write(mdl_aspect, k.domain.name, a.aspect, t, i); e != nil {
						err = e
						break
					}
				}
			}
		}
	}
	return
}

// uses the ancestry phase because it generates kinds ( one per aspect. )
// the assembly statement generates new ephemera for the aspect phase
// ( to fill the aspect's kind with bool fields representing the traits. )
func (op *EphAspects) Phase() Phase { return AncestryPhase }

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
		kid.AddRequirement(KindsOfAspect)
		if len(traits) > 0 {
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
	}
	return
}
