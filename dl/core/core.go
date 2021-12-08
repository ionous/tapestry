package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/rt"
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

func B(b bool) rt.BoolEval          { return &literal.BoolValue{b} }
func I(n int) rt.NumberEval         { return &literal.NumValue{float64(n)} }
func F(n float64) rt.NumberEval     { return &literal.NumValue{n} }
func P(p string) value.PatternName  { return value.PatternName{Str: p} }
func N(v string) value.VariableName { return value.VariableName{Str: v} }
func T(s string) *literal.TextValue { return &literal.TextValue{W(s)} }
func V(i string) *GetVar            { return &GetVar{N(i)} }
func W(v string) string             { return v }
