package story

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/ephemera/eph"
	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/iffy/tables"
)

// a stub so that we can record the pattern and its arguments as referenced
// note: "send" should be doing something similar, and it isnt.
// a simpler way of recording this -- handwaving something about template parsing -- would be nice.
func (op *Determine) ImportStub(k *Importer) (ret interface{}, err error) {
	if p, args, e := importCall(k, "patterns", op.Name, op.Arguments); e != nil {
		err = &OpError{Op: op, Err: e}
	} else {
		ret = &core.CallPattern{Pattern: value.PatternName{Str: p.String()}, Arguments: args}
	}
	return
}

func (op *Determine) ImportPhrase(k *Importer) (err error) {
	return errutil.New("determine should be transformed, not imported")
}

func (op *Make) ImportStub(k *Importer) (ret interface{}, err error) {
	// fix: add a reference to the kind.
	// fix: not recording this against a "pattern" name, but it could be recorded against a kind
	if args, e := importArgs(k, eph.Named{}, op.Arguments); e != nil {
		err = &OpError{Op: op, Err: e}
	} else {
		ret = &core.CallMake{Kind: op.Name, Arguments: args}
	}
	return
}

func (op *Send) ImportStub(k *Importer) (ret interface{}, err error) {
	pn := value.PatternName{Str: op.Event}
	if p, args, e := importCall(k, "actions", pn, op.Arguments); e != nil {
		err = &OpError{Op: op, Err: e}
	} else {
		// event, path ( list ), args
		ret = &core.CallSend{Event: p.String(), Path: op.Path, Arguments: args}
	}
	return
}

func importCall(k *Importer, slot string, n value.PatternName, stubs *Arguments) (retName eph.Named, retArgs core.CallArgs, err error) {
	if p, e := NewPatternName(k, n); e != nil {
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

func importArgs(k *Importer, p eph.Named, stubs *Arguments,
) (ret core.CallArgs, err error) {
	if stubs != nil {
		var argList []core.CallArg
		for _, stub := range stubs.Args {
			paramName := k.NewName(stub.Name, tables.NAMED_ARGUMENT, stub.At.String())

			if aff := stub.From.Affinity(); p.IsValid() && len(aff) > 0 {
				// fix: this shouldnt be "eval" here.
				// see buildPatternCache
				paramType := k.NewName(string(aff)+"_eval", tables.NAMED_TYPE, stub.At.String())
				k.NewPatternRef(p, paramName, paramType, "")
			}
			// after recording the "fact" of the parameter...
			// copy the stubbed argument data into the real argument list.
			newArg := core.CallArg{Name: stub.Name, From: stub.From}
			argList = append(argList, newArg)
		}
		if err == nil {
			ret = core.CallArgs{Args: argList}
		}
	}
	return
}
