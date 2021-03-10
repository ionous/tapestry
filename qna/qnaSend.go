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
func (run *Runner) Send(pat string, up []string, args []rt.Arg) (ret g.Value, err error) {
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
				for i, cnt := 0, len(rules); i < cnt && flags != 0; i++ {
					if phase := rt.Flags(1 << i); phase&flags != 0 {
						for _, el := range rules[i] {
							currentNoun.SetValue(g.StringOf(el.Noun))
							if next, e := rw.ApplyRule(run, el.Rule, flags); e != nil {
								err = errutil.New(e, "applying phase", phase)
								break AllPhases
							} else {
								flags = next
								// the first result from a rule kicks us out of this phase.
								if flags&phase == 0 {
									break
								}
							}
						}
						// note: if we have a return... we know its a bool for events.
						if rw.ComputedResult() {
							if res, e := rw.GetResult(); e != nil {
								err = errutil.New(e, "resulting from phase", phase)
								break AllPhases
							} else {
								ret = res
								// if the return is false, we end all the remaining phases
								if !res.Bool() {
									break AllPhases
								}
								// otherwise, we move on to the next.
								rw.ResetResult()
							}
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
	}
	return
}
