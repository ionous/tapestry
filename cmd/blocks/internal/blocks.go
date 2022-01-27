package blocks

import (
	"io/fs"

	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
)

func readSpec(files fs.FS, fileName string) (ret *spec.TypeSpec, err error) {
	if b, e := fs.ReadFile(files, fileName); e != nil {
		err = e
	} else {
		// the outer one is always (supposed to be) a group
		var blockType spec.TypeSpec
		if e := cin.Decode(&blockType, b, cin.Signatures{
			spec.Signatures, // we are reading specs so we need the spec signatures
			prim.Signatures, // and some of those commands use the primitive types.
		}); e != nil {
			err = e
		} else if e := importTypes(&blockType); e != nil {
			err = e
		} else {
			ret = &blockType
		}
	}
	return
}

func importTypes(types *spec.TypeSpec) error {
	var currGroups []string
	enc := chart.MakeEncoder()
	return enc.Marshal(types, story.Map(&enc, story.BlockMap{
		spec.TypeSpec_Type: story.KeyMap{
			story.BlockStart: func(b jsn.Block, _ interface{}) (err error) {
				if flow, ok := b.(jsn.FlowBlock); ok {
					if blockType, ok := flow.GetFlow().(*spec.TypeSpec); ok {
						switch blockType.Spec.Choice {
						case spec.UsesSpec_Group_Opt:
							// the block is group: push it
							currGroups = append(currGroups, blockType.Name)
						default:
							lookup[blockType.Name] = blockType
							// add in all of the parent groups --
							// they take precedence over the extra groups listed
							blockType.Groups = append(currGroups, blockType.Groups...)
						}
					}
				}
				return
			},
			story.BlockEnd: func(b jsn.Block, _ interface{}) (err error) {
				if flow, ok := b.(jsn.FlowBlock); ok {
					if blockType, ok := flow.GetFlow().(*spec.TypeSpec); ok {
						switch blockType.Spec.Choice {
						case spec.UsesSpec_Group_Opt:
							// the block *was* a group: pop it.
							currGroups = currGroups[:len(currGroups)-1]
						}
					}
				}
				return
			},
		},
	}))
}

var lookup = make(map[string]*spec.TypeSpec)
