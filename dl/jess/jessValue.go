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

// panics if unmatched
func (op *SingleValue) Assignment() (ret rt.Assignment) {
	if n := op.QuotedText; n != nil {
		ret = n.Assignment()
	} else if n := op.MatchingNumber; n != nil {
		ret = n.Assignment()
	} else {
		panic("unmatched assignment")
	}
	return
}

func (op *SingleValue) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	Optional(q, &next, &op.QuotedText) ||
		Optional(q, &next, &op.MatchingNumber) {
		*input, okay = next, true
	}
	return
}

func (op *QuotedText) Assignment() (ret rt.Assignment) {
	str := op.Matched.String()
	if v, e := ConvertTextTemplate(str); e == nil {
		ret = &assign.FromText{Value: v}
	} else {
		ret = text(str, "")
	}
	return
}

// match combines double quoted and backtick text:
// generating a leading "QuotedText" indicator
// and a single "word" containing the entire quoted text.
func (op *QuotedText) Match(q Query, input *InputState) (okay bool) {
	if width := input.MatchWord(match.Keywords.QuotedText); width > 0 {
		next := input.Skip(width) // skip over the quote indicator (1 word)
		op.Matched, *input, okay = next.Cut(1), next.Skip(1), true
	}
	return
}

func (op *MatchingNumber) Assignment() rt.Assignment {
	return number(op.Number, "")
}

func (op *MatchingNumber) Match(q Query, input *InputState) (okay bool) {
	if ws := input.Words(); len(ws) > 0 {
		word := ws[0].String()
		if v, ok := WordsToNum(word); ok && v > 0 {
			const width = 1
			op.Number = float64(v)
			*input, okay = input.Skip(width), true
		}
	}
	return
}

// tbd: i'm not sold on the idea that registar takes assignments
// maybe it'd make more sense to pass in generic "any" values,
// to have add factory functions to Registrar,
// or to have individual methods for the necessary types
// ( maybe just three: trait, text, number )
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
