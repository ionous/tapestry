package story

import "git.sr.ht/~ionous/tapestry/dl/eph"

func (op *MakePlural) ImportPhrase(k *Importer) (noerr error) {
	k.WriteEphemera(&eph.EphPlurals{
		Singular: op.Singular,
		Plural:   op.Plural,
	})
	return
}

func (op *MakeOpposite) ImportPhrase(k *Importer) (noerr error) {
	k.WriteEphemera(&eph.EphOpposites{
		Opposite: op.Opposite,
		Word:     op.Word,
	})
	return
}
