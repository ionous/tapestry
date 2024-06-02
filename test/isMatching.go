package test

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/test/testpat"
)

type MatchGroups struct {
	A       GroupSettings
	B       GroupSettings
	Matches bool
}

// a pattern for matching groups --
// we add rules that if things arent equal we return false
var matchGroups = testpat.Pattern{
	Name:   "match groups",
	Labels: []string{"a", "b"},
	Fields: []rt.Field{
		{Name: "a", Affinity: affine.Record, Type: "group settings"},
		{Name: "b", Affinity: affine.Record, Type: "group settings"},
		{Name: "matches", Affinity: affine.Bool},
	},
	Return: "matches",
	// rules are evaluated in reverse order
	Rules: []rt.Rule{
		makeRule(
			nil, matches(true)),
		makeRule(
			&math.CompareText{
				A:       object.Variable("a", "label"),
				Compare: math.C_Comparison_OtherThan,
				B:       object.Variable("b", "label"),
			}, matches(false)),
		makeRule(
			&math.CompareText{
				A:       object.Variable("a", "group options"),
				Compare: math.C_Comparison_OtherThan,
				B:       object.Variable("b", "group options"),
			}, matches(false)),
	},
}

func matches(b bool) rt.Execute {
	return &object.SetValue{
		Target: object.Variable("matches"),
		Value:  &call.FromBool{Value: B(b)},
	}
}

// for tests
func makeRule(filter rt.BoolEval, exe ...rt.Execute) (ret rt.Rule) {
	if filter == nil {
		ret.Exe = exe
	} else {
		ret = rt.Rule{Exe: []rt.Execute{
			&logic.ChooseBranch{
				Condition: filter,
				Exe:       exe,
			}}}
	}
	return
}
