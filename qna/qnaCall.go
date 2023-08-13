package qna

import (
	"golang.org/x/exp/slices"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/action"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"git.sr.ht/~ionous/tapestry/rt/safe"
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
		} else if !pat.Implements(action.AtEvent.Kind().String()) {
			// note: this doesnt positively affirm kindsOf.Pattern:
			// some tests use golang structs as faux patterns.
			if ret, err = run.call(rec, pat, aff, 0); err != nil {
				// .call can return both a value and an error.
				err = errutil.Fmt("%w calling %q", err, name)
			}
		} else {
			// get the request event target, and the chain of objects up to the root.
			if target, e := rec.GetIndexedField(action.Target.Index()); e != nil {
				err = e
			} else if els, e := run.Call(action.CapturePattern, affine.TextList, nil, []g.Value{target}); e != nil {
				err = e
			} else {
				chain, order := els.Strings(), action.Bubbles // fix? capture actual returns bubble order right now.
				prev, base := rec, 0
				for evt := action.FirstEvent; evt < action.NumEvents; evt++ {
					if kindOfEvent, e := run.getKind(evt.Name(name)); e != nil {
						err = e
					} else if next, e := copyEventObject(kindOfEvent.Kind, prev); e != nil {
						err = e
						break
					} else {
						// call or send depending on the flow:
						if flow := evt.Flow(); flow == action.Targets {
							if _, e := run.call(next, kindOfEvent, affine.None, base); e != nil {
								err = e
								break
							}
						} else {
							if flow != order {
								slices.Reverse(chain)
								order = order.Reverse()
							}
							if allDone, e := run.send(next, kindOfEvent, chain, base); e != nil {
								err = errutil.Fmt("%w calling %q", e, kindOfEvent.Name())
								break
							} else if allDone {
								break
							}
						}
						// re: base. once the initial fields have been initialized we dont have to initialize them again.
						prev, base = next, action.NumFields
					}
				}
			}
		}
	}
	return
}

// copies any fields in src to dst
// the action pattern declaration implies that the initial set of field sare the same
// ( even though the author can extend those patterns with non-overlapping locals )
// fix: it'd be nicer to backdoor record creation to slice the values across
func copyEventObject(dstKind *g.Kind, src *g.Record) (ret *g.Record, err error) {
	dst, srcKind := dstKind.NewRecord(), src.Kind()
	max := srcKind.NumField()
	if dstcnt := dstKind.NumField(); dstcnt < max {
		max = dstcnt
	}
	for i := 0; i < max; i++ {
		if i >= action.NumFields && srcKind.Field(i) != dstKind.Field(i) {
			// the first NumFields are guaranteed to be the same;
			// after those, the author might have customized the individual patterns.
			// even if rare, its not an error; its simply the end of the shared fields.
			break
		} else if src.HasValue(i) { // no need to set fields that were never set
			if v, e := src.GetIndexedField(i); e != nil {
				err = e
				break
			} else if e := dst.SetIndexedField(i, g.CopyValue(v)); e != nil {
				err = e
				break
			}
		}
	}
	if err == nil {
		ret = dst
	}
	return
}

// return true if all done ( canceled )
func (run *Runner) send(rec *g.Record, kind cachedKind, chain []string, base int) (retDone bool, err error) {
	var flags rt.Flags
	oldScope := run.Stack.ReplaceScope(pattern.NewResults(rec, nil, affine.None))
	run.currentPatterns.startedPattern(kind.Name())
	// run any record init statements
	// ( with the scope of the pattern so that locals can refer to parameters )
	if e := kind.recordInit(run, rec, base); e != nil {
		err = e
	} else {
		var stopPropogation bool // stopPropogation -- finish the current event level, and no more.
		for i, cnt := 0, len(chain); i < cnt && err == nil && !stopPropogation; i++ {
			tgt := chain[i]
			if rec.SetIndexedField(action.CurrentTarget.Index(), g.StringOf(tgt)); e != nil {
				err = e
				break
			} else if rules, e := run.GetRules(kind.Name(), "", &flags); e != nil {
				err = e
				break
			} else {
				for _, rule := range rules {
					// end if there are no flags left, and we didn't want to filter everything.
					if _, e := pattern.ApplyRule(run, rule, flags); e != nil {
						err = errutil.New(e, "while applying", rule.Name)
						break
					}

					// setting cancel blocks the rest of processing
					// setting it false lets the handler set continue
					// ( so it can be controlled via stopPropogation/Immediate similar to the dom )
					if rec.HasValue(action.Cancel.Index()) {
						retDone = true
						if cancel, e := rec.GetIndexedField(action.Cancel.Index()); e != nil {
							err = e
							break
						} else if cancel.Bool() {
							stopPropogation = true
							break //
						}
					}

					// setting stopPropogation blocks rest of the bubble / cancel
					// setting it true stops all other events on this noun.
					if rec.HasValue(action.Interupt.Index()) {
						stopPropogation = true
						if interrupt, e := rec.GetIndexedField(action.Interupt.Index()); e != nil {
							err = e
						} else if interrupt.Bool() {
							break
						}
					}
				}
			}
		}
	}
	// restore
	run.currentPatterns.stoppedPattern(kind.Name())
	run.Stack.ReplaceScope(oldScope)
	return
}

func (run *Runner) call(rec *g.Record, kind cachedKind, aff affine.Affinity, base int) (ret g.Value, err error) {
	name := kind.Name()
	if labels, e := run.GetField(meta.PatternLabels, name); e != nil {
		err = e
	} else {
		var flags rt.Flags
		res := pattern.NewResults(rec, labels.Strings(), aff)
		oldScope := run.Stack.ReplaceScope(res)
		run.currentPatterns.startedPattern(name)
		if e := kind.recordInit(run, rec, base); e != nil {
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
		run.Stack.ReplaceScope(oldScope)
	}
	return
}

// a copy of call using the new flags......
func (run *Runner) newCall(rec *g.Record, kind cachedKind, aff affine.Affinity, base int) (ret g.Value, err error) {
	name := kind.Name()
	if labels, e := run.GetField(meta.PatternLabels, name); e != nil {
		err = e
	} else {
		res := pattern.NewResults(rec, labels.Strings(), aff)
		oldScope := run.Stack.ReplaceScope(res)
		run.currentPatterns.startedPattern(name)
		if e := kind.recordInit(run, rec, base); e != nil {
			err = e
		} else if e := run.applyRules(name); e != nil {
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
		run.Stack.ReplaceScope(oldScope)
	}
	return
}

func (run *Runner) applyRules(name string) (err error) {
	if rs, e := run.getRules(lang.Normalize(name), ""); e != nil {
		err = e
	} else {
		alwaysUpdate, allowRun := rs.alwaysUpdate, true
		for _, rule := range rs.rules {
			if ok, e := safe.GetOptionalBool(run, rule.Filter, true); e != nil {
				err = e
			} else if allowRun && ok.Bool() {
				if e := safe.RunAll(run, rule.Execute); e != nil {
					err = e
				} else if rule.Terminates {
					if !alwaysUpdate {
						break
					}
					allowRun = false // turn off the run flag, but keep going
				}
			}
		}
	}
	return
}
