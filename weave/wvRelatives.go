package weave

import (
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

// validate that the pattern for the rule exists then add the rule to the *current* domain
// ( rules are de/activated based on domain, they can be part some child of the domain where the pattern was defined. )
func (cat *Catalog) AssertRelative(opRel, opNoun, opOtherNoun string) error {
	return cat.Schedule(assert.RelativePhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		if noun, ok := UniformString(opNoun); !ok {
			err = InvalidString(opNoun)
		} else if otherNoun, ok := UniformString(opOtherNoun); !ok {
			err = InvalidString(opOtherNoun)
		} else if rel, ok := UniformString(opRel); !ok {
			err = InvalidString(opRel)
		} else if first, e := d.GetClosestNoun(noun); e != nil {
			err = e
		} else if second, e := d.GetClosestNoun(otherNoun); e != nil {
			err = e
		} else {
			rel := d.singularize(rel)
			err = cat.writer.Pair(d.name, rel, first.name, second.name, at)
		}
		return
	})
}
