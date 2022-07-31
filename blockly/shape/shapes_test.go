package shape_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io/fs"
	"testing"

	"git.sr.ht/~ionous/tapestry/blockly/shape"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/idl"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// fix: generate and compare just testdl?
// right now this just tests that the shapes are well formed.
func TestBlocklyTypes(t *testing.T) {
	if str, e := shape.FromSpecs(idl.Specs); e != nil {
		t.Fatal(e)
	} else if out, e := indent(str); e != nil {
		t.Fatal(e)
	} else {
		t.Log(out)
	}
}

func TestRepeatingContainers(t *testing.T) {
	// reads all of the files in the passed filesystem as ifspecs and returns them as one big json array of shapes
	if reps, e := findRepeatingContainers(idl.Specs); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(reps,
		[]repeatingContainer{
			{"field_list", "field_value"},
			{"paragraph", "story_statement"},
		}); len(diff) > 0 {
		for _, rep := range reps {
			t.Log(rep.outer, rep.inner)
		}
		t.Fatal("dont want to use any repeating containers")
	}
}

type repeatingContainer struct {
	outer, inner string
}

func findRepeatingContainers(files fs.FS) (ret []repeatingContainer, err error) {
	if ts, e := rs.FromSpecs(files); e != nil {
		err = e
	} else {
		for _, k := range ts.Keys() {
			blockType := ts.Types[k]
			// search for flows...
			if flow, ok := blockType.Spec.Value.(*spec.FlowSpec); ok {
				// that have a term...
				for _, t := range flow.Terms {
					// that isnt a special internal term...
					if n := t.TypeName(); !t.Private {
						if ref, ok := ts.Types[n]; !ok {
							err = errutil.New("couldnt find", n)
						} else {
							// which is a flow...
							if flow, ok := ref.Spec.Value.(*spec.FlowSpec); ok {
								// containing a single repeating term
								if len(flow.Terms) == 1 && flow.Terms[0].Repeats {
									ret = append(ret, repeatingContainer{n, flow.Terms[0].TypeName()})
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

// make sure that story file has no output and one stacked input.
func TestStoryFileShape(t *testing.T) {
	expect := `{
  "type": "story_file",
  "colour": "%{BKY_COLOUR_HUE}",
  "tooltip": "top level node, currently just for blockly  might eventually contain story metadata  ex. author, description...",
  "extensions": [
    "tapestry_generic_mixin",
    "tapestry_generic_extension"
  ],
  "mutator": "tapestry_generic_mutation",
  "customData": {
    "shapeDef": [
      {
        "label": "Tapestry"
      },
      {
        "name": "STORY_LINES",
        "type": "input_statement",
        "checks": [
          "_story_statement_stack"
        ]
      }
    ]
  }
}`
	if ts, e := rs.ReadSpec(idl.Specs, "story.ifspecs"); e != nil {
		t.Fatal(e)
	} else if x, ok := ts.Types["story_file"]; !ok {
		t.Fatal("missing story file type")
	} else {
		var out js.Builder
		w := shape.ShapeWriter{ts}
		w.WriteShape(&out, x)
		//
		if str, e := indent(out.String()); e != nil {
			t.Fatal(e, str)
		} else if diff := pretty.Diff(str, expect); len(diff) > 0 {
			t.Log(str)
			t.Fatal("ng", diff)
		}
	}
}

// make sure that story file has no output and one stacked input.
func TestStoryTextShape(t *testing.T) {
	expect := `{
  "type": "text_field",
  "output": [
    "text_field",
    "field"
  ],
  "colour": "%{BKY_COLOUR_HUE}",
  "extensions": [
    "tapestry_generic_mixin",
    "tapestry_generic_extension"
  ],
  "mutator": "tapestry_generic_mutation",
  "customData": {
    "shapeDef": [
      {
        "label": "Text"
      },
      {
        "name": "NAME",
        "type": "mosaic_text_field",
        "text": "name"
      },
      {
        "label": "kind",
        "name": "TYPE",
        "type": "mosaic_text_field",
        "text": "kind",
        "optional": true
      },
      {
        "label": "initially",
        "name": "INITIALLY",
        "type": "input_value",
        "checks": [
          "text_eval"
        ],
        "optional": true
      }
    ]
  }
}`
	if ts, e := rs.FromSpecs(idl.Specs); e != nil {
		t.Fatal(e)
	} else if x, ok := ts.Types["text_field"]; !ok {
		t.Fatal("missing story file type")
	} else {
		var out js.Builder
		w := shape.ShapeWriter{ts}
		w.WriteShape(&out, x)
		//
		if str, e := indent(out.String()); e != nil {
			t.Fatal(e, str)
		} else if diff := pretty.Diff(str, expect); len(diff) > 0 {
			t.Log(str)
			t.Fatal("ng", diff)
		}
	}
}

func indent(str string) (ret string, err error) {
	var indent bytes.Buffer
	if e := json.Indent(&indent, []byte(str), "", "  "); e != nil {
		ret = str
	} else {
		ret = indent.String()
	}
	return
}
