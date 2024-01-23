package inspect

import (
	"log"
	r "reflect"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// Provides iteration of tapestry commands.
// The returned Iter starts pointing to the passed container,
// calling Next advances to each element of the container in turn.
// Containers can be a single command, a slot for a command, or slices of commands or slots;
// other values have undefined results and may panic.
func Walk(i typeinfo.Inspector) Iter {
	v := r.ValueOf(i)
	t := i.Inspect()
	switch t.(type) {
	case *typeinfo.Flow:
		// the autogenerated info for slices are slices, and slices are valid containers.
		// each flow implements TypeInfo(), and its *struct* is a valid container.
		// ( that matches flow members of another flow )
		if k := v.Kind(); k != r.Slice {
			v = v.Elem()
		}
	case *typeinfo.Slot:
		// the autogenerated info for slices are slices, and slices are valid containers.
		// this wants the interface interface value embedded in the autogenerated struct
		// ( that matches slot members of a flow )
		if k := v.Kind(); k == r.Ptr {
			v = v.Elem().Field(0)
		}
	case *typeinfo.Str, *typeinfo.Num:
		panic("not a valid container type")
	}
	return Iter{curr: v, currType: t}
}

// If a command, a slot for a command, and slices of commands or slots,
// can all be thought of as containers of values,
// where values are primitive values or other containers,
// then Iter provides a depth first traversal of those values.
// In this context, a slice of primitives values is handled treated as a single value.
type Iter struct {
	curr r.Value
	// track the type manually because slice and slot members dont have typeinf
	currType typeinfo.T
	index    int // *next* index in a slice, or field in a struct
}

// advance the focus within the current collection.
func (w *Iter) Next() (okay bool) {
	if cnt := w.Len(); w.index < cnt {
		w.index++
		okay = true
	}
	return
}

// returns a new walker for the currently focused element;
// panics if the focused element isn't a valid container.
func (w *Iter) Walk() (ret Iter) {
	if at := w.index - 1; at < 0 {
		log.Printf("container of %s %T hasn't advanced", w.curr.Type(), w.currType)
		panic("invalid index")
	} else if v, t := w.getFocus(); !v.IsValid() {
		log.Printf("container of %s %T invalid at %d", w.curr.Type(), w.currType, w.index)
		panic("can't walk an invalid value")
	} else {
		ret = Iter{curr: v, currType: t}
	}
	return
}

// the autogenerated description of the currently focused element.
func (w *Iter) TypeInfo() typeinfo.T {
	_, t := w.getFocus()
	return t
}

// metadata for the currently focused element.
// only valid flow; returns nil otherwise
func (w *Iter) Markup() (ret map[string]any) {
	v, _ := w.getFocus() // ugly.
	if m, ok := v.Addr().Interface().(typeinfo.FlowInspector); ok {
		ret = m.GetMarkup(true)
	}
	return
}

// only valid for the members of a flow; panics otherwise
func (w *Iter) Term() (ret typeinfo.Term) {
	if p, ok := w.currType.(*typeinfo.Flow); !ok || w.index == 0 {
		log.Printf("container of %T at index: %d", w.currType, w.index)
		panic("terms only make sense for structs")
	} else {
		ret = p.Terms[w.index-1]
	}
	return
}

// returns the value of the current focus as used by go.
// enums use $STRING_KEY, while bool uses true/false.
// falls back to the container itself if Next() has yet to be called.
func (w *Iter) ZeroValue() (okay bool) {
	v, t := w.getFocus()
	// temp: unpack Str struct;
	if _, ok := t.(*typeinfo.Str); ok && v.Kind() == r.Struct {
		str := v.Field(0).String() // Str holds $KEY in memory.
		// fix? the optional marshal treated zero value as "nothing specified"
		// not, as you might expect, the zeroth enum value.
		okay = len(str) == 0
	} else {
		// fix: $TRUE is listed as the zeroth value. while this treats false as zero.
		// $FALSE seems like it should be the zeroth value.
		okay = v.IsZero()
	}
	return
}

// returns the value of the current focus as used by go.
// enums use $STRING_KEY, while bool uses true/false.
// falls back to the container itself if Next() has yet to be called.
func (w *Iter) GoValue() (ret any) {
	v, t := w.getFocus()
	// temp: unpack Str struct
	if _, ok := t.(*typeinfo.Str); ok && v.Kind() == r.Struct {
		ret = v.Field(0).String() // Str enums holds $KEY in memory.
	} else {
		ret = v.Interface()
	}
	return
}

// returns the value of the current focus as it would appear in file.
// enums use lowercase strings, while bool uses true/false.
// falls back to the container itself if Next() has yet to be called.
func (w *Iter) CompactValue() (ret any) {
	v, t := w.getFocus()
	// temp: unpack Str struct
	if t, ok := t.(*typeinfo.Str); ok && v.Kind() == r.Struct {
		str := v.Field(0).String()
		if opt := t.Options; len(opt) > 0 {
			ret = strings.ToLower(str[1:])
		} else {
			ret = str
		}
	} else {
		ret = v.Interface()
	}
	return
}

// returns the value of the current focus as described by the spec.
// for example: enums, including boolean values, use $STRING_KEY format.
// falls back to the container itself if Next() has yet to be called.
func (w *Iter) NormalizedValue() (ret any) {
	v, t := w.getFocus()
	// temp: unpack Str struct
	if _, ok := t.(*typeinfo.Str); !ok {
		ret = v.Interface()
	} else {
		switch v.Kind() {
		case r.Struct:
			ret = v.Field(0).Interface()
		case r.Bool:
			if v.Bool() {
				ret = "$TRUE"
			} else {
				ret = "$FALSE"
			}
		default:
			ret = v.String()
		}
	}
	return
}

// write a value into the target of an iterator.
// ( SetSlot can't be on the slot itself since the slot is often a bare member )
func (w *Iter) SetSlot(val typeinfo.Inspector) (okay bool) {
	v, _ := w.getFocus()
	newVal := r.ValueOf(v)
	if newVal.Type().AssignableTo(v.Type()) {
		v.Set(newVal)
		okay = true
	}
	return
}

// write a value into the target of an iterator.
// returns false if the value is incompatible
// ( uses go rules of conversion when needed to complete the assignment )
// func (w *Iter) SetValue(val any) (okay bool) {
// 	if out, val := w.Value(), r.ValueOf(val); out.Kind() == val.Kind() {
// 		out.Set(val)
// 		okay = true
// 	} else if t := out.Type(); val.CanConvert(t) {
// 		out.Set(val.Convert(t))
// 	}
// 	return
// }

// change the size of the container
func (w *Iter) Resize(cnt int) {
	if w.curr.Kind() != r.Slice {
		log.Printf("can't resize a %s(%s)", w.curr.Kind(), w.curr.Type())
	}
	w.curr.Grow(cnt)
	w.curr.SetLen(cnt)
}

// Returns the number of repeated elements in the current container;
// doesn't change over the course of iteration.
// Filled slots have one element; empty slots zero elements.
func (w *Iter) Len() (ret int) {
	switch w.curr.Kind() {
	default:
		panic("can't measure the length of primitive values")
	case r.Slice:
		ret = w.curr.Len()
	case r.Struct:
		p := w.currType.(*typeinfo.Flow)
		ret = len(p.Terms)
	case r.Interface:
		if !w.curr.IsNil() {
			ret = 1
		}
	}
	return
}

func (w *Iter) getFocus() (rv r.Value, rt typeinfo.T) {
	if at := w.index - 1; at < 0 {
		rv, rt = w.curr, w.currType
	} else {
		switch k := w.curr.Kind(); k {
		case r.Slice:
			rv, rt = w.curr.Index(at), w.currType
		case r.Struct:
			p := w.currType.(*typeinfo.Flow)
			rv, rt = w.curr.Field(at), p.Terms[at].Type
		case r.Interface:
			// slots are filled with flows
			// the first Elem() gets the pointer, the second the struct.
			ptr := w.curr.Elem()
			inspect := ptr.Interface().(typeinfo.Inspector)
			rv, rt = ptr.Elem(), inspect.Inspect().TypeInfo()
		}
	}
	return
}
