package story

import (
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

// ex. "cats are a kind of record"
func (op *KindsOfRecord) ImportPhrase(k *Importer) (err error) {
	if kind, e := NewRecordPlural(k, op.RecordPlural); e != nil {
		err = e
	} else {
		record := k.NewName("record", tables.NAMED_KIND, op.RecordPlural.At.String())
		k.NewKind(kind, record)
	}
	return
}

// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func (op *RecordsPossessProperties) ImportPhrase(k *Importer) (err error) {
	if kind, e := NewRecordPlural(k, op.RecordPlural); e != nil {
		err = e
	} else {
		for _, n := range op.PropertyDecl {
			if e := n.ImportProperty(k, kind); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}
