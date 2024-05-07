package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"git.sr.ht/~ionous/tapestry/rt/scope"
	"github.com/ionous/errutil"
)

// note: this is mirrored/mimicked in package testpat
// note: in order to generate appropriate defaults ( ex. a record of the right type )
// can return both a both meaningful value *and* an error
func (run *Runner) Call(name string, aff affine.Affinity, keys []string, vals []rt.Value) (ret rt.Value, err error) {
	// create a container to hold results of args, locals, and the pending return value
	if pat, e := run.getKind(name); e != nil {
		err = e
	} else if rec, e := pattern.InitRecord(run, pat, keys, vals); e != nil {
		err = e
	} else {
		switch pattern.Categorize(pat.Ancestors()) {
		case pattern.Initializes:
			ret = rt.RecordOf(rec)

		case pattern.Calls:
			ret, err = run.call(pat, rec, aff)

		case pattern.Sends:
			// fix: currently, if the pattern contains "noun"
			// then use that as the event target;
			// otherwise, use the first parameter.
			// alt: maybe putting noun first, or renaming noun to target
			var target int
			if i := pat.FieldIndex("noun"); i >= 0 {
				target = i
			}
			if tgt, e := rec.GetIndexedField(target); e != nil {
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

func (run *Runner) call(kind *rt.Kind, rec *rt.Record, aff affine.Affinity) (ret rt.Value, err error) {
	if field, e := pattern.GetResultField(run, kind); e != nil {
		err = e
	} else {
		name := kind.Name()
		newScope := scope.FromRecord(run, rec)
		oldScope := run.scope.ReplaceScope(newScope)
		run.currentPatterns.startedPattern(name)
		if rules, e := run.getRules(name); e != nil {
			err = e
		} else if res, e := rules.Calls(run, rec, field); e != nil {
			err = e
		} else {
			ret, err = res.GetResult(run, aff)
		}
		run.currentPatterns.stoppedPattern(name)
		run.scope.RestoreScope(oldScope)
	}
	return
}

func (run *Runner) send(name string, rec *rt.Record, tgt rt.Value) (ret rt.Value, err error) {
	// sets the base pattern first since that's where most initialization lives
	callState := run.saveCallState(rec)
	if path, e := run.newPathForTarget(tgt); e != nil {
		err = e // ^ the bubble capture chain
	} else if evtObj, e := newEventRecord(name, tgt); e != nil {
		err = e // ^ create the "event object" sent to each event phase pattern.
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
				// create the phase record, set it into scope.
				// doesn't initialize locals in the subphase
				// fix? but im hoping to get rid of the subphases
				// its all *very* confusing
				// ( the event object isnt available till the main pattern body )
				phaseRec := rt.NewRecord(rec.Kind)
				if e := copyPhase(phaseRec, prevRec); e != nil {
					err = e
					break
				} else {
					// replace the previous phase data
					run.scope.ReplaceScope(scope.FromRecord(run, phaseRec))
					// add the event object, accessed via the name "event" ( ex. "event.target" )
					// pushed last, it takes precedence over locals if any have the same name.
					run.PushScope(scope.NewReadOnlyValue(event.Object, rt.RecordOf(evtObj)))
					if ok, e := rules.Sends(run, evtObj, path.slice(phase)); e != nil {
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
		ret = rt.BoolOf(!canceled)
	}
	callState.restore()
	return
}

func copyPhase(dst *rt.Record, src *rt.Record) (err error) {
	for i, srcCnt, dstCnt := 0, src.NumField(), dst.NumField(); i < srcCnt && i < dstCnt; i++ {
		// copy until the fields are mismatched
		if src.Field(i).Name != dst.Field(i).Name {
			break
		} else if v, e := src.GetIndexedField(i); e != nil {
			if !rt.IsNilRecord(e) {
				err = e
				break
			}
		} else if e := dst.SetIndexedField(i, rt.CopyValue(v)); e != nil {
			err = e
			break
		}
	}
	return
}
