package box

import (
	"io/fs"
	"strings"

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
									// only write flows to the toolbox
									if flow, ok := blockType.Spec.Value.(*spec.FlowSpec); ok {
										stacks, outputs := slotStacks(&blockType)

										terms := func() {
											var fields bool
											for _, term := range flow.Terms {
												var placeholder string
												if a, ok := ts.Types[term.TypeName()]; ok {
													if str, ok := a.Spec.Value.(*spec.StrSpec); ok {
														if len(str.Uses) == 0 {
															if !term.IsAnonymous() {
																placeholder = term.Label
															} else {
																placeholder = term.Name
															}
														}
													}
												}
												if len(placeholder) > 0 {
													els.R(js.Comma)
													if !fields {
														fields = true
														els.Q("fields").R(js.Colon).R(js.Obj[0])
													}
													fieldName := strings.ToUpper(term.Field())
													els.Kv(fieldName, strings.Replace(placeholder, "_", "", -1))
												}

											}
											// close if opened
											if fields {
												els.R(js.Obj[1])
											}
										}
										if stacks {
											b.Write(els, bconst.StackedName(blockType.Name), terms)
										}
										if outputs {
											b.Write(els, blockType.Name, terms)
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

func (b *block) Write(els *js.Builder, name string, hack func()) {
	if *b = (*b) + 1; *b > 1 {
		els.R(js.Comma)
	}
	els.Brace(js.Obj, func(el *js.Builder) {
		el.
			Kv("kind", "block").R(js.Comma).
			Kv("type", name)
		hack()
	})
}

// split the slots that this type supports into "stacks" and "values"
func slotStacks(blockType *spec.TypeSpec) (stacks, outputs bool) {
	for _, s := range blockType.Slots {
		slotRule := bconst.FindSlotRule(s)
		if slotRule.Stack {
			stacks = true
		} else {
			outputs = true
		}
	}
	return
}
