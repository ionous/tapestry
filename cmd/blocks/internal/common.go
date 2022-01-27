package blocks

import (
	"strings"
)

// write a number of fields, followed by the input that they merge down into
// using blockly json's message interpolation syntax.
func writeDummy(args *Js, name, label string, fields ...func(*Js)) (ret int) {
	// write any leading text as a field
	if n := label; len(n) > 0 {
		writeLabel(args, n)
		ret++
	}
	// write the explicit fields
	if len(fields) > 0 {
		args.R(comma)
		ret += writeFields(args, fields...)
	}
	// write the input the fields are a part of
	args.R(comma).
		Brace(obj, func(tail *Js) {
			if n := name; len(n) > 0 {
				tail.Kv("name", strings.ToUpper(n)).R(comma)
			}
			tail.Kv("type", InputDummy)
			ret++
		})
	return
}

func writeFields(out *Js, fields ...func(*Js)) (ret int) {
	for i, field := range fields {
		if i > 0 {
			out.R(comma)
		}
		out.Brace(obj, field)
		ret++
	}
	return
}

func writeLabel(out *Js, n string) *Js {
	return out.Brace(obj, func(lab *Js) {
		lab.
			Kv("type", FieldLabel).R(comma).
			Kv("text", n)
	})
}
