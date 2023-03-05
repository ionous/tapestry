package eph

import (
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

// write the kinds in a reasonable order
func (c *Catalog) WriteKinds(w Writer) (err error) {
	if deps, e := c.ResolveKinds(); e != nil {
		err = e
	} else {
		for _, dep := range deps {
			k, ancestors := dep.Leaf().(*ScopedKind), dep.Strings(true)
			if e := w.Write(mdl.Kind, k.domain.name, k.name, ancestors, k.at); e != nil {
				err = e
				break
			}
		}
	}
	return
}

var AncestryActions = PhaseAction{
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

// Kinds, From string, Contain []EphParams
func (op *EphKinds) Assemble(c *Catalog, d *Domain, at string) (err error) {
	// tbd: are the determiners of kinds useful for anything?
	kind, ancestor := op.Kind, op.Ancestor
	_, newName := d.StripDeterminer(kind)
	if newName, ok := UniformString(newName); !ok {
		err = InvalidString(kind)
	} else if newName, e := d.Pluralize(newName); e != nil {
		err = InvalidString(kind)
	} else {
		// add the kind we're talking about
		kid := d.EnsureKind(newName, at)
		// remember these fields for future evaluation
		if e := op.addFields(kid, at); e != nil {
			err = e
		} else {
			// if a parent kind is specified, make the kid dependent on it.
			if _, from := d.StripDeterminer(ancestor); len(from) > 0 {
				// note: a singular to plural (if needed ) gets handled by the dependency resolver's kindFinder and GetPluralKind()
				if parentKind, ok := UniformString(from); !ok {
					err = InvalidString(ancestor)
				} else {
					// we can only add requirements to the kind in the same domain that it was declared
					if kid.domain == d {
						kid.AddRequirement(parentKind) // fix? maybe it'd make sense for requirements to have origin at?
					} else {
						// if in a different domain: the kinds have to match up
						if pk, ok := d.GetPluralKind(parentKind); !ok {
							err = errutil.New("unknown parent kind", ancestor)
						} else if !kid.Requires.HasAncestor(pk.name) {
							err = KindError{kind, errutil.Fmt("can't redefine parent as %q", ancestor)}
						} else {
							LogWarning(KindError{kind, errutil.New("duplicate parent definition at", at)})
						}
					}
				}
			}
		}
	}
	return
}

// generated (for instance) from DefineFields...
// these make new ephemera which are processed during the PropertyPhase.
func (op *EphKinds) addFields(k *ScopedKind, at string) (err error) {
	for _, p := range op.Contain {
		if uf, e := p.unify(at); e != nil {
			err = e
		} else {
			k.pendingFields = append(k.pendingFields, uf)
		}
	}
	return
}
