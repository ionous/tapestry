package story

import (
	r "reflect"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/ephemera"
	"git.sr.ht/~ionous/iffy/ephemera/decode"
	"git.sr.ht/~ionous/iffy/export"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

func (op *PatternActions) ImportPhrase(k *Importer) (err error) {
	if patternName, e := NewPatternName(k, op.Name); e != nil {
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
	if patternName, e := NewPatternName(k, op.Name); e != nil {
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
	if patternName, e := NewPatternName(k, op.PatternName); e != nil {
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
		for _, el := range els {
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
	} else if flags != 0 && tgtFlags != 0 {
		// ensure flags were only set via the rule or via the pattern
		err = errutil.New("unexpected continuation flags in", pattern.String())
	} else {
		if tgtFlags != 0 {
			flags = tgtFlags
		} else if flags == 0 {
			flags = rt.Infix
		}
		// check if this rule is declared inside a specific domain
		guard := op.Guard
		if searchForCounters(r.ValueOf(guard)) {
			flags |= rt.Filter
		}
		// check if this rule is declared inside a specific domain
		if dom := k.Current.Domain.String(); len(dom) > 0 {
			guard = &core.AllTrue{[]rt.BoolEval{
				&core.HasDominion{value.Text{Str: dom}},
				guard,
			}}
		}
		// a token stream sure would be nice here -- then we could just strstr for countOf
		rule := &rt.Rule{Filter: guard, Execute: hook, RawFlags: flags}
		if patternProg, e := k.NewGob("rule", rule); e != nil {
			err = e
		} else {
			// currentDomain returns "entire_game" when k.Current.Domain is the empty string.
			k.NewPatternRule(pattern, target, k.currentDomain(), patternProg)
		}
	}
	return
}

func (op *PatternFlags) ReadFlags() (ret rt.Flags, err error) {
	switch str := op.Str; str {
	case PatternFlags_Before:
		// run other matching patterns, and then run this pattern. other...this.
		ret = rt.Postfix
	case PatternFlags_After:
		// keep going after running the current pattern. this...others.
		ret = rt.Prefix
	case PatternFlags_Terminate:
		ret = rt.Infix
	default:
		if len(str) > 0 {
			err = errutil.Fmt("unknown pattern flags %q", str)
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

func searchForCounters(rval r.Value) bool {
	return searchForType(rval.Elem(), r.TypeOf((*core.CallTrigger)(nil)).Elem())
}

func searchForType(rval r.Value, match r.Type) (ret bool) {
	if rtype := rval.Type(); rtype == match {
		ret = true
	} else {
		ret = export.WalkProperties(rtype, func(f *r.StructField, path []int) (ret bool) {
			switch ftype := f.Type; ftype.Kind() {
			case r.Slice:
				k := ftype.Elem().Kind()
				if k == r.Interface {
					f := rval.FieldByIndex(path)
					for i, cnt := 0, f.Len(); i < cnt; i++ {
						if el := f.Index(i); checkInterfaceType(el, match) {
							ret = true
							break
						}
					}
				}

			case r.Ptr:
				f := rval.FieldByIndex(path)
				ret = checkPtrType(f, match)
			case r.Interface:
				f := rval.FieldByIndex(path)
				ret = checkInterfaceType(f, match)
			}
			return
		})
	}
	return
}

func checkPtrType(f r.Value, match r.Type) (ret bool) {
	if !f.IsNil() {
		if el := f.Elem(); el.Kind() == r.Struct {
			ret = searchForType(el, match)
		}
	}
	return
}

func checkInterfaceType(f r.Value, match r.Type) (ret bool) {
	if !f.IsNil() {
		el := f.Elem()
		if k := el.Kind(); k == r.Ptr {
			ret = checkPtrType(el, match)
		}
	}
	return
}
