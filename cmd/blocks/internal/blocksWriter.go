package blocks

import (
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/tile/bc"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// return any fields which need mutation
// tbd: "helpUrl"
func writeBlock(block *js.Builder, blockType *spec.TypeSpec) (okay bool) {
	stacks, values := SlotStacks(blockType)
	switch blockType.Spec.Choice {
	case spec.UsesSpec_Flow_Opt:
		flow := blockType.Spec.Value.(*spec.FlowSpec)
		// we write to partial so that we can potentially have two blocks
		// one if we are stackable, one if we output a value
		// ( ex. something that is executable or returns bool )
		var partial js.Builder
		// write the label for the block itself; aka the lede.
		partial.Q("message0").R(js.Colon)
		if lede := flow.Name; len(lede) > 0 {
			partial.Q(lede)
		} else {
			partial.Q(blockType.Name)
		}
		// color
		var colour string = bc.COLOUR_HUE // default
		if len(values) > 0 {              // we take on the color of the first slot specified
			slot := bc.FindSlotRule(values[0])
			colour = slot.Colour
		} else if len(stacks) > 0 {
			slot := bc.FindSlotRule(stacks[0])
			colour = slot.Colour
		}
		partial.R(js.Comma).Kv("colour", colour)
		// comment
		if cmt := blockType.UserComment; len(cmt) > 0 {
			partial.R(js.Comma).Kv("tooltip", cmt)
		}
		partial.R(js.Comma)
		writeCustomData(&partial, blockType, flow)
		// are we stackable? ( ex. story statement or executable )
		if len(stacks) > 0 {
			block.Brace(js.Obj, func(out *js.Builder) {
				checks := js.QuotedStrings(stacks)
				out.
					Kv("type", "stacked_"+blockType.Name).R(js.Comma).
					Q("nextStatement").R(js.Colon).S(checks).R(js.Comma).
					Q("prevStatement").R(js.Colon).S(checks).R(js.Comma).
					S(partial.String())
			})
		}
		block.Brace(js.Obj, func(out *js.Builder) {
			out.
				Kv("type", blockType.Name).R(js.Comma).
				Q("output").R(js.Colon).Brace(js.Array, func(checks *js.Builder) {
				// add the flow itself as a possible output type
				// (useful for cases where the its used directly by other flows)
				checks.Q(blockType.Name)
				for _, el := range values {
					checks.R(js.Comma).Q(el)
				}
			}).
				R(js.Comma).
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
		slotRule := bc.FindSlotRule(s)
		if slotRule.Stack {
			retStack = append(retStack, slotRule.SlotType())
		} else {
			retValue = append(retValue, s)
		}
	}
	return
}
