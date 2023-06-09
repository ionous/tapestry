package weave

import (
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

func (cat *Catalog) AssertRelation(opRel, a, b string, amany, bmany bool) error {
	// uses ancestry because it defines kinds for each relation
	return cat.Schedule(assert.AncestryPhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		// like aspects, we dont try to singularize these.
		if rel, ok := UniformString(opRel); !ok {
			err = InvalidString(opRel)
		} else if acls, ok := UniformString(a); !ok {
			err = InvalidString(a)
		} else if bcls, ok := UniformString(b); !ok {
			err = InvalidString(b)
		} else if card := makeCard(amany, bmany); len(card) == 0 {
			err = errutil.New("unknown cardinality")
		} else {
			rel := d.singularize(rel)
			if e := d.addKind(rel, kindsOf.Relation.String(), at); e != nil {
				err = e
			} else {
				err = cat.Schedule(assert.MemberPhase, func(ctx *Weaver) (err error) {
					if e := cat.writer.Rel(d.name, rel, acls, bcls, card, at); e != nil {
						err = e
					}
					return
				})
			}
		}
		return
	})
}

func makeCard(amany, bmany bool) (ret string) {
	switch {
	case !amany && !bmany:
		ret = tables.ONE_TO_ONE
	case !amany && bmany:
		ret = tables.ONE_TO_MANY
	case amany && !bmany:
		ret = tables.MANY_TO_ONE
	case amany && bmany:
		ret = tables.MANY_TO_MANY
	}
	return
}
