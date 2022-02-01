package btypes

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// write a number of fields, followed by the input that they merge down into
// using blockly json's message interpolation syntax.
func writeDummy(args *js.Builder, name, label string, fields ...func(*js.Builder)) (ret int) {
	// write any leading text as a field
	if n := label; len(n) > 0 {
		writeLabel(args, n)
		ret++
	}
	// write the explicit fields
	if len(fields) > 0 {
		args.R(js.Comma)
		ret += writeFields(args, fields...)
	}
	// write the input the fields are a part of
	args.R(js.Comma).
		Brace(js.Obj, func(tail *js.Builder) {
			if n := name; len(n) > 0 {
				tail.Kv("name", strings.ToUpper(n)).R(js.Comma)
			}
			tail.Kv("type", bconst.InputDummy)
			ret++
		})
	return
}

func writeFields(out *js.Builder, fields ...func(*js.Builder)) (ret int) {
	for i, field := range fields {
		if i > 0 {
			out.R(js.Comma)
		}
		out.Brace(js.Obj, field)
		ret++
	}
	return
}

func writeLabel(out *js.Builder, n string) *js.Builder {
	return out.Brace(js.Obj, func(lab *js.Builder) {
		lab.
			Kv("type", bconst.FieldLabel).R(js.Comma).
			Kv("text", n)
	})
}
