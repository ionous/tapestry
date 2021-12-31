package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"github.com/ionous/errutil"
)

func (op *KindOfRelation) ImportPhrase(k *Importer) (err error) {
	if card, e := op.Cardinality.GetCardinality(); e != nil {
		err = e
	} else {
		k.Write(&eph.EphRelations{op.Relation.String(), card})
	}
	return
}

func (op *RelationCardinality) GetCardinality() (ret eph.EphCardinality, err error) {
	if c, ok := op.Value.(cardinalityImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Value))
	} else {
		ret = c.GetCardinality()
	}
	return
}

func (op *OneToOne) GetCardinality() eph.EphCardinality {
	return eph.EphCardinality{
		eph.EphCardinality_OneOne_Opt,
		&eph.OneOne{op.Kind.String(), op.OtherKind.String()},
	}
}
func (op *OneToMany) GetCardinality() eph.EphCardinality {
	return eph.EphCardinality{
		eph.EphCardinality_OneMany_Opt,
		&eph.OneMany{op.Kind.String(), op.Kinds.String()},
	}
}
func (op *ManyToOne) GetCardinality() eph.EphCardinality {
	return eph.EphCardinality{
		eph.EphCardinality_ManyOne_Opt,
		&eph.ManyOne{op.Kinds.String(), op.Kind.String()},
	}
}
func (op *ManyToMany) GetCardinality() eph.EphCardinality {
	return eph.EphCardinality{
		eph.EphCardinality_ManyMany_Opt,
		&eph.ManyMany{op.Kinds.String(), op.OtherKinds.String()},
	}
}

type cardinalityImporter interface {
	GetCardinality() eph.EphCardinality
}
