package walk

import (
	r "reflect"

	"git.sr.ht/~ionous/tapestry/support/tag"
)

// info about a member of a flow
type Field struct {
	fieldType r.Type
	tag       tag.StructTag // cached tag of the current field
}

func (f *Field) Type() r.Type {
	return f.fieldType
}

func (f *Field) Optional() bool {
	return f.tag.Exists("optional")
}

func (f *Field) Internal() bool {
	return f.tag.Exists("internal")
}

// true if the container is a slice of commands.
func (f *Field) Repeats() (okay bool) {
	k := f.fieldType.Kind()
	if k == r.Slice {
		t := sliceType(f.fieldType.Elem())
		okay = t != Value
	}
	return
}

func (f *Field) Label() (ret string, okay bool) {
	if label, ok := f.tag.Find("label"); ok {
		if label != "_" {
			ret = label
		}
		okay = true
	}
	return
}
