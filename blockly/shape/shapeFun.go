package shape

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
)

func fieldWriter(t typeinfo.T) func(*ShapeWriter, *js.Builder, typeinfo.Term) {
	switch t.(type) {
	case *typeinfo.Flow:
		return writeFlowField
	case *typeinfo.Slot:
		return writeSlotField
	case *typeinfo.Num:
		return writeNumField
	case *typeinfo.Str:
		return writeStrField
	}
	return nil
}

//func writeLabelField(out *js.Builder, term typeinfo.Term, typeSpec typeinfo.T) {
//out.Kv("label", term.Label)
//}

func writeNumField(w *ShapeWriter, out *js.Builder, term typeinfo.Term) {
	// termType := term.Type.(*typeinfo.NUm)
	out.Kv("type", bconst.FieldNumber)
}

func writeStrField(w *ShapeWriter, out *js.Builder, term typeinfo.Term) {
	// other options possible: spellcheck: true/false; text: the default value.
	if str := term.Type.(*typeinfo.Str); len(str.Options) == 0 {
		fieldType := bconst.MosaicTextField
		if str.Name == "lines" {
			fieldType = bconst.MosaicMultilineField
		}
		out.Kv("type", fieldType)
		// every field needs to be labeled
		var placeholder string
		if !term.IsAnonymous() {
			placeholder = term.Label
		} else {
			placeholder = term.Name
		}
		if len(placeholder) > 0 {
			out.R(js.Comma).Kv("text", strings.Replace(placeholder, "_", " ", -1))
		}
	} else {
		fieldType := bconst.FieldDropdown
		// all enums are limited to their string list now
		// if !str.Exclusively {
		// 	fieldType = bconst.MosaicStrField
		// }
		out.Kv("type", fieldType).R(js.Comma).
			Q("options").R(js.Colon).
			Brace(js.Array,
				func(options *js.Builder) {
					for i, pick := range str.Options {
						if i > 0 {
							options.R(js.Comma)
						}
						options.Brace(js.Array, func(opt *js.Builder) {
							label := strings.Replace(pick, "_", " ", -1)
							key := bconst.KeyName(pick)
							opt.Q(label).R(js.Comma).Q(key)
						})
					}
				})
	}
	return
}

func writeFlowField(w *ShapeWriter, out *js.Builder, term typeinfo.Term) {
	termType := term.Type.(*typeinfo.Flow)
	writeInput(out, bconst.InputValue, "", []string{termType.Name})
}

func writeSlotField(w *ShapeWriter, out *js.Builder, term typeinfo.Term) {
	// inputType might be a statement_input stack, or a single ( maybe repeatable ) input
	termType := term.Type.(*typeinfo.Slot)
	slot := bconst.MakeSlotRule(termType)
	writeInput(out, slot.InputType(), "", []string{slot.SlotType()})
}

// write an input for the fields that need it.
func writeInput(out *js.Builder, inputType, shadow string, checks []string) {
	out.Kv("type", inputType)
	appendChecks(out, "checks", checks)
	if len(shadow) > 0 {
		out.R(js.Comma).Kv("shadow", shadow)
	}
}

func getShadow(t typeinfo.T) (ret string) {
	switch t.(type) {
	case *typeinfo.Num, *typeinfo.Flow:
		ret = t.TypeName()
	}
	return
}
