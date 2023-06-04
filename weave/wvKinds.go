package weave

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

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
						if pk, ok := d.findPluralKind(ancestor); !ok {
							err = errutil.New("unknown parent kind", opAncestor)
						} else if !kid.Requires.HasAncestor(pk) {
							err = errutil.Fmt("kind %q can't redefine parent as %q", kind, opAncestor)
						} else {
							e := errutil.Fmt("kind %q duplicate parent definition at %v", kind, at)
							LogWarning(e)
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
	return cat.Schedule(assert.MemberPhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		_, newName := d.StripDeterminer(kind)
		if newName, ok := UniformString(newName); !ok {
			err = InvalidString(kind)
		} else if kid, ok := d.findPluralKind(newName); !ok {
			err = errutil.Fmt("unknown kind %q at %v", kind, at)
		} else if uf, e := MakeUniformField(aff, fieldName, class, at); e != nil {
			err = e
		} else if e := cat.writer.Member(d.name, kid, uf.Name, uf.Affinity, uf.Type, at); e != nil {
			err = e
		} else if init != nil {
			return cat.Schedule(assert.DefaultsPhase, func(ctx *Weaver) (err error) {
				return cat.writer.Default(d.name, kid, uf.Name, init)
			})
		}
		return
	})
}
