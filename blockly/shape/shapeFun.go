package shape

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/web/js"
)

var fieldWriter = map[string]func(*js.Builder, spec.TermSpec, *spec.TypeSpec){
	spec.UsesSpec_Flow_Opt: writeFlowField,
	spec.UsesSpec_Slot_Opt: writeSlotField,
	spec.UsesSpec_Swap_Opt: writeSwapField,
	spec.UsesSpec_Num_Opt:  writeNumField,
	spec.UsesSpec_Str_Opt:  writeStrField,
	// fields dont use groups, so we COULD highjack it for labels
	//spec.UsesSpec_Group_Opt: writeLabelField,
}

//func writeLabelField(out *js.Builder, term spec.TermSpec, typeSpec *spec.TypeSpec) {
//out.Kv("label", term.Label)
//}

func writeNumField(out *js.Builder, term spec.TermSpec, typeSpec *spec.TypeSpec) {
	out.Kv("type", bconst.FieldNumber)
}

func writeStrField(out *js.Builder, term spec.TermSpec, typeSpec *spec.TypeSpec) {
	// other options possible: spellcheck: true/false; text: the default value.
	if str := typeSpec.Spec.Value.(*spec.StrSpec); len(str.Uses) == 0 {
		out.Kv("type", bconst.FieldText)
	} else {
		fieldType := bconst.FieldDropdown
		if !str.Exclusively {
			fieldType = bconst.CustomFieldDropdown
		}
		out.Kv("type", fieldType).R(js.Comma).
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

func writeFlowField(out *js.Builder, term spec.TermSpec, typeSpec *spec.TypeSpec) {
	writeInput(out, bconst.InputValue, "", []string{typeSpec.Name})
}

func writeSlotField(out *js.Builder, term spec.TermSpec, typeSpec *spec.TypeSpec) {
	// inputType might be a statement_input stack, or a single ( maybe repeatable ) input
	slot := bconst.FindSlotRule(typeSpec.Name)
	writeInput(out, slot.InputType(), "", []string{slot.SlotType()})
}

// swaps are output as two items: one for the combo box editor, and one for the input plug.
// they are associated in mosaic by giving the blockly field and the blockly input the same name.
func writeSwapField(out *js.Builder, term spec.TermSpec, typeSpec *spec.TypeSpec) {
	swap := typeSpec.Spec.Value.(*spec.SwapSpec)
	out.
		Kv("type", bconst.FieldDropdown).R(js.Comma).
		// write the blockly dropdown data: a displayed label and a resulting key value.
		// the idea is that different translations could display different labels for the same value.
		// ex: [ "flow", "$FLOW" ]
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
		// write info on how to change the block with different selections.
		// ties the key to the specific type:  ex. "$FLOW": "flow_spec".
		// while each key can only appear once, different keys can theoretically map to a single type.
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
	// close off the dropdown item def, and open an input item def; the caller will close it out.
	//
	// fix: not every swap needs to generate an input: a swap of strings for instance.
	// it'd probably make more sense for the blockly side to determine the need based on the swap.
	// that'd save having to generate standalone blocks for open strs and nums.
	out.R(js.Obj[1]).R(js.Comma).R(js.Obj[0])

	fieldName := strings.ToUpper(term.Field())

	out.Kv("name", fieldName).R(js.Comma)
	writeInput(out, bconst.InputValue, getShadow(typeSpec), checks)
	return
}

// write an input for the fields that need it.
func writeInput(out *js.Builder, inputType, shadow string, checks []string) {
	out.Kv("type", inputType)
	appendChecks(out, "checks", checks)
	if len(shadow) > 0 {
		out.R(js.Comma).Kv("shadow", shadow)
	}
}

func getShadow(typeSpec *spec.TypeSpec) (ret string) {
	switch typeSpec.Spec.Choice {
	case spec.UsesSpec_Num_Opt, spec.UsesSpec_Flow_Opt:
		ret = typeSpec.Name
	}
	return
}
