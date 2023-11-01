package chart

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/charmed"
	"git.sr.ht/~ionous/tapestry/template/postfix"
	"git.sr.ht/~ionous/tapestry/template/types"
)

type Runes = charm.Runes
type State = charm.State

var Step = charm.Step
var Parallel = charm.Parallel
var Parse = charm.Parse
var RunStep = charm.RunStep
var Self = charm.Self
var OnExit = charm.OnExit
var Statement = charm.Statement

var Terminal = charm.Terminal
var spaces = charm.Optional(isSpace)

// fix: apparently this is handling more than just bool
// if it doesnt read a whole ident, or if it tries to return error ( instead of a nil function )
// then various tests fail ( ex. {cycle}a{end} )
type BooleanParser struct{ IdentParser }

func (p *BooleanParser) GetOperand() (ret postfix.Function, err error) {
	switch id := p.IdentParser.Identifier(); id {
	case "true":
		ret = types.Bool(true)
	case "false":
		ret = types.Bool(false)
	default:
		// err = errutil.New("mismatched bool")
	}
	return
}

type NumParser struct {
	charmed.NumParser
}

func (p *NumParser) GetOperand() (ret postfix.Function, err error) {
	if n, e := p.GetFloat(); e != nil {
		err = e
	} else {
		ret = types.Number(n)
	}
	return
}

type QuoteParser struct {
	charmed.QuoteParser
}

// NewRune starts with the leading quote mark; it finishes just after the matching quote mark.
func (p *QuoteParser) NewRune(r rune) (ret State) {
	if isQuote(r) {
		ret = p.ScanQuote(r)
	}
	return
}

func (p *QuoteParser) GetOperand() (ret postfix.Function, err error) {
	if r, e := p.GetString(); e != nil {
		err = e
	} else {
		ret = types.Quote(r)
	}
	return
}
