package story_test

import (
	"fmt"
	"testing"

	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/dl/story"
	"github.com/kr/pretty"
)

func TestImportNamedNouns(t *testing.T) {
	var els []eph.Ephemera
	k := story.NewImporter(collectEphemera(&els), storyMarshaller)
	//
	for i, nouns := 0, []string{
		"our", "Trevor",
		"an", "apple",
		"3", "triangles",
		"one", "square",
		"a gaggle of", "robot sheep",
	}; i < len(nouns); i += 2 {
		n := makeNoun(nouns[i], nouns[i+1])
		if e := n.ImportNouns(k); e != nil {
			t.Fatal(e, "at", i)
		}
	}
	expect := concat([]eph.Ephemera{
		// implicitly generated object type
		&eph.EphKinds{
			Kinds: "objects",
			From:  "kind",
		},
		// various implicitly defined aspects
		&eph.EphAspects{
			Aspects: "noun_types",
			Traits:  []string{"common_named", "proper_named", "counted"},
		},
		&eph.EphKinds{
			Kinds:   "objects",
			Contain: []eph.EphParams{eph.AspectParam("noun_types")},
		},
		&eph.EphAspects{
			Aspects: "private_names",
			Traits:  []string{"publicly_named", "privately_named"},
		},
		&eph.EphKinds{
			Kinds:   "objects",
			Contain: []eph.EphParams{eph.AspectParam("private_names")},
		},
		&eph.EphKinds{
			Kinds:   "objects",
			Contain: []eph.EphParams{textParam("indefinite_article")},
		},
		// the printed name field is implicitly generated by the counted nouns
		&eph.EphKinds{
			Kinds:   "objects",
			Contain: []eph.EphParams{textParam("printed_name")},
		},
		// note: we dont expect to see &EphNoun;
		// the phrases we are using for the test only add values.
		// ( *except* for counted nouns because they're implicitly generated )
		&eph.EphValues{
			Noun:  "Trevor",
			Field: "indefinite_article",
			Value: T("our"),
		},
		&eph.EphValues{
			Noun:  "Trevor",
			Field: "proper_named",
			Value: B(true),
		},
		&eph.EphValues{
			Noun:  "apple",
			Field: "indefinite_article",
			Value: T("an"),
		}},
		countedNouns("triangles", 3),
		countedNouns("square", 1),
		// finally, our last object...
		[]eph.Ephemera{
			&eph.EphValues{
				Noun:  "robot sheep",
				Field: "indefinite_article",
				Value: T("a gaggle of"),
			},
		},
	)
	els = append(k.Queued(), els...)
	if diff := pretty.Diff(els, expect); len(diff) > 0 {
		t.Log(diff)
		t.Error(pretty.Sprint(els))
	}
}

func makeNoun(det, name string) story.NamedNoun {
	return story.NamedNoun{
		Determiner: story.Determiner{Str: det},
		Name:       story.NounName{Str: name},
	}
}

func concat(els ...[]eph.Ephemera) (out []eph.Ephemera) {
	for _, k := range els {
		out = append(out, k...)
	}
	return
}

// fix? previously the name would have been "triangle_1" now its "triangles_1"
func countedNouns(kind string, cnt int) (out []eph.Ephemera) {
	out = append(out, &eph.EphKinds{
		Kinds: kind,
		From:  "thing",
	})
	for i := 1; i <= cnt; i++ {
		n := fmt.Sprintf("%s_%d", kind, i)
		out = append(out,
			&eph.EphNouns{Noun: n, Kind: kind},
			&eph.EphValues{
				Noun:  n,
				Field: "counted",
				Value: B(true),
			},
			&eph.EphValues{
				Noun:  n,
				Field: "printed_name",
				Value: T(kind),
			})
	}
	return out
}

func textParam(name string) eph.EphParams {
	return eph.EphParams{
		Affinity: eph.Affinity{eph.Affinity_Text},
		Name:     name,
	}
}
