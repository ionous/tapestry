package story

import (
	"git.sr.ht/~ionous/iffy/ephemera/eph"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

func (op *KindOfRelation) ImportPhrase(k *Importer) (err error) {
	// rec.NewRelation(r, k, q, c)
	if rel, e := NewRelation(k, op.Relation); e != nil {
		err = e
	} else if card, e := op.Cardinality.ImportCardinality(k); e != nil {
		err = e
	} else {
		k.NewRelation(rel, card.firstKind, card.secondKind, card.cardinality)
	}
	return
}

type importedCardinality struct {
	cardinality           string // tables.ONE_TO_ONE
	firstKind, secondKind eph.Named
}

func (op *RelationCardinality) ImportCardinality(k *Importer) (ret importedCardinality, err error) {
	type cardinalityImporter interface {
		ImportCardinality(k *Importer) (importedCardinality, error)
	}
	if c, ok := op.Value.(cardinalityImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Value))
	} else {
		ret, err = c.ImportCardinality(k)
	}
	return
}

func (op *OneToOne) ImportCardinality(k *Importer) (ret importedCardinality, err error) {
	if first, e := FixPlurals(k, op.Kind); e != nil {
		err = e
	} else if second, e := FixPlurals(k, op.OtherKind); e != nil {
		err = e
	} else {
		ret = importedCardinality{tables.ONE_TO_ONE, first, second}
	}
	return
}
func (op *OneToMany) ImportCardinality(k *Importer) (ret importedCardinality, err error) {
	if first, e := FixPlurals(k, op.Kind); e != nil {
		err = e
	} else if second, e := NewPluralKinds(k, op.Kinds); e != nil {
		err = e
	} else {
		ret = importedCardinality{tables.ONE_TO_MANY, first, second}
	}
	return
}
func (op *ManyToOne) ImportCardinality(k *Importer) (ret importedCardinality, err error) {
	if first, e := NewPluralKinds(k, op.Kinds); e != nil {
		err = e
	} else if second, e := FixPlurals(k, op.Kind); e != nil {
		err = e
	} else {
		ret = importedCardinality{tables.MANY_TO_ONE, first, second}
	}
	return
}
func (op *ManyToMany) ImportCardinality(k *Importer) (ret importedCardinality, err error) {
	if first, e := NewPluralKinds(k, op.Kinds); e != nil {
		err = e
	} else if second, e := NewPluralKinds(k, op.OtherKinds); e != nil {
		err = e
	} else {
		ret = importedCardinality{tables.MANY_TO_MANY, first, second}
	}
	return
}
