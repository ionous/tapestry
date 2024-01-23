package box

import (
	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// reads all specs from the passed file system
// returns them in blockly toolbox format
func FromTypeSet(types []*typeinfo.TypeSet) (ret string, err error) {
	ret = js.Embrace(js.Obj, func(out *js.Builder) {
		out.
			Kv("kind", "categoryToolbox").R(js.Comma).
			Q("contents").R(js.Colon).Brace(js.Array, func(box *js.Builder) {
			// range over groups
			for i, group := range types {
				groupName := typeinfo.FriendlyName(group.Name)
				if i > 0 {
					box.R(js.Comma)
				}
				box.Brace(js.Obj, func(cat *js.Builder) {
					cat.
						Kv("kind", "category").R(js.Comma).
						Kv("name", groupName).R(js.Comma).
						Q("contents").R(js.Colon).Brace(js.Array, func(els *js.Builder) {
						var b block // only flows
						for _, blockType := range group.Flow {
							stacks, outputs := slotStacks(blockType)
							if stacks {
								b.Write(els, bconst.StackedName(blockType.Name))
							}
							if outputs {
								b.Write(els, blockType.Name)
							}
						}
					})
				})
			}
		})
	})
	return
}

type block int

func (b *block) Write(els *js.Builder, name string) {
	if *b = (*b) + 1; *b > 1 {
		els.R(js.Comma)
	}
	els.Brace(js.Obj, func(el *js.Builder) {
		el.
			Kv("kind", "block").R(js.Comma).
			Kv("type", name)
	})
}

// split the slots that this type supports into "stacks" and "outputs"
func slotStacks(blockType *typeinfo.Flow) (stacks, outputs bool) {
	for _, s := range blockType.Slots {
		slotRule := bconst.MakeSlotRule(s)
		if slotRule.Stack {
			stacks = true
		} else {
			outputs = true
		}
	}
	return
}
