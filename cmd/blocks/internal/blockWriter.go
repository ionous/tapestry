package blocks

import (
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// write the args0 and message0 key-values.
func writeCustomData(out *js.Builder, blockType *spec.TypeSpec, flow *spec.FlowSpec) {
	out.WriteString(`"mutator": "tapestry_generic_mutation",` +
		`"extensions":["tapestry_mutation_mixin","tapestry_mutation_extension"],`)
	out.Q("customData").R(js.Colon).
		Brace(js.Obj, func(custom *js.Builder) {
			custom.Q("muiData").R(js.Colon).
				Brace(js.Array, func(mui *js.Builder) {
					var csv int
					for _, term := range flow.Terms {
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
	name, label := term.Field(), term.Label()
	// write the label for this term.
	writeLabel(args, label)
	// write other fields while collecting information for the trailing input:
	var checks []string
	var inputType = InputDummy
	//
	switch kind := termType.Spec.Choice; kind {
	case spec.UsesSpec_Flow_Opt:
		// a flow goes here: tbd: but probably a shadow
		// it only has the input, no special fields
		inputType, checks = InputValue, []string{termType.Name}

	case spec.UsesSpec_Slot_Opt:
		// inputType might be a statement_input stack, or a single ( maybe repeatable ) input
		// regardless, it only has the input, no special fields.
		slot := slotRules.FindSlot(termType.Name)
		inputType, checks = slot.InputType(), []string{slot.SlotType()}
		// if we are stack, we want to force a non-repeating input; one stack can already handle multiple blocks.
		// fix? we dont handle the case of a stack of one element; not sure that it exists in practice.
		if slot.Stack {
			term.Repeats = false
		}

	case spec.UsesSpec_Swap_Opt:
		swap := termType.Spec.Value.(*spec.SwapSpec)
		inputType = InputValue
		for _, pick := range swap.Between {
			checks = append(checks, pick.TypeName())
		}
		args.R(js.Comma).
			Brace(js.Obj, func(field *js.Builder) {
				field.
					Kv("name", name).R(js.Comma). // for blockly serialization
					Kv("type", FieldDropdown).R(js.Comma).
					Q("option").R(js.Colon).Brace(js.Array,
					func(options *js.Builder) {
						for i, pick := range swap.Between {
							if i > 0 {
								options.R(js.Comma)
							}
							options.Brace(js.Array, func(opt *js.Builder) {
								opt.Kv(pick.FriendlyName(), pick.TypeName())
							})
						}
					})
			})

	case spec.UsesSpec_Num_Opt:
		args.R(js.Comma).
			Brace(js.Obj, func(field *js.Builder) {
				field.
					Kv("name", name).R(js.Comma). // for blockly serialization
					Kv("type", FieldNumber)
			})

	case spec.UsesSpec_Str_Opt:
		// FIX - write combo box for enums
		// fix: future? a combo box with custom entry
		// ( ex. something like a variable that is shared globally *if* variable categories are allowed. )
		// ( or. a combo box with an "other" entry, or a mui option -- to change from selected to free typing )
		args.R(js.Comma).
			Brace(js.Obj, func(field *js.Builder) {
				field.
					Kv("name", name).R(js.Comma). // for blockly serialization
					Kv("type", FieldText)
				// other options:
				// spellcheck: true/false
				// text: the default value
			})

	default:
		log.Fatalln("unknown spec type", kind)
	}
	// write the input all of the above fields are a part of:
	args.R(js.Comma).
		Brace(js.Obj, func(tail *js.Builder) {
			tail.Kv("name", strings.ToUpper(term.Field())).R(js.Comma)
			tail.Kv("type", inputType)
			if len(checks) > 0 {
				tail.R(js.Comma).
					Q("check").R(js.Colon).Brace(js.Array, func(check *js.Builder) {
					for i, c := range checks {
						if i > 0 {
							check.R(js.Comma)
						}
						check.Q(c)
					}
				})
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
