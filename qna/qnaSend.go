package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/evt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/pattern"
	"git.sr.ht/~ionous/iffy/rt/scope"
	"github.com/ionous/errutil"
)

// where args should be of the set actor, noun, other noun.
// and the return for the event pattern is always a bool.
// optionally, likely, the locals include a "cancel" bool.
// returns whether true if the event handling didnt return false
func (run *Runner) Send(pat string, up []string, args []rt.Arg) (ret g.Value, err error) {
	okay := true                // provisionally
	name := lang.Breakcase(pat) // gets replaced with the actual name by query
	var labels, result string   // fix? consider a cache for this info?
	if e := run.fields.patternOf.QueryRow(name).Scan(&name, &labels, &result); e != nil {
		err = e
	} else if rec, e := pattern.NewRecord(run, name, labels, args); e != nil {
		err = e
	} else {
		// we always expect a "bool" result.
		rw := pattern.NewResults(rec, result, affine.Bool)
		if oldScope, e := run.ReplaceScope(rw, true); e != nil {
			err = e
		} else {
			currentNoun := scope.NewSingleValue("current_noun", g.Empty)
			run.PushScope(currentNoun)
			// note: the scope has to be established before BuildPath gets called
			// ( suspiciously like initialize value )
			var flags rt.Flags
			if rules, e := evt.BuildPath(run, name, up, &flags); e != nil {
				err = e
			} else {
			AllPhases:
				for i, cnt := 0, len(rules); okay && i < cnt && flags != 0; i++ {
					if phase := rt.Flags(1 << i); phase&flags != 0 {
						for _, el := range rules[i] {
							currentNoun.SetValue(g.StringOf(el.Noun))
							// fix? would it make more sense to return the result here?
							// possibly as a pointer so that we can check "has result"
							if next, e := rw.ApplyRule(run, el.Rule, flags); e != nil {
								err = errutil.New(e, "applying phase", phase)
								break AllPhases
							} else if flags = next; flags&phase == 0 {
								break // we're done with this phase.
							}
						}
						if ok, e := rw.GetContinuation(); e != nil {
							err = errutil.New(e, "resulting from phase", phase)
						} else if !ok {
							okay = false
							break
						}
					}
				}
			}
			run.PopScope()
			run.ReplaceScope(oldScope, false)
		}
	}
	if err != nil {
		err = errutil.New(err, "sending", name)
	} else {
		ret = g.BoolOf(okay)
	}
	return
}
