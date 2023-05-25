package weave

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

func (c *Catalog) WriteRelations(m mdl.Modeler) (err error) {
	if ks, e := c.ResolveKinds(); e != nil {
		err = e
	} else {
		for _, kdep := range ks {
			if k := kdep.Leaf().(*ScopedKind); k.HasParent(kindsOf.Relation) && len(k.fields) > 0 {
				one := k.fields[0]   // a field of affinity text referencing some other kind.
				other := k.fields[1] // the name is the cardinality, and the class is the kind.
				card := makeCard(one.name, other.name)
				if e := m.Rel(k.domain.name, k.name, one.class, other.class, card, k.at); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

// return a string compatible with package table's cardinality
// see also: relKind.short() for how it specifies fields
func makeCard(one, other string) string {
	var s strings.Builder
	oneOrAny(&s, one)
	s.WriteRune('_')
	oneOrAny(&s, other)
	return s.String()
}

func oneOrAny(out *strings.Builder, s string) (ret string) {
	if s[len(s)-1] == 's' {
		out.WriteString("any")
	} else {
		out.WriteString("one")
	}
	return
}

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
			// add the cardinality as a definition
			// ( used by eph.Relatives to determine the cardinality )
			var conflict *Conflict
			if e := d.AddDefinition(MakeKey("rel", rel, "card"), at, card); e != nil && !errors.As(e, &conflict) && conflict.Reason != Duplicated {
				err = e
			} else {
				kid := d.EnsureKind(rel, at)
				kid.AddRequirement(kindsOf.Relation.String())

				if ua, e := MakeUniformField(a.affinity(), a.short(false), a.class, at); e != nil {
					err = e
				} else if ub, e := MakeUniformField(b.affinity(), b.short(true), b.class, at); e != nil {
					err = e
				} else {
					kid.pendingFields = append(kid.pendingFields, ua, ub)

					err = cat.Schedule(assert.FieldPhase, func(ctx *Weaver) (err error) {
						if e := cat.writeField(d.name, kid.name, ua); e != nil {
							err = e
						} else if e := cat.writeField(d.name, kid.name, ub); e != nil {
							err = e
						}
						return
					})

				}
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
