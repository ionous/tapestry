package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/express"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/template"
	"git.sr.ht/~ionous/tapestry/template/types"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// --------------------------------------------------------------
// SingleValue
// --------------------------------------------------------------

// match and apply a value
func TrySingleValue(q JessContext, in InputState,
	an ActualNoun,
	accept func(SingleValue, InputState),
	reject func(error),
) {
	q.Try(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) {
		var sv SingleValue
		if !sv.Match(q, &in) {
			reject(FailedMatch{"didn't understand the value", in})
		} else {
			accept(sv, in)
		}
	}, reject)
}

// panics if unmatched
func (op *SingleValue) Assignment() (ret rt.Assignment) {
	if n := op.QuotedText; n != nil {
		ret = n.Assignment()
	} else if n := op.MatchingNum; n != nil {
		ret = n.Assignment()
	} else if n := op.Noun; n != nil {
		ret = nounAsTextValue(n.actualNoun)
	} else if k := op.Kind; k != nil {
		ret = kindAsTextValue(k.actualKind)
	} else {
		panic("unmatched assignment")
	}
	return
}

func (op *SingleValue) Match(q Query, input *InputState) (okay bool) {
	// alt: could rework each of these as implementing a "single_value" interface
	// ( well: probably single_value and property_value )
	if next := *input; //
	Optional(q, &next, &op.QuotedText) ||
		Optional(q, &next, &op.MatchingNum) ||
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
	return &call.FromText{Value: op.TextEval()}
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
		op.Matched = v.String()
		*input, okay = input.Skip(1), true
	}
	return
}

// --------------------------------------------------------------
// MatchingNum
// --------------------------------------------------------------

func (op *MatchingNum) Assignment() rt.Assignment {
	return number(op.Value, "")
}

// matches a natural number in words, or a literal natural number.
func (op *MatchingNum) Match(q Query, input *InputState) (okay bool) {
	if ws := input.Words(); len(ws) > 0 && ws[0].Token == match.String {
		word := ws[0].String()
		if v, ok := WordsToNum(word); ok && v > 0 {
			const width = 1
			op.Value = float64(v)
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
	return &call.FromNum{
		Value: &literal.NumValue{Value: value},
	}
}

func text(value, kind string) rt.Assignment {
	return &call.FromText{
		Value: &literal.TextValue{Value: value, KindName: kind},
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
// func ConvertNumberTemplate(str string) (ret rt.NumEval, err error) {
// 	if xs, e := template.Parse(str); e != nil {
// 		err = e
// 	} else if v, ok := getSimpleValue(xs); ok {
// 		ret = &literal.NumValue{Value: v}
// 	} else {
// 		if got, e := express.Convert(xs); e != nil {
// 			err = e
// 		} else if eval, ok := got.(rt.NumEval); !ok {
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
