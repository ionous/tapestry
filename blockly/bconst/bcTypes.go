package bconst

import "git.sr.ht/~ionous/tapestry/lang/typeinfo"

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

func BlockColor(blockType typeinfo.T) (ret string) {
	m := blockType.TypeMarkup()
	if c, ok := m[ColorMarkup].(string); ok {
		ret = c
	}
	return
}

func RootBlock(blockType typeinfo.T) (ret bool) {
	m := blockType.TypeMarkup()
	c, ok := m[RootBlockMarkup].(bool)
	return c && ok
}

func InlineBlock(blockType typeinfo.T) (ret bool) {
	m := blockType.TypeMarkup()
	c, ok := m[InlineBlockMarkup].(bool)
	return c && ok
}
