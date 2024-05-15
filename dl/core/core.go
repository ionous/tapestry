package core

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
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
