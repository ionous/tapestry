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

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{Bool: b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{Num: float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{Num: n} }
func T(s string) *literal.TextValue { return &literal.TextValue{Text: s} }

func P(p string) PatternName  { return PatternName{Str: p} }
func N(v string) VariableName { return VariableName{Str: v} }
func V(i string) *GetVar      { return &GetVar{N(i)} }
func W(v string) string       { return v }
