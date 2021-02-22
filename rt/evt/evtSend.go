package evt

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

/*

name := lang.Breakcase(event)
	var labels, result string // fix? consider a cache for this info?
	if e := run.fields.patternOf.QueryRow(name).Scan(&name, &labels, &result); e != nil {
		err = e
	} else
*/

// where args should be always actor, noun, other noun.
// and the return for the event pattern is always a bool.
// optionally, likely, the locals include a "cancel" bool.
func Send(run rt.Runtime, name, labels, result string, up []string, args []rt.Arg) (ret *bool, err error) {
	if rec, e := pattern.NewRecord(run, name, labels, args); e != nil {
		err = e
	} else {
		// we always expect a "bool" result.
		rw := pattern.NewResults(rec, result, affine.Bool)
		// the scope will be the same for the whole event cycle.
		if oldScope, e := run.ReplaceScope(rw, true); e != nil {
			err = e
		} else {
			// note: the scope has to be established before BuildPath gets called
			var allFlags rt.Flags
			if rules, e := BuildPath(run, name, up, &allFlags); e != nil {
				err = e
			} else {
				for i, rules := range rules {
					if phase := rt.Flags(1 << i); phase&allFlags == 0 {
						continue
					} else {
						// the rules stop processing if someone sets a return
						if e := rw.ApplyRules(run, rules, allFlags); e != nil {
							err = errutil.New("error in phase", phase, e)
							break
						} else if rw.HasResults() {
							// if we have a return... we know its a bool
							if b, e := rw.GetResult(); e != nil {
								err = errutil.New("error in phase", phase, e)
								break
							} else {
								p := new(bool)
								*p = b.Bool()
								ret = p
								break
							}
						}
					}
				}
			}
			// only init can return an error
			run.ReplaceScope(oldScope, false)
		}
	}
	if err != nil {
		err = errutil.New("error calling", name, err)
	}
	return
}
