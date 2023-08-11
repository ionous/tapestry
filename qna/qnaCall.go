package qna

import (
	"golang.org/x/exp/slices"

	"git.sr.ht/~ionous/tapestry/affine"
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
	if kind, e := run.getKind(name); e != nil {
		err = e
	} else if rec, e := safe.FillRecord(run, kind.NewRecord(), keys, vals); e != nil {
		err = e
	} else {
		// Call is being used as a side effect to initialize records
		if kind.Implements(kindsOf.Record.String()) {
			if aff != affine.Record {
				err = errutil.Fmt("attempting to call a record %q with affine %q", name, aff)
			} else {
				ret = g.RecordOf(rec)
			}
		} else if !kind.Implements(kindsOf.Event.String()) { // <- fix should be action if anything
			// note: this doesnt positively affirm kindsOf.Pattern:
			// some tests use golang structs as faux patterns.
			if ret, err = run.call(rec, kind, aff); err != nil {
				// .call can return both a value and an error.
				err = errutil.Fmt("%w calling %q", err, name)
			}
		} else {
			// get the request event target
			if target, e := rec.GetIndexedField(action.Target.Index()); e != nil {
				err = e
			} else if els, e := run.Call(action.CapturePattern, affine.TextList, nil, []g.Value{target}); e != nil {
				err = e
			} else {
				targets := els.Strings()
				for i, event := range action.EventNames(name) {
					if kind, e := run.getKind(event); e != nil {
						err = e
					} else {
						ctx := rec
						if event == name {
							// FX FIX FIX CREATE THE PROPER REC:
						}
						// before, after capture/bubble through multiple nouns; the others do not.
						if i&1 == 0 {
							if _, e := run.call(ctx, kind, affine.None); e != nil {
								err = e
								break
							}
						} else {
							slices.Reverse(targets) // capture should (eventually) return top to bottom
							if allDone, e := run.send(ctx, kind, targets); e != nil {
								err = errutil.Fmt("%w calling %q", e, name)
								break
							} else if allDone {
								break
							}
						}
					}
				}
			}
		}
	}
	return
}

// return true if all done ( canceled )
func (run *Runner) send(rec *g.Record, kind cachedKind, chain []string) (retDone bool, err error) {
	var flags rt.Flags
	oldScope := run.Stack.ReplaceScope(pattern.NewResults(rec, nil, affine.None))
	run.currentPatterns.startedPattern(kind.Name())

	if e := kind.initializeRecord(run, rec); e != nil {
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

func (run *Runner) call(rec *g.Record, kind cachedKind, aff affine.Affinity) (ret g.Value, err error) {
	name := kind.Name()
	if labels, e := run.GetField(meta.PatternLabels, name); e != nil {
		err = e
	} else {
		var flags rt.Flags
		res := pattern.NewResults(rec, labels.Strings(), aff)
		oldScope := run.Stack.ReplaceScope(res)
		run.currentPatterns.startedPattern(name)
		if e := kind.initializeRecord(run, rec); e != nil {
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
