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
func WriteSpecs(realout io.Writer, files fs.FS) (err error) {
	if types, e := rs.FromSpecs(files); e != nil {
		err = e
	} else {
		var out bytes.Buffer
		ctx := &Context{types: types}
		if tps, e := newTemplates(ctx); e != nil {
			err = e
		} else {
			w := StructWriter{types, tps}
			for _, groupType := range types.Groups {
				groupSpec := groupType.Spec.Value.(*spec.GroupSpec)
				typeNames := make([]string, 0, len(groupSpec.Specs))
				for _, t := range groupSpec.Specs {
					typeNames = append(typeNames, t.Name)
				}
				sort.Strings(typeNames) // in-place
				ctx.currentGroup = groupType.Name
				for _, key := range typeNames {
					blockType := types.Types[key]
					if e := w.WriteStruct(&out, blockType); e != nil {
						err = errutil.Append(err, errutil.New(e, "couldnt process", key))
					}
				}
			}
		}
		// run "gofmt" on the source
		if me, e := format.Source(out.Bytes()); e != nil {
			err = e
		} else {
			realout.Write(me)
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
