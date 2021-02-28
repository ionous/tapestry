package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/evt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/pattern"
	"github.com/ionous/errutil"
)

// note: this is mirrored/mimicked in package testpat
func (run *Runner) Call(pat string, aff affine.Affinity, args []rt.Arg) (ret g.Value, err error) {
	name := lang.Breakcase(pat) // gets replaced with the actual name by query
	var labels, result string   // fix? consider a cache for this info?
	var rec *g.Record
	if e := run.fields.patternOf.QueryRow(name).Scan(&name, &labels, &result); e != nil {
		err = e
	} else if rec, e = pattern.NewRecord(run, name, labels, args); e != nil {
		err = e
	} else {
		// locals can ( and often do ) read arguments.
		results := pattern.NewResults(rec, result, aff)
		if oldScope, e := run.ReplaceScope(results, true); e != nil {
			err = e
		} else {
			var allFlags rt.Flags
			if rules, e := run.GetRules(name, "", &allFlags); e != nil {
				err = e
			} else {
				ret, err = results.Compute(run, rules, allFlags)
			}
			// only init can return an error
			run.ReplaceScope(oldScope, false)
		}
	}
	if err != nil {
		err = errutil.New(err, "calling", pat, g.RecordToValue(rec))
	}
	return
}

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
			// note: the scope has to be established before BuildPath gets called
			// ( suspiciously like initialize value )
			var allFlags rt.Flags
			if rules, e := evt.BuildPath(run, name, up, &allFlags); e != nil {
				err = e
			} else {
				for i, rules := range rules {
					if phase := rt.Flags(1 << i); phase&allFlags == 0 {
						continue
					} else {
						// the rules stop processing if someone sets a return
						if e := rw.ApplyRules(run, rules, allFlags); e != nil {
							err = errutil.New(e, "in phase", phase)
							break
						} else if rw.HasResults() {
							// if we have a return... we know its a bool
							if b, e := rw.GetResult(); e != nil {
								err = errutil.New(e, "in phase", phase)
							} else {
								ret = b
							}
							break
						}
					}
				}
			}
			// only init can return an error
			run.ReplaceScope(oldScope, false)
		}
	}
	if err != nil {
		err = errutil.New(err, "sending", name)
	}
	return
}

// by now the initializers for the kind will have been cached....
func (run *Runner) initializeLocals(rec *g.Record) (err error) {
	k := rec.Kind()
	if qk, ok := run.qnaKinds.kinds[k.Name()]; !ok {
		err = errutil.New("unknown kind", k.Name())
	} else {
		// run all the initializers
		for i, init := range qk.init {
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
