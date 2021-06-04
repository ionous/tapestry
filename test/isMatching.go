package test

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/test/testpat"
)

type MatchGroups struct {
	A       GroupSettings
	B       GroupSettings
	Matches bool `if:"bool"`
}

// a pattern for matching groups --
// we add rules that if things arent equal we return false
var matchGroups = testpat.Pattern{
	Name:   "match_groups",
	Labels: []string{"a", "b"},
	Fields: []g.Field{
		{Name: "a", Affinity: affine.Record, Type: "group_settings"},
		{Name: "b", Affinity: affine.Record, Type: "group_settings"},
		{Name: "matches", Affinity: affine.Bool},
	},
	Return: "matches",
	// rules are evaluated in reverse order ( see sortRules )
	Rules: []rt.Rule{{
		Filter:  &core.Always{},
		Execute: matches(true),
	}, {
		Filter: &core.CompareText{
			&core.GetAtField{
				From:  &core.FromVar{N("a")},
				Field: W("label"),
			},
			&core.Unequal{},
			&core.GetAtField{
				From:  &core.FromVar{N("b")},
				Field: W("label"),
			},
		},
		Execute: matches(false),
	}, {
		Filter: &core.CompareText{
			&core.GetAtField{
				From:  &core.FromVar{N("a")},
				Field: W("innumerable"),
			},
			&core.Unequal{},
			&core.GetAtField{
				From:  &core.FromVar{N("b")},
				Field: W("innumerable"),
			},
		},
		Execute: matches(false),
	}, {
		Filter: &core.CompareText{
			&core.GetAtField{
				From:  &core.FromVar{N("a")},
				Field: W("group_options"),
			},
			&core.Unequal{},
			&core.GetAtField{
				From:  &core.FromVar{N("b")},
				Field: W("group_options"),
			},
		},
		Execute: matches(false),
	}},
}

func matches(b bool) rt.Execute {
	return &core.Assign{N("matches"), &core.FromBool{B(b)}}
}
