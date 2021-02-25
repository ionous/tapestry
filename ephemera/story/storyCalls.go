package story

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/ephemera"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

// a stub so that we can record the pattern and its arguments as referenced
// note: send should be doing something similar, and it isnt.
// a simpler way of recording this -- handwaving something about template parsing -- would be nice.
func (op *Determine) ImportStub(k *Importer) (ret interface{}, err error) {
	if p, args, e := importCall(k, "patterns", op.Name, op.Arguments); e != nil {
		err = ImportError(op, op.At, e)
	} else {
		ret = &core.Determine{Pattern: pattern.PatternName(p.String()), Arguments: args}
	}
	return
}

func importCall(k *Importer, slot string, n PatternName, stubs *Arguments) (retName ephemera.Named, retArgs *core.Arguments, err error) {
	if p, e := n.NewName(k); e != nil {
		err = e
	} else if args, e := importArgs(k, p, stubs); e != nil {
		err = e
	} else {
		// fix: tests expect pattern type to be declared last :'(
		// fix: object type names will need adaption of some sort re plural_kinds
		patternType := k.NewName(slot, tables.NAMED_TYPE, n.At.String())
		k.NewPatternRef(p, p, patternType, "")
		retName, retArgs = p, args
	}
	return
}

func importArgs(k *Importer, p ephemera.Named, stubs *Arguments) (ret *core.Arguments, err error) {
	if stubs != nil {
		var argList []*core.Argument
		for _, stub := range stubs.Args {
			aff := stub.From.Affinity()
			if paramName, e := stub.Name.NewName(k, tables.NAMED_ARGUMENT); e != nil {
				err = errutil.Append(err, e)
			} else {
				if aff := string(aff); len(aff) > 0 {
					// fix: this shouldnt be "eval" here.
					// see buildPatternCache
					paramType := k.NewName(aff+"_eval", tables.NAMED_TYPE, stub.At.String())
					k.NewPatternRef(p, paramName, paramType, "")
				}
				// after recording the "fact" of the parameter...
				// copy the stubbed argument data into the real argument list.
				newArg := &core.Argument{Name: paramName.String(), From: stub.From}
				argList = append(argList, newArg)
			}
		}
		if err == nil {
			ret = &core.Arguments{Args: argList}
		}
	}
	return
}
