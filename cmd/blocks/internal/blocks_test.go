package blocks

import (
	"log"
	"strconv"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/idl"
)

func TestBlocks(t *testing.T) {
	// errutil.Panic = true
	if str, e := run(t); e != nil {
		t.Fatal(e)
	} else {
		t.Log(str)
	}
}

func run(t *testing.T) (ret string, err error) {
	if _, e := readSpec(idl.Specs, "prim.ifspecs"); e != nil {
		err = e
	} else if _, e := readSpec(idl.Specs, "literal.ifspecs"); e != nil {
		err = e
	} else {
		ret = Embrace(array, func(out *Js) {
			var sep Commas
			for _, tap := range lookup {
				if s := writeBlock(tap); len(s) > 0 {
					out.R(sep.Next()).S(s)
				}
			}
		})
	}
	return
}

func writeBlock(tap *spec.TypeSpec) string {
	var out Js
	stacks, values := SlotTypes(tap)
	switch tap.Spec.Choice {
	case spec.UsesSpec_Flow_Opt:
		flow := tap.Spec.Value.(*spec.FlowSpec)
		var partial Js
		partial.
			Kv("message0", MessageOf(tap.Name, flow)).R(comma).
			Q("args0").R(colon).S(Args(flow))
		if cmt := tap.UserComment; len(cmt) > 0 {
			partial.R(comma).
				Kv("tooltip", cmt)
		}
		// out.Kv("helpUrl", ?)

		if len(stacks) > 0 {
			out.Brace(obj, func(out *Js) {
				slot := slotRules.FindSlot(stacks[0])
				types := QuotedStrings(stacks)
				out.
					Kv("type", "stacked_"+tap.Name).R(comma).
					Kv("colour", slot.Colour).R(comma).
					Q("nextStatement").R(colon).S(types).R(comma).
					Q("prevStatement").R(colon).S(types).R(comma).
					S(partial.String())

			})
		}
		// add the flow itself as a possible output type
		// (useful for cases where the its used directly by other flows)
		var colour string = BKY_COLOUR_HUE // default
		if len(values) > 0 {
			slot := slotRules.FindSlot(values[0])
			colour = slot.Colour
		}
		if len(stacks) > 0 {
			out.R(comma)
		}
		out.Brace(obj, func(out *Js) {
			types := QuotedStrings(append(values, flow.Name))
			out.
				Kv("type", tap.Name).R(comma).
				Kv("colour", colour).R(comma).
				Q("output").R(colon).S(types).R(comma).
				S(partial.String())
		})
	}
	return out.String()
}

// split the slots that this type supports into "stacks" and "values"
func SlotTypes(spec *spec.TypeSpec) (retStack, retValue []string) {
	for _, s := range spec.Slots {
		slot := slotRules.FindSlot(s)
		if slot.Stack {
			retStack = append(retStack, slot.SlotType())
		} else {
			retValue = append(retValue, s)
		}
	}
	return
}

// returns a comma separated list of {} block arguments.
// ex. [{"type": "input_value", "name": "DELTA"}]
func Args(flow *spec.FlowSpec) string {
	return Embrace(array, func(out *Js) {
		var sep Commas
		for _, term := range flow.Terms {
			if term.Private {
				continue // skip private terms
			}
			out.R(sep.Next())
			// arg entry
			var name string = term.Key
			if n := term.Name; len(n) > 0 {
				name = n
			}
			var kind string = term.Type
			if kind == "" {
				kind = name
			}
			if termType, ok := lookup[kind]; !ok {
				log.Fatalln("missing named type", kind)
			} else {
				switch t := termType.Spec.Choice; t {

				case spec.UsesSpec_Flow_Opt:
					// fix? im not sure what's best here.
					// for now: using a block, and in the toolbox: a shadow so authors dont have to find it
					// for repeating fields of these, youd get a mutation with dangling inputs ( like blockly's lists )
					out.S(NewInput(name, InputValue, term.Type))

				case spec.UsesSpec_Slot_Opt:
					slot := slotRules.FindSlot(term.Type)
					out.S(NewInput(name, slot.InputType(), slot.SlotType()))

				case spec.UsesSpec_Swap_Opt:
					out.S(NewDropdown(name, []string{"first", "ITEM"}))

				case spec.UsesSpec_Num_Opt:
					out.S(NewNumber(name))

				case spec.UsesSpec_Str_Opt:
					out.S(NewText(name))

				default:
					log.Fatalln("unknown spec type", t)
				}
			}
		}
	})
}

// name: sepional display name; if not this then we use the name of the type
// phrase -- ignore
// trim: influences blockly msg ( first param anon )
// uses terms....
// - key: this is our label ( unless first param anon. )
// - name: if specified then its our field name ( we need to UPPER for blockly ), otherwise key is
// - type: our field or input pin depends on what kind of spec the named type is.
// ---
// - private: this term is skipped; invisible to blockly
// - sepional: FIX: becomes a mutator
// --- how does that work with messages?
// - repeats: depends on what we are repeating
func MessageOf(k string, flow *spec.FlowSpec) string {
	var msg strings.Builder
	var lede string = k
	if n := flow.Name; len(n) > 0 {
		lede = n
	}
	msg.WriteString(strings.ReplaceAll(lede, "_", " "))

	var el int
	for _, t := range flow.Terms {
		if t.Private {
			continue // skip private terms
		}
		msg.WriteRune(space)
		// only write labels for body elements and non-trimmed first elements
		if !flow.Trim || el > 0 {
			msg.WriteString(t.Key)
			msg.WriteRune(space)
		}
		el = el + 1 // pre-inc b/c msgs are one-indexed
		msg.WriteString("%" + strconv.Itoa(el))
	}
	return msg.String()
}
