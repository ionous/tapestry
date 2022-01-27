package blocks

import (
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/spec"
)

// helper to write a workspace block.
// a clone of muiWriter; not sure what to share...
type blockWriter struct {
	blockType *spec.TypeSpec
	flow      *spec.FlowSpec
	argCount  int
	muiData   []muiData
}

// when multiple is false, the mui generates a checkbox for the field, when true: a number.
// note: stacked repeat elements ( statement input ) have multiple false.
type muiData struct {
	term     spec.TermSpec
	termType *spec.TypeSpec
	multiple bool
}

// write the args0 and message0 key-values.
func writeBlockInternals(out *Js, blockType *spec.TypeSpec, flow *spec.FlowSpec) {
	w := blockWriter{blockType: blockType, flow: flow}
	out.
		Q("args0").R(colon).
		Brace(array, w.writeArgs).R(comma).
		Q("message0").R(colon).
		Brace(quotes, func(msg *Js) { writeMsg(msg, w.argCount) }).
		If(len(w.muiData) > 0, w.writeMuiData)
}

func (w *blockWriter) writeMuiData(out *Js) {
	out.
		R(comma).
		Kv("mutator", "tapestry_generic_mutation").R(comma).
		Q("customData").R(colon).
		Brace(obj, func(custom *Js) {
			custom.Q("muiData").R(colon).
				Brace(obj, func(mui *Js) {
					for i, cd := range w.muiData {
						if i > 0 {
							mui.R(comma)
						}
						mui.Q(strings.ToUpper(cd.term.Field())).R(colon).Brace(array, func(opt *Js) {
							writeTerm(opt, cd.term, cd.termType)
						})
					}
				})
		})
}

//
func (w *blockWriter) writeArgs(args *Js) {
	// write our
	label := w.flow.Name
	if len(label) == 0 {
		label = w.blockType.Name
	}
	w.argCount += writeInput(args, Input{Label: label})
	// write any mutating terms:
	for _, term := range w.flow.Terms {
		if term.Private {
			continue // skip private terms
		}
		//
		typeName := term.TypeName() // lookup spec
		if termType, ok := lookup[typeName]; !ok {
			log.Fatalln("missing named type", typeName)
		} else {
			// look at all of the slots the term type implements
			stacks, _ := SlotStacks(termType)
			if term.Optional && (!term.Repeats || len(stacks) > 0) {
				w.muiData = append(w.muiData, muiData{term, termType, true})
			} else if term.Repeats {
				w.muiData = append(w.muiData, muiData{term, termType, false})
			}
			// write the field if we need it.
			if !term.Optional {
				args.R(comma)
				// ugly, but simplifies counting of repeat fields in tapestry_generic_mutation
				if term.Repeats {
					term.Name = term.Field() + "0"
				}
				w.argCount += writeTerm(args, term, termType)
			}
		}
	}
}

func writeTerm(out *Js, term spec.TermSpec, termType *spec.TypeSpec) (ret int) {
	name := term.Field()
	label := term.Label()
	// ideally then, we'd be able to refactor to write a single term.
	switch kind := termType.Spec.Choice; kind {

	case spec.UsesSpec_Flow_Opt:
		// a flow goes here: tbd: but probably a shadow
		ret += writeInput(out, Input{
			Name:  name,
			Check: termType.Name,
			Type:  InputValue,
			Label: label,
		})

	case spec.UsesSpec_Slot_Opt:
		slot := slotRules.FindSlot(termType.Name)
		ret += writeInput(out, Input{
			Name:  name,
			Check: slot.SlotType(),
			Type:  slot.InputType(),
			Label: label,
		})

	case spec.UsesSpec_Swap_Opt:
		pairs := []string{"first", "ITEM"}
		ret += writeInput(out, Input{
			Name:  name,
			Label: label,
		},
			func(field *Js) {
				field.Kv("type", FieldDropdown).R(comma).
					Q("option").R(colon).Brace(array,
					func(options *Js) {
						for i, cnt := 0, len(pairs); i < cnt; i += 2 {
							if i > 0 {
								options.R(comma)
							}
							options.Brace(array, func(opt *Js) {
								// "first item", "ITEM1"
								opt.Q(pairs[i]).R(comma).Q(pairs[i+1])
							})
						}
					})
			})

	case spec.UsesSpec_Num_Opt:
		ret += writeInput(out, Input{Name: name, Label: label},
			func(field *Js) {
				field.Kv("type", FieldNumber)
			})

	case spec.UsesSpec_Str_Opt:
		ret += writeInput(out, Input{Name: name, Label: label},
			func(field *Js) {
				field.Kv("type", FieldText)
				// other options:
				// spellcheck: true/false
				// text: the default value
			})

	default:
		log.Fatalln("unknown spec type", kind)
	}
	return
}
