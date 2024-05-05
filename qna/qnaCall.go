package qna

import (
	"slices"

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
		switch pattern.Categorize(pat.Path()) {
		case pattern.Initializes:
			ret = rt.RecordOf(rec)

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
			// run through keys to see if the target was specified
			if ft := pat.Field(target); !slices.Contains(keys, ft.Name) {
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
		} else if res, e := rules.Calls(run, newScope, field); e != nil {
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
	// fix? setup the base state first since that's where most initialization lives
	// however, it's not satisfying to move from "action", to "before action", back to "action" again
	// tbd: add an explicit middle pattern ( "when action" )
	callState := run.saveCallState(rec)
	if path, e := run.newPathForTarget(tgt); e != nil {
		err = e // ^ the bubble capture chain
	} else if evtObj, e := newEventRecord(run, event.Object, tgt); e != nil {
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
					// push the event scope. we dont pop later; we replace.
					// ( the event gets pushed last so it cant be hidden by pattern locals )
					scoped := scope.NewReadOnlyValue(event.Object, rt.RecordOf(evtObj))
					run.PushScope(scoped)
					if ok, e := rules.Sends(run, scoped, path.slice(phase)); e != nil {
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

// fix: serialize and deserialize might be better
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
