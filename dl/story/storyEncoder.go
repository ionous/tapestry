package story

import (
	"strings"

	"git.sr.ht/~ionous/iffy/dl/story"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
)

func CompactEncoder(m jsn.Marshaler, flow jsn.FlowBlock) (err error) {
	switch i, typeName := flow.GetFlow(), flow.GetType(); typeName {
	case NamedNoun_Type:
		ptr := i.(*story.NamedNoun)
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
