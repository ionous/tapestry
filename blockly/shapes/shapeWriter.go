package shapes

import (
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// write the args0 and message0 key-values.
func writeShapeDef(out *js.Builder, blockType *spec.TypeSpec, terms []spec.TermSpec) {
	out.WriteString(`"extensions":["tapestry_generic_mixin","tapestry_generic_extension"],`)
	if blockType.Spec.Choice == spec.UsesSpec_Flow_Opt {
		out.Kv("mutator", "tapestry_generic_mutation").R(js.Comma)
	}
	out.Q("customData").R(js.Colon).
		Brace(js.Obj, func(custom *js.Builder) {
			custom.Kv("mui", bconst.MutatorName(blockType.Name)).R(js.Comma)
			custom.Q("shapeDef").R(js.Colon).
				Brace(js.Array, func(mui *js.Builder) {
					var csv int
					for _, term := range terms {
						if term.Private {
							continue // skip private terms
						}
						if csv = csv + 1; csv > 1 {
							mui.R(js.Comma)
						}
						mui.Brace(js.Array, func(args *js.Builder) {
							writeFieldDefs(args, term)
						})
					}
				})
		})
}

//
func writeFieldDefs(args *js.Builder, term spec.TermSpec) {
	typeName := term.TypeName() // lookup spec
	if termType, ok := lookup[typeName]; !ok {
		log.Fatalln("missing named type", typeName)
	} else {
		writeTerm(args, term, termType)
	}
}

func writeTerm(args *js.Builder, term spec.TermSpec, termType *spec.TypeSpec) {
	// write the label for this term.
	if l := term.FriendlyName(); len(l) > 0 {
		writeLabel(args, l)
		args.R(js.Comma)
	}
	// write other fields while collecting information for the trailing input:
	var checks []string
	var inputType = bconst.InputDummy
	var shadow string
	//
	switch kind := termType.Spec.Choice; kind {

	case spec.UsesSpec_Flow_Opt:
		// a flow goes here: tbd: but probably a shadow
		// it only has the input, no special fields
		inputType, checks = bconst.InputValue, []string{termType.Name}
		shadow = termType.Name

	case spec.UsesSpec_Slot_Opt:
		// inputType might be a statement_input stack, or a single ( maybe repeatable ) input
		// regardless, it only has the input, no special fields.
		slot := bconst.FindSlotRule(termType.Name)
		inputType, checks = slot.InputType(), []string{slot.SlotType()}
		// if we are stack, we want to force a non-repeating input; one stack can already handle multiple blocks.
		// fix? we dont handle the case of a stack of one element; not sure that it exists in practice.
		if slot.Stack {
			term.Repeats = false
		}

	case spec.UsesSpec_Swap_Opt:
		swap := termType.Spec.Value.(*spec.SwapSpec)
		inputType = bconst.InputValue
		// allows all the types and changes the swap depending on what gets connected
		for _, pick := range swap.Between {
			checks = append(checks, pick.TypeName())
		}
		args.Brace(js.Obj, func(field *js.Builder) {
			field.
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
		}).R(js.Comma)

	case spec.UsesSpec_Num_Opt:
		args.Brace(js.Obj, func(field *js.Builder) {
			field.Kv("type", bconst.FieldNumber)
		}).R(js.Comma)

	case spec.UsesSpec_Str_Opt:
		// other options:
		// spellcheck: true/false
		// text: the default value
		str := termType.Spec.Value.(*spec.StrSpec)
		// open:
		// closed:
		args.Brace(js.Obj, func(field *js.Builder) {
			if !str.Exclusively {
				field.Kv("type", bconst.FieldText)
			} else {
				field.Kv("type", bconst.FieldDropdown).R(js.Comma).
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
		}).R(js.Comma)

	default:
		log.Fatalln("unknown spec type", kind)
	}
	// write the input that all of the above fields are a part of:
	args.Brace(js.Obj, func(tail *js.Builder) {
		tail.Kv("name", strings.ToUpper(term.Field())).R(js.Comma)
		tail.Kv("type", inputType)
		appendChecks(tail, "checks", checks)

		if len(shadow) > 0 {
			tail.R(js.Comma).Kv("shadow", shadow)
		}
		if term.Optional {
			tail.R(js.Comma).Q("optional").R(js.Colon).S("true")
		}
		if term.Repeats {
			tail.R(js.Comma).Q("repeats").R(js.Colon).S("true")
		}
	})
	return
}

func appendChecks(tail *js.Builder, label string, checks []string) {
	if len(checks) > 0 {
		tail.R(js.Comma).
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
