package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/event"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/scope"
	"github.com/ionous/errutil"
)

// note: this is mirrored/mimicked in package testpat
// note: in order to generate appropriate defaults ( ex. a record of the right type )
// can return both a both meaningful value *and* an error
func (run *Runner) Call(name string, aff affine.Affinity, keys []string, vals []g.Value) (ret g.Value, err error) {
	// create a container to hold results of args, locals, and the pending return value
	if pat, e := run.getKind(name); e != nil {
		err = e
	} else if rec, e := safe.FillRecord(run, pat.NewRecord(), keys, vals); e != nil {
		err = e
	} else {
		switch pattern.Categorize(g.Ancestry(pat)) {
		case pattern.Initializes:
			ret = g.RecordOf(rec)
		case pattern.Calls:
			ret, err = run.call(pat, rec, aff)
		case pattern.Sends:
			// fix: this search isnt great; especially since its looking for "noun" instead of "target"
			// putting target first might be best, or renaming noun to target
			// but would need some rework of parser and of any explicit calls to the action
			var target int
			if i := pat.FieldIndex("noun"); i >= 0 {
				target = i
			}
			if !rec.HasValue(target) {
				err = errutil.New("no target value specified")
			} else if tgt, e := rec.GetIndexedField(target); e != nil {
				err = e
			} else {
				ret, err = run.send(name, rec, tgt)
			}

		default:
			err = errutil.Fmt("attempting to call %q directly", name)
		}
	}
	return
}

func (run *Runner) call(kind *g.Kind, rec *g.Record, aff affine.Affinity) (ret g.Value, err error) {
	if field, e := pattern.GetResultField(run, kind); e != nil {
		err = e
	} else {
		name := kind.Name()
		oldScope := run.scope.ReplaceScope(scope.FromRecord(run, rec))
		run.currentPatterns.startedPattern(name)
		if e := initRecord(run, rec); e != nil {
			err = e
		} else if rules, e := run.getRules(name); e != nil {
			err = e
		} else if res, e := rules.Call(run, rec, field); e != nil {
			err = e
		} else {
			ret, err = res.GetResult(run, aff)
		}
		run.currentPatterns.stoppedPattern(name)
		run.scope.RestoreScope(oldScope)
	}
	return
}

func (run *Runner) send(name string, rec *g.Record, tgt g.Value) (ret g.Value, err error) {
	// fix? setup the base state first since that's where most initialization lives
	// however, it's not satisfying to move from "action", to "before action", back to "action" again
	// tbd: add an explicit middle pattern ( "is action" )
	callState := run.saveCallState(rec)
	if e := initRecord(run, rec); e != nil {
		err = e
	} else if path, e := run.newPathForTarget(tgt); e != nil {
		err = e // ^ the bubble capture chain
	} else if evtObj, e := newEventRecord(run, event.Object, tgt); e != nil {
		err = e // ^ create the "event object" sent to each event phase pattern.
	} else {
		// fix: shouldnt this be in scope so that it can pull from its sibling variables during initialization?
		// ex. locals from parameters.
		if initRecord(run, rec); e != nil {
			err = e
		}
		var canceled bool
		for prevRec, i := rec, 0; i < event.NumPhases; i++ {
			phase := event.Phase(i)
			if kindForPhase, e := run.getKind(phase.PatternName(name)); e != nil {
				err = e
				break
			} else if rules, e := run.getRules(kindForPhase.Name()); e != nil {
				err = e
			} else {
				// set the name of the phase
				// ( so that authors can determine the name of the phase during initialization )
				callState.setPattern(kindForPhase.Name())
				// create the phase record, set it into scope, and initialize.
				// ( the event object isnt available till the main pattern body )
				if phaseRec, e := run.newPhase(kindForPhase, prevRec); e != nil {
					err = e
					break
				} else {
					// push the event scope. we dont pop later; we replace.
					// ( the event gets pushed last so it cant be hidden by pattern locals )
					run.PushScope(scope.NewReadOnlyValue(event.Object, g.RecordOf(evtObj)))
					if ok, e := rules.Send(run, evtObj, path.slice(phase)); e != nil {
						err = errutil.Fmt("%w calling %q", e, kindForPhase.Name())
						break
					} else if !ok {
						// stopping before, is the same as canceling
						// tbd: should this also check for explicit canceling in later phases?
						canceled = phase < event.TargetPhase
						break
					}
					prevRec = phaseRec
				}
			}
		}
		ret = g.BoolOf(!canceled)
	}
	callState.restore()
	return
}

// creates a new record for an event phase:
// copies any matching values from the previous phase, and initialize the rest.
// just like call, replaces the scope so init can see all of the parameters and locals
// the action pattern declaration implies that the initial set of field sare the same
// ( even though the author can extend those patterns with non-overlapping locals )
// fix: backdoor record creation to slice the values across?
func (run *Runner) newPhase(k *g.Kind, src *g.Record) (ret *g.Record, err error) {
	dst := k.NewRecord()
	_ = run.scope.ReplaceScope(scope.FromRecord(run, dst))
	if e := copyPhase(src, dst); e != nil {
		err = e
	} else if e := initRecord(run, dst); e != nil {
		err = e
	} else {
		ret = dst
	}
	return
}

func copyPhase(src *g.Record, dst *g.Record) (err error) {
	srcKind, dstKind := src.Kind(), dst.Kind()
	for i, srcCnt, dstCnt := 0, srcKind.NumField(), dstKind.NumField(); i < srcCnt && i < dstCnt; i++ {
		// copy until the fields are mismatched
		if srcKind.Field(i) != dstKind.Field(i) {
			break
		} else if src.HasValue(i) {
			// only copy fields that were set
			// (ex. from a previous phase, or via an argument to Call )
			if v, e := src.GetIndexedField(i); e != nil {
				err = e
				break
			} else if e := dst.SetIndexedField(i, g.CopyValue(v)); e != nil {
				err = e
				break
			}
		}
	}
	return
}
