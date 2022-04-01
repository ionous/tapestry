package gomake

import (
	"bytes"
	"io/fs"
	"sort"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"github.com/ionous/errutil"
)

// reads all of the files in the passed filesystem as ifspecs
// generate golang structs
func WriteSpecs(ifspecs fs.FS, onGroup func(string, []byte) error) (err error) {
	if types, e := rs.FromSpecs(ifspecs); e != nil {
		err = e
	} else {
		// generate slot names, unboxed types, and other useful info.
		ctx := &Context{types: types} // types is rs.TypeSpecs
		ctx.unbox = map[string]string{"text": "string", "bool": "bool"}
		for typeName, t := range types.Types {
			if t.Spec.Choice == spec.UsesSpec_Num_Opt {
				if n := t.Spec.Value.(*spec.NumSpec); !n.Exclusively {
					ctx.unbox[typeName] = "float64"
				}

			}
		}

		if tps, e := newTemplates(ctx); e != nil {
			err = e
		} else {
			for _, groupType := range types.Groups {
				groupName := groupType.Name
				ctx.currentGroup = groupName
				typeNames := getAllTypes(groupType.Spec.Value.(*spec.GroupSpec))
				//
				var out bytes.Buffer
				if e := tps.ExecuteTemplate(&out, "fileHeader.tmpl", map[string]any{
					"Package": groupName,
					"Imports": []string{},
				}); e != nil {
					err = errutil.Append(err, errutil.New(e, "in header", groupName))
				} else {
					// registration lists.
					reg := RegistrationLists{types: types}

					// all the types in this group:
					for _, key := range typeNames {
						t := types.Types[key]
						reg.AddType(t)
						//
						if name := specShortName(t); len(name) == 0 {
							err = errutil.New("unknown type", t.Spec.Choice)
						} else if e := tps.ExecuteTemplate(&out, name+".tmpl", t); e != nil {
							err = errutil.Append(err, errutil.New(e, "couldnt process", key))
						}
					}

					if e := reg.Write(&out, tps); e != nil {
						err = errutil.Append(err, errutil.New(e, "couldnt write registrations"))
					}

					// get whatever we can of the output errors or no.
					if e := onGroup(groupName, out.Bytes()); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
		}
	}
	return
}

// names of all the types in the passed group ( sorted, alpha least first )
func getAllTypes(groupSpec *spec.GroupSpec) []string {
	typeNames := make([]string, 0, len(groupSpec.Specs))
	for _, t := range groupSpec.Specs {
		typeNames = append(typeNames, t.Name)
	}
	sort.Strings(typeNames) // in-place
	return typeNames
}
