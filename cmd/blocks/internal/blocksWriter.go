package blocks

import (
	"git.sr.ht/~ionous/tapestry/dl/spec"
)

// return any fields which need mutation
// tbd: "helpUrl"
func writeBlock(out *Js, blockType *spec.TypeSpec) (okay bool) {
	stacks, values := SlotStacks(blockType)
	switch blockType.Spec.Choice {
	case spec.UsesSpec_Flow_Opt:
		flow := blockType.Spec.Value.(*spec.FlowSpec)
		// we write to partial so that we can potentially have two blocks
		// one if we are stackable, one if we output a value
		// ( ex. something that is executable or returns bool )
		var partial Js
		writeBlockInternals(&partial, blockType, flow)
		// comment
		if cmt := blockType.UserComment; len(cmt) > 0 {
			partial.R(comma).Kv("tooltip", cmt)
		}
		// are we stackable? ( ex. story statement or executable )
		var colour string = BKY_COLOUR_HUE // default
		if len(stacks) > 0 {
			out.Brace(obj, func(out *Js) {
				slot := slotRules.FindSlot(stacks[0])
				types := quotedStrings(stacks)
				colour = slot.Colour //
				out.
					Kv("type", "stacked_"+blockType.Name).R(comma).
					Q("nextStatement").R(colon).S(types).R(comma).
					Q("prevStatement").R(colon).S(types).R(comma).
					S(partial.String())
			})
		}
		if len(values) > 0 { // we take on the color of the first slot specified
			slot := slotRules.FindSlot(values[0])
			colour = slot.Colour
		}
		if len(stacks) > 0 {
			out.R(comma)
		}
		out.Brace(obj, func(out *Js) {
			out.
				Kv("type", blockType.Name).R(comma).
				Kv("colour", colour).R(comma).
				Q("output").R(colon).Brace(array, func(out *Js) {
				// add the flow itself as a possible output type
				// (useful for cases where the its used directly by other flows)
				out.Q(blockType.Name)
				for _, el := range values {
					out.R(comma).Q(el)
				}
			}).
				R(comma).
				S(partial.String())
		})
		okay = true
	}
	return
}

// split the slots that this type supports into "stacks" and "values"
func SlotStacks(blockType *spec.TypeSpec) (retStack, retValue []string) {
	var slots []string
	if blockType.Spec.Choice == spec.TypeSpec_Field_Slots {
		slots = []string{blockType.Name}
	} else {
		slots = blockType.Slots
	}
	for _, s := range slots {
		slotRule := slotRules.FindSlot(s)
		if slotRule.Stack {
			retStack = append(retStack, slotRule.SlotType())
		} else {
			retValue = append(retValue, s)
		}
	}
	return
}
