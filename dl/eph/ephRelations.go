package eph

import (
	"strings"

	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

// rather than creating an entirely new hierarchy of relations
// and an entirely new verification mechanism...
// we create a relation kind with fields of an appropriate type
// and use the kind resolving mechanism
func (el *EphRelations) Phase() Phase { return AncestryPhase }

func (el *EphRelations) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if rel, e := c.Singularize(d.name, strings.TrimSpace(el.Relations)); e != nil {
		err = e
	} else if a, b, card := el.getCard(); len(card) < 0 {
		err = errutil.New("unknown cardinality")
	} else if ak, e := a.getKind(c, d); e != nil {
		err = e
	} else if bk, e := b.getKind(c, d); e != nil {
		err = e
	} else {
		kid := d.EnsureKind(rel, at)
		kid.AddRequirement(RelationKinds)
		err = c.AddEphemera(
			EphAt{at, &EphFields{Kinds: rel, Affinity: Affinity{Affinity_Text}, Name: a.short(), Class: ak}},
			EphAt{at, &EphFields{Kinds: rel, Affinity: Affinity{Affinity_Text}, Name: b.short(), Class: bk}})
		// we can still write a separate relation table by looping over domain kinds looking at relation types
	}
	return
}

type relKind struct {
	name   string
	plural bool
}

func (k *relKind) short() (ret string) {
	if k.plural {
		ret = "any"
	} else {
		ret = "one"
	}
	return
}
func (k *relKind) getKind(c *Catalog, d *Domain) (ret string, err error) {
	if n := strings.TrimSpace(k.name); k.plural {
		ret, err = c.Singularize(d.name, n)
	} else {
		ret = n
	}
	return
}

func (el *EphRelations) getCard() (first, second relKind, ret string) {
	switch c := el.Cardinality.Value.(type) {
	case OneOne:
		first = relKind{c.Kind, false}
		second = relKind{c.OtherKind, false}
		ret = tables.ONE_TO_ONE
	case OneMany:
		first = relKind{c.Kind, false}
		second = relKind{c.OtherKinds, true}
		ret = tables.ONE_TO_MANY
	case ManyOne:
		first = relKind{c.Kinds, true}
		second = relKind{c.OtherKind, false}
		ret = tables.MANY_TO_ONE
	case ManyMany:
		first = relKind{c.Kinds, true}
		second = relKind{c.OtherKinds, true}
		ret = tables.MANY_TO_MANY
	}
	return
}
