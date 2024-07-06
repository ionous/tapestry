package story

import (
	"fmt"

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
)

// Adds a new pattern declaration and optionally some associated pattern parameters.
func (op *DefinePattern) Weave(cat *weave.Catalog) (err error) {
	return cat.ScheduleCmd(op, weaver.PropertyPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
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

func (op *DefineAction) Weave(cat *weave.Catalog) error {
	return cat.ScheduleCmd(op, weaver.VerbPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if act, e := safe.GetText(run, op.PatternName); e != nil {
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

func (op *DefinePatternProvides) Weave(cat *weave.Catalog) (err error) {
	return cat.ScheduleCmd(op, weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if act, e := safe.GetText(run, op.PatternName); e != nil {
			err = e
		} else if ks, e := run.GetField(meta.KindAncestry, act.String()); e != nil {
			err = e // ^ verify the kind exists
		} else if ks := ks.Strings(); len(ks) == 0 {
			err = fmt.Errorf("%w kind %s", weaver.ErrMissing, act)
		} else {
			pb := mdl.NewPatternSubtype(ks[0], kindsOf.Action.String())
			pb.AddLocals(reduceFields(run, op.Provides))
			err = w.AddPattern(pb.Pattern)
		}
		return
	})
}

func (op *DefineRule) Weave(cat *weave.Catalog) (err error) {
	return cat.ScheduleCmd(op, weaver.VerbPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if phrase, e := safe.GetText(run, op.RuleTiming); e != nil {
			err = e
		} else if label, e := safe.GetOptionalText(run, op.RuleName, ""); e != nil {
			err = e
		} else if desc, e := rules.ReadPhrase(run, phrase.String(), label.String()); e != nil {
			err = e
		} else if rule, e := desc.GetRuleInfo(); e != nil {
			err = e
		} else {
			err = cat.ScheduleCmd(op, weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) error {
				return rule.WeaveRule(w, nil, op.Exe)
			})
		}
		return
	})
}

func (op *DefineNounRule) Weave(cat *weave.Catalog) (err error) {
	return cat.ScheduleCmd(op, weaver.VerbPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if noun, e := safe.GetText(run, op.NounName); e != nil {
			err = e
		} else if noun, e := run.GetField(meta.ObjectId, noun.String()); e != nil {
			err = e
		} else if phrase, e := safe.GetText(run, op.RuleTiming); e != nil {
			err = e
		} else if label, e := safe.GetOptionalText(run, op.RuleName, ""); e != nil {
			err = e
		} else if desc, e := rules.ReadPhrase(run, phrase.String(), label.String()); e != nil {
			err = e
		} else if rule, e := desc.GetRuleInfo(); e != nil {
			err = e
		} else {
			filters := rules.AddNounFilter(noun.String(), nil)
			if desc.IsEvent() {
				if !desc.ExcludesPlayer {
					filters = rules.AddPlayerFilter(filters)
				}
				filters = rules.AddEventFilters(filters)
			}
			err = cat.ScheduleCmd(op, weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) error {
				return rule.WeaveRule(w, filters, op.Exe)
			})
		}
		return
	})
}

func (op *DefineKindRule) Weave(cat *weave.Catalog) (err error) {
	return cat.ScheduleCmd(op, weaver.VerbPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if kind, e := safe.GetText(run, op.KindName); e != nil {
			err = e
		} else if ks, e := run.GetField(meta.KindAncestry, kind.String()); e != nil {
			err = e // ^ verify the kind exists
		} else if ks := ks.Strings(); len(ks) == 0 {
			err = fmt.Errorf("%w kind %s", weaver.ErrMissing, kind)
		} else if phrase, e := safe.GetText(run, op.RuleTiming); e != nil {
			err = e
		} else if label, e := safe.GetOptionalText(run, op.RuleName, ""); e != nil {
			err = e
		} else if desc, e := rules.ReadPhrase(run, phrase.String(), label.String()); e != nil {
			err = e
		} else if rule, e := desc.GetRuleInfo(); e != nil {
			err = e
		} else {
			filters := rules.AddKindFilter(ks[0], nil)
			if desc.IsEvent() {
				if !desc.ExcludesPlayer {
					filters = rules.AddPlayerFilter(filters)
				}
				filters = rules.AddEventFilters(filters)
			}
			err = cat.ScheduleCmd(op, weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) error {
				return rule.WeaveRule(w, filters, op.Exe)
			})
		}
		return
	})
}
