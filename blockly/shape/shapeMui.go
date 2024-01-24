package shape

import (
	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// write the args0 and message0 for a mutator ui block
func (w *ShapeWriter) writeMuiMsgArgs(out *js.Builder, flow *typeinfo.Flow) {
	var argCount int
	out.
		Q("args0").R(js.Colon).Brace(js.Array, func(args *js.Builder) {
		header := typeinfo.FriendlyName(flow.Lede)
		argCount += writeDummy(args, "", header)
		for _, term := range flow.Terms {
			if !term.Private {
				argCount += w.writeMuiInput(args, term)
			}
		}
	}).
		R(js.Comma).
		Q("message0").R(js.Colon).Brace(js.Quotes, func(msg *js.Builder) {
		// ex. "message0": "%1%2%3%4%5%6%7%8"
		for i := 1; i <= argCount; i++ {
			out.R(js.Percent).N(i)
		}
	})
}

func (w *ShapeWriter) writeMuiInput(args *js.Builder, term typeinfo.Term) (ret int) {
	return w.writeMuiTerm(args, term, term.Type)
}

// note: writes a leading comma :/
func (w *ShapeWriter) writeMuiTerm(args *js.Builder, term typeinfo.Term, termType typeinfo.T) (ret int) {
	label, name := term.Label, term.Name
	// stacked elements dont need to repeat inputs: one input allows multiple blocks.
	// ( and if they are optional, we'll want to use a checkbox )
	stacks, _ := slotStacks(w.Types, termType)
	if term.Repeats && len(stacks) == 0 {
		// for the mui: insist on having a label.
		if len(label) == 0 {
			label = term.Name
		}
		args.R(js.Comma)
		ret += writeDummy(args, name, label,
			func(field *js.Builder) {
				const (
					zero = "0"
					one  = "1"
				)
				var min = one
				if term.Optional {
					min = zero
				}
				field.
					Kv("type", bconst.FieldNumber).R(js.Comma).
					Q("min").R(js.Colon).Raw(min).R(js.Comma).
					Q("precision").R(js.Colon).Raw(one).R(js.Comma).
					// unique name needed for blockly undo
					Q("name").R(js.Colon).Brace(js.Quotes,
					func(val *js.Builder) {
						val.Str(name).R(js.Score).Raw("edit")
					})
			})

	} else if term.Optional {
		args.R(js.Comma)
		ret += writeDummy(args, name, label,
			func(field *js.Builder) {
				field.
					Kv("type", bconst.FieldCheckbox).R(js.Comma).
					// unique name needed for blockly undo:
					Q("name").R(js.Colon).Brace(js.Quotes,
					func(val *js.Builder) {
						val.Str(name).R(js.Score).Raw("edit")
					})
			})
	}
	return
}
