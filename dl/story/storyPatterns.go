package story

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/safe"
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
			name := lang.Normalize(name.String())
			pb := mdl.NewPatternBuilder(name)
			if e := addRequiredFields(w, pb, op.Requires); e != nil {
				err = e
			} else if e := addProvidingFields(w, pb, op.Provides); e != nil {
				err = e
			} else {
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
			act := lang.Normalize(act.String())
			for i := 0; i < event.NumPhases; i++ {
				phase := event.Phase(i)
				pb := mdl.NewPatternSubtype(phase.PatternName(act), phase.PatternKind())
				// note: actions dont have an explicit return
				if e := addRequiredFields(w, pb, op.Requires); e != nil {
					err = e
				} else if e := addFields(w, pb, mdl.PatternLocals, op.Provides); e != nil {
					err = e
				} else {
					err = cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) error {
						return w.Pin().AddPattern(pb.Pattern)
					})
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
		} else {
			if k, e := w.Pin().GetKind(lang.Normalize(act.String())); e != nil {
				err = e // ^ verify the kind exists
			} else {
				for i := 0; i < event.NumPhases; i++ {
					phase := event.Phase(i)
					pb := mdl.NewPatternSubtype(phase.PatternName(k), phase.PatternKind())
					if e := addFields(w, pb, mdl.PatternLocals, op.Provides); e != nil {
						err = e
					} else {
						err = cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) error {
							return w.Pin().AddPattern(pb.Pattern)
						})
					}
				}
			}
		}
		return
	})

}

func (op *RuleForPattern) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *RuleForPattern) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) (err error) {
		if act, e := safe.GetText(w, op.PatternName); e != nil {
			err = e
		} else if rule, e := safe.GetOptionalText(w, op.RuleName, ""); e != nil {
			err = e
		} else if e := weaveRule(w, act.String(), rule.String(), nil, op.Exe); e != nil {
			err = errutil.Fmt("%w weaving a rule", e)
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
		} else if act, e := safe.GetText(w, op.PatternName); e != nil {
			err = e
		} else if rule, e := safe.GetOptionalText(w, op.RuleName, ""); e != nil {
			err = e
		} else {
			filter := &core.CompareText{
				A:  core.Variable(event.Object, event.Target.String()),
				Is: core.Equal,
				B:  &literal.TextValue{Value: noun},
			}
			if e := weaveRule(w, act.String(), rule.String(), filter, op.Exe); e != nil {
				err = errutil.Fmt("%w weaving a rule", e)
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
		} else if k, e := w.Pin().GetKind(lang.Normalize(kind.String())); e != nil {
			err = e // ^ verify the kind exists
		} else if exact, e := safe.GetOptionalBool(w, op.Exactly, false); e != nil {
			err = e
		} else if act, e := safe.GetText(w, op.PatternName); e != nil {
			err = e
		} else if rule, e := safe.GetOptionalText(w, op.RuleName, ""); e != nil {
			err = e
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
			if e := weaveRule(w, act.String(), rule.String(), filter, op.Exe); e != nil {
				err = errutil.Fmt("%w weaving a rule", e)
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

func weaveRule(w *weave.Weaver, pat, rule string, filter rt.BoolEval, exe []rt.Execute) (err error) {
	if prefix, e := rules.ReadName(w, lang.Normalize(pat)); e != nil {
		err = errutil.New("determining prefix", e)
	} else {
		pat, rank := prefix.Name, prefix.Rank
		if k, e := w.GetKindByName(pat); e != nil {
			err = errutil.Fmt("finding base pattern %q %q %s", pat, rule, e)
		} else if canFilterActor := rules.CanFilterActor(k); prefix.ExcludesPlayer && !canFilterActor {
			rules.CanFilterActor(k)
			err = errutil.Fmt("only actor events can filter by actor for pattern %q rule %q", pat, rule)
		} else {
			updates := rules.DoesUpdate(exe)
			// term := rules.DoesTerminate(exe)

			if rule = lang.Normalize(rule); len(rule) == 0 {
				rule = rules.FindNamedResponse(exe)
			}

			// add additional filters:
			filters := make([]rt.BoolEval, 0, 3)
			if filter != nil {
				filters = append(filters, filter)
			}
			// by default: all event handlers are filtered to the player and the innermost target.
			eventLike := k.Implements(kindsOf.Event.String()) || k.Implements(kindsOf.Action.String())
			if eventLike {
				// if the focus of the event involves an actor;
				// then we automatically filter for the player
				if !prefix.ExcludesPlayer {
					filters = append(filters,
						&core.CompareText{
							A:  core.Variable(event.Actor),
							Is: core.Equal,
							B:  T("self"),
						})
				}
				// filter to the innermost target.
				filters = append(filters,
					&core.CompareText{
						A:  core.Variable(event.Object, event.CurrentTarget.String()),
						Is: core.Equal,
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
			pb := mdl.NewPatternBuilder(pat)
			pb.AppendRule(rank, rt.Rule{
				Name:    rule,
				Exe:     exe,
				Updates: updates,
				Stop:    prefix.Stop,
				Jump:    prefix.Jump,
			})
			err = w.Catalog.Schedule(weave.RequirePatterns, func(w *weave.Weaver) error {
				return w.Pin().ExtendPattern(pb.Pattern)
			})
		}
	}
	return
}

func addRequiredFields(run rt.Runtime, pb *mdl.PatternBuilder, fields []FieldDefinition) (err error) {
	return addFields(run, pb, mdl.PatternParameters, fields)
}

// assumes the first field ( if any ) is the return value
// i'm sure i'll come to hate this, but for now i like that it simplifies specification
func addProvidingFields(run rt.Runtime, pb *mdl.PatternBuilder, fields []FieldDefinition) (err error) {
	if len(fields) > 0 {
		if e := addOptionalField(run, pb, mdl.PatternResults, fields[0]); e != nil {
			err = e
		} else if e := addFields(run, pb, mdl.PatternLocals, fields[1:]); e != nil {
			err = e
		}
	}
	return
}

func addOptionalField(run rt.Runtime, pb *mdl.PatternBuilder, ft mdl.FieldType, field FieldDefinition) (err error) {
	if field != nil {
		var empty mdl.FieldInfo // the Nothing type generates a blank field info
		if f, e := field.FieldInfo(run); e != nil {
			err = errutil.Append(err, e)
		} else if f != empty {
			pb.AddField(ft, f)
		}
	}
	return
}

func addFields(run rt.Runtime, pb *mdl.PatternBuilder, ft mdl.FieldType, fields []FieldDefinition) (err error) {
	// fix; should probably be an error if nothing is used for locals
	// or if nothing exists in a list of more than one nothing parameter
	for _, field := range fields {
		var empty mdl.FieldInfo // the Nothing type generates a blank field info
		if f, e := field.FieldInfo(run); e != nil {
			err = errutil.Append(err, e)
		} else if f != empty {
			pb.AddField(ft, f)
		}
	}
	return
}
