package gomake

import (
	"bytes"
	"go/format"
	"io"
	"io/fs"
	"sort"
	"text/template"

	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"github.com/ionous/errutil"
)

// reads all of the files in the passed filesystem as ifspecs
// generate golang structs
func WriteSpecs(ifspecs fs.FS, onGroup func(string, []byte)) (err error) {
	if types, e := rs.FromSpecs(ifspecs); e != nil {
		err = e
	} else {
		ctx := &Context{types: types}
		if tps, e := newTemplates(ctx); e != nil {
			err = e
		} else {
			w := StructWriter{types, tps}
			for _, groupType := range types.Groups {
				groupName := groupType.Name
				ctx.currentGroup = groupName
				groupSpec := groupType.Spec.Value.(*spec.GroupSpec)
				typeNames := make([]string, 0, len(groupSpec.Specs))
				for _, t := range groupSpec.Specs {
					typeNames = append(typeNames, t.Name)
				}
				sort.Strings(typeNames) // in-place
				var out bytes.Buffer
				tps.ExecuteTemplate(&out, "fileHeader.tmpl", map[string]any{
					"Package": groupName,
					"Imports": []string{},
				})
				for _, key := range typeNames {
					blockType := types.Types[key]
					if e := w.WriteStruct(&out, blockType); e != nil {
						err = errutil.Append(err, errutil.New(e, "couldnt process", key))
					}
				}

				// get whatever we can of the output errors or no.
				res := out.Bytes()
				// if the writing worked okay, run "gofmt" on the source
				if err == nil {
					if gofmt, e := format.Source(res); e != nil {
						err = errutil.Append(err, errutil.New(e, "while formating", groupName))
					} else {
						res = gofmt
					}
				}
				onGroup(groupName, res)
			}
		}
	}
	return
}

func (w *StructWriter) WriteStruct(out io.Writer, blockType *spec.TypeSpec) (err error) {
	switch t := blockType.Spec.Choice; t {
	case spec.UsesSpec_Flow_Opt:
		err = w.writeFlow(out, blockType)
	case spec.UsesSpec_Slot_Opt:
		err = w.writeSlot(blockType)
	case spec.UsesSpec_Swap_Opt:
		err = w.writeSwap(blockType)
	case spec.UsesSpec_Num_Opt:
		err = w.writeNum(blockType)
	case spec.UsesSpec_Str_Opt:
		err = w.writeStr(blockType)
	default:
		err = errutil.New("unknown type", blockType.Spec.Choice)
	}
	return
}

type StructWriter struct {
	rs.TypeSpecs
	tps *template.Template
}

func (w *StructWriter) writeFlow(out io.Writer, blockType *spec.TypeSpec) error {
	// todo: groups.
	// each group is going to go to a different out so...
	return w.tps.ExecuteTemplate(out, "flow.tmpl", blockType)
}

func (w *StructWriter) writeSlot(blockType *spec.TypeSpec) error {
	return nil
}

func (w *StructWriter) writeSwap(blockType *spec.TypeSpec) error {
	return nil
}

func (w *StructWriter) writeNum(blockType *spec.TypeSpec) error {
	return nil
}

func (w *StructWriter) writeStr(blockType *spec.TypeSpec) error {
	return nil
}
