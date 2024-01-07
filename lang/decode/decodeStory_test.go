package decode_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/test/debug"
	"github.com/kr/pretty"
)

// given a container ( as defined by walker )
// and a json like object tree
// read the json like data into the container;
// overwriting any current data.
func TestDecodeStory(t *testing.T) {
	var file story.StoryFile // fix: it'd be nice to use hand rolled structs or test data instead of story
	var m map[string]any
	if e := json.Unmarshal([]byte(debug.FactorialJs), &m); e != nil {
		t.Fatal(e)
	} else {
		d := decode.Decoder{
			SignatureTable: story.AllSignatures,
			CustomDecoder:  core.CoreDecoder,
			PatternDecoder: story.TryPattern,
		}
		if e := d.Unmarshal(&file, m); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(debug.FactorialStory, &file); len(diff) != 0 {
			pretty.Print(file)
			t.Fatal(diff)
		}
	}
}
