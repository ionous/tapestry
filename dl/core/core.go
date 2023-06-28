package core

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"github.com/ionous/errutil"
)

func P(patternName string) string  { return patternName }
func N(variableName string) string { return variableName }
func W(plainText string) string    { return plainText }

var cmdError = assign.CmdError       // backwards compat
var cmdErrorCtx = assign.CmdErrorCtx // backwards compat

var (
	B           = literal.B
	F           = literal.F
	I           = literal.I
	T           = literal.T
	Ts          = literal.Ts
	CmdError    = assign.CmdError    // backwards compat
	CmdErrorCtx = assign.CmdErrorCtx // backwards compat
)

func Object(name, field string, path ...any) *assign.ObjectRef {
	return &assign.ObjectRef{
		Name:  literal.T(name),
		Field: literal.T(field),
		Dot:   MakeDot(path...),
	}
}

// generate a statement which extracts a variable's value.
// path can include strings ( for reading from records ) or integers ( for reading from lists )
func Variable(name string, path ...any) *assign.VariableRef {
	return &assign.VariableRef{
		Name: literal.T(name),
		Dot:  MakeDot(path...),
	}
}

func MakeDot(path ...any) []assign.Dot {
	out := make([]assign.Dot, len(path))
	for i, p := range path {
		switch el := p.(type) {
		case string:
			out[i] = &assign.AtField{Field: literal.T(el)}
		case int:
			out[i] = &assign.AtIndex{Index: literal.I(el)}
		case assign.Dot:
			out[i] = el
		default:
			panic(errutil.Fmt("expected an int or string element; got %T", el))
		}
	}
	return out
}
