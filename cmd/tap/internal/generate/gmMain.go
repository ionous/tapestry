package generate

import (
	"bytes"
	"io"
	"io/fs"
	"sort"
	"text/template"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/support/distill"
	"github.com/ionous/errutil"
)

// read all of the files in the passed filesystem as interface description files,
// and generate golang structs
func WriteSpecs(specs fs.FS, onGroup func(string, []byte) error) (err error) {
	if types, e := rs.FromSpecs(specs); e != nil {
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
		Loop:
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
					err = errutil.New(e, "in header", groupName)
					break Loop
				} else {
					// registration lists per group
					reg := distill.MakeRegistry(types.Types)

					// all the types in this group:
					for _, key := range typeNames {
						t := types.Types[key]
						if t == nil {
							err = errutil.Fmt("groups in groups need work still ( %q in %q )", key, groupName)
							break Loop
						} else {
							reg.AddType(t)
							//
							if name := specShortName(t); len(name) == 0 {
								err = errutil.New("unknown type", t.Spec.Choice)
								break Loop
							} else if e := tps.ExecuteTemplate(&out, name+".tmpl", t); e != nil {
								err = errutil.New(e, "couldnt process", key)
								break Loop
							}
						}
					}

					if e := writeLists(&out, reg, tps); e != nil {
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

func writeLists(w io.Writer, reg distill.Registry, tps *template.Template) (err error) {
	// sort registration lists ( in place )
	reg.Sort()

	// write registration lists
	if e := tps.ExecuteTemplate(w, "regList.tmpl", map[string]any{
		"Name": "Slots",
		"List": reg.Slots,
		"Type": "interface{}",
	}); e != nil {
		err = errutil.New(e, "couldnt process slots")
	} else if e := tps.ExecuteTemplate(w, "regList.tmpl", map[string]any{
		"Name": "Slats",
		"List": reg.Slats,
		"Type": "composer.Composer",
	}); e != nil {
		err = errutil.New(e, "couldnt process slats")
	} else if e := tps.ExecuteTemplate(w, "sigList.tmpl", reg.Sigs); e != nil {
		err = errutil.New(e, "couldnt process signatures")
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
