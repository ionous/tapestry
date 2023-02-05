package core

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/literal"
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
func V(i string) *GetVar      { return &GetVar{Name: N(i)} }
func W(v string) string       { return v }

// MakeGetFromVar - generate a statement which extracts a variable's value.
// path can include strings ( for reading from records ) or integers ( for reading from lists )
func MakeGetFromVar(v string, path ...any) *GetFromVar {
	var dot []Dot
	for _, p := range path {
		switch el := p.(type) {
		case string:
			dot = append(dot, &AtField{Field: T(el)})
		case int:
			dot = append(dot, &AtIndex{Index: I(el)})
		default:
			panic("invalid path")
		}
	}
	return &GetFromVar{
		Name: T(v),
		Dot:  dot,
	}
}
