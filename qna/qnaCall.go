package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
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
			if ret, err = run.call(rec, pat, aff); err != nil {
				// .call can return both a value and an error.
				err = errutil.Fmt("%w calling %q", err, name)
			}
		} else {
			callState := run.saveCallState()
			if len(keys) <= 0 {
				err = errutil.Fmt("attempting to call an event %q with no target  %q", name, aff)
			} else if tgt, e := rec.GetNamedField(keys[0]); e != nil {
				err = e
			} else if path, e := run.newPathForTarget(tgt); e != nil {
				err = e
			} else if evtObj, e := newEventRecord(run, event.Object, tgt); e != nil {
				err = e // ^ create the "event" object sent to each event phase pattern.
			} else {
				for prevRec, i := rec, 0; i < event.NumPhases; i++ {
					phase := event.Phase(i)
					if kindForPhase, e := run.getKind(phase.PatternName(name)); e != nil {
						err = e
						break
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
							if allDone, e := run.send(evtObj, kindForPhase, path.slice(phase)); e != nil {
								err = errutil.Fmt("%w calling %q", e, kindForPhase.Name())
								break
							} else if allDone {
								break
							}
							prevRec = phaseRec
						}
					}
				}
			}
			callState.restore()
		}
	}
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
	_ = run.replaceScope(scope.FromRecord(dst)) // assumes the caller handles scope restoration

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

// return true if all done ( canceled )
func (run *Runner) send(evtObj *g.Record, kind cachedKind, chain []string) (retDone bool, err error) {
	var stopPropogation bool // stopPropogation -- finish the current event level, and no more.
	for tgtIdx, cnt := 0, len(chain); tgtIdx < cnt && err == nil && !stopPropogation; tgtIdx++ {
		tgt := chain[tgtIdx]
		if e := evtObj.SetIndexedField(event.CurrentTarget.Index(), g.StringOf(tgt)); e != nil {
			err = e
			break
		} else if rs, e := run.getRules(kind.Name(), ""); e != nil {
			err = e // fix? get rules earlier to skip setup if they dont exist?
			break
		} else {
			for ruleIdx := range rs.rules {
				if done, e := rs.applyRule(run, ruleIdx); e != nil || done {
					err = e
					break
				} else {
					// setting cancel blocks the rest of processing
					// clearing it lets the handler set continue
					// ( so it can be controlled via stopPropogation/Immediate similar to the dom )
					if evtObj.HasValue(event.Cancel.Index()) {
						retDone = true
						if cancel, e := evtObj.GetIndexedField(event.Cancel.Index()); e != nil {
							err = e
							break
						} else if cancel.Bool() {
							stopPropogation = true
							break // hard exit targets and rules
						}
					}

					// setting interrupt stops all handlers in this flow;
					// clearing this blocks allows this level of the hierarchy to finish, then stops.
					if evtObj.HasValue(event.Interupt.Index()) {
						stopPropogation = true // stops target loop, but not rule loop.
						if interrupt, e := evtObj.GetIndexedField(event.Interupt.Index()); e != nil {
							err = e
						} else if interrupt.Bool() {
							break
						}
					}
				}
			}
		}
	}
	return
}

func (run *Runner) call(rec *g.Record, kind cachedKind, aff affine.Affinity) (ret g.Value, err error) {
	name := kind.Name()
	if labels, e := run.GetField(meta.PatternLabels, name); e != nil {
		err = e
	} else {
		var flags rt.Flags
		res := pattern.NewResults(rec, labels.Strings(), aff)
		oldScope := run.replaceScope(res)
		run.currentPatterns.startedPattern(name)
		if e := kind.recordInit(run, rec); e != nil {
			err = e
		} else if rules, e := run.GetRules(name, "", &flags); e != nil {
			err = e
		} else if e := res.ApplyRules(run, rules, flags); e != nil {
			err = e
		} else if v, e := res.GetResult(); e != nil {
			err = e
		} else {
			// warning: in order to generate appropriate defaults ( ex. a record of the right type )
			// while still informing the caller of lack of pattern decision in a concise manner
			// can return both a valid value and an error
			ret = v
			if !res.ComputedResult() {
				err = errutil.Fmt("%w computing %s", rt.NoResult, aff)
			}
		}
		run.currentPatterns.stoppedPattern(name)
		run.restoreScope(oldScope)
	}
	return
}

func (rs *ruleSet) applyRule(run rt.Runtime, i int) (done bool, err error) {
	rule := rs.rules[i]
	if ok, e := safe.GetOptionalBool(run, rule.Filter, true); e != nil {
		err = e
	} else if ok.Bool() && !rs.skipRun {
		if e := safe.RunAll(run, rule.Execute); e != nil {
			err = e
		} else if rule.Terminates {
			if !rs.updateAll {
				done = true
			}
			rs.skipRun = true
		}
	}
	return
}
