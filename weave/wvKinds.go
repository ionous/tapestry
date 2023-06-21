package weave

import (
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

// if a parent kind was specified, make the kid dependent on it.
// note: a singular to plural (if needed ) gets handled by the dependency resolver's kindFinder and GetPluralKind()
func (cat *Catalog) AssertAncestor(opKind, opAncestor string) error {
	return cat.Schedule(assert.RequirePlurals, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		// tbd: are the determiners of kinds useful for anything?
		if _, kind := d.UniformDeterminer(opKind); len(kind) == 0 {
			err = InvalidString(kind)
		} else if _, ancestor := d.UniformDeterminer(opAncestor); len(ancestor) == 0 && len(opAncestor) > 0 {
			err = InvalidString(opAncestor)
		} else {
			err = cat.Schedule(assert.RequireDeterminers, func(ctx *Weaver) error {
				return d.addKind(kind, ancestor, at)
			})
		}
		return
	})
}
