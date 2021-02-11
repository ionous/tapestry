package test

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// a pattern for matching groups --
// we add rules that if things arent equal we return false
var matchGroups = pattern.Pattern{
	Name:   "matchGroups",
	Labels: []string{"a", "b"},
	Fields: []g.Field{
		{Name: "a", Affinity: affine.Record, Type: "GroupSettings"},
		{Name: "b", Affinity: affine.Record, Type: "GroupSettings"},
		{Name: "matches", Affinity: affine.Bool},
	},
	Return: "matches",
	// rules are evaluated in reverse order ( see sortRules )
	Rules: []*pattern.Rule{{
		Filter:  &core.Always{},
		Execute: matches(true),
	}, {
		Filter: &core.CompareText{
			&core.GetAtField{
				From:  &core.FromVar{N("a")},
				Field: "Label",
			},
			&core.NotEqualTo{},
			&core.GetAtField{
				From:  &core.FromVar{N("b")},
				Field: "Label",
			},
		},
		Execute: matches(false),
	}, {
		Filter: &core.CompareText{
			&core.GetAtField{
				From:  &core.FromVar{N("a")},
				Field: "Innumerable",
			},
			&core.NotEqualTo{},
			&core.GetAtField{
				From:  &core.FromVar{N("b")},
				Field: "Innumerable",
			},
		},
		Execute: matches(false),
	}, {
		Filter: &core.CompareText{
			&core.GetAtField{
				From:  &core.FromVar{N("a")},
				Field: "GroupOptions",
			},
			&core.NotEqualTo{},
			&core.GetAtField{
				From:  &core.FromVar{N("b")},
				Field: "GroupOptions",
			},
		},
		Execute: matches(false),
	}},
}

func matches(b bool) rt.Execute {
	return &core.Assign{core.Variable{Str: "matches"}, &core.FromBool{&core.Bool{b}}}
}
