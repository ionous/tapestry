package pattern

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
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

// assumes rec is already in scope and initialized
func (rs *RuleSet) Call(run rt.Runtime, rec *g.Record, resultField int) (res Result, err error) {
	stopJump := stopJump{jump: rt.JumpLater}
	for i, rule := range rs.rules {
		if ranRule, e := rs.tryRule(run, i); e != nil {
			err = e
			break
		} else {
			// if there's a return value, it must be set for the pattern to be considered done.
			if ranRule && (resultField < 0 || rec.HasValue(resultField)) {
				stopJump.ranRule(rule.Stop, rule.Jump)
			}
			if stopJump.jump == rt.JumpNow {
				break
			}
		}
	}
	if err == nil {
		res = Result{rec: rec, field: resultField, hasResult: stopJump.runCount > 0}
	}
	return
}

// trigger a pattern for each of the targets in the passed chain.
// return false if stopped
func (rs *RuleSet) Send(run rt.Runtime, evtObj *g.Record, chain []string) (okay bool, err error) {
	stopJump := stopJump{jump: rt.JumpLater}
	for tgtIdx, cnt := 0, len(chain); tgtIdx < cnt && stopJump.jump == rt.JumpLater; tgtIdx++ {
		tgt := chain[tgtIdx]
		if e := evtObj.SetIndexedField(event.CurrentTarget.Index(), g.StringOf(tgt)); e != nil {
			err = e
			break
		} else if e := rs.send(run, evtObj, &stopJump); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		okay = !stopJump.stop
	}
	return
}

func (rs *RuleSet) send(run rt.Runtime, evtObj *g.Record, stopJump *stopJump) (err error) {
	for i, rule := range rs.rules {
		if ranRule, e := rs.tryRule(run, i); e != nil {
			err = e
			break
		} else {
			// event handlers dont have return values, so when they match it may be the end
			// ( depending on values in the db determined during weave based on phase )
			if ranRule {
				stopJump.ranRule(rule.Stop, rule.Jump)
			}
			if e := stopJump.mergeEvent(evtObj); e != nil {
				err = e
				break
			} else if stopJump.jump == rt.JumpNow {
				break
			}
		}
	}
	return
}

func (rs *RuleSet) tryRule(run rt.Runtime, i int) (okay bool, err error) {
	var pushes int
	var prog []rt.Execute
	rule := rs.rules[i]
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
		if e := safe.RunAll(run, prog); e != nil {
			err = e
		} else {
			okay = true
		}
	}
	safe.PopSeveral(run, pushes)
	return
}
