package pattern_test

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/scope"
	"git.sr.ht/~ionous/iffy/rt/writer"
	"git.sr.ht/~ionous/iffy/test/testutil"
)

func B(i bool) rt.BoolEval     { return &core.Bool{i} }
func I(i int) rt.NumberEval    { return &core.Number{float64(i)} }
func T(i string) rt.TextEval   { return &core.Text{i} }
func V(i string) *core.Var     { return &core.Var{Name: i} }
func N(n string) core.Variable { return core.Variable{Str: n} }

type baseRuntime struct {
	testutil.PanicRuntime
}

type patternRuntime struct {
	baseRuntime
	scope.ScopeStack    // parameters are pushed onto the stack.
	testutil.PatternMap // holds pointers to patterns
}

func (patternRuntime) Writer() writer.Output {
	return writer.NewStdout()
}
