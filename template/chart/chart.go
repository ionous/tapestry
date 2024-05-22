package chart

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/template/postfix"
	"git.sr.ht/~ionous/tapestry/template/types"
	"github.com/ionous/errutil"
	"github.com/ionous/tell/charm"
	"github.com/ionous/tell/charmed"
)

type State = charm.State

var OnExit = charm.OnExit
var Parallel = charm.Parallel
var Parse = charm.ParseEof
var RunStep = charm.RunStep
var Self = charm.Self
var Statement = charm.Statement
var Step = charm.Step

var spaces = charm.Optional(isSpace)

// for the very next rune, return nil ( unhandled )
// it may be the end of parsing, or some parent state might be taking over from here on out.
func Finished(reason string) State {
	return Statement(reason, func(rune) (none State) {
		return
	})
}

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

func (p *NumParser) NewRune(r rune) charm.State {
	return p.Decode().NewRune(r)
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
	err error
	res string
}

// NewRune starts with the leading quote mark;
// it finishes just after the matching quote mark.
func (p *QuoteParser) NewRune(r rune) (ret State) {
	var str strings.Builder
	if decoder, ok := charmed.DecodeQuote(r, &str); ok {
		p.err = errutil.New("unclosed quote") // provisionally
		ret = charm.Step(decoder, charm.OnExit("post quote", func() {
			p.res, p.err = str.String(), nil
		}))
	}
	return
}

func (p *QuoteParser) GetOperand() (ret postfix.Function, err error) {
	if p.err != nil {
		err = p.err
	} else {
		ret = types.Quote(p.res)
	}
	return
}
