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
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *DefinePattern) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// Adds a new pattern declaration and optionally some associated pattern parameters.
func (op *DefinePattern) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weave.RequireAncestry, func(w *weave.Weaver) (err error) {
		if name, e := safe.GetText(cat.Runtime(), op.PatternName); e != nil {
			err = e
		} else {
			name := inflect.Normalize(name.String())
			pb := mdl.NewPatternBuilder(name)
			pb.AddParams(reduceFields(w, op.Requires))
			// assumes the first field ( if any ) is the result
			if ps := op.Provides; len(ps) > 0 {
				pb.AddResult(ps[0].GetFieldInfo(w))
				pb.AddLocals(reduceFields(w, ps[1:]))
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
			err = cat.Schedule(weave.RequireAncestry, func(w *weave.Weaver) error {
				return w.Pin().AddPattern(pb.Pattern)
			})
		}
		return
	})
}

func (op *DefineAction) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineAction) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) (err error) {
		if act, e := safe.GetText(w, op.Action); e != nil {
			err = e
		} else {
			act := mdl.NewPatternSubtype(inflect.Normalize(act.String()), kindsOf.Action.String())
			// note: actions dont have an explicit return
			act.AddParams(reduceFields(w, op.Requires))
			act.AddLocals(reduceFields(w, op.Provides))
			if e := cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) error {
				return w.Pin().AddPattern(act.Pattern)
			}); e != nil {
				err = e
			} else {
				// derive the before and after phases
				for _, phase := range []event.Phase{event.BeforePhase, event.AfterPhase} {
					pb := mdl.NewPatternSubtype(phase.PatternName(act.Name()), act.Name())
					if e := cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) error {
						return w.Pin().AddPattern(pb.Pattern)
					}); e != nil {
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
	return cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) (err error) {
		if act, e := safe.GetText(w, op.PatternName); e != nil {
			err = e
		} else if act, e := w.Pin().GetKind(inflect.Normalize(act.String())); e != nil {
			err = e // ^ verify the kind exists
		} else {
			pb := mdl.NewPatternSubtype(act, kindsOf.Action.String())
			pb.AddLocals(reduceFields(w, op.Provides))
			err = cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) error {
				return w.Pin().AddPattern(pb.Pattern)
			})
		}
		return
	})
}

func (op *RuleForPattern) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *RuleForPattern) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) (err error) {
		if phrase, e := safe.GetText(w, op.PatternName); e != nil {
			err = e
		} else if label, e := safe.GetOptionalText(w, op.RuleName, ""); e != nil {
			err = e
		} else {
			rule := rules.ReadPhrase(phrase.String(), label.String())
			if rule.IsDomainEvent() {
				// are we in the domain?
				domainName, eventName := rule.Short, rule.EventName()
				if v, e := w.GetField(meta.Domain, domainName); e == nil && v.Bool() {
					// cheat by adding the pattern as if it were in the root domain
					// regardless of where we are.
					pin := cat.Modeler.Pin(domainName, w.At)
					pb := mdl.NewPatternBuilder(eventName)
					err = pin.AddPattern(pb.Pattern)
				}
			}
			if err == nil {
				if e := weaveRule(w, rule, nil, op.Exe); e != nil {
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
	return cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) (err error) {
		if noun, e := safe.GetText(w, op.NounName); e != nil {
			err = e
		} else if noun, e := w.Pin().GetClosestNoun(noun.String()); e != nil {
			err = e
		} else if phrase, e := safe.GetText(w, op.PatternName); e != nil {
			err = e
		} else if label, e := safe.GetOptionalText(w, op.RuleName, ""); e != nil {
			err = e
		} else if rule := rules.ReadPhrase(phrase.String(), label.String()); rule.IsDomainEvent() {
			err = errutil.New("can't target nouns for domain events")
		} else {
			filter := &core.CompareText{
				A:  core.Variable(event.Object, event.Target.String()),
				Is: core.C_Comparison_EqualTo,
				B:  &literal.TextValue{Value: noun},
			}
			if e := weaveRule(w, rule, filter, op.Exe); e != nil {
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
	return cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) (err error) {
		if kind, e := safe.GetText(w, op.KindName); e != nil {
			err = e
		} else if k, e := w.Pin().GetKind(inflect.Normalize(kind.String())); e != nil {
			err = e // ^ verify the kind exists
		} else if exact, e := safe.GetOptionalBool(w, op.Exactly, false); e != nil {
			err = e
		} else if phrase, e := safe.GetText(w, op.PatternName); e != nil {
			err = e
		} else if label, e := safe.GetOptionalText(w, op.RuleName, ""); e != nil {
			err = e
		} else if rule := rules.ReadPhrase(phrase.String(), label.String()); rule.IsDomainEvent() {
			err = errutil.New("can't target kinds for domain events")
		} else {
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
			if e := weaveRule(w, rule, filter, op.Exe); e != nil {
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

func weaveRule(w *weave.Weaver, rule rules.RuleName, filter rt.BoolEval, exe []rt.Execute) (err error) {
	if info, e := rule.GetRuleInfo(w); e != nil {
		err = e
	} else if k, e := w.GetKindByName(info.Name); e != nil {
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
		err = w.Catalog.Schedule(weave.RequirePatterns, func(w *weave.Weaver) error {
			return w.Pin().ExtendPattern(pb.Pattern)
		})
	}
	return
}
