package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/express"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/template"
	"git.sr.ht/~ionous/tapestry/template/types"
)

// --------------------------------------------------------------
// SingleValue
// --------------------------------------------------------------

// panics if unmatched
func (op *SingleValue) Assignment() (ret rt.Assignment) {
	if n := op.QuotedText; n != nil {
		ret = n.Assignment()
	} else if n := op.MatchingNumber; n != nil {
		ret = n.Assignment()
	} else if n := op.Noun; n != nil {
		ret = nounAsTextValue(n.ActualNoun)
	} else if k := op.Kind; k != nil {
		ret = kindAsTextValue(k.ActualKind)
	} else {
		panic("unmatched assignment")
	}
	return
}

func (op *SingleValue) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	Optional(q, &next, &op.QuotedText) ||
		Optional(q, &next, &op.MatchingNumber) ||
		Optional(q, &next, &op.Noun) ||
		Optional(q, &next, &op.Kind) {
		*input, okay = next, true
	}
	return
}

func nounAsTextValue(noun ActualNoun) rt.Assignment {
	return text(noun.Name, noun.Kind)
}

func kindAsTextValue(kind ActualKind) rt.Assignment {
	return text(kind.Name, "") // tbd: should these be typed? ex. as "kinds" or something?
}

// --------------------------------------------------------------
// QuotedText
// --------------------------------------------------------------

func (op *QuotedText) String() string {
	return op.Matched
}

func (op *QuotedText) Assignment() rt.Assignment {
	return &assign.FromText{Value: op.TextEval()}
}

func (op *QuotedText) TextEval() (ret rt.TextEval) {
	str := op.Matched
	if v, e := ConvertTextTemplate(str); e == nil {
		ret = v
	} else {
		ret = &literal.TextValue{Value: str}
	}
	return
}

// match combines double quoted and backtick text:
// generating a leading "QuotedText" indicator
// and a single "word" containing the entire quoted text.
func (op *QuotedText) Match(q Query, input *InputState) (okay bool) {
	if v, ok := input.GetNext(match.Quoted); ok {
		op.Matched = v.(string)
		*input, okay = input.Skip(1), true
	}
	return
}

// --------------------------------------------------------------
// MatchingNumber
// --------------------------------------------------------------

func (op *MatchingNumber) Assignment() rt.Assignment {
	return number(op.Number, "")
}

// matches a natural number in words, or a literal natural number.
func (op *MatchingNumber) Match(q Query, input *InputState) (okay bool) {
	if ws := input.Words(); len(ws) > 0 && ws[0].Token == match.String {
		word := ws[0].Value.(string)
		if v, ok := WordsToNum(word); ok && v > 0 {
			const width = 1
			op.Number = float64(v)
			*input, okay = input.Skip(width), true
		}
	}
	return
}

// --------------------------------------------------------------
// support
// --------------------------------------------------------------

// tbd: i'm not sold on the idea that weave takes assignments
// maybe it'd make more sense to pass in generic "any" values
// but... note: text templates.
// or to have individual methods for the necessary types
func number(value float64, kind string) rt.Assignment {
	return &assign.FromNumber{
		Value: &literal.NumValue{Value: value, Kind: kind},
	}
}

func text(value, kind string) rt.Assignment {
	return &assign.FromText{
		Value: &literal.TextValue{Value: value, Kind: kind},
	}
}

// returns a string or a FromText assignment as a slice of bytes
func ConvertTextTemplate(str string) (ret rt.TextEval, err error) {
	if xs, e := template.Parse(str); e != nil {
		err = e
	} else if v, ok := getSimpleString(xs); ok {
		ret = &literal.TextValue{Value: v}
	} else {
		if got, e := express.Convert(xs); e != nil {
			err = e
		} else if eval, ok := got.(rt.TextEval); !ok {
			// todo: could probably fix this now; passing expected aff maybe
			// ( or maybe via unpackPatternArg? )
			err = fmt.Errorf("render template has unknown expression %T", got)
		} else {
			ret = eval
		}
	}
	return
}

// return true if the expression contained only a string, or was empty
func getSimpleString(xs template.Expression) (ret string, okay bool) {
	switch len(xs) {
	case 0:
		okay = true
	case 1:
		if quote, ok := xs[0].(types.Quote); ok {
			ret, okay = quote.Value(), true
		}
	}
	return
}

// // returns a string or a FromText assignment as a slice of bytes
// func ConvertNumberTemplate(str string) (ret rt.NumberEval, err error) {
// 	if xs, e := template.Parse(str); e != nil {
// 		err = e
// 	} else if v, ok := getSimpleValue(xs); ok {
// 		ret = &literal.NumValue{Value: v}
// 	} else {
// 		if got, e := express.Convert(xs); e != nil {
// 			err = e
// 		} else if eval, ok := got.(rt.NumberEval); !ok {
// 			err = fmt.Errorf("render template has unknown expression %T", got)
// 		} else {
// 			ret = eval
// 		}
// 	}
// 	return
// }

// // return true if the expression contained only a string, or was empty
// func getSimpleValue(xs template.Expression) (ret float64, okay bool) {
// 	switch len(xs) {
// 	case 0:
// 		okay = true
// 	case 1:
// 		if quote, ok := xs[0].(types.Number); ok {
// 			ret, okay = quote.Value(), true
// 		}
// 	}
// 	return
// }
