package generate

import (
	"log"
	"strings"
)

type Signature struct {
	Type     string // ex. DiffOf
	Slot     string // ex. number_eval
	Sig      string // ex. "number_eval=Dec:by: "
	Hash     uint64 // ex. 10788210406716082593
	Optional bool
}

func (a *Signature) Body() (ret string) {
	parts := a.parts()
	if cnt := len(parts); cnt > 0 {
		ret = parts[cnt-1]
	}
	return
}

func (a *Signature) IsLessThan(b Signature) (okay bool) {
	as, bs := a.parts(), b.parts()
	if ac, bc := len(as), len(bs); ac < bc {
		okay = true
	} else if ac == bc {
		multipart := ac > 1
		if multipart && as[1] < bs[1] {
			okay = true
		} else if !multipart || as[1] == bs[1] {
			okay = as[0] < bs[0]
		}
	}
	return
}

func (a *Signature) parts() []string {
	return strings.Split(a.Sig, "=")
}

func makeSig(t specData, sig string, slots []string) (ret []Signature) {
	name := t.Name
	if strings.Contains(sig, "::") ||
		strings.Contains(sig, "_") ||
		strings.Contains(sig, ": ") {
		log.Fatalln("bad signature for", name, sig)
	}
	// we dont generally need signatures for structs
	// b/c we aren't trying to create those types dynamically from the signature
	// we already have the type, and we're simply deserializing the fields into that type.
	// ( still it's nice to see them )
	if len(slots) == 0 {
		h := Hash(sig, "")
		ret = append(ret, Signature{
			Type:     name,
			Sig:      h.String,
			Hash:     h.Value,
			Optional: true,
		})
	} else {
		for _, slotName := range slots {
			h := Hash(sig, slotName)
			ret = append(ret, Signature{
				Type: name,
				Slot: slotName,
				Sig:  h.String,
				Hash: h.Value,
			})
		}
	}
	return
}

// loop over a subset of parameters generating signatures
// where each signature an array of parts.
func sigTerms(flow flowData) [][]string {
	commandName := Pascal(flow.Lede)
	var sets = [][]string{{commandName}}
	for _, term := range flow.Terms {
		if term.Private {
			continue
		}
		var sel string
		if term.Label != "_" {
			sel = Camelize(term.Label)
		}
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
	}
	return sets
}
