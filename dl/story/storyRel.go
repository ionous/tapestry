package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/imp"
	"github.com/ionous/errutil"
)

func (op *KindOfRelation) PostImport(k *imp.Importer) (err error) {
	if card, e := op.Cardinality.GetCardinality(); e != nil {
		err = e
	} else {
		k.WriteEphemera(&eph.EphRelations{Rel: op.Relation.String(), Cardinality: card})
	}
	return
}

func (op *RelationCardinality) GetCardinality() (ret eph.EphCardinality, err error) {
	if c, ok := op.Value.(cardinalityImporter); !ok {
		err = ImportError(op, errutil.Fmt("%w for %T", UnhandledSwap, op.Value))
	} else {
		ret = c.GetCardinality()
	}
	return
}

func (op *OneToOne) GetCardinality() eph.EphCardinality {
	return eph.EphCardinality{
		eph.EphCardinality_OneOne_Opt,
		&eph.OneOne{Kind: op.Kind.String(), OtherKind: op.OtherKind.String()},
	}
}
func (op *OneToMany) GetCardinality() eph.EphCardinality {
	return eph.EphCardinality{
		eph.EphCardinality_OneMany_Opt,
		&eph.OneMany{Kind: op.Kind.String(), OtherKinds: op.Kinds.String()},
	}
}
func (op *ManyToOne) GetCardinality() eph.EphCardinality {
	return eph.EphCardinality{
		eph.EphCardinality_ManyOne_Opt,
		&eph.ManyOne{Kinds: op.Kinds.String(), OtherKind: op.Kind.String()},
	}
}
func (op *ManyToMany) GetCardinality() eph.EphCardinality {
	return eph.EphCardinality{
		eph.EphCardinality_ManyMany_Opt,
		&eph.ManyMany{Kinds: op.Kinds.String(), OtherKinds: op.OtherKinds.String()},
	}
}

type cardinalityImporter interface {
	GetCardinality() eph.EphCardinality
}
