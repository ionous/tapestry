package core

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"github.com/ionous/errutil"
)

type Say = SayText            // backwards compat
var cmdError = CmdError       // backwards compat
var cmdErrorCtx = CmdErrorCtx // backwards compat

func CmdError(op composer.Composer, err error) error {
	return cmdErrorCtx(op, "", err)
}

func CmdErrorCtx(op composer.Composer, ctx string, err error) error {
	e := &composer.CommandError{Cmd: op, Ctx: ctx}
	return errutil.Append(e, err)
}

func B(b bool) *literal.BoolValue       { return &literal.BoolValue{Value: b} }
func I(n int) *literal.NumValue         { return &literal.NumValue{Value: float64(n)} }
func F(n float64) *literal.NumValue     { return &literal.NumValue{Value: n} }
func T(s string) *literal.TextValue     { return &literal.TextValue{Value: s} }
func Ts(s []string) *literal.TextValues { return &literal.TextValues{Values: s} }

func P(p string) PatternName  { return PatternName{Str: p} }
func N(v string) VariableName { return VariableName{Str: v} }
func W(v string) string       { return v }

// generate a statement which extracts a variable's value.
// path can include strings ( for reading from records ) or integers ( for reading from lists )
func GetVariable(name string, path ...any) *GetValue {
	return &GetValue{Source: Variable(name, path...)}
}

func Object(name, field string, path ...any) Address {
	return Address{
		Choice: Address_Object_Opt,
		Value: &ObjectRef{
			Name:  T(name),
			Field: T(field),
			Dot:   MakeDot(path...),
		},
	}
}

func Variable(name string, path ...any) Address {
	return Address{
		Choice: Address_Variable_Opt,
		Value: &VariableRef{
			Name: T(name),
			Dot:  MakeDot(path...),
		},
	}
}

func MakeDot(path ...any) []Dot {
	out := make([]Dot, len(path))
	for i, p := range path {
		switch el := p.(type) {
		case string:
			out[i] = &AtField{Field: T(el)}
		case int:
			out[i] = &AtIndex{Index: I(el)}
		case Dot:
			out[i] = el
		default:
			panic(errutil.Fmt("expected an int or string element; got %T", el))
		}
	}
	return out
}
