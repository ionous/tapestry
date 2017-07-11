package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// Objects with ids, findable by the game.
type Objects struct {
	ObjectMap
	classes *Classes
}
type ObjectMap map[string]*RefObject

// NewObject from the passed class the object is anonymous and cant be found by id.
// Compatible with rt.Runtime.
func (or *Objects) NewObject(class string) (ret rt.Object, err error) {
	if cls, ok := or.classes.GetClass(class); !ok {
		err = errutil.New("no such class", class)
	} else {
		ret = or.newObject(cls.(*RefClass))
	}
	return
}

func (or *Objects) newObject(cls *RefClass) *RefObject {
	rval := r.New(cls.rtype).Elem()
	return &RefObject{"", rval, cls, or}
}

// Emplace wraps the passed value as an anonymous object.
// Compatible with rt.Runtime.
func (or *Objects) Emplace(i interface{}) (ret rt.Object, err error) {
	if rval, e := unique.ValuePtr(i); e != nil {
		err = e
	} else if cls, e := or.classes.GetByType(rval.Type()); e != nil {
		err = e
	} else {
		ret = &RefObject{"", rval, cls, or}
	}
	return
}

// GetObject is compatible with rt.Runtime. The map can also be used directly.
func (or *Objects) GetObject(name string) (ret rt.Object, okay bool) {
	id := id.MakeId(name)
	ret, okay = or.ObjectMap[id]
	return
}

// GetByValue expects a pointer to a value, and it returns the ref object which wraps it.
// WARNING: it can return nil without error
func (or *Objects) GetByValue(rval r.Value) (ret *RefObject, err error) {
	if !rval.IsNil() {
		rval := rval.Elem()
		if id, e := idFromValue(rval); e != nil {
			err = e
		} else if obj, ok := or.ObjectMap[id]; !ok {
			err = errutil.Fmt("object not found '%s'", id)
		} else if obj.rval.Interface() != rval.Interface() {
			err = errutil.Fmt("conflicting objects '%s' %T %T", id, obj, rval)
		} else {
			ret = obj
		}
	}
	return
}

// FIX: this is going to be way too slow for *casual use.
// an mru might of type might help,
// better might be caching the id path in the class,
// best might be forcing all classes to carry an explict id field as their first member.
// good for serialization would be store ids as much as possible.
func IdFromValue(rval r.Value) (ret string, err error) {
	if !rval.IsNil() {
		ret, err = idFromValue(rval.Elem())
	}
	return
}

func idFromValue(rval r.Value) (ret string, err error) {
	rtype := rval.Type()
	if field, ok := FieldPathOfId(rtype); !ok {
		err = errutil.New("couldnt find id for", rtype)
	} else {
		field := rval.FieldByIndex(field.FullPath())
		name := field.String()
		ret = id.MakeId(name)
	}
	return
}

// FieldPathOfId extracts the index of the id field
func FieldPathOfId(rtype r.Type) (ret *unique.FieldInfo, okay bool) {
	for fw := unique.Fields(rtype); fw.HasNext(); {
		field := fw.GetNext()
		tag := unique.Tag(field.Tag)
		if _, ok := tag.Find("id"); ok {
			if field.Type.Kind() == r.String {
				ret, okay = &field, true
			}
			break
		}
	}
	return
}
