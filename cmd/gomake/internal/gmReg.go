package gomake

import (
	"io"
	"sort"
	"strings"
	"text/template"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"github.com/ionous/errutil"
)

type RegistrationLists struct {
	slots, slats []string
	sigs         []Sig
	types        rs.TypeSpecs
}

func (reg *RegistrationLists) AddType(t *spec.TypeSpec) {
	switch t.Spec.Choice {
	case spec.UsesSpec_Group_Opt:
		// skip
	case spec.UsesSpec_Slot_Opt:
		reg.slots = append(reg.slots, t.Name)
	default:
		reg.slats = append(reg.slats, t.Name)

		// add signatures:
		switch t.Spec.Choice {
		case spec.UsesSpec_Flow_Opt:
			reg.addFlow(t)

		case spec.UsesSpec_Swap_Opt:
			reg.addSwap(t)
		}
	}
}

func (reg *RegistrationLists) addSwap(t *spec.TypeSpec) {
	swap := t.Spec.Value.(*spec.SwapSpec)
	commandName := pascal(t.Name)
	for _, pick := range swap.Between {
		sel := camelize(pick.Name)
		reg.sigs = append(reg.sigs, makeSig(t, commandName+" "+sel+":")...)
	}
}

func (reg *RegistrationLists) addFlow(t *spec.TypeSpec) {
	flow := t.Spec.Value.(*spec.FlowSpec)
	lede := flow.Name
	if len(lede) == 0 {
		lede = t.Name
	}
	sets := sigParts(flow, pascal(lede), reg.types)
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
		reg.sigs = append(reg.sigs, makeSig(t, sig)...)
	}
}

func (reg *RegistrationLists) Write(w io.Writer, tps *template.Template) (err error) {
	// sort registration lists ( in place )
	sort.Strings(reg.slots)
	sort.Strings(reg.slats)
	sort.Slice(reg.sigs, func(i, j int) bool {
		a, b := reg.sigs[i], reg.sigs[j]
		as := strings.Split(a.Sig, "=")
		bs := strings.Split(b.Sig, "=")
		return as[1] < bs[1] || (as[1] == bs[1] && (as[0] < bs[0]))
	})
	// write registration lists
	if e := tps.ExecuteTemplate(w, "regList.tmpl", map[string]any{
		"Name": "Slots",
		"List": reg.slots,
		"Type": "interface{}",
	}); e != nil {
		err = errutil.New(e, "couldnt process slots")
	} else if e := tps.ExecuteTemplate(w, "regList.tmpl", map[string]any{
		"Name": "Slats",
		"List": reg.slats,
		"Type": "composer.Composer",
	}); e != nil {
		err = errutil.New(e, "couldnt process slats")
	} else if e := tps.ExecuteTemplate(w, "sigList.tmpl", reg.sigs); e != nil {
		err = errutil.New(e, "couldnt process signatures")
	}
	return
}
