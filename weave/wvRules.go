package weave

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/rules"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// fix: move to jess or break into parts, some in rule
func (cat *Catalog) WeaveRule(rule rules.RuleInfo, filter rt.BoolEval, exe []rt.Execute) (err error) {
	run := cat.GetRuntime()
	if k, e := run.GetKindByName(rule.Name); e != nil {
		err = fmt.Errorf("finding base pattern %q %q %s", rule.Name, rule.Label, e)
	} else if canFilterActor := rules.CanFilterActor(k); rule.ExcludesPlayer && !canFilterActor {
		rules.CanFilterActor(k)
		err = fmt.Errorf("only actor events can filter by actor for pattern %q rule %q", rule.Name, rule.Label)
	} else {
		updates := rules.DoesUpdate(exe)
		// term := rules.DoesTerminate(exe)

		label := rule.Label
		if len(label) == 0 {
			label = rules.FindNamedResponse(exe)
		}

		// add additional filters:
		filters := make([]rt.BoolEval, 0, 3)
		if filter != nil {
			filters = append(filters, filter)
		}
		// by default: all event handlers are filtered to the player and the innermost target.
		eventLike := k.Implements(kindsOf.Action.String())
		if eventLike {
			// if the focus of the event involves an actor;
			// then we automatically filter for the player
			if !rule.ExcludesPlayer {
				filters = append(filters,
					&core.CompareText{
						A:  core.Variable(event.Actor),
						Is: core.C_Comparison_EqualTo,
						B:  core.T("self"),
					})
			}
			// filter to the innermost target.
			filters = append(filters,
				&core.CompareText{
					A:  core.Variable(event.Object, event.CurrentTarget.String()),
					Is: core.C_Comparison_EqualTo,
					B:  core.Variable(event.Object, event.Target.String()),
				})
		}
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
			Updates: updates,
			Stop:    rule.Stop,
			Jump:    rule.Jump,
		})
		// tbd: is scheduling needed here?
		// callers are in the VerbPhrase
		err = cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) error {
			return w.ExtendPattern(pb.Pattern)
		})
	}
	return
}
