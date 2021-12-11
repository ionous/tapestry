package eph

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

func (c *Catalog) WriteRelations(w Writer) (err error) {
	if ks, e := c.ResolveKinds(); e != nil {
		err = e
	} else {
		for _, kdep := range ks {
			if k := kdep.Leaf().(*ScopedKind); k.HasParent(KindsOfRelation) && len(k.fields) > 0 {
				one := k.fields[0]   // a field of affinity text referencing some other kind.
				other := k.fields[1] // the name is the cardinality, and the class is the kind.
				card := makeCard(one.name, other.name)
				if e := w.Write(mdl_rel, k.domain.name, k.name, one.class, card, other.class, k.at); e != nil {
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
func (op *EphRelations) Phase() Phase { return AncestryPhase }

func (op *EphRelations) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if rel, e := d.Singularize(strings.TrimSpace(op.Rel)); e != nil {
		err = e
	} else if a, b, card := op.getCard(); len(card) == 0 {
		err = errutil.New("unknown cardinality")
	} else if ak, e := a.getKind(c, d); e != nil {
		err = e
	} else if bk, e := b.getKind(c, d); e != nil {
		err = e
	} else {
		// add the cardinality as a definition
		// ( used by EphRelatives to determine the cardinality )
		var conflict *Conflict
		if e := d.AddDefinition(rel+"?card", at, card); e != nil && !errors.As(e, &conflict) && conflict.Reason != Duplicated {
			err = e
		} else {
			kid := d.EnsureKind(rel, at)
			kid.AddRequirement(KindsOfRelation)
			err = d.AddEphemera(
				EphAt{at, &EphKinds{
					Kinds: rel,
					Contain: []EphParams{{
						Affinity: a.affinity(),
						Name:     a.short(false),
						Class:    ak,
					}, {
						Affinity: b.affinity(),
						Name:     b.short(true),
						Class:    bk,
					}},
				}})
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
func (k *relKind) getKind(c *Catalog, d *Domain) (ret string, err error) {
	if n := strings.TrimSpace(k.name); len(n) == 0 {
		err = errutil.New("missing name")
	} else if k.plural {
		ret, err = d.Singularize(n)
	} else {
		ret = n
	}
	return
}

func (op *EphRelations) getCard() (first, second relKind, ret string) {
	switch c := op.Cardinality.Value.(type) {
	case *OneOne:
		first = relKind{c.Kind, false}
		second = relKind{c.OtherKind, false}
		ret = tables.ONE_TO_ONE
	case *OneMany:
		first = relKind{c.Kind, false}
		second = relKind{c.OtherKinds, true}
		ret = tables.ONE_TO_MANY
	case *ManyOne:
		first = relKind{c.Kinds, true}
		second = relKind{c.OtherKind, false}
		ret = tables.MANY_TO_ONE
	case *ManyMany:
		first = relKind{c.Kinds, true}
		second = relKind{c.OtherKinds, true}
		ret = tables.MANY_TO_MANY
	}
	return
}
