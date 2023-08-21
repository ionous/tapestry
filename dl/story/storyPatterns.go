package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
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
			pb := mdl.NewPatternBuilder(name.String())
			addFields(pb, mdl.PatternLocals, op.Locals)
			addFields(pb, mdl.PatternParameters, op.Params)
			addOptionalField(pb, mdl.PatternResults, op.Result)
			addRules(pb, "", op.Rules)
			err = cat.Schedule(weave.RequireAncestry, func(w *weave.Weaver) error {
				return w.Pin().AddPattern(pb.Pattern)
			})
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *ExtendPattern) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *ExtendPattern) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) (err error) {
		if name, e := safe.GetText(cat.Runtime(), op.PatternName); e != nil {
			err = e
		} else if name, e := w.Pin().GetKind(lang.Normalize(name.String())); e != nil {
			err = e
		} else if k, e := w.GetKindByName(name); e != nil {
			err = e
		} else {
			pb := mdl.NewPatternBuilder(name)
			addFields(pb, mdl.PatternLocals, op.Locals)
			addRules(pb, "", op.Rules)
			err = cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) (err error) {
				if !k.Implements(kindsOf.Action.String()) {
					err = w.Pin().ExtendPattern(pb.Pattern)
				} else {
					for i := 0; i < event.NumPhases; i++ {
						// fix: we copy all the initialization too
						// maybe an explicit ExtendPatternSet that in mdl to do better things.
						name := event.Phase(i).PatternName(name)
						if e := w.Pin().ExtendPattern(pb.Copy(name)); e != nil {
							err = e
							break
						}
					}
				}
				return
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
			act := lang.Normalize(act.String())
			for i := 0; i < event.NumPhases; i++ {
				phase := event.Phase(i)
				pb := mdl.NewPatternSubtype(phase.PatternName(act), phase.PatternKind())
				addFields(pb, mdl.PatternParameters, op.Params)
				addFields(pb, mdl.PatternLocals, op.Locals)
				err = cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) error {
					return w.Pin().AddPattern(pb.Pattern)
				})
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
			err = errutil.New("finding base pattern", e)
		} else {
			updates := rules.DoesUpdate(exe)
			cancels := prefix.Cancels
			interrupts := cancels || rules.DoesTerminate(exe)

			// add additional filters:
			filters := make([]rt.BoolEval, 0, 3)
			if filter != nil {
				filters = append(filters, filter)
			}
			// by default: all event handlers are filtered to the player and the innermost target.
			if k.Implements(kindsOf.Event.String()) || k.Implements(kindsOf.Action.String()) {
				// if the focus of the event involves an actor;
				// then we automatically filter for the player
				if k.NumField() > 0 {
					if f := k.Field(0); f.Type == event.Actors {
						filters = append(filters,
							&core.CompareText{
								A:  core.Variable(f.Name),
								Is: core.Equal,
								B:  T("self"),
							})
					}
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
			pb.AppendRule(mdl.Rule{
				Name: rule,
				Rank: rank,
				Prog: assign.Prog{
					Exe:        exe,
					Updates:    updates,
					Interrupts: interrupts,
					Cancels:    cancels,
				},
			})
			err = w.Catalog.Schedule(weave.RequirePatterns, func(w *weave.Weaver) error {
				return w.Pin().ExtendPattern(pb.Pattern)
			})
		}
	}
	return
}

// note:  statements can set flags for a bunch of rules at once or within each rule separately, but not both.
func ImportRules(pb *mdl.PatternBuilder, target string, els []PatternRule) {
	// write in reverse order because within a given pattern, earlier rules take precedence.
	for i := len(els) - 1; i >= 0; i-- {
		els[i].addRule(pb, target)
	}
}

func (op *PatternRule) addRule(pb *mdl.PatternBuilder, target string) {
	pb.AppendRule(mdl.Rule{Target: target, Prog: assign.Prog{
		Filter:  op.Guard,
		Exe:     op.Exe,
		Updates: rules.FilterHasCounter(op.Guard),
	}})
}

func addOptionalField(pb *mdl.PatternBuilder, ft mdl.FieldType, field FieldDefinition) {
	if field != nil {
		var empty mdl.FieldInfo // the Nothing type generates a blank field info
		if f := field.FieldInfo(); f != empty {
			pb.AddField(ft, f)
		}
	}
	return
}

func addFields(pb *mdl.PatternBuilder, ft mdl.FieldType, fields []FieldDefinition) {
	// fix; should probably be an error if nothing is used for locals
	// or if nothing exists in a list of more than one nothing parameter
	for _, field := range fields {
		var empty mdl.FieldInfo // the Nothing type generates a blank field info
		if f := field.FieldInfo(); f != empty {
			pb.AddField(ft, f)
		}
	}
}

// note:  statements can set flags for a bunch of rules at once or within each rule separately, but not both.
func addRules(pb *mdl.PatternBuilder, target string, els []PatternRule) {
	// write in reverse order because within a given pattern, earlier rules take precedence.
	for i := len(els) - 1; i >= 0; i-- {
		els[i].addRule(pb, target)
	}
}
