package blocks

import (
	"log"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
)

// fix: might eventually want to add to this as we are reading files....
var slotRules = SlotRules{{
	Name:   literal.LiteralValue_Type,
	Colour: BKY_MATH_HUE,
}, {
	Name:   rt.Execute_Type,
	Colour: BKY_PROCEDURES_HUE,
}, {
	Name:   rt.Assignment_Type,
	Colour: BKY_PROCEDURES_HUE,
}, {
	Name:   rt.BoolEval_Type,
	Colour: BKY_LOGIC_HUE,
}, {
	Name:   rt.NumberEval_Type,
	Colour: BKY_MATH_HUE,
}, {
	Name:   rt.TextEval_Type,
	Colour: BKY_TEXTS_HUE,
}, {
	Name:   rt.NumListEval_Type,
	Colour: BKY_MATH_HUE,
}, {
	Name:   rt.TextListEval_Type,
	Colour: BKY_TEXTS_HUE,
}, {
	Name:   rt.RecordEval_Type,
	Colour: BKY_LISTS_HUE,
}, {
	Name:   rt.RecordListEval_Type,
	Colour: BKY_LISTS_HUE,
}}

type SlotRules []SlotRule

func (slots SlotRules) FindSlot(name string) (ret SlotRule) {
	var found bool
	for _, s := range slots {
		if s.Name == name {
			ret, found = s, true
			break
		}
	}
	if !found {
		log.Fatalln("couldnt find slot", name)
	}
	return
}

// maybe "internal" could indicate shadows...
type SlotRule struct {
	Name   string
	Stack  bool // if false, then: input_value, if true: input_statement
	Colour string
}

// slots are referenced by terms of a flow
// blockly needs to know whether they stacked or produce a single value
// note: in blockly, a *block* can only be one or the other
// while tapestry allows you to produce and discard a value.
// so, there are some *flows* which can be stackable *or* produce a value.
// ( fwiw: most languages allow this, including go. )
// what we care about here though, is just the slot.
func (slot *SlotRule) InputType() (ret string) {
	if slot.Stack {
		ret = InputStatement
	} else {
		ret = InputValue
	}
	return
}

// ex. "stacked_story_statement"
func (slot *SlotRule) SlotType() (ret string) {
	ret = slot.Name
	if slot.Stack {
		ret = "stacked_" + ret
	}
	return
}
