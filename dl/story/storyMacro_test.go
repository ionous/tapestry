package story_test

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
	"testing"
)

func TestMacroImport(t *testing.T) {
	var els []eph.Ephemera
	k := imp.NewImporter(collectEphemera(&els))
	if e := story.ImportStory(k, t.Name(), macroStory); e != nil {
		t.Fatal(e)
	} else {
		t.Log(els)
	}
}

// tbd: is there anyway to make the first parameter anonymous? does that even matter in blockly?
// Carrier actor:carries: [X, Ys]
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
					// tbd: allow variables to have determiners? (and strip them off during import or weave)
					Nouns: core.Variable("actor"),
					Kind:  literal.T("an actor"),
				},
				// 2. each thing is a thing ( the relation is any object so... )
				&story.DefineNouns{
					Nouns: core.Variable("carries"),
					Kind:  literal.T("things"),
				},
				// 3. use whereabouts
				&story.DefineRelatives{
					Nouns:      core.Variable("carries"),
					Relation:   literal.T("whereabouts"),
					OtherNouns: core.Variable("actor"),
				},
			},
		},
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
}
