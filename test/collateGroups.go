package test

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/test/testpat"
)

var runCollateGroups = list.ListReduce{
	FromList:     V("settings"),
	IntoValue:    W("collation"),
	UsingPattern: W("collate_groups")}

type CollateGroups struct {
	Settings  GroupSettings
	Collation GroupCollation
	Idx       float64
	Groups    []GroupedObjects
	Group     GroupedObjects
	Names     []string
}

var collateGroups = testpat.Pattern{
	Name:   "collate_groups",
	Labels: []string{"settings", "collation"},
	Return: "collation",
	Fields: []g.Field{
		{Name: "settings", Affinity: affine.Record, Type: "group_settings"},
		{Name: "collation", Affinity: affine.Record, Type: "group_collation"},
		{Name: "idx", Affinity: affine.Number},
		{Name: "groups", Affinity: affine.RecordList, Type: "grouped_objects"},
		{Name: "group", Affinity: affine.Record, Type: "grouped_objects"},
		{Name: "names", Affinity: affine.TextList},
	},
	Rules: []rt.Rule{{
		Execute: core.MakeActivity(
			// walk collation.Groups for matching settings
			&core.Assign{
				Var:  N("groups"),
				From: V("collation", "groups"),
			},
			&list.ListEach{
				List: V("groups"),
				As:   &list.AsRec{Var: N("el")},
				Does: core.MakeActivity(
					&core.ChooseAction{
						If: &core.CallPattern{
							Pattern: P("match_groups"),
							Arguments: core.Args(
								V("settings"),
								V("el", "settings"),
							),
						},
						Does: core.MakeActivity(
							&core.Assign{
								Var:  N("idx"),
								From: V("index"),
							},
							// implement a "break" for the each that returns a constant error?
						),
					},
				)}, // end go-each
			&core.ChooseAction{
				If: &core.CompareNum{
					A:  V("idx"),
					Is: core.Equal,
					B:  F(0),
				},
				// havent found a matching group?
				// pack the object and its settings into it,
				// push the group into the groups.
				Does: core.MakeActivity(
					&list.PutEdge{
						Into: &list.IntoTxtList{Var: N("names")},
						From: V("setting", "name"),
					},
					SetVar("group", "objects", V("names")),
					SetVar("group", "settings", V("settings")),
					&list.PutEdge{Into: &list.IntoRecList{Var: N("groups")}, From: V("group")},
				), // end true
				// found a matching group?
				// unpack it, add the object to it, then pack it up again.
				Else: &core.ChooseNothingElse{Does: core.MakeActivity(
					SetVar("group", &list.ListAt{List: V("groups"), Index: V("idx")}),
					SetVar("names", V("group", "objects")),
					&list.PutEdge{
						Into: &list.IntoTxtList{Var: N("names")},
						V("settings", "name"),
					},
					SetVar("group", "objects", V("names")),
					&list.ListSet{List: W("groups"), Index: V("idx"), From: V("group")},
				)}, // end false
			},
			SetVar("collation", "groups", V("groups")),
		)},
	},
}
