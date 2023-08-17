package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"github.com/ionous/errutil"
)

// where args should be of the set actor, noun, other noun.
// and the return for the event pattern is always a bool.
// optionally, likely, the locals include a "cancel" bool.
// returns whether true if the event handling didnt return false
func (run *Runner) Send(rec *g.Record, up []string) (ret g.Value, err error) {
	okay := true // provisionally
	name := rec.Kind().Name()
	if labels, e := run.GetField(meta.PatternLabels, name); e != nil {
		err = e
	} else {
		res := pattern.NewResults(rec, labels.Strings(), affine.Bool)
		oldScope := run.replaceScope(res)
		if cached, e := run.getKindOf(name, kindsOf.Pattern.String()); e != nil {
			err = e
		} else if e := cached.recordInit(run, rec); e != nil {
			err = e
		} else {
			// note: the scope has to be established before BuildPath gets called
			// ( suspiciously like initialize value )
			if rules, e := BuildPath(run, name, up); e != nil {
				err = e
			} else {
				var skip bool
				for _, el := range rules {
					if next, e := res.ApplyRule(run, el.Rule, skip); e != nil {
						err = errutil.New(e, "applying phase")
						break
					} else {
						skip = next
					}
				}
				if ok, e := res.GetContinuation(); e != nil {
					err = errutil.New(e, "resulting from phase")
				} else if !ok {
					okay = false
				}
			}

		}
		run.restoreScope(oldScope)
	}
	if err != nil {
		err = errutil.New(err, "sending", name)
	} else {
		ret = g.BoolOf(okay)
	}
	return
}
