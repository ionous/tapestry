package weave

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
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

// Kinds, From string, Contain []eph.Params
func (cat *Catalog) AssertAncestor(opKind, opAncestor string) error {
	return cat.Schedule(assert.AncestryPhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		// tbd: are the determiners of kinds useful for anything?
		_, kind := d.StripDeterminer(opKind)
		if kind, ok := UniformString(kind); !ok {
			err = InvalidString(kind)
		} else {
			// ensure the kind into existence
			kid := d.EnsureKind(kind, at)
			// if a parent kind is specified, make the kid dependent on it.
			if _, ancestor := d.StripDeterminer(opAncestor); len(ancestor) > 0 {
				// note: a singular to plural (if needed ) gets handled by the dependency resolver's kindFinder and GetPluralKind()
				if ancestor, ok := UniformString(ancestor); !ok {
					err = InvalidString(opAncestor)
				} else {
					// we can only add requirements to the kind in the same domain that it was declared
					if kid.domain == d {
						kid.AddRequirement(ancestor) // fix? maybe it'd make sense for requirements to have origin at?
					} else {
						// otherwise, if in a different domain: the kinds have to match up
						if pk, ok := d.GetPluralKind(ancestor); !ok {
							err = errutil.New("unknown parent kind", opAncestor)
						} else if !kid.Requires.HasAncestor(pk.name) {
							err = KindError{kind, errutil.Fmt("can't redefine parent as %q", opAncestor)}
						} else {
							LogWarning(KindError{kind, errutil.New("duplicate parent definition at", at)})
						}
					}
				}
			}
		}
		return
	})
}

// generated (for instance) from DefineFields...
// these make new ephemera which are processed during the PropertyPhase.
func (cat *Catalog) AssertField(kind, fieldName, class string, aff affine.Affinity, init assign.Assignment) error {
	return cat.Schedule(assert.AncestryPhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		_, newName := d.StripDeterminer(kind)
		if newName, ok := UniformString(newName); !ok {
			err = InvalidString(kind)
		} else if kid, ok := d.GetPluralKind(newName); !ok {
			err = KindError{kind, errutil.New("unknown kind at", at)}
		} else if uf, e := MakeUniformField(aff, fieldName, class, at); e != nil {
			err = e
		} else {
			kid.pendingFields = append(kid.pendingFields, uf)
		}
		return
	})
}
