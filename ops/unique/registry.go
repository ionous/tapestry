package unique

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	r "reflect"
)

type Types map[string]r.Type

// RegisterBlock registers a structure containing pointers to commands.
func (reg Types) Find(name string) (r.Type, bool) {
	id := id.MakeId(name)
	rtype, ok := reg[id]
	return rtype, ok
}

// RegisterBlock registers a structure containing pointers to commands.
func (reg Types) RegisterBlock(block interface{}) (err error) {
	if blockType := r.TypeOf(block); blockType.Kind() != r.Ptr {
		err = errutil.New("expected (nil) pointer (to a struct).")
	} else if structType := blockType.Elem(); structType.Kind() != r.Struct {
		err = errutil.New("expected a struct pointer.")
	} else {
		for i, cnt := 0, structType.NumField(); i < cnt; i++ {
			field := structType.Field(i)
			if e := reg.registerType(field.Type); e != nil {
				err = errutil.New(field.Name, e)
				break
			}
		}
	}
	return
}

// RegisterType registers a single pointer to a command.
func (reg Types) RegisterType(cmd interface{}) (err error) {
	if e := reg.registerType(r.TypeOf(cmd)); e != nil {
		err = errutil.New("command", e)
	}
	return
}

// rtype should be a struct ptr.
func (reg Types) registerType(cmdType r.Type) (err error) {
	if ptrType := cmdType; ptrType.Kind() != r.Ptr {
		err = errutil.New("expected (nil) pointer (to a struct).")
	} else if rtype := ptrType.Elem(); rtype.Kind() != r.Struct {
		err = errutil.New("expected a struct pointer.")
	} else {
		id := id.MakeId(rtype.Name())
		// println("regsiter", id)
		if was, exists := reg[id]; exists && was != rtype {
			err = errutil.New("has conflicting names, id:", id, "was:", was, "type:", rtype)
		} else {
			reg[id] = rtype
		}
	}
	return
}
