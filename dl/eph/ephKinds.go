package eph

import (
	"strings"

	"github.com/ionous/errutil"
)

const (
	KindsOfAspect   = "aspect"
	KindsOfRecord   = "record"
	KindsOfRelation = "relation"
)

// traverse the domains and then kinds in a reasonable order
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
func (el *EphKinds) Phase() Phase { return AncestryPhase }

func (el *EphKinds) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if singleKind, e := d.Singularize(strings.TrimSpace(el.Kinds)); e != nil {
		err = e
	} else if newKind, ok := UniformString(singleKind); !ok {
		err = InvalidString(el.Kinds)
	} else {
		kid := d.EnsureKind(newKind, at)
		if parentKind, ok := UniformString(el.From); !ok && len(el.From) > 0 {
			err = InvalidString(el.From)
		} else if kid.domain == d {
			// we can only add requirements to the kind in the same domain that it was declared
			// if in a different domain: the kinds have to match up
			if len(parentKind) > 0 {
				kid.AddRequirement(parentKind)
			}
		} else if ok, e := kid.HasAncestor(parentKind); e != nil {
			err = e
		} else if !ok {
			err = KindError{newKind, errutil.Fmt("can't redefine parent as %q", parentKind)}
		} else {
			e := errutil.New("duplicate parent definition at", at)
			LogWarning(KindError{newKind, e})
		}
	}
	return
}
