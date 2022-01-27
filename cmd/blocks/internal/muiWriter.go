package blocks

import (
	"git.sr.ht/~ionous/tapestry/dl/spec"
)

type muiWriter struct {
	blockType *spec.TypeSpec
	flow      *spec.FlowSpec
	argCount  int
}

// write the args0 and message0 key-values.
func writeMuiArgs(out *Js, blockType *spec.TypeSpec, flow *spec.FlowSpec) bool {
	w := muiWriter{blockType: blockType, flow: flow}
	out.
		Q("args0").R(colon).
		Brace(array, w.writeArgs).R(comma).
		Q("message0").R(colon).
		Brace(quotes, func(msg *Js) { writeMsg(out, w.argCount) })
	return w.argCount > 0
}

func (w *muiWriter) writeArgs(args *Js) {
	// write a header for the mutation block
	w.argCount += writeInput(args, Input{Label: w.blockType.Name})
	// write any mutating terms:
	for _, term := range w.flow.Terms {
		if !term.Private {
			name, label := term.Field(), term.Label()
			if term.Repeats {
				// for the mui: insist on having a label.
				if len(label) == 0 {
					if len(term.Name) > 0 {
						label = term.Name
					} else {
						label = term.Type
					}
				}
				args.R(comma)
				w.argCount += writeInput(args, Input{Name: name, Label: label},
					func(field *Js) {
						const (
							zero = "0"
							one  = "1"
						)
						var min = one
						if term.Optional {
							min = zero
						}
						field.
							Kv("type", FieldNumber).R(comma).
							Q("min").R(colon).S(min).R(comma).
							Q("precision").R(colon).S(one).R(comma).
							Q("name").R(colon).Brace(quotes,
							func(val *Js) {
								// unique name needed for blockly undo
								val.S(term.Field()).R(score).S("edit")
							})
					})

			} else if term.Optional {
				args.R(comma)
				w.argCount += writeInput(args, Input{Name: name, Label: label},
					func(field *Js) {
						field.
							Kv("type", FieldCheckbox).R(comma).
							Q("name").R(colon).Brace(quotes,
							func(val *Js) {
								// unique name needed for blockly undo
								val.S(term.Field()).R(score).S("edit")
							})
					})
			}
		}
	}
}
