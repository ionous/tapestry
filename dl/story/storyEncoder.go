package story

import (
	"encoding/json"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"github.com/ionous/errutil"
)

// customized writer of compact data
func CompactEncoder(m jsn.Marshaler, flow jsn.FlowBlock) (err error) {
	// typeName := flow.GetType()
	// switch ptr := flow.GetFlow().(type) {
	// case *NamedNoun:

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

func Decode(dst jsn.Marshalee, msg json.RawMessage, sig cin.Signatures) error {
	return cin.NewDecoder(sig).
		SetFlowDecoder(CompactFlowDecoder).
		SetSlotDecoder(CompactSlotDecoder).
		Decode(dst, msg)
}

var CompactSlotDecoder = core.CompactSlotDecoder

// customized reader of compact data
func CompactFlowDecoder(flow jsn.FlowBlock, msg json.RawMessage) (err error) {
	// switch typeName, ptr := flow.GetType(), flow.GetFlow(); ptr.(type) {
	// default:
	// 	err = chart.Unhandled(typeName)

	// case *NamedNoun:
	// 	var str string
	// 	if e := json.Unmarshal(msg, &str); e != nil {
	// 		err = chart.Unhandled(typeName)
	// 	} else {
	// 		var out NamedNoun
	// 		if space := strings.Index(str, " "); space < 0 {
	// 			out.Name.Str = str
	// 		} else {
	// 			jsn.MakeEnum(&out.Determiner, &out.Determiner.Str).SetValue(str[:space])
	// 			out.Name.Str = str[space+1:]
	// 		}
	// 		if !flow.SetFlow(&out) {
	// 			err = errutil.New("could set result to flow", typeName, flow)
	// 		}
	// 	}
	// }
	switch typeName := flow.GetType(); typeName {
	default:
		err = chart.Unhandled("CustomFlow")

	case NamedNoun_Type:
		var str string
		if e := json.Unmarshal(msg, &str); e != nil {
			err = chart.Unhandled(typeName)
		} else {
			var out NamedNoun
			if space := strings.Index(str, " "); space < 0 {
				out.Name.Str = str
			} else {
				jsn.MakeEnum(&out.Determiner, &out.Determiner.Str).SetValue(str[:space])
				out.Name.Str = str[space+1:]
			}
			if !flow.SetFlow(&out) {
				err = errutil.New("could set result to flow", typeName, flow)
			}
		}
	}
	return
}
