package story

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/rules"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *DefinePattern) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// Adds a new pattern declaration and optionally some associated pattern parameters.
func (op *DefinePattern) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weaver.VerbPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if name, e := safe.GetText(run, op.PatternName); e != nil {
			err = e
		} else {
			name := inflect.Normalize(name.String())
			pb := mdl.NewPatternBuilder(name)
			pb.AddParams(reduceFields(run, op.Requires))
			// assumes the first field ( if any ) is the result
			if ps := op.Provides; len(ps) > 0 {
				pb.AddResult(ps[0].GetFieldInfo(run))
				pb.AddLocals(reduceFields(run, ps[1:]))
			}
			if len(op.Exe) > 0 {
				pb.AppendRule(0, rt.Rule{
					Name:    fmt.Sprintf("the default %s rule", name),
					Exe:     op.Exe,
					Updates: rules.DoesUpdate(op.Exe),
					// Stop/Jump is 0/0 by default;
					// and so is the rule Rank
				})
			}
			err = w.AddPattern(pb.Pattern)
		}
		return
	})
}

func (op *DefineAction) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineAction) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.VerbPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if act, e := safe.GetText(run, op.Action); e != nil {
			err = e
		} else {
			act := mdl.NewPatternSubtype(inflect.Normalize(act.String()), kindsOf.Action.String())
			// note: actions dont have an explicit return
			act.AddParams(reduceFields(run, op.Requires))
			act.AddLocals(reduceFields(run, op.Provides))
			if e := w.AddPattern(act.Pattern); e != nil {
				err = e
			} else {
				// derive the before and after phases
				for _, phase := range []event.Phase{event.BeforePhase, event.AfterPhase} {
					pb := mdl.NewPatternSubtype(phase.PatternName(act.Name()), act.Name())
					if e := w.AddPattern(pb.Pattern); e != nil {
						err = e
						break
					}
				}
			}
		}
		return
	})
}

func (op *RuleProvides) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *RuleProvides) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weaver.VerbPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if act, e := safe.GetText(run, op.PatternName); e != nil {
			err = e
		} else if ks, e := run.GetField(meta.KindAncestry, act.String()); e != nil {
			err = e // ^ verify the kind exists
		} else {
			ks := ks.Strings()
			act := ks[len(ks)-1] // get the kind's real name (ex. plural fixup)
			pb := mdl.NewPatternSubtype(act, kindsOf.Action.String())
			pb.AddLocals(reduceFields(run, op.Provides))
			err = w.AddPattern(pb.Pattern)
		}
		return
	})
}

func (op *RuleForPattern) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *RuleForPattern) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weaver.VerbPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if phrase, e := safe.GetText(run, op.PatternName); e != nil {
			err = e
		} else if label, e := safe.GetOptionalText(run, op.RuleName, ""); e != nil {
			err = e
		} else {
			rule := rules.ReadPhrase(phrase.String(), label.String())
			if rule.IsDomainEvent() {
				// are we in the domain?
				domainName, eventName := rule.Short, rule.EventName()
				if v, e := run.GetField(meta.Domain, domainName); e == nil && v.Bool() {
					// cheat by adding the pattern as if it were in the root domain
					// regardless of where we are.
					pb := mdl.NewPatternBuilder(eventName)
					err = w.AddPattern(pb.Pattern)
				}
			}
			if err == nil {
				if e := weaveRule(cat, w, run, rule, nil, op.Exe); e != nil {
					err = errutil.Fmt("%w weaving a rule for a pattern", e)
				}
			}
		}
		return
	})
}

func (op *RuleForNoun) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *RuleForNoun) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weaver.VerbPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if noun, e := safe.GetText(run, op.NounName); e != nil {
			err = e
		} else if noun, e := run.GetField(meta.ObjectId, noun.String()); e != nil {
			err = e
		} else if phrase, e := safe.GetText(run, op.PatternName); e != nil {
			err = e
		} else if label, e := safe.GetOptionalText(run, op.RuleName, ""); e != nil {
			err = e
		} else if rule := rules.ReadPhrase(phrase.String(), label.String()); rule.IsDomainEvent() {
			err = errutil.New("can't target nouns for domain events")
		} else {
			filter := &core.CompareText{
				A:  core.Variable(event.Object, event.Target.String()),
				Is: core.C_Comparison_EqualTo,
				B:  &literal.TextValue{Value: noun.String()},
			}
			if e := weaveRule(cat, w, run, rule, filter, op.Exe); e != nil {
				err = errutil.Fmt("%w weaving a rule for a noun", e)
			}
		}
		return
	})
}

func (op *RuleForKind) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *RuleForKind) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weaver.VerbPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if kind, e := safe.GetText(run, op.KindName); e != nil {
			err = e
		} else if ks, e := run.GetField(meta.KindAncestry, kind.String()); e != nil {
			err = e // ^ verify the kind exists
		} else if exact, e := safe.GetOptionalBool(run, op.Exactly, false); e != nil {
			err = e
		} else if phrase, e := safe.GetText(run, op.PatternName); e != nil {
			err = e
		} else if label, e := safe.GetOptionalText(run, op.RuleName, ""); e != nil {
			err = e
		} else if rule := rules.ReadPhrase(phrase.String(), label.String()); rule.IsDomainEvent() {
			err = errutil.New("can't target kinds for domain events")
		} else {
			ks := ks.Strings()
			k := ks[len(ks)-1] // get the kind's real name (ex. plural fixup)
			var filter rt.BoolEval
			if exact.Bool() {
				filter = &core.IsExactKindOf{Object: core.Variable(event.Object, event.Target.String()),
					Kind: k,
				}
			} else {
				filter = &core.IsKindOf{Object: core.Variable(event.Object, event.Target.String()),
					Kind: k,
				}
			}
			if e := weaveRule(cat, w, run, rule, filter, op.Exe); e != nil {
				err = errutil.Fmt("%w weaving a rule for a kind", e)
			}
		}
		return
	})
}

type ruleNoun string

type ruleKind struct {
	name    string
	exactly bool
}

func weaveRule(cat *weave.Catalog, w weaver.Weaves, run rt.Runtime, rule rules.RuleName, filter rt.BoolEval, exe []rt.Execute) (err error) {
	if info, e := rule.GetRuleInfo(run); e != nil {
		err = e
	} else if k, e := run.GetKindByName(info.Name); e != nil {
		err = errutil.Fmt("finding base pattern %q %q %s", info.Name, info.Label, e)
	} else if canFilterActor := rules.CanFilterActor(k); info.ExcludesPlayer && !canFilterActor {
		rules.CanFilterActor(k)
		err = errutil.Fmt("only actor events can filter by actor for pattern %q rule %q", info.Name, info.Label)
	} else {
		updates := rules.DoesUpdate(exe)
		// term := rules.DoesTerminate(exe)

		label := info.Label
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
			if !info.ExcludesPlayer {
				filters = append(filters,
					&core.CompareText{
						A:  core.Variable(event.Actor),
						Is: core.C_Comparison_EqualTo,
						B:  T("self"),
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
		pb := mdl.NewPatternBuilder(info.Name)
		pb.AppendRule(info.Rank, rt.Rule{
			Name:    label,
			Exe:     exe,
			Updates: updates,
			Stop:    info.Stop,
			Jump:    info.Jump,
		})
		err = cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) error {
			return w.ExtendPattern(pb.Pattern)
		})
	}
	return
}
