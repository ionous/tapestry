package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/kindsOf"
	"git.sr.ht/~ionous/iffy/rt/meta"
	"git.sr.ht/~ionous/iffy/rt/pattern"
	"github.com/ionous/errutil"
)

func (run *Runner) ReplaceScope(s rt.Scope, init bool) (ret rt.Scope, err error) {
	ret = run.Stack.ReplaceScope(s)
	if init {
		// fix... yeah, possibly this needs work.
		// Runner is the thing calling replace scope
		// so it could initialize locals, but... also depends on what happens with "Send"
		// cant really put locals into g. because rt depends on g --
		// but could put an initializer function per field maybe.
		if res, ok := s.(*pattern.Results); !ok {
			err = errutil.New("can only initialize records")
		} else {
			err = run.initializeLocals(res.Record())
		}
	}
	// fix our errors
	if err != nil {
		run.Stack.ReplaceScope(ret)
		ret = nil
	}
	return
}

// get the initializer and ... init them.
func (run *Runner) initializeLocals(rec *g.Record) (err error) {
	k := rec.Kind()
	if cached, e := run.getKindOf(k.Name(), kindsOf.Pattern.String()); e != nil {
		err = e
	} else {
		for fieldIndex, init := range cached.init {
			if init != nil { // not every field necessarily has an initializer
				ft := k.Field(fieldIndex)
				if src, e := init.GetAssignedValue(run); e != nil {
					err = errutil.New("error determining local", k.Name(), ft.Name, e)
					break
				} else if val, e := autoConvert(run, ft, src); e != nil {
					err = e
				} else if e := rec.SetIndexedField(fieldIndex, val); e != nil {
					err = errutil.New("error setting local", k.Name(), ft.Name, e)
					break
				}
			}
		}
	}
	return
}

// if the target field (ex. a pattern local) requires text of a certain type
// and the incoming value is untyped: convert it.
// FIX: see pattern.NewRecord()
func autoConvert(run rt.Runtime, ft g.Field, val g.Value) (ret g.Value, err error) {
	if needsConversion := ft.Affinity == affine.Text && len(ft.Type) > 0 &&
		val.Affinity() == affine.Text && len(val.Type()) == 0; !needsConversion {
		ret = val
	} else {
		// set indexed field validates the ft.Type and the val.Type match
		// we just have to give it the proper value in the first place.
		ret, err = run.GetField(meta.ObjectId, val.String())
	}
	return
}
