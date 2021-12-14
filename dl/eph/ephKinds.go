package eph

import (
	"strings"

	"github.com/ionous/errutil"
)

const (
	// default kinds
	KindsOfAction   = "action"
	KindsOfAspect   = "aspect"
	KindsOfEvent    = "event"
	KindsOfKind     = "kind"
	KindsOfPattern  = "pattern"
	KindsOfRecord   = "record"
	KindsOfRelation = "relation"
)

// write the kinds in a reasonable order
func (c *Catalog) WriteKinds(w Writer) (err error) {
	if deps, e := c.ResolveKinds(); e != nil {
		err = e
	} else {
		for _, dep := range deps {
			k, ancestors := dep.Leaf().(*ScopedKind), dep.Strings(true)
			if e := w.Write(mdl_kind, k.domain.name, k.name, ancestors, k.at); e != nil {
				err = e
				break
			}
		}
	}
	return
}

var AncestryPhaseActions = PhaseAction{
	PhaseFlags{NoDuplicates: true},
	func(d *Domain) error {
		_, e := d.ResolveKinds()
		return e
	},
}

type KindError struct {
	Kind string
	Err  error
}

func (n KindError) Error() string {
	return errutil.Sprintf("%v for kind %q", n.Err, n.Kind)
}
func (n KindError) Unwrap() error {
	return n.Err
}

// name of a kind to assembly info
func (op *EphKinds) Phase() Phase { return AncestryPhase }

func (op *EphKinds) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if singleKind, e := d.Singularize(strings.TrimSpace(op.Kinds)); e != nil {
		err = e
	} else if newKind, ok := UniformString(singleKind); !ok {
		err = InvalidString(op.Kinds)
	} else if e := op.addFields(d, newKind, at); e != nil {
		err = e
	} else {
		kid := d.EnsureKind(newKind, at)
		// is there a "parent kind specified?"
		if trim := strings.TrimSpace(op.From); len(trim) > 0 {
			if parentKind, ok := UniformString(trim); !ok {
				err = InvalidString(trim)
			} else {
				// we can only add requirements to the kind in the same domain that it was declared
				if kid.domain == d {
					kid.AddRequirement(parentKind) // fix? maybe it'd make sense for requirements to have origin at?
				} else {
					// if in a different domain: the kinds have to match up
					if !kid.HasAncestor(parentKind) {
						err = KindError{newKind, errutil.Fmt("can't redefine parent as %q", parentKind)}
					} else {
						e := errutil.New("duplicate parent definition at", at)
						LogWarning(KindError{newKind, e})
					}
				}
			}
		}
	}
	return
}

func (op *EphKinds) addFields(d *Domain, n, at string) (err error) {
	// fix? for backwards compat, every param becomes a "field" command
	// probably better to loop inside of the fields command rather than here.
	for _, p := range op.Contain {
		if e := d.AddEphemera(EphAt{at, &ephFields{Kinds: n, EphParams: p}}); e != nil {
			err = e
			break
		}
	}
	return
}
