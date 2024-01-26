package generate

import (
	"sort"
	"strings"
)

// Accumulates tapestry type specs to build lists of slots and slats and their signatures.
// Doesnt currently support multiple groups
type Registry []Signature

func (reg Registry) addPrim(t specData) Registry {
	commandName := Pascal(t.Name)
	sigs := makeSig(t, commandName+":", nil)
	return append(reg, sigs...)
}

func (reg Registry) addFlow(t flowData) Registry {
	sets := sigTerms(t)
	for _, set := range sets {
		sig, params := set[0], set[1:] // index 0 is the command name itself
		if len(params) > 0 {
			var next int // if the first parameter is named, it comes before the first colon.
			if first := strings.TrimSpace(params[0]); len(first) > 0 {
				sig += " " + first + ":"
				next++
			}
			// add the rest of the parameters
			if rest := params[next:]; len(rest) > 0 {
				sig += strings.Join(rest, ":") + ":"
			}
		}
		sigs := makeSig(t.specData, sig, t.Slots)
		reg = append(reg, sigs...)
	}
	return reg
}

func (reg Registry) Sort() {
	sort.Slice(reg, func(i, j int) bool {
		a, b := reg[i], reg[j]
		return a.IsLessThan(b)
	})
}
