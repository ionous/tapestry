package shapes

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io/fs"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/idl"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
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

func TestRepeatingContainers(t *testing.T) {
	// reads all of the files in the passed filesystem as ifspecs and returns them as one big json array of shapes
	if reps, e := findRepeatingContainers(idl.Specs); e != nil {
		t.Fatal(e)
	} else {
		for _, rep := range reps {
			t.Log(rep.outer, rep.inner)
		}
	}
}

type repeatingContainer struct {
	outer, inner string
}

func findRepeatingContainers(files fs.FS) (ret []repeatingContainer, err error) {
	// first, read into global "lookup"
	if e := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, e error) (err error) {
		if e != nil {
			err = e // can happen if it failed to read the contents of a director
		} else if !d.IsDir() { // the first dir we get is "."
			println("reading", path)
			if _, e := readSpec(files, path); e != nil {
				err = errutil.New(e, "reading", path)
			}
		}
		return
	}); e != nil {
		err = e
	} else {
		for _, blockType := range lookup {
			// search for flows...
			if flow, ok := blockType.Spec.Value.(*spec.FlowSpec); ok {
				// that have a term...
				for _, t := range flow.Terms {
					// that isnt a special internal term...
					if n := t.TypeName(); !t.Private {
						if ref, ok := lookup[n]; !ok {
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
	lookup = make(TypeSpecs) // reset
	expect := `{
  "type": "story_file",
  "message0": "story_file",
  "colour": "%{BKY_COLOUR_HUE}",
  "tooltip": "top level node, currently just for blockly  might eventually contain story metadata  ex. author, description...",
  "extensions": [
    "tapestry_generic_mixin",
    "tapestry_generic_extension"
  ],
  "mutator": "tapestry_generic_mutation",
  "customData": {
    "mui": "_story_file_mutator",
    "shapeDef": [
      [
        {
          "type": "field_label",
          "text": "story_lines"
        },
        {
          "name": "STORY_LINES",
          "type": "input_value",
          "checks": [
            "story_lines"
          ],
          "shadow": "story_lines"
        }
      ]
    ]
  }
}`
	if _, e := readSpec(idl.Specs, "story.ifspecs"); e != nil { // reads into the global lookup
		t.Fatal(e)
	} else {
		var out js.Builder
		writeShape(&out, lookup["story_file"])
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
