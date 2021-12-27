package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/kindsOf"
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
		if v, ok := s.(*pattern.Results).Scope.(g.Value); !ok || v.Affinity() != affine.Record {
			err = errutil.New("can only initialize records")
		} else {
			err = run.initializeLocals(v.Record())
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
	if cached, e := run.getKindOf(k.Name(), kindsOf.Pattern); e != nil {
		err = e
	} else {
		for i, init := range cached.init {
			if init != nil { // not every field necessarily has an initializer
				if v, e := init.GetAssignedValue(run); e != nil {
					err = errutil.New("error determining local", k.Name(), k.Field(i).Name, e)
					break
				} else if e := rec.SetIndexedField(i, v); e != nil {
					err = errutil.New("error setting local", k.Name(), k.Field(i).Name, e)
					break
				}
			}
		}
	}
	return
}
