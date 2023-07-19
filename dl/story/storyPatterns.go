package story

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *ExtendPattern) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *ExtendPattern) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(assert.RequireDependencies, func(w *weave.Weaver) (err error) {
		if name, e := safe.GetText(cat.Runtime(), op.PatternName); e != nil {
			err = e
		} else {
			pb := mdl.NewPatternBuilder(name.String())
			if e := addFields(pb, mdl.PatternLocals, op.Locals); e != nil {
				err = e
			} else if e := addRules(pb, "", op.Rules, assert.DefaultTiming); e != nil {
				err = e
			} else {
				err = w.Pin().ExtendPattern(pb.Pattern)
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefinePattern) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// Adds a new pattern declaration and optionally some associated pattern parameters.
func (op *DefinePattern) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(assert.RequireDependencies, func(w *weave.Weaver) (err error) {
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
			} else if e := addRules(pb, "", op.Rules, assert.DefaultTiming); e != nil {
				err = e
			} else {
				err = w.Pin().AddPattern(pb.Pattern)
			}
		}
		return
	})
}

// note:  statements can set flags for a bunch of rules at once or within each rule separately, but not both.
func ImportRules(pb *mdl.PatternBuilder, target string, els []PatternRule, flags assert.EventTiming) (err error) {
	// write in reverse order because within a given pattern, earlier rules take precedence.
	for i := len(els) - 1; i >= 0; i-- {
		if e := els[i].addRule(pb, target, flags); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (op *PatternRule) addRule(pb *mdl.PatternBuilder, target string, tgtFlags assert.EventTiming) (err error) {
	act := op.Does
	if flags, e := op.Flags.ReadFlags(); e != nil {
		err = e
	} else if flags > 0 && tgtFlags > 0 {
		// ensure flags were only set via the rule or via the pattern
		err = errutil.New("unexpected continuation flags", pb.Name())
	} else {
		if tgtFlags > 0 {
			flags = tgtFlags
		} else if flags == 0 {
			flags = assert.During
		}

		// check if this rule is declared inside a specific domain
		if guard, ok := op.Guard.(jsn.Marshalee); !ok {
			err = errutil.New("missing guard", pb.Name())
		} else {
			// fix? could we instead just strstr for countOf
			// also might be cool to augment or replace the serialized type
			// with our own that has an pre-calced field ( at import, via state parser )
			if SearchForCounters(guard) {
				flags |= assert.RunAlways
			}
			// fix via runtime? check if this rule is declared inside a specific domain
			// if domain != k.Env().Game.Domain {
			// 	guard = &core.AllTrue{[]rt.BoolEval{
			// 		&core.HasDominion{domain.String()},
			// 		guard,
			// 	}}
			pb.AddRule(target, op.Guard, flags, act)
		}
	}
	return
}

func (op *PatternFlags) ReadFlags() (ret assert.EventTiming, err error) {
	switch str := op.Str; str {
	case PatternFlags_Before:
		ret = assert.After // run other matching patterns, and then run this pattern. other...this.
	case PatternFlags_After:
		ret = assert.Before // keep going after running the current pattern. this...others.
	case PatternFlags_Terminate:
		ret = assert.During
	default:
		if len(str) > 0 {
			err = errutil.Fmt("unknown pattern flags %q", str)
		}
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
func addRules(pb *mdl.PatternBuilder, target string, els []PatternRule, flags assert.EventTiming) (err error) {
	// write in reverse order because within a given pattern, earlier rules take precedence.
	for i := len(els) - 1; i >= 0; i-- {
		if e := els[i].addRule(pb, target, flags); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
