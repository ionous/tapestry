package shape

import (
	"log"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// write the args0 and message0 key-values.
func (w *ShapeWriter) writeShapeDef(out *js.Builder, lede string, blockType *spec.TypeSpec, terms []spec.TermSpec) {
	out.WriteString(`"extensions":["tapestry_generic_mixin","tapestry_generic_extension"],`)
	hasMutator := blockType.Spec.Choice == spec.UsesSpec_Flow_Opt
	if hasMutator {
		out.Kv("mutator", "tapestry_generic_mutation").R(js.Comma)
	}
	out.Q("customData").R(js.Colon).
		Brace(js.Obj, func(custom *js.Builder) {
			if hasMutator {
				custom.Kv("mui", bconst.MutatorName(blockType.Name)).R(js.Comma)
			}
			custom.Q("shapeDef").R(js.Colon).
				Brace(js.Array, func(out *js.Builder) {
					// an initial item containing just the lede
					out.Brace(js.Obj, func(out *js.Builder) {
						out.Kv("label", lede)
					})
					// now any following terms as their own items
					for _, term := range terms {
						if term.Private {
							continue // skip private terms
						} else if fd, e := w.newFieldDef(term); e != nil {
							log.Fatalln(e) // exit if we couldnt create the field def
						} else if fn, ok := writeFn[fd.termType()]; !ok {
							log.Fatalln("unknown term type", fd.termType())
						} else if fn != nil {
							out.R(js.Comma)
							out.Brace(js.Obj, func(out *js.Builder) {
								// add a label for non-anonymous fields
								if label := fd.blocklyLabel(); len(label) > 0 {
									out.Kv("label", label).R(js.Comma)
								}
								out.Kv("name", fd.name()).R(js.Comma)
								fn(out, fd)
								if fd.term.Optional {
									out.R(js.Comma).Q("optional").R(js.Colon).S("true")
								}
								// if we are stack, we want to force a non-repeating input; one stack can already handle multiple blocks.
								// fix? we dont handle the case of a stack of one element; not sure that it exists in practice.
								if !fd.slot.Stack && fd.term.Repeats {
									out.R(js.Comma).Q("repeats").R(js.Colon).S("true")
								}
							})
						}
					}
				})
		})
}

func writeNum(out *js.Builder, fd *fieldDef) {
	out.Kv("type", bconst.FieldNumber)
}

func writeStr(out *js.Builder, fd *fieldDef) {
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

func writeFlow(out *js.Builder, fd *fieldDef) {
	writeInput(out, fd, bconst.InputValue, []string{fd.typeSpec.Name})
}

func writeSlot(out *js.Builder, fd *fieldDef) {
	// inputType might be a statement_input stack, or a single ( maybe repeatable ) input
	writeInput(out, fd, fd.slot.InputType(), []string{fd.slot.SlotType()})
}

// swaps are output as two items: one for the combo box editor, and one for the input plug.
// they are associated in mosaic by giving the blockly field and the blockly input the same name.
func writeSwap(out *js.Builder, fd *fieldDef) {
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
	out.R(js.Obj[1]).R(js.Comma).R(js.Obj[0])
	writeInput(out, fd, bconst.InputValue, checks)
	return
}

// write an input for the fields that need it.
func writeInput(out *js.Builder, fd *fieldDef, inputType string, checks []string) {
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
