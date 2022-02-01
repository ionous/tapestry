package shapes

import (
	"log"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/web/js"
)

var writeBlockType = map[string]func(*js.Builder, *spec.TypeSpec) bool{
	spec.UsesSpec_Flow_Opt: writeFlowBlock,
	spec.UsesSpec_Slot_Opt: writeSlotBlock,
	spec.UsesSpec_Swap_Opt: writeStandalone,
	spec.UsesSpec_Num_Opt:  writeStandalone,
	spec.UsesSpec_Str_Opt:  writeStandalone,
}

// return any fields which need mutation
func writeBlock(block *js.Builder, blockType *spec.TypeSpec) (okay bool) {
	if cb, ok := writeBlockType[blockType.Spec.Choice]; !ok {
		log.Fatalln("unknown type", blockType.Spec.Choice)
	} else {
		okay = cb(block, blockType)
	}
	return
}

func writeSlotBlock(block *js.Builder, blockType *spec.TypeSpec) bool {
	return false // slots dont have have corresponding blocks
}

// ideally: we'd only need a standalone block for strs, etc. when they implement some slot.
// however, if they are used by a slot -- then we need a block for them too.
// fix: maybe consider writing an "inputDef" object {} as the value of "swaps"
// ( for simple types or maybe all of them ) and change the block's input on selection.
func writeStandalone(block *js.Builder, blockType *spec.TypeSpec) (okay bool) {
	// we simply pretend we're a flow of one anonymous member.
	okay = _writeBlock(block, blockType.Name, blockType, []spec.TermSpec{{
		Key:  "",
		Name: blockType.Name,
		Type: blockType.Name,
	}})
	return okay
}

func writeFlowBlock(block *js.Builder, blockType *spec.TypeSpec) bool {
	flow := blockType.Spec.Value.(*spec.FlowSpec)
	name := blockType.Name
	if n := flow.Name; len(n) > 0 {
		name = n
	}
	return _writeBlock(block, name, blockType, flow.Terms)
}

// writes one or possible two blocks to represent the blockType.
// will always generate an output block because every type outputs itself
// it may also generate a stackable block if any of the slots implemented have a stackable SlotRule.
// ( ex. a type that implements rt.BoolEval and rt.Execute will write both types of blocks )
func _writeBlock(block *js.Builder, name string, blockType *spec.TypeSpec, terms []spec.TermSpec) bool {
	stacks, values := SlotStacks(blockType)
	// we write to partial so that we can potentially have two blocks
	var partial js.Builder
	// write the label for the block itself; aka the lede.
	partial.Kv("message0", name)
	// color
	var colour string = bconst.COLOUR_HUE // default
	if len(values) > 0 {                  // we take on the color of the first slot specified
		slot := bconst.FindSlotRule(values[0])
		colour = slot.Colour
	} else if len(stacks) > 0 {
		slot := bconst.FindSlotRule(stacks[0])
		colour = slot.Colour
	}
	partial.R(js.Comma)
	if len(colour) > 0 {
		partial.Kv("colour", colour)
	} else {
		partial.Kv("colour", bconst.COLOUR_HUE)
	}
	// comment
	if cmt := blockType.UserComment; len(cmt) > 0 {
		partial.R(js.Comma).Kv("tooltip", cmt)
	}
	partial.R(js.Comma)

	// write the terms:
	writeBlockDef(&partial, blockType, terms)

	// are we stackable? ( ex. story statement or executable )
	if len(stacks) > 0 {
		block.Brace(js.Obj, func(out *js.Builder) {
			out.Kv("type", "stacked_"+blockType.Name)
			appendChecks(out, "nextStatement", stacks)
			appendChecks(out, "prevStatement", stacks)
			appendString(out, partial.String())
		}).R(js.Comma)
	}
	if rootBlock := rootBlocks.IsRoot(blockType.Name); !rootBlock {
		values = append([]string{blockType.Name}, values...)
	}
	block.Brace(js.Obj, func(out *js.Builder) {
		out.Kv("type", blockType.Name)
		appendChecks(out, "output", values)
		appendString(out, partial.String())
	})
	return true
}

func appendString(out *js.Builder, s string) {
	if len(s) > 0 {
		out.R(js.Comma).S(s)
	}
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
		slotRule := bconst.FindSlotRule(s)
		if slotRule.Stack {
			retStack = append(retStack, slotRule.SlotType())
		} else {
			retValue = append(retValue, s)
		}
	}
	return
}
