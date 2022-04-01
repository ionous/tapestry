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
	var pi int
	for _, p := range flow.Terms {
		if p.Private {
			continue
		}
		var sel string
		if pi = pi + 1; pi > 1 || !flow.Trim {
			sel = camelize(p.Key)
		}
		pt := types.Types[p.TypeName()]
		if simpleSwap := !p.Repeats && pt.Spec.Choice == spec.UsesSpec_Swap_Opt; !simpleSwap {
			var rest [][]string
			for _, a := range sets {
				rest = append(rest, append(a, sel))
			}
			if !p.Optional {
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
					mul = append(mul, append(a, choice))
				}
			}
			sets = mul
		}
	}
	return sets
}
