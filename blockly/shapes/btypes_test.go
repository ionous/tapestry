package btypes

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
)

// fix: generate and test just testdl?
func TestBlocklyTypes(t *testing.T) {
	if str, e := run(t); e != nil {
		t.Fatal(e)
	} else {
		var out bytes.Buffer
		if e := json.Indent(&out, []byte(str), "", "  "); e != nil {
			t.Fatal(str)
		} else {
			t.Log(out.String())
		}
	}
}

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
