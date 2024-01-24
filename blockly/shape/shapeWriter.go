package shape

import (
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/lang/markup"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
)

type ShapeWriter struct {
	Types TypeMap
}

func NewShapeWriter(ts TypeMap) *ShapeWriter {
	return &ShapeWriter{ts}
}

// return any fields which need mutation
func (w *ShapeWriter) WriteShape(block *js.Builder, blockType typeinfo.T) (okay bool) {
	switch blockType := blockType.(type) {
	case *typeinfo.Flow:
		okay = w.writeFlowBlock(block, blockType)
	case *typeinfo.Slot:
		okay = w.writeSlotBlock(block, blockType)
	case *typeinfo.Num:
		okay = w.writeStandalone(block, blockType)
	case *typeinfo.Str:
		okay = w.writeStandalone(block, blockType)
	default:
		log.Fatalf("unknown type %T", blockType)
	}
	return
}

func (w *ShapeWriter) writeSlotBlock(block *js.Builder, blockType *typeinfo.Slot) bool {
	return false // slots dont have have corresponding blocks
}

// ideally: we'd only need a standalone block for strs, etc. when they implement some slot.
// however, if they are used by a slot -- then we need a block for them too.
// fix: maybe consider writing an "inputDef" object {} as the value of "swaps"
// ( for simple types or maybe all of them ) and change the block's input on selection.
func (w *ShapeWriter) writeStandalone(block *js.Builder, blockType typeinfo.T) (okay bool) {
	// we simply pretend we're a flow of one anonymous member.
	n := blockType.TypeName()
	name := typeinfo.FriendlyName(n)
	okay = w._writeShape(block, name, blockType, []typeinfo.Term{{
		Label: "",
		Name:  n,
		Type:  w.Types[n],
	}})
	return okay
}

func (w *ShapeWriter) writeFlowBlock(block *js.Builder, flow *typeinfo.Flow) bool {
	name := typeinfo.FriendlyName(flow.Lede)
	return w._writeShape(block, name, flow, flow.Terms)
}

// writes one or possible two blocks to represent the blockType.
// will always generate an output block because every type outputs itself
// it may also generate a stackable block if any of the slots implemented have a stackable SlotRule.
// ( ex. a type that implements rt.BoolEval and rt.Execute will write both types of blocks )
func (w *ShapeWriter) _writeShape(block *js.Builder, name string, blockType typeinfo.T, terms []typeinfo.Term) bool {
	stacks, outSlots := slotStacks(w.Types, blockType)
	// we write to partial so that we can potentially have two blocks
	var partial js.Builder
	// color
	var colour string // default
	if c := bconst.BlockColor(blockType); len(c) > 0 {
		colour = c
	} else if len(outSlots) > 0 { // we take on the color of the first slot specified
		n := outSlots[0]
		colour = bconst.BlockColor(n.Type)
	} else if len(stacks) > 0 {
		n := stacks[0]
		colour = bconst.BlockColor(n.Type)
	}
	if len(colour) == 0 {
		colour = bconst.DefaultColor
	}
	partial.Q("colour").Raw(`:"%{BKY_`).Raw(colour).Raw(`}"`)

	// comment
	if cmt := comment(blockType.TypeMarkup()); len(cmt) > 0 {
		partial.R(js.Comma).Kv("tooltip", cmt)
	}
	partial.R(js.Comma)

	// write the terms:
	w.writeShapeDef(&partial, name, blockType, publicTerms(terms))

	// are we stackable? ( ex. story statement or executable )
	if len(stacks) > 0 {
		block.Brace(js.Obj, func(out *js.Builder) {
			out.Kv("type", bconst.StackedName(blockType.TypeName()))
			checks := slotTypes(stacks)
			appendChecks(out, "nextStatement", checks)
			appendChecks(out, "previousStatement", checks)
			appendString(out, partial.String())
		}).R(js.Comma)
	}
	outputs := slotTypes(outSlots)
	if rootBlock := bconst.RootBlock(blockType); !rootBlock {
		outputs = append([]string{blockType.TypeName()}, outputs...)
	}
	block.Brace(js.Obj, func(out *js.Builder) {
		out.Kv("type", blockType.TypeName())
		appendChecks(out, "output", outputs)
		appendString(out, partial.String())
	})
	return true
}

