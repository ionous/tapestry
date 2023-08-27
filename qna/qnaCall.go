package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/event"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
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
		// Call is being used as a side effect to initialize records
		if pat.Implements(kindsOf.Record.String()) {
			if aff != affine.Record {
				err = errutil.Fmt("attempting to call a record %q with affine %q", name, aff)
			} else {
				ret = g.RecordOf(rec)
			}
		} else if pat.Implements(kindsOf.Event.String()) {
			err = errutil.Fmt("attempting to call an event  %q directly", name)
		} else if !pat.Implements(kindsOf.Action.String()) {
			// note: this doesnt positively affirm kindsOf.Pattern:
			// some tests use golang structs as faux patterns.
			ret, err = run.call(rec, pat, aff)
		} else {
			if len(vals) <= 0 {
				err = errutil.Fmt("attempting to call an event %q with no target  %q", name, aff)
			} else {
				ret, err = run.send(name, rec, vals[0])
			}
		}
	}
	return
}

func (run *Runner) call(rec *g.Record, kind cachedKind, aff affine.Affinity) (ret g.Value, err error) {
	if field, e := pattern.GetResultField(run, kind.Kind); e != nil {
		err = e
	} else {
		name := kind.Name()
		oldScope := run.scope.ReplaceScope(scope.FromRecord(rec))
		run.currentPatterns.startedPattern(name)
		if e := kind.recordInit(run, rec); e != nil {
			err = e
		} else if rules, e := run.getRules(name); e != nil {
			err = e
		} else if res, e := rules.Call(run, rec, field); e != nil {
			err = e
		} else {
			ret, err = res.GetResult(aff)
		}
		run.currentPatterns.stoppedPattern(name)
		run.scope.RestoreScope(oldScope)
	}
	return
}

func (run *Runner) send(name string, rec *g.Record, tgt g.Value) (ret g.Value, err error) {
	callState := run.saveCallState()
	if path, e := run.newPathForTarget(tgt); e != nil {
		err = e
	} else if evtObj, e := newEventRecord(run, event.Object, tgt); e != nil {
		err = e // ^ create the "event" object sent to each event phase pattern.
	} else {
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
func (run *Runner) newPhase(k cachedKind, src *g.Record) (ret *g.Record, err error) {
	dst, srcKind := k.NewRecord(), src.Kind()
	_ = run.scope.ReplaceScope(scope.FromRecord(dst)) // assumes the caller handles scope restoration

	i := 0
	// copy until the fields are mismatched
	for srcCnt, dstCnt := srcKind.NumField(), k.NumField(); i < srcCnt && i < dstCnt; i++ {
		if srcKind.Field(i) != k.Field(i) {
			break
		} else if src.HasValue(i) { // no need to copy fields that were never set
			if v, e := src.GetIndexedField(i); e != nil {
				err = e
				break
			} else if e := dst.SetIndexedField(i, g.CopyValue(v)); e != nil {
				err = e
				break
			}
		} else if initCnt := len(k.init); i < initCnt { // init any unset fieds
			if e := k.initIndex(run, dst, i); e != nil {
				err = e
				break
			}
		}
	}
	// init any remaining fields
	for initCnt := len(k.init); i < initCnt; i++ {
		if e := k.initIndex(run, dst, i); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		ret = dst
	}
	return
}
