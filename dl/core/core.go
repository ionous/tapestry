package core

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

type Say = SayText // backwards compat

func cmdError(op composer.Composer, err error) error {
	return cmdErrorCtx(op, "", err)
}

func cmdErrorCtx(op composer.Composer, ctx string, err error) error {
	// avoid triggering errutil panics for break statements
	if _, ok := err.(DoInterrupt); !ok {
		e := &composer.CommandError{Cmd: op, Ctx: ctx}
		err = errutil.Append(err, e)
	}
	return err
}

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{Value: b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{Value: float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{Value: n} }
func T(s string) *literal.TextValue { return &literal.TextValue{Value: s} }

func P(p string) PatternName  { return PatternName{Str: p} }
func N(v string) VariableName { return VariableName{Str: v} }
func W(v string) string       { return v }

// fix: rename to GetVar ( once GetVar{} is gone )
// generate a statement which extracts a variable's value.
// path can include strings ( for reading from records ) or integers ( for reading from lists )
func V(v string, path ...any) *GetFromVar {
	return &GetFromVar{
		Name: T(v),
		Dot:  MakeDot(path...),
	}
}

func GetName(v string, path ...any) *GetFromName {
	return &GetFromName{
		Name: T(v),
		Dot:  MakeDot(path...),
	}
}

func SetVar(name string, patheval ...any) (ret rt.Execute) {
	n := T(name)
	cnt := len(patheval)
	eval := patheval[cnt-1]
	dots := MakeDot(patheval[:cnt-1]...)
	switch eval := eval.(type) {
	case rt.BoolEval:
		ret = &SetVarFromBool{Name: n, Bool: eval, Dot: dots}
	case rt.NumberEval:
		ret = &SetVarFromNumber{Name: n, Number: eval, Dot: dots}
	case rt.TextEval:
		ret = &SetVarFromText{Name: n, Text: eval, Dot: dots}
	case rt.ListEval:
		ret = &SetVarFromList{Name: n, List: eval, Dot: dots}
	case rt.RecordEval:
		ret = &SetVarFromRecord{Name: n, Record: eval, Dot: dots}
	default:
		panic("unknown eval type")
	}
	return
}

func MakeDot(path ...any) []Dot {
	out := make([]Dot, len(path))
	for i, p := range path {
		switch el := p.(type) {
		case string:
			out[i] = &AtField{Field: T(el)}
		case int:
			out[i] = &AtIndex{Index: I(el)}
		default:
			panic("expected an int or string element")
		}
	}
	return out
}
