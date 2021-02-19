package test

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/list"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

var runGroupTogther = list.Map{
	FromList:     &core.Var{Name: "objects"},
	ToList:       "settings",
	UsingPattern: "assign_grouping"}

type AssignGrouping struct {
	In  string
	Out GroupSettings
}

// from a list of object names, build a list of group settings
var assignGrouping = pattern.Pattern{
	Name:   "assign_grouping",
	Return: "out",
	Labels: []string{"in"},
	Fields: []g.Field{
		{Name: "in", Affinity: affine.Text},
		{Name: "out", Affinity: affine.Record, Type: "group_settings"},
	},
	Rules: []*pattern.Rule{
		{Execute: &core.Activity{[]rt.Execute{
			Put("out", "name", V("in")),
			&core.ChooseAction{
				If: &core.Matches{
					Text:    &core.Var{Name: "in"},
					Pattern: "^thing"},
				Do: core.MakeActivity(
					Put("out", "label", &core.FromText{&core.Text{"thingies"}}),
				),
			},
		}}},
	},
}
