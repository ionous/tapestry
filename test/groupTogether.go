package test

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/test/testpat"
)

var runGroupTogther = list.ListMap{
	FromList:     V("objects"),
	ToList:       W("settings"),
	UsingPattern: W("assign_grouping"),
}

type AssignGrouping struct {
	In  string
	Out GroupSettings
}

// from a list of object names, build a list of group settings
var assignGrouping = testpat.Pattern{
	Name:   "assign_grouping",
	Return: "out",
	Labels: []string{"in"},
	Fields: []g.Field{
		{Name: "in", Affinity: affine.Text},
		{Name: "out", Affinity: affine.Record, Type: "group_settings"},
	},
	Rules: []rt.Rule{
		{Execute: &core.Activity{[]rt.Execute{
			Put("out", "name", V("in")),
			&core.ChooseAction{
				If: &core.Matches{
					Text:    V("in"),
					Pattern: "^thing"},
				Do: core.MakeActivity(
					Put("out", "label", &core.FromText{T("thingies")}),
				),
			},
		}}},
	},
}
