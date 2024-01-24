package shape_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/blockly/shape"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/frame"
	"git.sr.ht/~ionous/tapestry/dl/game"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rel"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// fix: generate and compare just testdl?
// right now this just tests that the shapes are well formed.
func TestBlocklyTypes(t *testing.T) {
	if str, e := shape.FromTypes(blocks); e != nil {
		t.Fatal(e)
	} else if !json.Valid([]byte(str)) {
		t.Fatal(e)
	} else {
		t.Log(str)
	}
}

var blocks = []*typeinfo.TypeSet{
	&assign.Z_Types,
	&core.Z_Types,
	&debug.Z_Types,
	&frame.Z_Types,
	&game.Z_Types,
	&grammar.Z_Types,
	&list.Z_Types,
	&literal.Z_Types,
	// &play.Z_Types,
	&prim.Z_Types,
	&rel.Z_Types,
	&render.Z_Types,
	&rtti.Z_Types,
	// &spec.Z_Types,
	&story.Z_Types,
	// &testdl.Z_Types,
}

func TestRepeatingContainers(t *testing.T) {
	// read all tapestry idl files from the filesystem
	m := shape.MakeTypeMap(blocks)
	if reps, e := findRepeatingContainers(m); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(reps,
		[]repeatingContainer{
			{"field_list", "field_value"},
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

func findRepeatingContainers(ts shape.TypeMap) (ret []repeatingContainer, err error) {
	for _, k := range ts.Keys() {
		blockType := ts[k]
		// search for flows...
		if flow, ok := blockType.(*typeinfo.Flow); ok {
			// that have a term...
			for _, t := range flow.Terms {
				// that isnt a special internal term...
				if !t.Private {
					n := t.Type.TypeName()
					if ref, ok := ts[n]; !ok {
						err = errutil.New("couldnt find", n)
					} else {
						// which is a flow...
						if flow, ok := ref.(*typeinfo.Flow); ok {
							// containing a single repeating term
							if len(flow.Terms) == 1 && flow.Terms[0].Repeats {
								ret = append(ret, repeatingContainer{n, flow.Terms[0].Type.TypeName()})
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
  "colour": "%{BKY_TAP_HUE_ROOT}",
  "tooltip": "top level node, currently just for blockly might eventually contain story metadata  ex. author, description...",
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
        "name": "STORY_STATEMENTS",
        "type": "input_statement",
        "checks": [
          "_story_statement_stack"
        ]
      }
    ]
  }
}`
	x := &story.Zt_StoryFile
	ts := shape.TypeMap{x.Name: x}

	var out js.Builder
	w := shape.ShapeWriter{ts}
	w.WriteShape(&out, x)
	//
	got := jsn.Indent(out.String())
	if got != expect {
		t.Log("have: \n", got)
		t.Log("want: \n", expect)
		t.Fatal("ng", len(got), len(expect))
	}
}

// make sure that story file has no output and one stacked input.
func TestStoryTextShape(t *testing.T) {
	expect := `{
  "type": "text_field",
  "output": [
    "text_field",
    "field_definition"
  ],
  "colour": "%{BKY_TAP_HUE}",
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
        "type": "input_value",
        "checks": [
          "text_eval"
        ]
      },
      {
        "label": "kind",
        "name": "TYPE",
        "type": "input_value",
        "checks": [
          "text_eval"
        ],
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
	x := &story.Zt_TextField
	ts := shape.TypeMap{x.Name: x}

	var out js.Builder
	w := shape.ShapeWriter{ts}
	w.WriteShape(&out, x)
	//
	str := jsn.Indent(out.String())
	if diff := pretty.Diff(str, expect); len(diff) > 0 {
		t.Log(str)
		t.Fatal("ng", diff)
	}
}

// test the generation of an enumeration
func TestStrEnum(t *testing.T) {
	expect := `{
  "type": "compare_text",
  "output": [
    "compare_text",
    "bool_eval"
  ],
  "colour": "%{BKY_LOGIC_HUE}",
  "tooltip": "True if eq,ne,gt,lt,ge,le two strings ( lexical. ).",
  "extensions": [
    "tapestry_generic_mixin",
    "tapestry_generic_extension"
  ],
  "mutator": "tapestry_generic_mutation",
  "customData": {
    "shapeDef": [
      {
        "label": "Is"
      },
      {
        "name": "A",
        "type": "input_value",
        "checks": [
          "text_eval"
        ]
      },
      {
        "label": "matching",
        "name": "IS",
        "type": "field_dropdown",
        "options": [
          [
            "equal to",
            "$EQUAL_TO"
          ],
          [
            "other than",
            "$OTHER_THAN"
          ],
          [
            "greater than",
            "$GREATER_THAN"
          ],
          [
            "less than",
            "$LESS_THAN"
          ],
          [
            "at least",
            "$AT_LEAST"
          ],
          [
            "at most",
            "$AT_MOST"
          ]
        ]
      },
      {
        "label": "text",
        "name": "B",
        "type": "input_value",
        "checks": [
          "text_eval"
        ]
      }
    ]
  }
}`
	x := &core.Zt_CompareText
	ts := shape.TypeMap{x.Name: x}

	var out js.Builder
	w := shape.ShapeWriter{ts}
	w.WriteShape(&out, x)
	//
	str := jsn.Indent(out.String())
	if diff := pretty.Diff(str, expect); len(diff) > 0 {
		t.Log(str)
		t.Fatal("ng", diff)
	}
}
