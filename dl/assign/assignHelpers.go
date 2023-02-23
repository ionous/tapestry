package assign

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"github.com/ionous/errutil"
)

func B(b bool) *literal.BoolValue       { return &literal.BoolValue{Value: b} }
func I(n int) *literal.NumValue         { return &literal.NumValue{Value: float64(n)} }
func F(n float64) *literal.NumValue     { return &literal.NumValue{Value: n} }
func T(s string) *literal.TextValue     { return &literal.TextValue{Value: s} }
func Ts(s []string) *literal.TextValues { return &literal.TextValues{Values: s} }

func P(patternName string) string  { return patternName }
func N(variableName string) string { return variableName }
func W(plainText string) string    { return plainText }

func Object(name, field string, path ...any) *ObjectRef {
	return &ObjectRef{
		Name:  T(name),
		Field: T(field),
		Dot:   MakeDot(path...),
	}
}

// generate a statement which extracts a variable's value.
// path can include strings ( for reading from records ) or integers ( for reading from lists )
func Variable(name string, path ...any) *VariableRef {
	return &VariableRef{
		Name: T(name),
		Dot:  MakeDot(path...),
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
