package shape

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// "append" b/c preceeded by a comma if needed
func appendString(out *js.Builder, s string) {
	if len(s) > 0 {
		out.R(js.Comma).Raw(s)
	}
}

// "append" b/c preceeded by a comma if needed
func appendKv(out *js.Builder, k, v string) {
	if len(v) > 0 {
		out.R(js.Comma).Kv(k, v)
	}
}

// "append" b/c preceeded by a comma if needed
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

func slotTypes(rules []bconst.SlotRule) []string {
	str := make([]string, len(rules))
	for i, s := range rules {
		str[i] = s.SlotType()
	}
	return str
}

func writeShapeDef(out *js.Builder, cb func(*js.Builder)) {
	out.WriteString(`,"extensions":["tapestry_generic_mixin","tapestry_generic_extension"]`)
	out.R(js.Comma).Q("customData").R(js.Colon).
		Brace(js.Obj, func(custom *js.Builder) {
			out.Q("shapeDef").R(js.Colon).
				Brace(js.Array, cb)
		})
}

// write a number of fields followed by their dummy input owner
// returning the number of items written.
func writeDummy(out *js.Builder, name, label string, fields ...func(*js.Builder)) (ret int) {
	// write any leading text as a field
	if n := label; len(n) > 0 {
		writeLabel(out, n)
		out.R(js.Comma)
		ret++
	}
	// write the explicit fields
	for _, field := range fields {
		out.Brace(js.Obj, field)
		out.R(js.Comma)
		ret++
	}
	// write the input the fields are a part of
	out.Brace(js.Obj, func(tail *js.Builder) {
		if n := name; len(n) > 0 {
			tail.Kv("name", strings.ToUpper(n)).R(js.Comma)
		}
		tail.Kv("type", bconst.InputDummy)
		ret++
	})
	return
}

func writeLabel(out *js.Builder, n string) *js.Builder {
	return out.Brace(js.Obj, func(lab *js.Builder) {
		lab.
			Kv("type", bconst.FieldLabel).R(js.Comma).
			Kv("text", n)
	})
}
