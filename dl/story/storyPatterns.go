package story

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/jsn"
	"github.com/ionous/errutil"
)

func (op *PatternActions) ImportPhrase(k *Importer) (err error) {
	if res, e := convertRes(op.PatternReturn); e != nil {
		err = e
	} else {
		patternName := op.Name.String()
		var locals []eph.EphParams
		if els := op.PatternLocals; els != nil {
			locals, err = els.ImportLocals(k, patternName)
		}
		if err == nil {
			k.Write(&eph.EphPatterns{Name: patternName, Locals: locals, Result: res})
			// write the rules last ( order doesnt matter except for tests )
			err = op.PatternRules.ImportRules(k, patternName, "", eph.EphTiming{})
		}
	}
	return
}

// Adds a new pattern declaration and optionally some associated pattern parameters.
func (op *PatternDecl) ImportPhrase(k *Importer) (err error) {
	if res, e := convertRes(op.PatternReturn); e != nil {
		err = e
	} else {
		patternName := op.Name.String()

		// tell the system about the pattern subtype ( if any )
		if patternType, ok := composer.FindChoice(&op.Type, op.Type.Str); ok {
			// fix: subtypes ( do we really need them? ) end with (s)
			// but the actual kinds do not :/
			patternType = patternType[:len(patternType)-1]
			if patternType != eph.KindsOfPattern {
				k.Write(&eph.EphKinds{Kinds: patternName, From: patternType})
			}
		}
		//
		var params []eph.EphParams
		if els := op.Optvars; els != nil && err == nil {
			if ps, e := reduceParams(els.VariableDecl); e != nil {
				err = e
			} else {
				params = ps
			}
		}
		if err == nil {
			k.Write(&eph.EphPatterns{Name: patternName, Result: res, Params: params})
		}
	}
	return
}

func (op *PatternVariablesDecl) ImportPhrase(k *Importer) (err error) {
	if ps, e := reduceParams(op.VariableDecl); e != nil {
		err = e
	} else {
		k.Write(&eph.EphPatterns{Name: op.PatternName.String(), Params: ps})
	}
	return
}

func (op *PatternRules) ImportRules(k *Importer, pattern, target string, flags eph.EphTiming) (err error) {
	if els := op.PatternRule; els != nil {
		for _, el := range els {
			if e := el.importRule(k, pattern, target, flags); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

func (op *PatternRule) importRule(k *Importer, pattern, target string, tgtFlags eph.EphTiming) (err error) {
	if act, e := op.Hook.ImportProgram(k); e != nil {
		err = e
	} else if flags, e := op.Flags.ReadFlags(); e != nil {
		err = e
	} else if len(flags.Str) > 0 && len(tgtFlags.Str) > 0 {
		// ensure flags were only set via the rule or via the pattern
		err = errutil.New("unexpected continuation flags in", pattern)
	} else {
		if len(tgtFlags.Str) > 0 {
			flags = tgtFlags
		} else if len(flags.Str) == 0 {
			flags = eph.EphTiming{eph.EphTiming_During}
		}

		var always eph.EphAlways
		// check if this rule is declared inside a specific domain
		if guard, ok := op.Guard.(jsn.Marshalee); !ok {
			err = errutil.New("missing guard in", pattern, "at", op.Hook.At.String())
		} else {
			// fix? could we instead just strstr for countOf
			// also might be cool to augment or replace the serialized type
			// with our own that has an pre-calced field ( at import, via state parser )
			if SearchForCounters(guard) {
				always = eph.EphAlways{eph.EphAlways_Always}
			}

			// fix via runtime? check if this rule is declared inside a specific domain
			// if domain != k.Env().Game.Domain {
			// 	guard = &core.AllTrue{[]rt.BoolEval{
			// 		&core.HasDominion{domain.String()},
			// 		guard,
			// 	}}
			k.Write(&eph.EphRules{Name: pattern, Filter: op.Guard, When: flags, Exe: &act, Touch: always})
		}
	}
	return
}

func (op *PatternFlags) ReadFlags() (ret eph.EphTiming, err error) {
	switch str := op.Str; str {
	case PatternFlags_Before:
		// run other matching patterns, and then run this pattern. other...this.
		ret = eph.EphTiming{eph.EphTiming_After}
	case PatternFlags_After:
		// keep going after running the current pattern. this...others.
		ret = eph.EphTiming{eph.EphTiming_Before}
	case PatternFlags_Terminate:
		ret = eph.EphTiming{eph.EphTiming_During}
	default:
		if len(str) > 0 {
			err = errutil.Fmt("unknown pattern flags %q", str)
		}
	}
	return
}

func (op *PatternLocals) ImportLocals(k *Importer, patternName string) (ret []eph.EphParams, err error) {
	var locals []eph.EphParams
	for _, el := range op.LocalDecl {
		if p, e := el.VariableDecl.GetParam(); e != nil {
			err = e
			break
		} else {
			if init := el.Value; init != nil {
				p.Initially = init.Value
			}
			locals = append(locals, p)
		}
	}
	if err == nil {
		ret = locals
	}
	return
}

func (op *PatternType) ImportType(k *Importer) (ret string, err error) {
	if t, found := composer.FindChoice(op, op.Str); !found {
		err = errutil.Fmt("choice %s not found in %T", op.Str, op)
	} else {
		ret = t
	}
	return
}

func convertRes(res *PatternReturn) (ret *eph.EphParams, err error) {
	if res != nil {
		if p, e := res.Result.GetParam(); e != nil {
			err = e
		} else {
			ret = &p
		}
	}
	return
}

func reduceParams(els []VariableDecl) (ret []eph.EphParams, err error) {
	var ps []eph.EphParams
	for _, el := range els {
		if p, e := el.GetParam(); e != nil {
			err = errutil.Append(err, e)
		} else {
			ps = append(ps, p)
		}
	}
	if err == nil {
		ret = ps
	}
	return
}
