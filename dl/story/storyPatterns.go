package story

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story/internal/rules"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
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
			if e := addFields(pb, mdl.PatternLocals, op.Locals); e != nil {
				err = e
			} else if e := addFields(pb, mdl.PatternParameters, op.Params); e != nil {
				err = e
			} else if e := addOptionalField(pb, mdl.PatternResults, op.Result); e != nil {
				err = e
			} else if e := addRules(pb, "", op.Rules); e != nil {
				err = e
			} else if e := cat.Schedule(weave.RequireAncestry, func(w *weave.Weaver) error {
				return w.Pin().AddPattern(pb.Pattern)
			}); e != nil {
				err = e
			}
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
			if e := addFields(pb, mdl.PatternLocals, op.Locals); e != nil {
				err = e
			} else if e := addRules(pb, "", op.Rules); e != nil {
				err = e
			} else {
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
				if e := addFields(pb, mdl.PatternParameters, op.Params); e != nil {
					err = e
				} else if e := addFields(pb, mdl.PatternLocals, op.Locals); e != nil {
					err = e
				} else if e := cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) error {
					return w.Pin().AddPattern(pb.Pattern)
				}); e != nil {
					err = errutil.Fmt("%w defining action %q", e, act)
					break
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
		} else if e := weaveRule(w, act.String(), rule.String(), ruleNoun(noun), op.Exe); e != nil {
			err = errutil.Fmt("%w weaving a rule", e)
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
		} else if e := weaveRule(w, act.String(), rule.String(), ruleKind{k, exact.Bool()}, op.Exe); e != nil {
			err = errutil.Fmt("%w weaving a rule", e)
		}
		return
	})
}

type ruleNoun string

type ruleKind struct {
	name    string
	exactly bool
}

func weaveRule(w *weave.Weaver, pat, rule string, extra any, exe []rt.Execute) (err error) {
	if prefix, e := rules.ReadName(w, lang.Normalize(pat)); e != nil {
		err = errutil.New("determining prefix", e)
	} else {
		pat, rank := prefix.Name, prefix.Rank
		if k, e := w.GetKindByName(pat); e != nil {
			err = errutil.New("finding base pattern", e)
		} else {
			equals := make([]rt.BoolEval, 0, 3)

			// by default: all event handlers are filtered to the player and the innermost target.
			if k.Implements(kindsOf.Event.String()) || k.Implements(kindsOf.Action.String()) {

				// if the focus of the event involves an actor;
				// then we automatically filter for the player
				if k.NumField() > 0 {
					if f := k.Field(0); f.Type == event.Actors {
						equals = append(equals,
							&core.CompareText{
								A:  core.Variable(f.Name),
								Is: core.Equal,
								B:  T("self"),
							})
					}
				}
				equals = append(equals,
					&core.CompareText{
						A:  core.Variable(event.Object, event.CurrentTarget.String()),
						Is: core.Equal,
						B:  core.Variable(event.Object, event.Target.String()),
					})

				switch extra := extra.(type) {
				case nil:
				case ruleNoun:
					equals = append(equals,
						&core.CompareText{
							A:  core.Variable(event.Object, event.Target.String()),
							Is: core.Equal,
							B:  &literal.TextValue{Value: string(extra)},
						})
				case ruleKind:
					if extra.exactly {
						equals = append(equals,
							&core.IsExactKindOf{Object: core.Variable(event.Object, event.Target.String()),
								Kind: extra.name,
							})
					} else {
						equals = append(equals,
							&core.IsKindOf{Object: core.Variable(event.Object, event.Target.String()),
								Kind: extra.name,
							})
					}

				default:
					panic(errutil.Sprint("unknown guard %T", extra))
				}

			}
			exe = []rt.Execute{&core.ChooseAction{
				If:  &core.AllTrue{Test: equals},
				Exe: exe,
			}}
		}
		updates := rules.DoesUpdate(exe)
		terminates := /*prefix.Terminates ||*/ rules.DoesTerminate(exe)
		pb := mdl.NewPatternBuilder(pat)
		pb.AddNewRule(lang.Normalize(rule), rank, updates, terminates, exe)
		err = w.Catalog.Schedule(weave.RequirePatterns, func(w *weave.Weaver) (err error) {
			if e := w.Pin().ExtendPattern(pb.Pattern); e != nil {
				err = errutil.New("extending pattern", e)
			}
			return
		})
	}
	return
}

// note:  statements can set flags for a bunch of rules at once or within each rule separately, but not both.
func ImportRules(pb *mdl.PatternBuilder, target string, els []PatternRule) (err error) {
	// write in reverse order because within a given pattern, earlier rules take precedence.
	for i := len(els) - 1; i >= 0; i-- {
		if e := els[i].addRule(pb, target); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (op *PatternRule) addRule(pb *mdl.PatternBuilder, target string) (err error) {
	act := op.Exe
	// check if this rule is declared inside a specific domain
	if guard, ok := op.Guard.(jsn.Marshalee); !ok {
		err = errutil.New("missing guard", pb.Name())
	} else {
		update := rules.SearchForCounters(guard)
		pb.AddRule(target, op.Guard, update, act)
	}
	return
}

func addOptionalField(pb *mdl.PatternBuilder, ft mdl.FieldType, field FieldDefinition) (_ error) {
	if field != nil {
		var empty mdl.FieldInfo // the Nothing type generates a blank field info
		if f := field.FieldInfo(); f != empty {
			pb.AddField(ft, f)
		}
	}
	return
}

func addFields(pb *mdl.PatternBuilder, ft mdl.FieldType, fields []FieldDefinition) (_ error) {
	// fix; should probably be an error if nothing is used for locals
	// or if nothing exists in a list of more than one nothing parameter
	for _, field := range fields {
		var empty mdl.FieldInfo // the Nothing type generates a blank field info
		if f := field.FieldInfo(); f != empty {
			pb.AddField(ft, f)
		}
	}
	return
}

// note:  statements can set flags for a bunch of rules at once or within each rule separately, but not both.
func addRules(pb *mdl.PatternBuilder, target string, els []PatternRule) (err error) {
	// write in reverse order because within a given pattern, earlier rules take precedence.
	for i := len(els) - 1; i >= 0; i-- {
		if e := els[i].addRule(pb, target); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
