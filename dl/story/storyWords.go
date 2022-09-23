package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/imp"
)

func (op *MakePlural) PostImport(k *imp.Importer) (noerr error) {
	k.WriteEphemera(&eph.EphPlurals{
		Singular: op.Singular,
		Plural:   op.Plural,
	})
	return
}

func (op *MakeOpposite) PostImport(k *imp.Importer) (noerr error) {
	k.WriteEphemera(&eph.EphOpposites{
		Opposite: op.Opposite,
		Word:     op.Word,
	})
	return
}
