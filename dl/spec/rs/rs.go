// Package rs reads specs from file(system)s
package rs

import (
	"io/fs"
	"sort"

	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"github.com/ionous/errutil"
)

// type name to type spec lookup
type TypeSpecs struct {
	Types  map[string]*spec.TypeSpec
	Groups []*spec.TypeSpec
}

// return a list of sorted keys
func (ts *TypeSpecs) Keys() []string {
	keys := make([]string, 0, len(ts.Types))
	for k, _ := range ts.Types {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

type typeMap map[string]*spec.TypeSpec

// reads all of the files in the passed filesystem as ifspecs and returns them as one big map
func FromSpecs(files fs.FS) (ret TypeSpecs, err error) {
	ts := TypeSpecs{Types: make(typeMap)}
	if e := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, e error) (err error) {
		if e != nil {
			err = e // can happen if it failed to read the contents of a director
		} else if !d.IsDir() { // the first dir we get is "."
			if _, e := readSpec(&ts, files, path); e != nil {
				err = errutil.New(e, "reading", path)
			}
		}
		return
	}); e != nil {
		err = e
	} else {
		ret = ts
	}
	return
}

func ReadSpec(files fs.FS, fileName string) (ret TypeSpecs, err error) {
	ts := TypeSpecs{Types: make(typeMap)}
	if _, e := readSpec(&ts, files, fileName); e != nil {
		err = e
	} else {
		ret = ts
	}
	return
}

// reads a single typespec from the named file from the passed filesystem into the lookup.
// ( usually a group containing other yet still other typespecs )
func readSpec(ts *TypeSpecs, files fs.FS, fileName string) (ret *spec.TypeSpec, err error) {
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
		} else if e := importTypes(ts, &blockType); e != nil {
			err = e
		} else {
			ret = &blockType
		}
	}
	return
}

func importTypes(ts *TypeSpecs, types *spec.TypeSpec) error {
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
							ts.Groups = append(ts.Groups, blockType)

						default:
							ts.Types[blockType.Name] = blockType
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
