package importer_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/ephemera/importer"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"github.com/kr/pretty"
)

func TestEventually(t *testing.T) {
	// Aspect      Aspect      `if:"label=_"`
	// TraitPhrase TraitPhrase `if:"label=trait_phrase"`
	asp := story.AspectTraits{
		Aspect: story.Aspect{Str: "test"},
		TraitPhrase: story.TraitPhrase{
			// AreEither: "canbe"|"either"
			Trait: []story.Trait{{
				Str: "yes",
			}, {
				Str: "no",
			}},
		},
	}

	// first, read the name "test aspect"... and once we have it
	// pass it to the traits, at every trait append
	var name string
	var out []string
	ts := chart.MakeEncoder(nil)
	if e := ts.Marshal(&asp, importer.Block(
		story.AspectTraits_Type,
		func(b jsn.BlockType) (err error) {
			name = asp.Aspect.Str
			ts.ChangeState(importer.Eventually(
				importer.Value(story.Trait_Type, func(v interface{}) (err error) {
					str := *v.(*string)
					out = append(out, name+":"+str)
					return
				})))
			return
		})); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{"test:yes", "test:no"}); len(diff) > 0 {
		t.Fatal(diff)
	} else {
		t.Log("okay")
	}
}
