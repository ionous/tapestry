package test

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/list"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/dl/term"
	"git.sr.ht/~ionous/iffy/rt"
)

var runGroupTogther = list.Map{
	FromList:     &core.Var{Name: "Objects"},
	ToList:       "Settings",
	UsingPattern: "assignGrouping"}

// from a list of object names, build a list of group settings
var assignGrouping = pattern.Pattern{
	Name: "assignGrouping",
	Params: []term.Preparer{
		&term.Text{Name: "in"},
		&term.Record{Name: "out", Kind: "GroupSettings"},
	},
	Rules: []*pattern.Rule{
		{Execute: &core.Activity{[]rt.Execute{
			Put("out", "Name", V("in")),
			&core.ChooseAction{
				If: &core.Matches{
					Text:    &core.Var{Name: "in"},
					Pattern: "^thing"},
				Do: core.MakeActivity(
					Put("out", "Label", &core.FromText{&core.Text{"thingies"}}),
				),
			},
		}}},
	},
}
