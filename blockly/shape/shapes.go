package shape

import (
	"io/fs"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// read all tapestry idl files from the filesystem and return them as an json array of shapes
func FromSpecs(files fs.FS) (ret string, err error) {
	if ts, e := rs.FromSpecs(files); e != nil {
		err = e
	} else {
		w := NewShapeWriter(ts)
		ret = js.Embrace(js.Array, func(out *js.Builder) {
			var comma bool
			for _, key := range ts.Keys() {
				blockType := ts.Types[key]
				if comma {
					out.R(js.Comma)
					comma = false
				}
				if w.WriteShape(out, blockType) {
					comma = true
					if flow, ok := blockType.Spec.Value.(*spec.FlowSpec); ok {
						out.R(js.Comma)
						w.writeMutator(out, blockType, flow)
					}
				}
			}
		})
	}
	return
}
