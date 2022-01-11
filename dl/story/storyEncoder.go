package story

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
)

// story encode writes story data without the surrounding Paragraph markers.
// ie. instead of: `{"Story:":[{"Paragraph:":["Example"]},{"Paragraph:":["Example","Example"]}]}`
// it writes: `[["Example"],["Example","Example"]]`
// an outer array of paragraphs, where each direct child is a array of story statements.
func Encode(src *Story) (interface{}, error) {
	return CustomEncode(src, CompactEncoder)
}

// helper exposed as public for some tests
func CustomEncode(src *Story, encoder func(jsn.Marshaler, jsn.FlowBlock) error) (ret interface{}, err error) {
	x := Paragraph_Slice(src.Paragraph)
	if a, e := cout.Encode(&x, encoder); e != nil {
		err = e
	} else {
		var out []interface{}
		for _, x := range a.([]interface{}) {
			switch els := x.(type) {
			case string:
				// empty paragraph. i think its fine to elide it.
			case map[string]interface{}:
				for _, v := range els {
					out = append(out, v)
				}
			}
		}
		ret = out
	}
	return
}

// customized writer of compact data
var CompactEncoder = core.CompactEncoder
