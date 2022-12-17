package distill

import (
	"sort"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/spec"
)

// Accumulates types to build signatures and lists of slots and slats
// Doesnt currently support multiple groups
type Registry struct {
	Types        []string
	Slots, Slats []string
	Sigs         []Sig
	types        map[string]*spec.TypeSpec
}

func MakeRegistry(types map[string]*spec.TypeSpec) Registry {
	return Registry{types: types}
}

func (reg *Registry) AddType(t *spec.TypeSpec) {
	reg.Types = append(reg.Types, t.Name)
	switch t.Spec.Choice {
	case spec.UsesSpec_Group_Opt:
		// skip
	case spec.UsesSpec_Slot_Opt:
		reg.Slots = append(reg.Slots, t.Name)
	default:
		reg.Slats = append(reg.Slats, t.Name)

		// add signatures:
		switch v := t.Spec.Value.(type) {
		case *spec.FlowSpec:
			reg.addFlow(t, v)

		case *spec.SwapSpec:
			reg.addSwap(t, v)

		case *spec.StrSpec:
			reg.addPrim(t, v.Name)

		case *spec.NumSpec:
			reg.addPrim(t, v.Name)
		}
	}
}

func (reg *Registry) addPrim(t *spec.TypeSpec, lede string) {
	if len(lede) == 0 {
		lede = t.Name
	}
	commandName := Pascal(lede)
	reg.Sigs = append(reg.Sigs, makeSig(t, commandName+":")...)
}

func (reg *Registry) addSwap(t *spec.TypeSpec, swap *spec.SwapSpec) {
	lede := swap.Name
	if len(lede) == 0 {
		lede = t.Name
	}
	commandName := Pascal(t.Name)
	for _, pick := range swap.Between {
		sel := Camelize(pick.Name)
		reg.Sigs = append(reg.Sigs, makeSig(t, commandName+" "+sel+":")...)
	}
}

func (reg *Registry) addFlow(t *spec.TypeSpec, flow *spec.FlowSpec) {
	lede := flow.Name
	if len(lede) == 0 {
		lede = t.Name
	}
	sets := sigParts(flow, Pascal(lede), reg.types)
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
		reg.Sigs = append(reg.Sigs, makeSig(t, sig)...)
	}
}

func (reg *Registry) Sort() {
	sort.Strings(reg.Types)
	sort.Strings(reg.Slots)
	sort.Strings(reg.Slats)
	sort.Slice(reg.Sigs, func(i, j int) bool {
		a, b := reg.Sigs[i], reg.Sigs[j]
		return a.IsLessThan(b)
	})
}
