package bconst

import (
	"git.sr.ht/~ionous/tapestry/dl/spec"
)

// transform the passed block name into the name of its corresponding stacked block.
// ex. _name_stack
func StackedName(blockType string) string {
	return "_" + blockType + "_stack"
}

// transform the passed block name into the name of its corresponding mutator block.
// ex. _name_mutator
func MutatorName(blockType string) string {
	return "_" + blockType + "_mutator"
}

func BlockColor(blockType *spec.TypeSpec) (ret string) {
	if c, ok := blockType.Markup[ColorMarkup].(string); ok {
		ret = c
	}
	return
}
