package story

import (
	"git.sr.ht/~ionous/iffy/dl/eph"
	"github.com/ionous/errutil"
)

// ex. "cats are a kind of record"
func (op *KindsOfRecord) ImportPhrase(k *Importer) (err error) {
	kinds := op.RecordPlural.String()
	k.Write(&eph.EphKinds{Kinds: kinds, From: "record"})
	return
}

// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func (op *RecordsPossessProperties) ImportPhrase(k *Importer) (err error) {
	kinds := op.RecordPlural.String()
	for _, n := range op.PropertyDecl {
		if e := n.ImportProperty(k, kinds); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}