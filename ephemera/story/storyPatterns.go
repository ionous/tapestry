package story

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/ephemera"
	"git.sr.ht/~ionous/iffy/ephemera/decode"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

func (op *PatternActions) ImportPhrase(k *Importer) (err error) {
	if patternName, e := op.Name.NewName(k); e != nil {
		err = e
	} else if e := op.PatternRules.ImportRules(k, patternName, ephemera.Named{}, 0); e != nil {
		err = e
	} else if e := op.PatternReturn.ImportReturn(k, patternName); e != nil {
		err = e
	} else {
		// import each local if they exist
		if els := op.PatternLocals; els != nil {
			err = els.ImportLocals(k, patternName)
		}
	}
	return
}

func (op *PatternReturn) ImportReturn(k *Importer, patternName ephemera.Named) (err error) {
	if op != nil { // pattern returns are optional
		if val, e := op.Result.ImportVariable(k, tables.NAMED_RETURN); e != nil {
			err = errutil.Append(err, e)
		} else {
			k.NewPatternDecl(patternName, val.name, val.typeName, val.affinity)
		}
	}
	return
}

// Adds a new pattern declaration and optionally some associated pattern parameters.
func (op *PatternDecl) ImportPhrase(k *Importer) (err error) {
	if patternName, e := op.Name.NewName(k); e != nil {
		err = e
	} else if patternType, e := op.Type.ImportType(k); e != nil {
		err = e
	} else {
		k.NewPatternDecl(patternName, patternName, patternType, "")

		if res := op.PatternReturn; res != nil {
			err = res.ImportReturn(k, patternName)
		}
		//
		if els := op.Optvars; els != nil {
			for _, el := range els.VariableDecl {
				if val, e := el.ImportVariable(k, tables.NAMED_PARAMETER); e != nil {
					err = errutil.Append(err, e)
				} else {
					k.NewPatternDecl(patternName, val.name, val.typeName, val.affinity)
				}
			}
		}
	}
	return
}

func (op *PatternVariablesDecl) ImportPhrase(k *Importer) (err error) {
	if patternName, e := op.PatternName.NewName(k); e != nil {
		err = e
	} else {
		// fix: shouldnt this be called pattern parameters?
		for _, el := range op.VariableDecl {
			if val, e := el.ImportVariable(k, tables.NAMED_PARAMETER); e != nil {
				err = errutil.Append(err, e)
			} else {
				k.NewPatternDecl(patternName, val.name, val.typeName, val.affinity)
			}
		}
	}
	return
}

func (op *PatternRules) ImportRules(k *Importer, pattern, target ephemera.Named, flags rt.Flags) (err error) {
	if els := op.PatternRule; els != nil {
		for _, el := range *els {
			if e := el.ImportRule(k, pattern, target, flags); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

func (op *PatternRule) ImportRule(k *Importer, pattern, target ephemera.Named, tgtFlags rt.Flags) (err error) {
	if hook, e := op.Hook.ImportProgram(k); e != nil {
		err = e
	} else if flags, e := op.Flags.ReadFlags(); e != nil {
		err = e
	} else {
		// make sure only one set of flags is ... set
		if flags != 0 && tgtFlags != 0 {
			err = errutil.New("unexpected continuation flags in", pattern.String())
		}
		if tgtFlags != 0 {
			flags = tgtFlags
		}
		if err == nil {
			guard := op.Guard
			// if this rule is declared inside a specific domain, add a check for that.
			if dom := k.Current.Domain.String(); len(dom) > 0 {
				guard = &core.AllTrue{[]rt.BoolEval{
					&core.HasDominion{dom},
					guard,
				}}
			}
			rule := &rt.Rule{Filter: guard, Execute: hook, RawFlags: flags}
			if patternProg, e := k.NewGob("rule", rule); e != nil {
				err = e
			} else {
				// currentDomain returns "entire_game" when k.Current.Domain is the empty string.
				k.NewPatternRule(pattern, target, k.currentDomain(), patternProg)
			}
		}
	}
	return
}

func (op *PatternFlags) ReadFlags() (ret rt.Flags, err error) {
	if op != nil {
		switch str := op.Str; str {
		case "$BEFORE":
			// run other matching patterns, and then run this pattern. other...this.
			ret = rt.Postfix
		case "$AFTER":
			// keep going after running the current pattern. this...others.
			ret = rt.Prefix
		case "$TERMINATE":
			ret = rt.Infix
		default:
			err = errutil.New("unknown pattern flags", str)
		}
	}
	return
}

func (op *PatternLocals) ImportLocals(k *Importer, patternName ephemera.Named) (err error) {
	for _, el := range op.LocalDecl {
		if val, e := el.VariableDecl.ImportVariable(k, tables.NAMED_LOCAL); e != nil {
			err = e
			break
		} else {
			var prog ephemera.Prog
			if init := el.Value; init != nil {
				if p, e := k.NewGob("assignment", init); e != nil {
					err = e
					break
				} else {
					prog = p
				}
			}
			k.NewPatternInit(patternName, val.name, val.typeName, val.affinity, prog)
		}
	}
	return
}

func (op *PatternType) ImportType(k *Importer) (ret ephemera.Named, err error) {
	if t, found := decode.FindChoice(op, op.Str); !found {
		err = errutil.Fmt("choice %s not found in %T", op.Str, op)
	} else {
		ret = k.NewName(t, tables.NAMED_TYPE, op.At.String())
	}
	return
}
