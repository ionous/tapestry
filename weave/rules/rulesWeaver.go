package rules

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
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

// ensure that the target of an action is the player
// ( jess apples this filter to actor actions unless specifically asked not to )
func AddPlayerFilter(filters []rt.BoolEval) (ret []rt.BoolEval) {
	return append(filters,
		&core.CompareText{
			A:  assign.Variable(event.Actor),
			Is: core.C_Comparison_EqualTo,
			B:  core.T("self"),
		})
}

// filter to the innermost target.
func AddEventFilters(filters []rt.BoolEval) (ret []rt.BoolEval) {
	return append(filters,
		&core.CompareText{
			A:  assign.Variable(event.Object, event.CurrentTarget.String()),
			Is: core.C_Comparison_EqualTo,
			B:  assign.Variable(event.Object, event.Target.String()),
		})
}

func AddNounFilter(noun string, filters []rt.BoolEval) (ret []rt.BoolEval) {
	return append(filters,
		&core.CompareText{
			A:  assign.Variable(event.Object, event.Target.String()),
			Is: core.C_Comparison_EqualTo,
			B:  &literal.TextValue{Value: noun},
		})
}

func AddKindFilter(kind string, filters []rt.BoolEval) (ret []rt.BoolEval) {
	return append(filters,
		&core.IsKindOf{
			Object: assign.Variable(event.Object, event.Target.String()),
			Kind:   kind,
		})
}
