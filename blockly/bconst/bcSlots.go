package bconst

import (
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

type Types interface {
	FindType(string) (*spec.TypeSpec, bool)
}

func MakeSlotRule(slot *typeinfo.Slot) SlotRule {
	stacks, _ := slot.Markup[StackMarkup].(bool)
	color, _ := slot.Markup[ColorMarkup].(string)
	return SlotRule{
		Name:   slot.Name,
		Stack:  stacks,
		Colour: color,
	}
}

// pass the name of the slot
func FindSlotRule(types Types, name string) (ret SlotRule) {
	if slot, ok := types.FindType(name); !ok {
		ret = SlotRule{Name: name}
	} else {
		stacks, _ := slot.Markup[StackMarkup].(bool)
		ret = SlotRule{
			Name:   name,
			Stack:  stacks,
			Colour: BlockColor(slot),
		}
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
		ret = StackedName(ret)
	}
	return
}
