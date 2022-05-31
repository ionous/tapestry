package shape

import (
	"log"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/web/js"
)

type ShapeWriter struct {
	rs.TypeSpecs
}

// return any fields which need mutation
func (w *ShapeWriter) WriteShape(block *js.Builder, blockType *spec.TypeSpec) (okay bool) {
	switch t := blockType.Spec.Choice; t {
	case spec.UsesSpec_Flow_Opt:
		okay = w.writeFlowBlock(block, blockType)
	case spec.UsesSpec_Slot_Opt:
		okay = w.writeSlotBlock(block, blockType)
	case spec.UsesSpec_Swap_Opt:
		okay = w.writeStandalone(block, blockType)
	case spec.UsesSpec_Num_Opt:
		okay = w.writeStandalone(block, blockType)
	case spec.UsesSpec_Str_Opt:
		okay = w.writeStandalone(block, blockType)
	default:
		log.Fatalln("unknown type", blockType.Spec.Choice)
	}
	return
}

func (w *ShapeWriter) writeSlotBlock(block *js.Builder, blockType *spec.TypeSpec) bool {
	return false // slots dont have have corresponding blocks
}

// ideally: we'd only need a standalone block for strs, etc. when they implement some slot.
// however, if they are used by a slot -- then we need a block for them too.
// fix: maybe consider writing an "inputDef" object {} as the value of "swaps"
// ( for simple types or maybe all of them ) and change the block's input on selection.
func (w *ShapeWriter) writeStandalone(block *js.Builder, blockType *spec.TypeSpec) (okay bool) {
	// we simply pretend we're a flow of one anonymous member.
	name := spec.FriendlyName(blockType.Name, false)
	okay = w._writeShape(block, name, blockType, []spec.TermSpec{{
		Label: "",
		Name:  blockType.Name,
		Type:  blockType.Name,
	}})
	return okay
}

func (w *ShapeWriter) writeFlowBlock(block *js.Builder, blockType *spec.TypeSpec) bool {
	flow := blockType.Spec.Value.(*spec.FlowSpec)
	name := flow.FriendlyLede(blockType)
	return w._writeShape(block, name, blockType, flow.Terms)
}

// writes one or possible two blocks to represent the blockType.
// will always generate an output block because every type outputs itself
// it may also generate a stackable block if any of the slots implemented have a stackable SlotRule.
// ( ex. a type that implements rt.BoolEval and rt.Execute will write both types of blocks )
func (w *ShapeWriter) _writeShape(block *js.Builder, name string, blockType *spec.TypeSpec, terms []spec.TermSpec) bool {
	stacks, values := slotStacks(blockType)
	// we write to partial so that we can potentially have two blocks
	var partial js.Builder
	// MOD-stravis... try without this.
	// write the label for the block itself; aka the lede.
	// partial.Kv("message0", name).R(js.Comma)
	// color
	var colour string = bconst.COLOUR_HUE // default
	if len(values) > 0 {                  // we take on the color of the first slot specified
		slot := bconst.FindSlotRule(values[0])
		colour = slot.Colour
	} else if len(stacks) > 0 {
		slot := bconst.FindSlotRule(stacks[0])
		colour = slot.Colour
	}
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
	w.writeShapeDef(&partial, name, blockType, terms)

	// are we stackable? ( ex. story statement or executable )
	if len(stacks) > 0 {
		block.Brace(js.Obj, func(out *js.Builder) {
			out.Kv("type", bconst.StackedName(blockType.Name))
			appendChecks(out, "nextStatement", stacks)
			appendChecks(out, "previousStatement", stacks)
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
func slotStacks(blockType *spec.TypeSpec) (retStack, retValue []string) {
	var slots []string
	if blockType.Spec.Choice == spec.UsesSpec_Slot_Opt {
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
