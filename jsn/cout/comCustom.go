package cout

import (
	"strings"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
)

func (m *xEncoder) customFlow(flow jsn.FlowBlock) (err error) {
	switch i, typeName := flow.GetFlow(), flow.GetType(); typeName {
	default:
		err = chart.Unhandled("CustomFlow")

	case core.BoolValue_Type:
		var out bool = i.(*core.BoolValue).Bool
		err = m.MarshalValue(typeName, out)

	case core.NumValue_Type:
		var out float64 = i.(*core.NumValue).Num
		err = m.MarshalValue(typeName, out)

	case core.Numbers_Type:
		var out []float64 = i.(*core.Numbers).Values
		err = m.MarshalValue(typeName, out)

		// write text as a raw string
	case core.TextValue_Type:
		str := i.(*core.TextValue).Text
		// if the text starts with an @, add another @
		if len(str) > 0 && str[0] == '@' {
			str = "@" + str
		}
		err = m.MarshalValue(typeName, str)

	case core.Texts_Type:
		var out []string = i.(*core.Texts).Values
		err = m.MarshalValue(typeName, out)

	case story.NamedNoun_Type:
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

	// write variables as a string prepended by @
	case core.GetVar_Type:
		ptr := i.(*core.GetVar)
		str := ptr.Name.Str
		// a leading ampersand would conflict with @@ escaped text serialization.
		if leadingAmp := len(str) > 0 && str[0] == '@'; !leadingAmp {
			err = m.MarshalValue(typeName, "@"+str)
		} else {
			err = chart.Unhandled(typeName)
		}
	}
	return
}
