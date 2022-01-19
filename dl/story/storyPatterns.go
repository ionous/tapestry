package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/jsn"
	"github.com/ionous/errutil"
)

func (op *PatternActions) ImportPhrase(k *Importer) (err error) {
	patternName := op.Name.String()
	var locals []eph.EphParams
	if els := op.Provides; els != nil {
		locals = els.ImportLocals(k, patternName)
	}
	if len(locals) > 0 {
		k.WriteEphemera(&eph.EphPatterns{Name: patternName, Locals: locals})
	}
	// write the rules last ( order doesnt matter except for tests )
	return op.PatternRules.ImportRules(k, patternName, "", eph.EphTiming{})
}

// Adds a new pattern declaration and optionally some associated pattern parameters.
func (op *PatternDecl) ImportPhrase(k *Importer) (err error) {
	patternName := op.Name.String()
	ps := op.reduceProps()
	res := convertRes(op.PatternReturn)
	k.WriteEphemera(&eph.EphPatterns{Name: patternName, Result: res, Params: ps})
	return
}

func (op *PatternDecl) reduceProps() (ret []eph.EphParams) {
	if els := op.Params; els != nil {
		ret = reduceProps(els.Props)
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
			err = errutil.New("missing guard in", pattern)
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
			k.WriteEphemera(&eph.EphRules{
				Name:   pattern,
				Target: target, // fix: this should become part of the guards i think, even if its less slightly less efficient
				Filter: op.Guard,
				When:   flags,
				Exe:    act,
				Touch:  always})
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

func (op *PatternLocals) ImportLocals(k *Importer, patternName string) []eph.EphParams {
	var out []eph.EphParams
	for _, el := range op.Locals {
		out = append(out, el.GetParam())
	}
	return out
}

func convertRes(res *PatternReturn) (ret *eph.EphParams) {
	if res != nil {
		p := res.Result.GetParam()
		ret = &p
	}
	return
}

func reduceProps(els []PropertySlot) []eph.EphParams {
	var out []eph.EphParams
	for _, el := range els {
		out = append(out, el.GetParam())
	}
	return out
}
