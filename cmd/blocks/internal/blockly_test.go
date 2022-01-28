package blocks

import (
	"bytes"
	"encoding/json"
	"sort"
	"testing"

	_ "embed"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/idl"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/kr/pretty"
)

func TestBlocklyOutput(t *testing.T) {
	// errutil.Panic = true
	if str, e := run(t); e != nil {
		t.Fatal(e)
	} else {
		var out bytes.Buffer
		if e := json.Indent(&out, []byte(str), "", "  "); e != nil {
			t.Fatal(str)
		} else {
			if diff := pretty.Diff(out.Bytes(), testOutput); len(diff) > 0 {
				t.Log(diff)
				t.Fatal(out.String())
			}
		}
	}
}

func run(t *testing.T) (ret string, err error) {
	if _, e := readSpec(idl.Specs, "prim.ifspecs"); e != nil {
		err = e
	} else if _, e := readSpec(idl.Specs, "literal.ifspecs"); e != nil {
		err = e
	} else /*if _, e := readSpec(idl.Specs, "testdl.ifspecs"); e != nil {
		err = e
	} else */{
		var keys []string
		for k, _ := range lookup {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		ret = js.Embrace(js.Array, func(out *js.Builder) {
			var csv bool
			for _, key := range keys {
				blockType := lookup[key]
				if csv {
					out.R(js.Comma)
				}
				if !writeBlock(out, blockType) {
					csv = false // good enough for testing.
				} else {
					csv = true
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

//go:embed testOutput.json
var testOutput []byte
