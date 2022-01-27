package blocks

import (
	"log"
	"strconv"

	"git.sr.ht/~ionous/tapestry/dl/spec"
)

// write the args0 and message0 for a mutator ui block
func writeMuiMsgArgs(out *Js, blockType *spec.TypeSpec, flow *spec.FlowSpec) {
	var argCount int
	out.
		Q("args0").R(colon).Brace(array, func(args *Js) {
		var header string
		if lede := flow.Name; len(lede) > 0 {
			header = lede
		} else {
			header = blockType.Name
		}
		argCount += writeDummy(args, "", header)
		for _, term := range flow.Terms {
			if !term.Private {
				argCount += writeMuiInput(args, term)
			}
		}
	}).
		R(comma).
		Q("message0").R(colon).Brace(quotes, func(msg *Js) {
		// ex. "message0": "%1%2%3%4%5%6%7%8"
		for i := 1; i <= argCount; i++ {
			out.R(percent).S(strconv.Itoa(i))
		}
	})
}

func writeMuiInput(args *Js, term spec.TermSpec) (ret int) {
	typeName := term.TypeName() // lookup spec
	if termType, ok := lookup[typeName]; !ok {
		log.Fatalln("missing named type", typeName)
	} else {
		ret = writeMuiTerm(args, term, termType)
	}
	return
}

// note: writes a leading comma :/
func writeMuiTerm(args *Js, term spec.TermSpec, termType *spec.TypeSpec) (ret int) {
	name, label := term.Field(), term.Label()
	// stacked elements dont need to repeat inputs: one input allows multiple blocks.
	// ( and if they are optional, we'll want to use a checkbox )
	stacks, _ := SlotStacks(termType)
	if term.Repeats && len(stacks) == 0 {
		// for the mui: insist on having a label.
		if len(label) == 0 {
			if len(term.Name) > 0 {
				label = term.Name
			} else {
				label = term.Type
			}
		}
		args.R(comma)
		ret += writeDummy(args, name, label,
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
					// unique name needed for blockly undo
					Q("name").R(colon).Brace(quotes, func(val *Js) {
					val.S(term.Field()).R(score).S("edit")
				})
			})

	} else if term.Optional {
		args.R(comma)
		ret += writeDummy(args, name, label,
			func(field *Js) {
				field.
					Kv("type", FieldCheckbox).R(comma).
					// unique name needed for blockly undo:
					Q("name").R(colon).Brace(quotes, func(val *Js) {
					val.S(term.Field()).R(score).S("edit")
				})
			})
	}
	return
}
