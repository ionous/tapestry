package composer

import (
	"github.com/ionous/errutil"

	r "reflect"
)

// note: if it were ever to turn out this is the only reflection needed
// could replace with function callbacks that new() rather than just a generic r.Type
type TypeRegistry map[string]r.Type

func (reg TypeRegistry) NewType(typeName string) (ret interface{}, okay bool) {
	if rtype, ok := reg[typeName]; ok {
		ret = r.New(rtype).Interface()
		okay = true
	}
	return
}

func (reg *TypeRegistry) RegisterTypes(cmds []Composer) (err error) {
	if *(reg) == nil {
		*(reg) = make(TypeRegistry)
	}
	for _, cmd := range cmds {
		if spec := cmd.Compose(); len(spec.Name) == 0 {
			e := errutil.Fmt("Missing type name %T", cmd)
			errutil.Append(err, e)
		} else if was, exists := (*reg)[spec.Name]; exists {
			e := errutil.Fmt("Duplicate type name %q now: %T, was: %s", spec.Name, cmd, was.String())
			errutil.Append(err, e)
			break
		} else {
			(*reg)[spec.Name] = r.TypeOf(cmd).Elem()
		}
	}
	return
}
