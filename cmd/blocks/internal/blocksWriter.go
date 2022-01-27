package blocks

import (
	"git.sr.ht/~ionous/tapestry/dl/spec"
)

// return any fields which need mutation
// tbd: "helpUrl"
func writeBlock(block *Js, blockType *spec.TypeSpec) (okay bool) {
	stacks, values := SlotStacks(blockType)
	switch blockType.Spec.Choice {
	case spec.UsesSpec_Flow_Opt:
		flow := blockType.Spec.Value.(*spec.FlowSpec)
		// we write to partial so that we can potentially have two blocks
		// one if we are stackable, one if we output a value
		// ( ex. something that is executable or returns bool )
		var partial Js
		// write the label for the block itself; aka the lede.
		partial.Q("message0").R(colon)
		if lede := flow.Name; len(lede) > 0 {
			partial.Q(lede)
		} else {
			partial.Q(blockType.Name)
		}
		// color
		var colour string = BKY_COLOUR_HUE // default
		if len(values) > 0 {               // we take on the color of the first slot specified
			slot := slotRules.FindSlot(values[0])
			colour = slot.Colour
		} else if len(stacks) > 0 {
			slot := slotRules.FindSlot(stacks[0])
			colour = slot.Colour
		}
		partial.R(comma).Kv("colour", colour)
		// comment
		if cmt := blockType.UserComment; len(cmt) > 0 {
			partial.R(comma).Kv("tooltip", cmt)
		}
		partial.R(comma)
		writeCustomData(&partial, blockType, flow)
		// are we stackable? ( ex. story statement or executable )
		if len(stacks) > 0 {
			block.Brace(obj, func(out *Js) {
				checks := quotedStrings(stacks)
				out.
					Kv("type", "stacked_"+blockType.Name).R(comma).
					Q("nextStatement").R(colon).S(checks).R(comma).
					Q("prevStatement").R(colon).S(checks).R(comma).
					S(partial.String())
			})
		}
		block.Brace(obj, func(out *Js) {
			out.
				Kv("type", blockType.Name).R(comma).
				Q("output").R(colon).Brace(array, func(checks *Js) {
				// add the flow itself as a possible output type
				// (useful for cases where the its used directly by other flows)
				checks.Q(blockType.Name)
				for _, el := range values {
					checks.R(comma).Q(el)
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
