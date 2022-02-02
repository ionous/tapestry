package shapes

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io/fs"
	"sort"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/idl"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// fix: generate and test just testdl?
func TestBlocklyTypes(t *testing.T) {
	if str, e := run(t); e != nil {
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
    "blockDef": [
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
		writeBlock(&out, lookup["story_lines"])
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

//   "type": "story_lines",
//   "message0": "story_lines",
//   "colour": "%{BKY_COLOUR_HUE}",
//   "extensions": [
//     "tapestry_generic_mixin",
//     "tapestry_generic_extension"
//   ],
//   "mutator": "tapestry_generic_mutation",
//   "customData": {
//     "blockDef": [
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

func run(t *testing.T) (ret string, err error) {
	if e := fs.WalkDir(idl.Specs, ".", func(path string, d fs.DirEntry, e error) (err error) {
		if e != nil {
			err = e // can happen if it failed to read the contents of a director
		} else if !d.IsDir() { // the first dir we get is "."
			println("reading", path)
			if _, e := readSpec(idl.Specs, path); e != nil { // reads into the global lookup
				err = errutil.New(e, "reading", path)
			}
		}
		return
	}); e != nil {
		err = e
	} else {
		var keys []string
		for k, _ := range lookup {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		ret = js.Embrace(js.Array, func(out *js.Builder) {
			var comma bool
			for _, key := range keys {
				blockType := lookup[key]
				if comma {
					out.R(js.Comma)
					comma = false
				}
				if writeBlock(out, blockType) {
					comma = true
					if flow, ok := blockType.Spec.Value.(*spec.FlowSpec); ok {
						out.R(js.Comma)
						writeMutator(out, blockType, flow)
					}
				}
			}
		})
	}
	return
}
