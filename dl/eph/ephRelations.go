package eph

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/imp/assert"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

func (c *Catalog) WriteRelations(w Writer) (err error) {
	if ks, e := c.ResolveKinds(); e != nil {
		err = e
	} else {
		for _, kdep := range ks {
			if k := kdep.Leaf().(*ScopedKind); k.HasParent(kindsOf.Relation) && len(k.fields) > 0 {
				one := k.fields[0]   // a field of affinity text referencing some other kind.
				other := k.fields[1] // the name is the cardinality, and the class is the kind.
				card := makeCard(one.name, other.name)
				if e := w.Write(mdl.Rel, k.domain.name, k.name, one.class, other.class, card, k.at); e != nil {
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

// rather than creating an entirely new hierarchy of relations
// and an entirely new verification mechanism...
// we create a relation kind with fields of an appropriate type
// and use the kind resolving mechanism
func (op *EphRelations) Phase() assert.Phase { return assert.AncestryPhase }

func (op *EphRelations) Weave(k assert.Assertions) (err error) {
	switch c := op.Cardinality.Value.(type) {
	case *OneOne:
		err = k.AssertRelation(op.Rel, c.Kind, c.OtherKind, false, false)
	case *OneMany:
		err = k.AssertRelation(op.Rel, c.Kind, c.OtherKinds, false, true)
	case *ManyOne:
		err = k.AssertRelation(op.Rel, c.Kinds, c.OtherKind, true, false)
	case *ManyMany:
		err = k.AssertRelation(op.Rel, c.Kinds, c.OtherKinds, true, true)
	}
	return
}

func (ctx *Context) AssertRelation(opRel, a, b string, amany, bmany bool) (err error) {
	d, at := ctx.d, ctx.at
	// like aspects, we dont try to singularize these.
	if rel, ok := UniformString(opRel); !ok {
		err = InvalidString(opRel)
	} else if a, b, card := makeKind(a, b, amany, bmany); len(card) == 0 {
		err = errutil.New("unknown cardinality")
	} else if len(a.name) == 0 {
		err = errutil.New("invalid name")
	} else if len(b.name) == 0 {
		err = errutil.New("invalid name")
	} else {
		// add the cardinality as a definition
		// ( used by EphRelatives to determine the cardinality )
		var conflict *Conflict
		if e := d.AddDefinition(MakeKey("rel", rel, "card"), at, card); e != nil && !errors.As(e, &conflict) && conflict.Reason != Duplicated {
			err = e
		} else {
			kid := d.EnsureKind(rel, at)
			kid.AddRequirement(kindsOf.Relation.String())
			err = d.QueueEphemera(at, &EphKinds{
				Kind: rel,
				Contain: []EphParams{{
					Affinity: a.affinity(),
					Name:     a.short(false),
					Class:    a.name,
				}, {
					Affinity: b.affinity(),
					Name:     b.short(true),
					Class:    b.name,
				}},
			})
		}
	}
	return
}

type relKind struct {
	name   string
	plural bool
}

func (k *relKind) affinity() (ret Affinity) {
	if k.plural {
		ret = Affinity{Affinity_TextList}
	} else {
		ret = Affinity{Affinity_Text}
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
