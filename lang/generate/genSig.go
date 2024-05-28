package generate

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

type Signature struct {
	Type     string // ex. DiffOf
	Slot     string // ex. number_eval
	Sig      string // ex. "number_eval=Dec:by: "
	Hash     uint64 // ex. 10788210406716082593
	Optional bool
}

// the right side of the equal sign
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

// returns lhs slot, and rhs signature body.
func (a *Signature) parts() []string {
	return strings.Split(a.Sig, "=")
}

func makeSig(t flowData, set sigTerm) (ret []Signature) {
	if sig := set.Signature(); strings.Contains(sig, "::") ||
		strings.Contains(sig, "_") ||
		strings.Contains(sig, ": ") {
		log.Fatalln("bad signature for", t.Name, sig)
	} else if slots := t.Slots; len(slots) == 0 {
		// we dont generally need signatures for structs
		// b/c we aren't trying to create those types dynamically from the signature
		// we already have the type, and we're simply deserializing the fields into that type.
		// ( still it's nice to see them )
		h := Hash(sig, "")
		ret = append(ret, Signature{
			Type:     t.Name,
			Sig:      h.String,
			Hash:     h.Value,
			Optional: true,
		})
	} else {
		for _, slotName := range slots {
			h := Hash(sig, slotName)
			ret = append(ret, Signature{
				Type: t.Name,
				Slot: slotName,
				Sig:  h.String,
				Hash: h.Value,
			})
		}
	}
	return
}

type TypeLink string

func (t TypeLink) Type() string {
	return string(t)
}

func (t TypeLink) Link() (ret string, err error) {
	if a, ok := hackForLinks.linkByName(t.Type()); !ok {
		err = fmt.Errorf("unknown type %q creating link", t)
	} else {
		ret = a
	}
	return
}

// a signature and the the terms it uses
type sigTerm struct {
	Flow  flowData // ugh. but useful for templates
	parts []string // to be separated by colons
	terms []termData
}

func (set sigTerm) Terms() []termData {
	return set.terms
}

// unique types used by the terms
func (set sigTerm) CommandLink() (ret []TypeLink) {
	for _, a := range set.terms {
		if typeName := a.Type; !slices.Contains(ret, TypeLink(typeName)) {
			ret = append(ret, TypeLink(typeName))
		}
	}
	return
}

// unique types used by the terms
func (set sigTerm) TypeLinks() (ret []TypeLink) {
	for _, a := range set.terms {
		if typeName := a.Type; !slices.Contains(ret, TypeLink(typeName)) {
			ret = append(ret, TypeLink(typeName))
		}
	}
	return
}

func (set sigTerm) TrimmedSignature() string {
	sig := set.Signature()
	end := len(sig) - 1
	if end > 0 && sig[end] == ':' {
		sig = sig[:end]
	}
	return sig
}

func (set sigTerm) Lede() string {
	return set.parts[0]
}

func (set sigTerm) Signature() string {
	sig, params := set.parts[0], set.parts[1:] // index 0 is the command name itself
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
	return sig
}

// loop over a subset of parameters generating signatures
// where each signature an array of parts.
func sigTerms(flow flowData) []sigTerm {
	commandName := Pascal(flow.Lede)
	var sets = []sigTerm{{Flow: flow, parts: []string{commandName}}}
	for _, term := range flow.Terms {
		if term.Private {
			continue
		}
		var sel string
		if term.Label != "_" {
			sel = Camelize(term.Label)
		}
		var rest []sigTerm
		for _, a := range sets {
			// without copy, the reserve gets re-used, causes a sharing of memory between slices
			// it feels like there should be some simpler way to trigger a reallocing append
			newParts := append([]string{}, append(a.parts, sel)...)
			newTerms := append([]termData{}, append(a.terms, term)...)
			rest = append(rest, sigTerm{
				Flow:  flow,
				parts: newParts,
				terms: newTerms,
			})
		}
		if !term.Optional {
			sets = rest
		} else {
			sets = append(sets, rest...)
		}
	}
	return sets
}
