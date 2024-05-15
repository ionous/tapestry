package test

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
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
		core.MakeRule(
			nil, matches(true)),
		core.MakeRule(
			&core.CompareText{
				A:  assign.Variable("a", "label"),
				Is: core.C_Comparison_OtherThan,
				B:  assign.Variable("b", "label"),
			}, matches(false)),
		core.MakeRule(
			&core.CompareText{
				A:  assign.Variable("a", "group options"),
				Is: core.C_Comparison_OtherThan,
				B:  assign.Variable("b", "group options"),
			}, matches(false)),
	},
}

func matches(b bool) rt.Execute {
	return &assign.SetValue{
		Target: assign.Variable("matches"),
		Value:  &assign.FromBool{Value: B(b)},
	}
}
