package generate

import (
	"sort"
)

// Accumulates tapestry type specs to build lists of slots and slats and their signatures.
// Doesnt currently support multiple groups
type Registry []Signature

func (reg Registry) addFlow(t flowData) Registry {
	sets := t.Signatures()
	for _, set := range sets {
		sigs := makeSig(t, set)
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
