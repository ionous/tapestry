package rules

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// fix: move to jess or break into parts, some in rule
func (rule RuleInfo) WeaveRule(w weaver.Weaves, filters []rt.BoolEval, exe []rt.Execute) error {
	label := rule.Label
	if len(label) == 0 {
		label = FindNamedResponse(exe)
	}
	// add the filters to the execute
	// fix? this is a little questionable; probably is mucking up the runtime some.
	// better might be to *extract* user booleans and put them into this filter list.
	if len(filters) > 0 {
		exe = []rt.Execute{&core.ChooseBranch{
			If:  &core.AllTrue{Test: filters},
			Exe: exe,
		}}
	}
	//
	pb := mdl.NewPatternBuilder(rule.Name)
	pb.AppendRule(rule.Rank, rt.Rule{
		Name:    label,
		Exe:     exe,
		Updates: DoesUpdate(exe),
		Stop:    rule.Stop,
		Jump:    rule.Jump,
	})
	return w.ExtendPattern(pb.Pattern)
}

func AddPlayerFilter(filters []rt.BoolEval) (ret []rt.BoolEval) {
	/*if k, e := run.GetKindByName(rule.Name); e != nil {
		err = fmt.Errorf("finding base pattern %q %q %s", rule.Name, rule.Label, e)
	} else
		disabled for now: hard test...

		if canFilterActor := CanFilterActor(k); rule.ExcludesPlayer && !canFilterActor {
			err = fmt.Errorf("only actor events can filter by actor for pattern %q rule %q", rule.Name, rule.Label)
		} else */

	// if the focus of the event involves an actor;
	// then we automatically filter for the player
	// if !rule.ExcludesPlayer {
	// 	filters =
	// }
	return append(filters,
		&core.CompareText{
			A:  core.Variable(event.Actor),
			Is: core.C_Comparison_EqualTo,
			B:  core.T("self"),
		})
}

// filter to the innermost target.
// eventLike := k.Implements(kindsOf.Action.String())
// if eventLike {
func AddEventFilters(filters []rt.BoolEval) (ret []rt.BoolEval) {
	return append(filters,
		&core.CompareText{
			A:  core.Variable(event.Object, event.CurrentTarget.String()),
			Is: core.C_Comparison_EqualTo,
			B:  core.Variable(event.Object, event.Target.String()),
		})
}
