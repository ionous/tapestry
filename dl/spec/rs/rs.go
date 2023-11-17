// Package rs reads specs from file(system)s
package rs

import (
	"encoding/json"
	"io/fs"
	"sort"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"github.com/ionous/errutil"
)

// type name to type spec lookup
// the map of types differs from what's on disk.
// it excludes groups as types and instead puts those in the groups array.
// it also adds any group containers as members of each type's group list
type TypeSpecs struct {
	Types  map[string]*spec.TypeSpec
	Groups []*spec.TypeSpec
}

func (types *TypeSpecs) FindType(name string) (ret *spec.TypeSpec, okay bool) {
	if t, ok := types.Types[name]; ok {
		ret, okay = t, true
	}
	return
}

// return a list of sorted keys
func (types *TypeSpecs) Keys() []string {
	keys := make([]string, 0, len(types.Types))
	for k, _ := range types.Types {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

type typeMap map[string]*spec.TypeSpec

// reads all of the files in the passed filesystem as ifspecs and returns them as one big map
func FromSpecs(files fs.FS) (ret TypeSpecs, err error) {
	types := TypeSpecs{Types: make(typeMap)}
	if e := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, e error) (err error) {
		if e != nil {
			err = e
		} else if !d.IsDir() && strings.HasSuffix(d.Name(), ".ifspecs") { // the first dir we get is "."
			if _, e := readSpec(&types, files, path); e != nil {
				err = errutil.New("reading", path, e)
			}
		}
		return
	}); e != nil {
		err = e
	} else {
		ret = types
	}
	return
}

func ReadSpec(files fs.FS, fileName string) (ret TypeSpecs, err error) {
	types := TypeSpecs{Types: make(typeMap)}
	if _, e := readSpec(&types, files, fileName); e != nil {
		err = e
	} else {
		ret = types
	}
	return
}

// reads a single typespec from the named file from the passed filesystem into the lookup.
// ( usually a group containing other yet still other typespecs )
func readSpec(types *TypeSpecs, files fs.FS, fileName string) (ret *spec.TypeSpec, err error) {
	if b, e := fs.ReadFile(files, fileName); e != nil {
		err = e
	} else {
		var msg map[string]any
		if e := json.Unmarshal(b, &msg); e != nil {
			err = e
		} else {
			// the outer one is always (supposed to be) a group
			var blockType spec.TypeSpec
			// note: we don't have to pass signatures, because .ifspecs always use concrete types.
			if e := cin.Decode(&blockType, msg, nil); e != nil {
				err = e
			} else if e := importTypes(types, &blockType); e != nil {
				err = e
			} else {
				ret = &blockType
			}
		}
	}
	return
}

func importTypes(types *TypeSpecs, block *spec.TypeSpec) error {
	var currGroups []string
	enc := chart.MakeEncoder()
	return enc.Marshal(block, chart.Map(&enc, chart.BlockMap{
		spec.TypeSpec_Type: chart.KeyMap{
			chart.BlockStart: func(b jsn.Block, _ interface{}) (err error) {
				if flow, ok := b.(jsn.FlowBlock); ok {
					if blockType, ok := flow.GetFlow().(*spec.TypeSpec); ok {
						switch blockType.Spec.Choice {
						case spec.UsesSpec_Group_Opt:
							// the block is group: push it
							currGroups = append(currGroups, blockType.Name)
							types.Groups = append(types.Groups, blockType)

						default:
							types.Types[blockType.Name] = blockType
							// add any parent groups; they take precedence over the inline ones.
							// ( append to an empty array to make sure each blockType gets its own copy )
							blockType.Groups = append([]string{}, append(currGroups, blockType.Groups...)...)
						}
					}
				}
				return
			},
			chart.BlockEnd: func(b jsn.Block, _ interface{}) (err error) {
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
