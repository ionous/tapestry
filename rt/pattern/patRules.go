package pattern

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

type RuleSet struct {
	rules     []rt.Rule
	updateAll bool
}

func (rs *RuleSet) AddRule(rule rt.Rule) {
	rs.rules = append(rs.rules, rule)
	rs.updateAll = rule.Updates
}

// assumes scope is initialized
func (rs *RuleSet) Calls(run rt.Runtime, scope rt.Scope, resultField string) (res Result, err error) {
	stopJump := stopJump{jump: rt.JumpLater}
	for i := range rs.rules {
		if res, e := rs.tryRule(run, i); e != nil {
			err = e
			break
		} else {
			gotResult := len(resultField) == 0 || scope.FieldChanged(resultField)
			if done, e := stopJump.update(res, nil, gotResult); e != nil || done {
				break
			}
		}
	}
	if err == nil {
		res = Result{scope: scope, field: resultField, hasResult: stopJump.runCount > 0}
	}
	return
}

// trigger a pattern for each of the targets in the passed chain.
// return false if stopped
func (rs *RuleSet) Sends(run rt.Runtime, scope rt.Scope, chain []string) (okay bool, err error) {
	stopJump := stopJump{jump: rt.JumpLater}
	for tgtIdx, cnt := 0, len(chain); tgtIdx < cnt && stopJump.jump == rt.JumpLater; tgtIdx++ {
		tgt := chain[tgtIdx]
		if e := scope.SetFieldByName(event.CurrentTarget.String(), rt.StringOf(tgt)); e != nil {
			err = e
			break
		} else if e := rs.send(run, scope, &stopJump); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		okay = !stopJump.stop
	}
	return
}

// event handlers dont have return values, so whenever they match it may be the end
// ( depending on values in the db determined during weave based on phase )
func (rs *RuleSet) send(run rt.Runtime, evtObj rt.Scope, stopJump *stopJump) (err error) {
	for i := range rs.rules {
		if res, e := rs.tryRule(run, i); e != nil {
			err = e
			break
		} else if done, e := stopJump.update(res, evtObj, true); e != nil || done {
			break
		}
	}
	return
}

func (rs *RuleSet) tryRule(run rt.Runtime, i int) (ret stopJump, err error) {
	var pushes int
	var prog []rt.Execute
	rule := rs.rules[i] // copies
	// scan for the first matching case
	// if none apply then the rule isn't considered to have been run
	if tree := core.PickTree(rule.Exe); tree == nil {
		prog = rule.Exe
	} else if branch, e := tree.PickBranch(run, &pushes); e != nil {
		err = e
	} else {
		prog = branch
	}
	if err == nil && prog != nil {
		// println("- ", rule.Name)
		var ri core.DoInterrupt
		switch e := safe.RunAll(run, prog); {
		case e == nil:
			ret = stopJump{
				runCount: 1,
				jump:     rule.Jump,
				stop:     rule.Stop,
			}
		case errors.As(e, &ri):
			if ri.KeepGoing {
				ret = stopJump{
					runCount: 1,
					jump:     rt.JumpLater,
					stop:     false,
				}
			} else {
				ret = stopJump{
					runCount: 1,
					jump:     rt.JumpNow,
					stop:     true,
				}
			}
		default:
			err = e
		}
	}
	safe.PopSeveral(run, pushes)
	return
}
