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
				Brace(js.Array, func(mui *js.Builder) {
					csv := OptionalComma{Builder: mui}
					// fake an initial term
					queue := []fieldDef{{
						term: spec.TermSpec{
							Label: lede,
						},
					}}
					for _, term := range terms {
						if term.Private {
							continue // skip private terms
						} else if fd, e := w.newFieldDef(term); e != nil {
							log.Fatalln(e) // exit if we couldnt create the field def
						} else if !fd.canCombine() {
							// write earlier items ( if any )
							if queue != nil {
								csv.Comma()
								flushQueue(out, queue)
								queue = nil
							}
							// write this term all on its own.
							csv.Comma()
							mui.Brace(js.Array, func(args *js.Builder) {
								if fd.writeField(args) {
									csv.Comma()
								}
								fd.writeInput(args)
							})
						} else {
							// allow this to join earlier terms
							queue = append(queue, fd)
							// but write all those terms and this one just added if we can't continue afterwards
							if !fd.canContinue() {
								csv.Comma()
								flushQueue(out, queue)
								queue = nil
							}
						}
					}
					// after all terms
					if queue != nil {
						csv.Comma()
						flushQueue(out, queue)
					}
				})
		})
}

// returns true if it wrote things
func (fd *fieldDef) writeField(args *js.Builder) (okay bool) {
	// write a leading label ( for non-anonymous fields )
	if label := fd.blocklyLabel(); len(label) > 0 {
		writeLabel(args, label)
		okay = true
	}
	// possibly write some other editable fields
	switch fd.termType() {
	case spec.UsesSpec_Swap_Opt:
		swap := fd.typeSpec.Spec.Value.(*spec.SwapSpec)
		if okay {
			args.R(js.Comma)
		}
		okay = true
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
		})

	case spec.UsesSpec_Num_Opt:
		if okay {
			args.R(js.Comma)
		}
		okay = true
		args.Kv("name", fd.name()).R(js.Comma)
		args.Brace(js.Obj, func(field *js.Builder) {
			field.Kv("type", bconst.FieldNumber)
		})

	case spec.UsesSpec_Str_Opt:
		// other options: spellcheck: true/false; text: the default value; open:; closed:
		if okay {
			args.R(js.Comma)
		}
		okay = true
		args.Brace(js.Obj, func(field *js.Builder) {
			field.Kv("name", fd.name()).R(js.Comma)
			if str := fd.typeSpec.Spec.Value.(*spec.StrSpec); !str.Exclusively {
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
		})
	}
	return
}

// write a blockly input to contain the passed field
func (fd *fieldDef) writeInput(args *js.Builder) {
	args.Brace(js.Obj, func(tail *js.Builder) {
		inputType, checks := fd.inputChecks()
		tail.Kv("name", fd.name()).R(js.Comma)
		tail.Kv("type", inputType)
		appendChecks(tail, "checks", checks)
		if shadow := fd.shadow(); len(shadow) > 0 {
			tail.R(js.Comma).Kv("shadow", shadow)
		}
		if fd.term.Optional {
			tail.R(js.Comma).Q("optional").R(js.Colon).S("true")
		}
		if fd.usesRepeatingInput() {
			tail.R(js.Comma).Q("repeats").R(js.Colon).S("true")
		}
	})

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
