package shapes

import (
	"io/fs"
	"sort"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
)

// reads all of the files in the passed filesystem as ifspecs and returns them as one big json array of shapes
func FromSpecs(files fs.FS) (ret string, err error) {
	if e := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, e error) (err error) {
		if e != nil {
			err = e // can happen if it failed to read the contents of a director
		} else if !d.IsDir() { // the first dir we get is "."
			println("reading", path)
			if _, e := readSpec(files, path); e != nil { // reads into the global lookup
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
				if writeShape(out, blockType) {
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
