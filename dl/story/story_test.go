package story_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/dl/story"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/test/debug"

	"git.sr.ht/~ionous/iffy/jsn/cout"
	"git.sr.ht/~ionous/iffy/jsn/din"
)

func TestImportStory(t *testing.T) {
	var curr story.Story
	if e := din.Decode(&curr, iffy.Registry(), []byte(debug.Blob)); e != nil {
		t.Fatal(e)
	} else {
		var els []eph.Ephemera
		k := story.NewImporter(collectEphemera(&els), storyMarshaller)
		if e := k.ImportStory(t.Name(), &curr); e != nil {
			t.Fatal("import", e)
		} else {
			t.Log("ok")
		}
	}
}

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{n} }
func T(s string) *literal.TextValue { return &literal.TextValue{s} }

func P(p string) core.PatternName  { return core.PatternName{Str: p} }
func N(v string) core.VariableName { return core.VariableName{Str: v} }
func V(i string) *core.GetVar      { return &core.GetVar{N(i)} }
func W(v string) string            { return v }

func storyMarshaller(m jsn.Marshalee) (string, error) {
	return cout.Marshal(m, story.CompactEncoder)
}

func collectEphemera(sink *[]eph.Ephemera) story.WriterFun {
	return func(el eph.Ephemera) {
		*sink = append(*sink, el)
	}
}
