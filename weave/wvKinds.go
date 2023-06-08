package weave

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

// if a parent kind was specified, make the kid dependent on it.
// note: a singular to plural (if needed ) gets handled by the dependency resolver's kindFinder and GetPluralKind()

func (cat *Catalog) AssertAncestor(opKind, opAncestor string) error {
	return cat.Schedule(assert.DeterminerPhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		// tbd: are the determiners of kinds useful for anything?
		_, kind := d.StripDeterminer(opKind)
		_, ancestor := d.StripDeterminer(opAncestor)
		//
		if kind, ok := UniformString(kind); !ok {
			err = InvalidString(kind)
		} else if ua, ok := UniformString(ancestor); !ok && len(ancestor) > 0 {
			err = InvalidString(opAncestor)
		} else {
			kind, ancestor := d.singularize(kind), d.singularize(ua)
			err = cat.Schedule(assert.AncestryPhase, func(ctx *Weaver) error {
				return d.addKind(kind, ancestor, at)
			})
		}
		return
	})
}

// generated (for instance) from DefineFields...
// these make new ephemera which are processed during the PropertyPhase.
func (cat *Catalog) AssertField(kind, fieldName, class string, aff affine.Affinity, init assign.Assignment) error {
	return cat.Schedule(assert.MemberPhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		_, name := d.StripDeterminer(kind)
		if name, ok := UniformString(name); !ok {
			err = InvalidString(kind)
		} else {
			kid := d.singularize(name)
			if uf, e := MakeUniformField(aff, fieldName, class, at); e != nil {
				err = e
			} else if e := cat.writer.Member(d.name, kid, uf.Name, uf.Affinity, uf.Type, at); e != nil {
				err = e
			} else if init != nil {
				return cat.Schedule(assert.DefaultsPhase, func(ctx *Weaver) (err error) {
					return cat.writer.Default(d.name, kid, uf.Name, init)
				})
			}
		}
		return
	})
}
