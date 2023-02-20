package test

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/test/testpat"
)

var runCollateGroups = list.ListReduce{
	Target:       core.Variable("collation"),
	List:         &assign.FromRecordList{Value: GetVariable("settings")},
	UsingPattern: W("collate_groups")}

// pattern:
type CollateGroups struct {
	// args:
	Settings  GroupSettings  // constant input
	Collation GroupCollation // accumulator; FIX: this should be []GroupedObjects -- no need for the wrapping structure.

	// locals
	Found bool `if:"bool"`
	Group GroupedObjects
}

// called multiple times: once for each "setting"
var collateGroups = testpat.Pattern{
	Name:   "collate_groups",
	Labels: []string{"settings", "collation"},
	Return: "collation",
	Fields: []g.Field{
		// arguments:
		{Name: "settings", Affinity: affine.Record, Type: "group_settings"},
		{Name: "collation", Affinity: affine.Record, Type: "group_collation"},
		// locals:
		{Name: "found", Affinity: affine.Bool},
		{Name: "group", Affinity: affine.Record, Type: "grouped_objects"},
	},
	Rules: []rt.Rule{{
		Execute: core.MakeActivity(
			// walk collation.groups, set .idx to whichever best matches.
			// fix: could this be list find? ( could list find take an optional pattern )
			&list.ListEach{
				List: &assign.FromRecordList{Value: GetVariable("collation", "groups")},
				As:   W("el"),
				Does: core.MakeActivity(
					&core.ChooseAction{
						If: &core.CallPattern{
							Pattern: P("match_groups"),
							Arguments: core.MakeArgs(
								&assign.FromRecord{Value: GetVariable("settings")},
								// "el" is specified by us during ListEach
								// fix: use a special $element name, and require "Up" scope to reach out.
								// ( or allow locals and force them to assign )
								&assign.FromRecord{Value: GetVariable("el", "settings")},
							),
						},
						Does: core.MakeActivity(
							&core.SetValue{
								Target: core.Variable("found"),
								Value:  &assign.FromBool{Value: B(true)},
							},
							// found a matching group? add the object to it group.
							&list.ListPush{
								// "index" is generated by ListEach -- fix: make the name more special?
								Target: core.Variable(
									"collation", "groups", &core.AtIndex{Index: GetVariable("index")}, "objects"),
								Value: &assign.FromText{Value: GetVariable("settings", "name")}},
							// todo: implement a "break"
						),
					},
				)}, // end go-each
			&core.ChooseAction{
				// didn't find a matching group?
				If: &core.Not{Test: GetVariable("found")},
				// pack the object and its settings into the local 'group'
				// then push the group into the groups.
				// FIX: a command to MakeRecord from args, and remove the local.
				Does: core.MakeActivity(
					&core.SetValue{
						Target: core.Variable("group", "settings"),
						Value:  &assign.FromRecord{Value: GetVariable("settings")}},
					&list.ListPush{
						Target: core.Variable("group", "objects"),
						Value:  &assign.FromText{Value: GetVariable("settings", "name")}},
					&list.ListPush{
						Target: core.Variable("collation", "groups"),
						Value:  &assign.FromRecord{Value: GetVariable("group")},
					},
				), // end true
			},
		)},
	},
}
