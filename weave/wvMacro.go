package weave

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// duplicates cachedKind.
type macroKind struct {
	*g.Kind
	init []assign.Assignment
	do   []rt.Execute
}

func (k macroKind) initializeRecord(run rt.Runtime, rec *g.Record) (err error) {
	for fieldIndex, init := range k.init {
		if init != nil {
			ft := k.Field(fieldIndex)
			if src, e := safe.GetAssignment(run, init); e != nil {
				err = errutil.New("error determining local", k.Name(), ft.Name, e)
				break
			} else if val, e := safe.AutoConvert(run, ft, src); e != nil {
				err = e
			} else if e := rec.SetIndexedField(fieldIndex, val); e != nil {
				err = errutil.New("error setting local", k.Name(), ft.Name, e)
				break
			}
		}
	}
	return
}

type macroReg map[string]macroKind
