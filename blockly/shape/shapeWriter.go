package shape

import (
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/web/js"
)

type ShapeWriter struct {
	rs.TypeSpecs
}

func NewShapeWriter(ts rs.TypeSpecs) *ShapeWriter {
	return &ShapeWriter{ts}
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
	// color
	var colour string    // default
	if len(values) > 0 { // we take on the color of the first slot specified
		slot := bconst.FindSlotRule(values[0])
		colour = slot.Colour
	} else if len(stacks) > 0 {
		slot := bconst.FindSlotRule(stacks[0])
		colour = slot.Colour
	}
	if len(colour) == 0 {
		colour = bconst.COLOUR_HUE
	}
	partial.Kv("colour", colour)

	// comment
	if cmt := comment(blockType.Markup); len(cmt) > 0 {
		partial.R(js.Comma).Kv("tooltip", cmt)
	}
	partial.R(js.Comma)

	// write the terms:
	w.writeShapeDef(&partial, name, blockType, publicTerms(terms))

	// are we stackable? ( ex. story statement or executable )
	if len(stacks) > 0 {
		block.Brace(js.Obj, func(out *js.Builder) {
			out.Kv("type", bconst.StackedName(blockType.Name))
			appendChecks(out, "nextStatement", stacks)
			appendChecks(out, "previousStatement", stacks)
			appendString(out, partial.String())
		}).R(js.Comma)
	}
	if rootBlock := blockType.InGroup(RootBlock); !rootBlock {
		values = append([]string{blockType.Name}, values...)
	}
	block.Brace(js.Obj, func(out *js.Builder) {
		out.Kv("type", blockType.Name)
		appendChecks(out, "output", values)
		appendString(out, partial.String())
	})
	return true
}

func comment(markup map[string]any) (ret string) {
	switch cmt := markup["comment"].(type) {
	case string:
		ret = cmt
	case []string:
		ret = strings.Join(cmt, "\n")
	}
	return
}

// write the args0 and message0 key-values.
func (w *ShapeWriter) writeShapeDef(out *js.Builder, lede string, blockType *spec.TypeSpec, terms []spec.TermSpec) {
	out.WriteString(`"extensions":["tapestry_generic_mixin","tapestry_generic_extension"],`)
	// note: currently if  excluding empty term sets causes a problem because the block output generates an extraState object that's empty: {}
	// blockly in blocks.js loadExtraState() believes that's valid data; sees we don't have the mutator loadExtraState function
	// and therefore tries to parse the empty object as xml... which raises an (uncaught) exception: stopping the whole page from loading.
	// ( so much power this one little line )
	hasMutator := /*len(terms) > 0 &&*/ blockType.Spec.Choice == spec.UsesSpec_Flow_Opt
	inlineBlock := blockType.InGroup(InlineBlock)
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
						} else if fn, ok := fieldWriter[fd.termType()]; !ok {
							log.Fatalln("unknown term type", fd.termType())
						} else if fn != nil {
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
								fn(out, fd.term, fd.typeSpec)
								// write optional, and repeating status
								if fd.term.Optional {
									out.R(js.Comma).Q("optional").R(js.Colon).S("true")
								}
								// if we are stack, we want to force a non-repeating input; one stack can already handle multiple blocks.
								// fix? we dont handle the case of a stack of one element; not sure that it exists in practice.
								if !fd.slot.Stack && fd.term.Repeats {
									out.R(js.Comma).Q("repeats").R(js.Colon).S("true")
								}
							})
						}
					}
				})
		})
}

func publicTerms(terms []spec.TermSpec) (ret []spec.TermSpec) {
	for _, t := range terms {
		if !t.Private {
			ret = append(ret, t)
		}
	}
	return
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
