package weave

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
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
		} else if a, b, card := makeKind(a, b, amany, bmany); len(card) == 0 {
			err = errutil.New("unknown cardinality")
		} else if len(a.class) == 0 {
			err = errutil.New("invalid lhs")
		} else if len(b.class) == 0 {
			err = errutil.New("invalid rhs")
		} else {
			rel := d.singularize(rel)
			if e := d.addKind(rel, kindsOf.Relation.String(), at); e != nil {
				err = e
			} else if ua, e := MakeUniformField(a.affinity(), a.short(false), a.class, at); e != nil {
				err = e
			} else if ub, e := MakeUniformField(b.affinity(), b.short(true), b.class, at); e != nil {
				err = e
			} else {
				err = cat.Schedule(assert.MemberPhase, func(ctx *Weaver) (err error) {
					if e := cat.writer.Member(d.name, rel, ua.Name, ua.Affinity, ua.Type, at); e != nil {
						err = e
					} else if e := cat.writer.Member(d.name, rel, ub.Name, ub.Affinity, ub.Type, at); e != nil {
						err = e
					} else if e := cat.writer.Rel(d.name, rel, ua.Type, ub.Type, card, at); e != nil {
						err = e // ^ not sure the best time to write this....
					}
					return
				})
			}
		}
		return
	})
}

type relKind struct {
	class  string
	plural bool
}

func (k *relKind) affinity() (ret affine.Affinity) {
	if k.plural {
		ret = affine.TextList
	} else {
		ret = affine.Text
	}
	return
}

// matches tables cardinality
// fix? could also name the fields after the specific kind
func (k *relKind) short(other bool) (ret string) {
	if other {
		if k.plural {
			ret = "other_kinds"
		} else {
			ret = "other_kind"
		}
	} else {
		if k.plural {
			ret = "kinds"
		} else {
			ret = "kind"
		}
	}
	return
}

func makeKind(a, b string, amany, bmany bool) (first, second relKind, ret string) {
	a, b = strings.TrimSpace(a), strings.TrimSpace(b)
	switch {
	case !amany && !bmany:
		first = relKind{a, false}
		second = relKind{b, false}
		ret = tables.ONE_TO_ONE
	case !amany && bmany:
		first = relKind{a, false}
		second = relKind{b, true}
		ret = tables.ONE_TO_MANY
	case amany && !bmany:
		first = relKind{a, true}
		second = relKind{b, false}
		ret = tables.MANY_TO_ONE
	case amany && bmany:
		first = relKind{a, true}
		second = relKind{b, true}
		ret = tables.MANY_TO_MANY
	}
	return
}
