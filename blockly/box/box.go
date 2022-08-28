package box

import (
	"io/fs"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// reads the passed filesystem of ifspecs
// returns them in blockly toolbox format
func FromSpecs(files fs.FS) (ret string, err error) {
	if ts, e := rs.FromSpecs(files); e != nil {
		err = e
	} else {
		ret = js.Embrace(js.Obj, func(out *js.Builder) {
			out.
				Kv("kind", "categoryToolbox").R(js.Comma).
				Q("contents").R(js.Colon).Brace(js.Array, func(box *js.Builder) {
				// range over groups
				for i, t := range ts.Groups {
					if group, ok := t.Spec.Value.(*spec.GroupSpec); ok {
						groupName := t.FriendlyName()
						if i > 0 {
							box.R(js.Comma)
						}
						box.Brace(js.Obj, func(cat *js.Builder) {
							cat.
								Kv("kind", "category").R(js.Comma).
								Kv("name", groupName).R(js.Comma).
								Q("contents").R(js.Colon).Brace(js.Array, func(els *js.Builder) {
								var b block
								for _, blockType := range group.Specs {
									// only flows
									if _, ok := blockType.Spec.Value.(*spec.FlowSpec); ok {
										stacks, outputs := slotStacks(&ts, &blockType)
										if stacks {
											b.Write(els, bconst.StackedName(blockType.Name))
										}
										if outputs {
											b.Write(els, blockType.Name)
										}
									}
								}
							})
						})
					}
				}
			})
		})
	}
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
func slotStacks(types bconst.Types, blockType *spec.TypeSpec) (stacks, outputs bool) {
	for _, s := range blockType.Slots {
		slotRule := bconst.FindSlotRule(types, s)
		if slotRule.Stack {
			stacks = true
		} else {
			outputs = true
		}
	}
	return
}
