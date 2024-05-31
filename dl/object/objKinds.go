package object

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// tbd: fields of kind could be expanded to full introspection
// if the Field structure was a record.
// ( maybe some core classes could be exposed using the test reflection?
// might make porting harder )
func (op *FieldsOfKind) GetTextList(run rt.Runtime) (ret rt.Value, err error) {
	if name, e := safe.GetText(run, op.KindName); e != nil {
		err = e
	} else {
		ret, err = run.GetField(meta.FieldsOfKind, name.String())
	}
	return
}
