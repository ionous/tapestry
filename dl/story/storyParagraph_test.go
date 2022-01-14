package story_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"github.com/kr/pretty"
)

// one possible way of compacting paragraphs in story.if files:
// they are implicitly arrays of arrays.
func TestCompactParagraph(t *testing.T) {
	body := `[["Example"],["Example","Example"]]`
	var out story.Story
	story.Decode(&out, []byte(body), cin.Signatures{{
		cin.Hash("Paragraph:"): (*story.Paragraph)(nil),
		cin.Hash("Story:"):     (*story.Story)(nil),
		cin.Hash("Example"):    (*Example)(nil),
	}})
	match := story.Story{
		Paragraph: []story.Paragraph{{
			StoryStatement: []story.StoryStatement{&Example{}},
		}, {
			StoryStatement: []story.StoryStatement{&Example{}, &Example{}},
		}}}
	if diff := pretty.Diff(out, match); len(diff) > 0 {
		pretty.Println("got", out)
		pretty.Println("got", match)
		t.Fatal("fail", diff)
	}
	if d, e := story.Encode(&out); e != nil {
		t.Fatal(e)
	} else if b, e := json.Marshal(d); e != nil {
		t.Fatal(e)
	} else {
		b := string(b)
		if body != b {
			t.Fatal(b)
		} else {
			// {"Story:":[{"Paragraph:":["Example"]},{"Paragraph:":["Example","Example"]}]}
			t.Log(b)
		}
	}
}

type Example struct{}

func (d *Example) ImportPhrase(*story.Importer) error { return nil }

func (op *Example) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(Example_Flow{op}); err == nil {
		m.EndBlock()
	}
	return
}

type Example_Flow struct{ ptr *Example }

func (n Example_Flow) GetType() string      { return "Example" }
func (n Example_Flow) GetLede() string      { return "Example" }
func (n Example_Flow) GetFlow() interface{} { return n.ptr }
func (n Example_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Example); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}
