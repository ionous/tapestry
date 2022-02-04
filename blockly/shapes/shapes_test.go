package shapes

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/idl"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/kr/pretty"
)

// fix: generate and test just testdl?
func TestBlocklyTypes(t *testing.T) {
	if str, e := FromSpecs(idl.Specs); e != nil {
		t.Fatal(e)
	} else if out, e := indent(str); e != nil {
		t.Fatal(e)
	} else {
		t.Log(out)
	}
}

// make sure that story lines has no output and one stacked input.
func TestStoryLineShape(t *testing.T) {
	lookup = make(TypeSpecs) // reset
	expect := `{
  "type": "story_lines",
  "message0": "story_lines",
  "colour": "%{BKY_COLOUR_HUE}",
  "extensions": [
    "tapestry_generic_mixin",
    "tapestry_generic_extension"
  ],
  "mutator": "tapestry_generic_mutation",
  "customData": {
    "mui": "_story_lines_mutator",
    "shapeDef": [
      [
        {
          "type": "field_label",
          "text": "lines"
        },
        {
          "name": "LINES",
          "type": "input_statement",
          "checks": [
            "_story_statement_stack"
          ]
        }
      ]
    ]
  }
}`
	if _, e := readSpec(idl.Specs, "story.ifspecs"); e != nil { // reads into the global lookup
		t.Fatal(e)
	} else {
		var out js.Builder
		writeShape(&out, lookup["story_lines"])
		//
		if str, e := indent(out.String()); e != nil {
			t.Fatal(e, str)
		} else if diff := pretty.Diff(str, expect); len(diff) > 0 {
			t.Log(str)
			t.Fatal("ng", diff)
		}
	}
}

// make sure that story lines has no output and one stacked input.
func TestStoryLineMui(t *testing.T) {
	lookup = make(TypeSpecs) // reset
	expect := `{
  "type": "_story_lines_mutator",
  "style": "logic_blocks",
  "inputsInline": false,
  "args0": [
    {
      "type": "field_label",
      "text": "story_lines"
    },
    {
      "type": "input_dummy"
    }
  ],
  "message0": "%1%2"
}`
	if _, e := readSpec(idl.Specs, "story.ifspecs"); e != nil { // reads into the global lookup
		t.Fatal(e)
	} else {
		blockType := lookup["story_lines"]
		if flow, ok := blockType.Spec.Value.(*spec.FlowSpec); !ok {
			t.Fatal("unexpected")
		} else {
			var out js.Builder
			writeMutator(&out, blockType, flow)
			if str, e := indent(out.String()); e != nil {
				t.Fatal(e, str)
			} else if diff := pretty.Diff(str, expect); len(diff) > 0 {
				t.Log(str)
				t.Fatal("ng", diff)
			}
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

//   "type": "story_lines",
//   "message0": "story_lines",
//   "colour": "%{BKY_COLOUR_HUE}",
//   "extensions": [
//     "tapestry_generic_mixin",
//     "tapestry_generic_extension"
//   ],
//   "mutator": "tapestry_generic_mutation",
//   "customData": {
//     "shapeDef": [
//       [
//         {
//           "type": "field_label",
//           "text": "lines"
//         },
//         {
//           "name": "LINES",
//           "type": "input_statement",
//           "checks": [
//             "stacked_story_statement"
//           ]
//         }
//       ]
//     ]
//   }
// }
