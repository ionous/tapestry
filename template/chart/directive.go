package chart

import (
	"bytes"

	"git.sr.ht/~ionous/tapestry/template/postfix"
	"git.sr.ht/~ionous/tapestry/template/types"
	"github.com/ionous/errutil"
)

// Directive containing the parsed content of a types.
// Both or either of the key and the expression can be empty.
// ex. {elseIf b = 5}
type Directive struct {
	Key string
	postfix.Expression
}

func UnknownDirective(v Directive) error {
	return errutil.Fmt("unknown directive '%s'", v.Key)
}

func UnexpectedExpression(v Directive) (err error) {
	if len(v.Expression) > 0 {
		err = errutil.Fmt("unexpected expression following key '%s'", v.Key)
	}
	return
}
func ExpectedExpression(v Directive) (err error) {
	if len(v.Expression) == 0 {
		err = errutil.Fmt("expected expression following key '%s'", v.Key)
	}
	return
}

// String of a directive in the format:
// {key:expression} or {expression}
func (d Directive) String() (ret string) {
	if len(d.Key) > 0 {
		ret = "{" + d.Key + ":" + d.Expression.String() + "}"
	} else if q, ok := d.isQuote(); ok {
		ret = string(q)
	} else {
		ret = "{" + d.Expression.String() + "}"
	}
	return
}

func (d Directive) isQuote() (ret types.Quote, okay bool) {
	if cnt := len(d.Expression); cnt == 1 {
		ret, okay = d.Expression[0].(types.Quote)
	}
	return
}

// Format a string from slice of directives.
func Format(ds []Directive) string {
	var buf bytes.Buffer
	for _, d := range ds {
		buf.WriteString(d.String())
	}
	return buf.String()
}
