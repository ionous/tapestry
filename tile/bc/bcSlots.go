package bc

import (
	"log"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
)

func FindSlotRule(name string) (ret SlotRule) {
	var found bool
	for _, s := range slots {
		if s.Name == name {
			ret, found = s, true
			break
		}
	}
	if !found {
		log.Fatalln("couldn't find slot", name)
	}
	return
}

type SlotRule struct {
	Name   string
	Stack  bool // if false, then: input_value, if true: input_statement
	Colour string
	// fix? maybe "internal" could indicate shadows...
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

// fix: might eventually want to add to this as we are reading files....
var slots = []SlotRule{{
	Name:   literal.LiteralValue_Type,
	Colour: MATH_HUE,
}, {
	Name:   rt.Execute_Type,
	Colour: PROCEDURES_HUE,
	Stack:  true,
}, {
	Name:   rt.Assignment_Type,
	Colour: PROCEDURES_HUE,
}, {
	Name:   rt.BoolEval_Type,
	Colour: LOGIC_HUE,
}, {
	Name:   rt.NumberEval_Type,
	Colour: MATH_HUE,
}, {
	Name:   rt.TextEval_Type,
	Colour: TEXTS_HUE,
}, {
	Name:   rt.NumListEval_Type,
	Colour: MATH_HUE,
}, {
	Name:   rt.TextListEval_Type,
	Colour: TEXTS_HUE,
}, {
	Name:   rt.RecordEval_Type,
	Colour: LISTS_HUE,
}, {
	Name:   rt.RecordListEval_Type,
	Colour: LISTS_HUE,
}, {
	Name:   story.StoryStatement_Type,
	Colour: VARIABLES_HUE,
	Stack:  true,
}}
