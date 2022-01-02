package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
)

// story encode writes story data without the surrounding Paragraph markers.
// ie. instead of: `{"Story:":[{"Paragraph:":["Example"]},{"Paragraph:":["Example","Example"]}]}`
// it writes: `[["Example"],["Example","Example"]]`
// an outer array of paragraphs, where each direct child is a array of story statements.
func Encode(src *Story) (ret interface{}, err error) {
	x := Paragraph_Slice(src.Paragraph)
	if a, e := cout.Encode(&x, CompactEncoder); e != nil {
		err = e
	} else {
		var out []interface{}
		for _, x := range a.([]interface{}) {
			for _, v := range x.(map[string]interface{}) {
				out = append(out, v)
			}
		}
		ret = out
	}
	return
}

// customized writer of compact data
func CompactEncoder(m jsn.Marshaler, flow jsn.FlowBlock) (err error) {
	switch i, typeName := flow.GetFlow(), flow.GetType(); typeName {
	case NamedNoun_Type:
		ptr := i.(*NamedNoun)
		// to shorten: if there's more than one word we also need a determiner
		det, words := strings.TrimSpace(ptr.Determiner.Str), strings.Fields(ptr.Name.Str)
		if hasDet, cnt := len(det) > 0, len(words); (cnt == 0) || (!hasDet && cnt > 1) {
			err = chart.Unhandled(typeName)
		} else {
			var out strings.Builder
			if hasDet {
				det := jsn.MakeEnum(&ptr.Determiner, &ptr.Determiner.Str).GetCompactValue().(string)
				out.WriteString(det)
			}
			for _, w := range words {
				if hasDet {
					out.WriteRune(' ')
				}
				out.WriteString(w)
				hasDet = true // reuse to add spacing
			}
			err = m.MarshalValue(typeName, out.String())
		}
	default:
		err = core.CompactEncoder(m, flow)
	}
	return
}
