package shape

import (
	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/web/js"
)

var writeFn = map[string]func(*js.Builder, *shapeField){
	spec.UsesSpec_Flow_Opt:  writeFlow,
	spec.UsesSpec_Slot_Opt:  writeSlot,
	spec.UsesSpec_Swap_Opt:  writeSwap,
	spec.UsesSpec_Num_Opt:   writeNum,
	spec.UsesSpec_Str_Opt:   writeStr,
	spec.UsesSpec_Group_Opt: nil,
}

func writeNum(out *js.Builder, fd *shapeField) {
	out.Kv("type", bconst.FieldNumber)
}

func writeStr(out *js.Builder, fd *shapeField) {
	// other options possible: spellcheck: true/false; text: the default value.
	if str := fd.typeSpec.Spec.Value.(*spec.StrSpec); !str.Exclusively {
		out.Kv("type", bconst.FieldText)
	} else {
		out.Kv("type", bconst.FieldDropdown).R(js.Comma).
			Q("options").R(js.Colon).
			Brace(js.Array,
				func(options *js.Builder) {
					for i, pick := range str.Uses {
						if i > 0 {
							options.R(js.Comma)
						}
						options.Brace(js.Array, func(opt *js.Builder) {
							opt.Q(pick.FriendlyName()).R(js.Comma).Q(pick.Value())
						})
					}
				})
	}
	return
}

func writeFlow(out *js.Builder, fd *shapeField) {
	writeInput(out, fd, bconst.InputValue, []string{fd.typeSpec.Name})
}

func writeSlot(out *js.Builder, fd *shapeField) {
	// inputType might be a statement_input stack, or a single ( maybe repeatable ) input
	writeInput(out, fd, fd.slot.InputType(), []string{fd.slot.SlotType()})
}

// swaps are output as two items: one for the combo box editor, and one for the input plug.
// they are associated in mosaic by giving the blockly field and the blockly input the same name.
func writeSwap(out *js.Builder, fd *shapeField) {
	swap := fd.typeSpec.Spec.Value.(*spec.SwapSpec)
	out.
		Kv("name", fd.name()).R(js.Comma).
		Kv("type", bconst.FieldDropdown).R(js.Comma).
		// write the blockly dropdown data
		Q("options").R(js.Colon).
		Brace(js.Array,
			func(options *js.Builder) {
				for i, pick := range swap.Between {
					if i > 0 {
						options.R(js.Comma)
					}
					options.Brace(js.Array, func(opt *js.Builder) {
						opt.Q(pick.FriendlyName()).R(js.Comma).Q(pick.Value())
					})
				}
			}).R(js.Comma).
		// write info on how to change the input when swapping types
		Q("swaps").R(js.Colon).
		Brace(js.Obj,
			func(swaps *js.Builder) {
				for i, pick := range swap.Between {
					if i > 0 {
						swaps.R(js.Comma)
					}
					swaps.Kv(pick.Value(), pick.TypeName())
				}
			})

	// allows all the types and changes the swap depending on what gets connected
	var checks []string
	for _, pick := range swap.Between {
		checks = append(checks, pick.TypeName())
	}
	// close off the dropdown item def, and open an input item def
	// the caller will close that out.
	out.R(js.Obj[1]).R(js.Comma).R(js.Obj[0])
	out.Kv("name", fd.name()).R(js.Comma)
	writeInput(out, fd, bconst.InputValue, checks)
	return
}

// write an input for the fields that need it.
func writeInput(out *js.Builder, fd *shapeField, inputType string, checks []string) {
	out.Kv("type", inputType)
	appendChecks(out, "checks", checks)
	if shadow := fd.shadow(); len(shadow) > 0 {
		out.R(js.Comma).Kv("shadow", shadow)
	}
}

func appendChecks(out *js.Builder, label string, checks []string) {
	if len(checks) > 0 {
		out.R(js.Comma).
			Q(label).R(js.Colon).Brace(js.Array, func(check *js.Builder) {
			for i, c := range checks {
				if i > 0 {
					check.R(js.Comma)
				}
				check.Q(c)
			}
		})
	}
}
