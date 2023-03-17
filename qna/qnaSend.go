package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"git.sr.ht/~ionous/tapestry/rt/scope"
	"github.com/ionous/errutil"
)

// where args should be of the set actor, noun, other noun.
// and the return for the event pattern is always a bool.
// optionally, likely, the locals include a "cancel" bool.
// returns whether true if the event handling didnt return false
func (run *Runner) Send(rec *g.Record, up []string) (ret g.Value, err error) {
	okay := true // provisionally
	name := rec.Kind().Name()
	if res, e := pattern.NewResults(run, rec, affine.Bool); e != nil {
		err = e
	} else {
		oldScope := run.Stack.ReplaceScope(res)
		if cached, e := run.getKindOf(name, kindsOf.Pattern.String()); e != nil {
			err = e
		} else {
			cached.initializeRecord(run, rec)
			// fix: nobody is using "current_noun" currently... so what does that say?
			// ( and what type should the current noun be here? )
			currentNoun := scope.NewSingleValue("current_noun", g.Empty)
			run.PushScope(currentNoun)
			// note: the scope has to be established before BuildPath gets called
			// ( suspiciously like initialize value )
			var flags rt.Flags
			if rules, e := BuildPath(run, name, up, &flags); e != nil {
				err = e
			} else {
			AllPhases:
				for i, cnt := 0, len(rules); okay && i < cnt && flags != 0; i++ {
					if phase := rt.Flags(1 << i); phase&flags != 0 {
						for _, el := range rules[i] {
							currentNoun.SetValue(g.StringFrom(el.Noun, el.Kind))
							// fix? would it make more sense to return the result here?
							// possibly as a pointer so that we can check "has result"
							if next, e := res.ApplyRule(run, el.Rule, flags); e != nil {
								err = errutil.New(e, "applying phase", phase)
								break AllPhases
							} else if flags = next; flags&phase == 0 {
								break // we're done with this phase.
							}
						}
						if ok, e := res.GetContinuation(); e != nil {
							err = errutil.New(e, "resulting from phase", phase)
						} else if !ok {
							okay = false
							break
						}
					}
				}
			}
			run.PopScope()
		}
		run.Stack.ReplaceScope(oldScope)
	}
	if err != nil {
		err = errutil.New(err, "sending", name)
	} else {
		ret = g.BoolOf(okay)
	}
	return
}
