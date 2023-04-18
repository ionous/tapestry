package story_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// generate ephemera for macros
func TestMacroEphemera(t *testing.T) {
	errutil.Panic = true
	var els []eph.Ephemera
	k := imp.NewImporter(collectEphemera(&els))
	if e := story.ImportStory(k, t.Name(), macroStory); e != nil {
		t.Fatal(e)
	} else {
		expect := []eph.Ephemera{
			&eph.EphValues{
				Noun:  "Hershel",
				Field: "proper_named",
				Value: literal.B(true),
			},
			&eph.EphNouns{
				Noun: "Hershel",
				Kind: "an actor",
			},
			&eph.EphNouns{
				Noun: "scissors",
				Kind: "things",
			},
			&eph.EphRelatives{
				Rel:       "whereabouts",
				Noun:      "Hershel",
				OtherNoun: "scissors",
			},
		}
		if diff := pretty.Diff(els, expect); len(diff) != 0 {
			t.Log(pretty.Sprint(els))
			t.Fatal(diff)
		}
	}
}

var macroStory = &story.StoryFile{
	StoryStatements: []story.StoryStatement{
		&story.DefineMacro{
			MacroName: core.T("carrier"),
			Params: []story.FieldDefinition{
				// todo: optional labels for fields so these can be "carrier", "nouns" internally
				// ( and maybe that would allow the anonymous first parameter )
				&story.TextField{
					Name: "actor",
				},
				&story.TextListField{
					Name: "carries",
				},
			},
			Result: &story.NothingField{},
			Locals: nil,
			MacroStatements: []rt.Execute{
				// assert:
				// 1. the actor is an actor
				&story.DefineNouns{
					// fix: autoconversions of text to text list?
					Nouns: &list.MakeTextList{
						// tbd: allow variables to have determiners? (and strip them off during import or weave)
						Values: []rt.TextEval{core.Variable("actor")},
					},
					Kind: literal.T("an actor"),
				},
				// 2. each thing is a thing ( the relation is any object so... )
				&story.DefineNouns{
					Nouns: core.Variable("carries"),
					Kind:  literal.T("things"),
				},
				// 3. use whereabouts
				&story.DefineRelatives{
					Nouns:    core.Variable("carries"),
					Relation: literal.T("whereabouts"),
					OtherNouns: &list.MakeTextList{
						Values: []rt.TextEval{core.Variable("actor")},
					},
				},
			},
		},
		&story.DefineScene{
			Scene: literal.T("carrier"),
			With: []story.StoryStatement{
				&story.CallMacro{
					MacroName: "carrier",
					Arguments: []assign.Arg{{
						Name:  "actor",
						Value: &assign.FromText{Value: literal.T("Hershel")},
					}, {
						Name:  "carries",
						Value: &assign.FromTextList{Value: literal.Ts("the scissors")},
					}}},
			},
		},
	},
}