func comment(m map[string]any) (ret string) {
	if lines := markup.UserComment(m); len(lines) > 0 {
		ret = strings.Join(lines, "\n")
	}
	return
}

// write the args0 and message0 key-outputs.
func (w *ShapeWriter) writeShapeDef(out *js.Builder, lede string, blockType typeinfo.T, terms []typeinfo.Term) {
	out.WriteString(`"extensions":["tapestry_generic_mixin","tapestry_generic_extension"],`)
	// note: currently if  excluding empty term sets causes a problem because the block output generates an extraState object that's empty: {}
	// blockly in blocks.js loadExtraState() believes that's valid data; sees we don't have the mutator loadExtraState function
	// and therefore tries to parse the empty object as xml... which raises an (uncaught) exception: stopping the whole page from loading.
	// ( so much power this one little line )
	_, hasMutator := /*len(terms) > 0 &&*/ blockType.(*typeinfo.Flow)
	inlineBlock := bconst.InlineBlock(blockType)
	if hasMutator {
		out.Kv("mutator", "tapestry_generic_mutation").R(js.Comma)
	}
	out.Q("customData").R(js.Colon).
		Brace(js.Obj, func(custom *js.Builder) {
			custom.Q("shapeDef").R(js.Colon).
				Brace(js.Array, func(out *js.Builder) {
					// an initial item containing just the lede
					var comma bool
					if !inlineBlock {
						out.Brace(js.Obj, func(out *js.Builder) {
							out.Kv("label", lede)
						})
						comma = true
					}
					// now any following terms as their own items
					for _, term := range terms {
						if fd, e := w.newFieldDef(term); e != nil {
							log.Fatalln(e) // exit if we couldnt create the field def
						} else if fn := fieldWriter(term.Type); fn != nil {
							if comma {
								out.R(js.Comma)
							}
							comma = true
							out.Brace(js.Obj, func(out *js.Builder) {
								// add a label for non-anonymous fields
								if label := fd.blocklyLabel(); !inlineBlock && len(label) > 0 {
									out.Kv("label", label).R(js.Comma)
								}
								// every term needs a name ( for blockly's sake )
								out.Kv("name", fd.name()).R(js.Comma)
								// write the contents of the term
								fn(w, out, fd.term)
								// write optional, and repeating status
								if fd.term.Optional {
									out.R(js.Comma).Q("optional").R(js.Colon).Raw("true")
								}
								// if we are stack, we want to force a non-repeating input; one stack can already handle multiple blocks.
								// fix? we dont handle the case of a stack of one element; not sure that it exists in practice.
								if !fd.slot.Stack && fd.term.Repeats {
									out.R(js.Comma).Q("repeats").R(js.Colon).Raw("true")
								}
							})
						}
					}
				})
		})
}

func publicTerms(terms []typeinfo.Term) (ret []typeinfo.Term) {
	for _, t := range terms {
		if !t.Private {
			ret = append(ret, t)
		}
	}
	return
}

// split the slots that this type supports into "stacks" and "outputs"
func slotStacks(types TypeMap, blockType typeinfo.T) (retStack, retOutput []bconst.SlotRule) {
	var slots []*typeinfo.Slot
	switch t := blockType.(type) {
	case *typeinfo.Flow:
		slots = t.Slots
	case *typeinfo.Slot:
		slots = []*typeinfo.Slot{t}
	}
	for _, s := range slots {
		slotRule := bconst.MakeSlotRule(s)
		if slotRule.Stack {
			retStack = append(retStack, slotRule)
		} else {
			retOutput = append(retOutput, slotRule)
		}
	}
	return
}
