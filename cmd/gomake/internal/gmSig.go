package gomake

import (
	"hash/fnv"
	"io"
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
)

type Sig struct {
	Type string // ex. Arg.
	Sig  string // ex. "Arg:from:"
	Hash uint64 // ex. 6291103735245333139
}

func makeSig(typeName string, sig string) Sig {
	if strings.Contains(sig, "::") ||
		strings.Contains(sig, "_") ||
		strings.Contains(sig, ": ") {
		log.Fatalln("bad signature for", typeName, sig)
	}
	w := fnv.New64a()
	io.WriteString(w, sig)
	return Sig{
		Hash: w.Sum64(),
		Type: typeName,
		Sig:  sig,
	}
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
			sel = camelize(term.Key)
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
