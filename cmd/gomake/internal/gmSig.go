package gomake

import (
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
)

type Sig struct {
	Type string // ex. Arg.
	Sig  string // ex. "Arg:from:"
	Hash uint64 // ex. 6291103735245333139
}

func makeSig(t *spec.TypeSpec, sig string) (ret []Sig) {
	if strings.Contains(sig, "::") ||
		strings.Contains(sig, "_") ||
		strings.Contains(sig, ": ") {
		log.Fatalln("bad signature for", t.Name, sig)
	}
	// we dont generally need signatures for structs
	// b/c we aren't trying to create those types dynamically from the signature
	// we already have the type, and we're simply deserializing the fields into that type.
	// if len(t.Slots) == 0 {
	// 	h := cin.Hash(sig, "")
	// 	ret = append(ret, Sig{
	// 		Type: t.Name,
	// 		Sig:  h.String,
	// 		Hash: h.Value,
	// 	})
	// }
	for _, slotType := range t.Slots {
		h := cin.Hash(sig, slotType)
		ret = append(ret, Sig{
			Type: t.Name,
			Sig:  h.String,
			Hash: h.Value,
		})
	}
	return
}

// loop over a subset of parameters generating signatures
// where each signature an array of parts.
func sigParts(flow *spec.FlowSpec, commandName string, types rs.TypeSpecs) [][]string {
	var sets = [][]string{{commandName}}
	for _, term := range flow.Terms {
		if term.Private {
			continue
		}
		var sel string
		if !term.IsAnonymous() {
			sel = camelize(term.Label)
		}
		pt := types.Types[term.TypeName()]
		if simpleSwap := !term.Repeats && pt.Spec.Choice == spec.UsesSpec_Swap_Opt; !simpleSwap {
			var rest [][]string
			for _, a := range sets {
				// without copy, the reserve gets re-used, causes a sharing of memory between slices
				// it feels like there should be some simpler way to trigger a reallocing append
				copy := append([]string{}, append(a, sel)...)
				rest = append(rest, copy)
			}
			if !term.Optional {
				sets = rest
			} else {
				sets = append(sets, rest...)
			}
		} else {
			// every choice in a swap gets its own selector for each existing set
			var mul [][]string
			for _, c := range pt.Spec.Value.(*spec.SwapSpec).Between {
				choice := sel + " " + camelize(c.Name)
				for _, a := range sets {
					copy := append([]string{}, append(a, choice)...)
					mul = append(mul, copy)
				}
			}
			sets = mul
		}
	}
	return sets
}