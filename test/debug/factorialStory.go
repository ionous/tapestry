package debug

import (
	_ "embed"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/rt"
)

func UserComment(s string) map[string]any {
	return map[string]any{compact.Comment: s}
}

//go:embed factorialStory.if
var FactorialJs string

// a program that can check factorials
var FactorialStory = story.StoryFile{
	Statements: []story.StoryStatement{
		&story.DefineTest{
			TestName: T("factorial"),
			Exe:      FactorialCheck,
		},
		&story.DefinePattern{
			PatternName: T("factorial"),
			Requires: []story.FieldDefinition{
				&story.NumField{
					Markup:    UserComment("just one argument, a number called 'num'"),
					FieldName: T("num"),
				}},
			Provides: []story.FieldDefinition{&story.NumField{
				Markup:    UserComment("the result uses the same variable as the pattern input does"),
				FieldName: T("num"),
			}},
			Exe: FactorialDefaultRule,
		},
		&story.DefineRule{
			RuleTiming: T("factorial"),
			Exe:        FactorialDecreaseRule,
		}},
}

// run 3! factorial
var FactorialCheck = []rt.Execute{
	&debug.Expect{
		Value: &math.CompareNum{
			A: F(6), Compare: math.C_Comparison_EqualTo, B: &call.CallPattern{
				PatternName: "factorial",
				Arguments: []call.Arg{{
					Name:  "num",
					Value: &call.FromNum{Value: F(3)},
					// fix: for some reason, the comment isn't appearing in the output.
					// Markup: UserComment("start the factorial with '3'"),
				}}},
		},
	},
}

// subtracts 1 from the num and multiples by one
var FactorialMulMinusOne = []rt.Execute{
	&object.SetValue{
		Target: object.Variable("num"),
		Value: &call.FromNum{Value: &math.MultiplyValue{
			A: object.Variable("num"),
			B: &math.SubtractValue{
				A: object.Variable("num"),
				B: I(1),
			},
		}}},
}

// override the default behavior:
var FactorialDecreaseRule = []rt.Execute{
	&logic.ChooseBranch{
		Condition: &math.CompareNum{
			Markup:  UserComment("above zero, subtract one"),
			A:       object.Variable("num"),
			Compare: math.C_Comparison_GreaterThan,
			B:       F(0)},
		Exe: FactorialMulMinusOne,
	},
}

// the default rule: use the number 1
var FactorialDefaultRule = []rt.Execute{
	&object.SetValue{
		Markup: UserComment("by default, return one"),
		Target: object.Variable("num"),
		Value:  &call.FromNum{Value: I(1)},
	},
}
