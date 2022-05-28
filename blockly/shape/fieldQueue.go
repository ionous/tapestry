package shape

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
)

type OptionalComma struct {
	*js.Builder
	c int
}

func (p *OptionalComma) Comma() {
	if (*p).c = (*p).c + 1; (*p).c > 1 {
		p.R(js.Comma)
	}
}

// everything in this list of fields goes into a single input
// the type of the input is dependent on the last field
// assumes that there's at least one thing in the list
func flushQueue(args *js.Builder, queue []fieldDef) {
	args.Brace(js.Array, func(out *js.Builder) {
		wrote := false
		for _, fd := range queue {
			if wrote {
				out.R(js.Comma)
				wrote = false
			}
			if fd.writeField(out) {
				wrote = true
			}
		}
		// all the terms "collapse" into the last input
		if wrote {
			out.R(js.Comma)
		}
		queue[len(queue)-1].writeInput(out)
	})
}

type fieldDef struct {
	term     spec.TermSpec
	typeSpec *spec.TypeSpec
	slot     bconst.SlotRule
}

func (w *ShapeWriter) newFieldDef(term spec.TermSpec) (ret fieldDef, err error) {
	typeName := term.TypeName() // lookup spec
	if typeSpec, ok := w.Types[typeName]; !ok {
		err = errutil.New("missing named type", typeName)
	} else {
		var slot bconst.SlotRule
		if typeSpec.Spec.Choice == spec.UsesSpec_Slot_Opt {
			// inputType might be a statement_input stack, or a single ( maybe repeatable ) input
			// regardless, it only has the input, no special fields.
			slot = bconst.FindSlotRule(typeSpec.Name)
		}
		ret = fieldDef{term, typeSpec, slot}
	}
	return
}

func (fd *fieldDef) name() string {
	return strings.ToUpper(fd.term.Field())
}

func (fd *fieldDef) usesRepeatingInput() bool {
	// if we are stack, we want to force a non-repeating input; one stack can already handle multiple blocks.
	// fix? we dont handle the case of a stack of one element; not sure that it exists in practice.
	return !fd.slot.Stack && fd.term.Repeats
}

// can this appear on a line with other fields?
func (fd *fieldDef) canCombine() (okay bool) {
	return !fd.term.Optional && !fd.usesRepeatingInput()
}

// can any other fields follow this one?
// ( assumes canCombine would return true )
func (fd *fieldDef) canContinue() (okay bool) {
	switch fd.termType() {
	case spec.UsesSpec_Num_Opt, spec.UsesSpec_Str_Opt:
		okay = true
	}
	return
}

// will we need a label for required anonymous terms?
// maybe at least for the mui?
func (fd *fieldDef) blocklyLabel() (ret string) {
	if !fd.term.IsAnonymous() {
		ret = fd.term.Label
	} /*else {
		ret = fd.term.Name
	}*/
	return
}

func (fd *fieldDef) shadow() (ret string) {
	switch fd.termType() {
	case spec.UsesSpec_Num_Opt, spec.UsesSpec_Flow_Opt:
		ret = fd.typeSpec.Name
	}
	return
}

// handle our fake field for the leading label.
func (fd *fieldDef) termType() (ret string) {
	if fd.typeSpec != nil {
		ret = fd.typeSpec.Spec.Choice
	}
	return
}

func (fd *fieldDef) inputChecks() (inputType string, checks []string) {
	switch fd.termType() {
	default:
		inputType = bconst.InputDummy // provisionally

	case spec.UsesSpec_Flow_Opt:
		inputType = bconst.InputValue
		checks = []string{fd.typeSpec.Name}

	case spec.UsesSpec_Slot_Opt:
		// inputType might be a statement_input stack, or a single ( maybe repeatable ) input
		// regardless, it only has the input, no special fields.
		inputType = fd.slot.InputType()
		checks = []string{fd.slot.SlotType()}

	case spec.UsesSpec_Swap_Opt:
		inputType = bconst.InputValue
		swap := fd.typeSpec.Spec.Value.(*spec.SwapSpec)
		// allows all the types and changes the swap depending on what gets connected
		for _, pick := range swap.Between {
			checks = append(checks, pick.TypeName())
		}
	}
	return
}
