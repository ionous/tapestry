package rules

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/dl/object"
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
		exe = []rt.Execute{&logic.ChooseBranch{
			Condition: &logic.IsAll{Test: filters},
			Exe:       exe,
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
		&math.CompareText{
			A:       object.Variable(event.Actor),
			Compare: math.C_Comparison_EqualTo,
			B:       literal.T("self"),
		})
}

// filter to the innermost target.
func AddEventFilters(filters []rt.BoolEval) (ret []rt.BoolEval) {
	return append(filters,
		&math.CompareText{
			A:       object.Variable(event.Object, event.CurrentTarget.String()),
			Compare: math.C_Comparison_EqualTo,
			B:       object.Variable(event.Object, event.Target.String()),
		})
}

func AddNounFilter(noun string, filters []rt.BoolEval) (ret []rt.BoolEval) {
	return append(filters,
		&math.CompareText{
			A:       object.Variable(event.Object, event.Target.String()),
			Compare: math.C_Comparison_EqualTo,
			B:       &literal.TextValue{Value: noun},
		})
}

func AddKindFilter(kind string, filters []rt.BoolEval) (ret []rt.BoolEval) {
	return append(filters,
		&object.IsKindOf{
			Object: object.Variable(event.Object, event.Target.String()),
			Kind:   kind,
		})
}
