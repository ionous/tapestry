package weave

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/rules"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// fix: move to jess or break into parts, some in rule
func (cat *Catalog) WeaveRule(rule rules.RuleInfo, filter rt.BoolEval, exe []rt.Execute) (err error) {
	run := cat.GetRuntime()
	if pb, e := rule.WeaveRule(run, filter, exe); e != nil {
		err = e
	} else {
		err = cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) error {
			return w.ExtendPattern(pb.Pattern)
		})
	}
	return
}
