package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *PatternActions) PostImport(k *imp.Importer) (err error) {
	patternName := op.PatternName
	if locals := ImportLocals(k, patternName, op.Locals); len(locals) > 0 {
		k.WriteEphemera(&eph.EphPatterns{PatternName: patternName, Locals: locals})
	}
	// write the rules last ( to help with test output consistency )
	return ImportRules(k, patternName, "", op.Rules, eph.EphTiming{})
}

func (op *ExtendPattern) PostImport(k *imp.Importer) (err error) {
	if name, e := safe.GetText(nil, op.PatternName); e != nil {
		err = e
	} else if e := ImportRules(k, name.String(), "", op.Rules, eph.EphTiming{}); e != nil {
		err = e
	} else if locals := ImportLocals(k, name.String(), op.Locals); len(locals) > 0 {
		k.WriteEphemera(&eph.EphPatterns{PatternName: name.String(), Locals: locals})
		// write the rules last ( order doesnt matter except it helps with test output consistency )
	}
	return
}

// Adds a new pattern declaration and optionally some associated pattern parameters.
func (op *DefinePattern) PostImport(k *imp.Importer) (err error) {
	ps := op.reduceProps()
	if name, e := safe.GetText(nil, op.PatternName); e != nil {
		err = e
	} else {
		var pres *eph.EphParams
		if opres := op.Result; opres != nil {
			if res, okay := opres.GetParam(); okay {
				pres = &res
			}
		}
		k.WriteEphemera(&eph.EphPatterns{PatternName: name.String(), Result: pres, Params: ps})
		// write the rules last ( order doesnt matter except it helps with test output consistency )
		err = ImportRules(k, name.String(), "", op.Rules, eph.EphTiming{})
	}
	return
}

func (op *DefinePattern) reduceProps() []eph.EphParams {
	return reduceProps(op.Params)
}

// note:  statements can set flags for a bunch of rules at once or within each rule separately, but not both.
func ImportRules(k *imp.Importer, pattern, target string, els []PatternRule, flags eph.EphTiming) (err error) {
	// write in reverse order because within a given pattern, earlier rules take precedence.
	for i := len(els) - 1; i >= 0; i-- {
		if e := els[i].importRule(k, pattern, target, flags); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (op *PatternRule) importRule(k *imp.Importer, pattern, target string, tgtFlags eph.EphTiming) (err error) {
	act := op.Does
	if flags, e := op.Flags.ReadFlags(); e != nil {
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
				PatternName: pattern,
				Target:      target, // fix: this should become part of the guards i think, even if its less slightly less efficient
				Filter:      op.Guard,
				When:        flags,
				Exe:         act,
				Touch:       always})
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

func ImportLocals(k *imp.Importer, patternName string, locals []FieldDefinition) (ret []eph.EphParams) {
	for _, el := range locals {
		if p, ok := el.GetParam(); ok {
			ret = append(ret, p)
		}
	}
	return
}

func reduceProps(els []FieldDefinition) []eph.EphParams {
	var out []eph.EphParams
	for _, el := range els {
		if p, ok := el.GetParam(); ok {
			out = append(out, p)
		}
	}
	return out
}
