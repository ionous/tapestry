package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *ExtendPattern) Execute(macro rt.Runtime) error {
	return weave.StoryStatement(macro, op)
}

func (op *ExtendPattern) Schedule(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireDependencies, func(w *weave.Weaver) (err error) {
		if name, e := safe.GetText(w, op.PatternName); e != nil {
			err = e
		} else {
			// tbd: assert declares it into existence
			// how do we instead simply say it should exist?
			pattern := name.String()
			if e := cat.AssertAncestor(pattern, kindsOf.Pattern.String()); e != nil {
				err = e
			} else if e := declareFields(op.Locals, func(name, class string, aff affine.Affinity, init assign.Assignment) (err error) {
				return cat.AssertField(pattern, name, class, aff, init)
			}); e != nil {
				err = e
			} else {
				// write the rules last to help with test output consistency
				err = ImportRules(cat, pattern, "", op.Rules, assert.DefaultTiming)
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefinePattern) Execute(macro rt.Runtime) error {
	return weave.StoryStatement(macro, op)
}

// Adds a new pattern declaration and optionally some associated pattern parameters.
func (op *DefinePattern) Schedule(cat *weave.Catalog) (err error) {
	return cat.Schedule(assert.RequireDependencies, func(w *weave.Weaver) (err error) {
		if name, e := safe.GetText(w, op.PatternName); e != nil {
			err = e
		} else {
			pattern := name.String()
			if e := cat.AssertAncestor(pattern, kindsOf.Pattern.String()); e != nil {
				err = e
			} else {
				// fix: probably always want to declare a result; even if its "nothing".
				var resList []FieldDefinition
				if res := op.Result; res != nil {
					resList = []FieldDefinition{res}
				}
				if e := declareFields(resList, func(name, class string, aff affine.Affinity, init assign.Assignment) error {
					return cat.AssertResult(pattern, name, class, aff, init)
				}); e != nil {
					err = e
				} else if e := declareFields(op.Params, func(name, class string, aff affine.Affinity, init assign.Assignment) error {
					return cat.AssertParam(pattern, name, class, aff, init)
				}); e != nil {
					err = e
				} else if e := declareFields(op.Locals, func(name, class string, aff affine.Affinity, init assign.Assignment) error {
					return cat.AssertField(pattern, name, class, aff, init)
				}); e != nil {
					err = e
				} else {
					// write the rules last to help with test output consistency
					err = ImportRules(cat, pattern, "", op.Rules, assert.DefaultTiming)
				}
			}
		}
		return
	})
}

// note:  statements can set flags for a bunch of rules at once or within each rule separately, but not both.
func ImportRules(cat *weave.Catalog, pattern, target string, els []PatternRule, flags assert.EventTiming) (err error) {
	// write in reverse order because within a given pattern, earlier rules take precedence.
	for i := len(els) - 1; i >= 0; i-- {
		if e := els[i].importRule(cat, pattern, target, flags); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (op *PatternRule) importRule(cat *weave.Catalog, pattern, target string, tgtFlags assert.EventTiming) (err error) {
	act := op.Does
	if flags, e := op.Flags.ReadFlags(); e != nil {
		err = e
	} else if flags > 0 && tgtFlags > 0 {
		// ensure flags were only set via the rule or via the pattern
		err = errutil.New("unexpected continuation flags in", pattern)
	} else {
		if tgtFlags > 0 {
			flags = tgtFlags
		} else if flags == 0 {
			flags = assert.During
		}

		// check if this rule is declared inside a specific domain
		if guard, ok := op.Guard.(jsn.Marshalee); !ok {
			err = errutil.New("missing guard in", pattern)
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
			err = cat.AssertRule(pattern, target, op.Guard, flags, act)
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

func declareFields(els []FieldDefinition, ft fieldType) (err error) {
	for _, el := range els {
		if e := el.DeclareField(ft); e != nil {
			err = e
			break
		}
	}
	return
}
