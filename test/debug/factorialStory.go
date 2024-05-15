package debug

import (
	_ "embed"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
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
			TestName: "factorial",
			Exe:      FactorialCheck,
		},
		&story.DefinePattern{
			PatternName: T("factorial"),
			Requires: []story.FieldDefinition{
				&story.NumberField{
					Markup:    UserComment("just one argument, a number called 'num'"),
					FieldName: T("num"),
				}},
			Provides: []story.FieldDefinition{&story.NumberField{
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
		Value: &core.CompareNum{
			A: F(6), Is: core.C_Comparison_EqualTo, B: &assign.CallPattern{
				PatternName: "factorial",
				Arguments: []assign.Arg{{
					Name:  "num",
					Value: &assign.FromNumber{Value: F(3)},
					// fix: for some reason, the comment isn't appearing in the output.
					// Markup: UserComment("start the factorial with '3'"),
				}}},
		},
	},
}

// subtracts 1 from the num and multiples by one
var FactorialMulMinusOne = []rt.Execute{
	&assign.SetValue{
		Target: assign.Variable("num"),
		Value: &assign.FromNumber{Value: &core.MultiplyValue{
			A: assign.Variable("num"),
			B: &core.SubtractValue{
				A: assign.Variable("num"),
				B: I(1),
			},
		}}},
}

// override the default behavior:
var FactorialDecreaseRule = []rt.Execute{
	&core.ChooseBranch{
		If: &core.CompareNum{
			Markup: UserComment("above zero, subtract one"),
			A:      assign.Variable("num"),
			Is:     core.C_Comparison_GreaterThan,
			B:      F(0)},
		Exe: FactorialMulMinusOne,
	},
}

// the default rule: use the number 1
var FactorialDefaultRule = []rt.Execute{
	&assign.SetValue{
		Markup: UserComment("by default, return one"),
		Target: assign.Variable("num"),
		Value:  &assign.FromNumber{Value: I(1)},
	},
}
