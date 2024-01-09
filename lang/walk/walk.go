package walk

import (
	"log"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/support/tag"
)

// Provides depth first iteration of tapestry commands;
// requires a call to Next() to advance to the first element of the passed container.
// Containers can be a single command, a slot for a command, or slices of commands or slots;
// other values have undefined results and may panic.
// The passed value should the instance value, not a pointer to one.
func Walk(container r.Value) Walker {
	return Walker{curr: container}
}

// If a command, a slot for a command, and slices of commands or slots,
// can all be thought of as containers of values,
// where values are primitive values or other containers,
// then Walker provides a depth first traversal of those values.
// In this context, a slice of primitives values is handled treated as a single value.
type Walker struct {
	curr, focus r.Value
	index       int // *next* index in a slice, or field in a struct
}

// Returns the number of repeated elements in the current container;
// doesn't change over the course of iteration.
// Filled slots have one element; empty slots zero elements.
func (w *Walker) Len() (ret int) {
	switch w.curr.Kind() {
	default:
		panic("can't measure the length of primitive values")
	case r.Interface:
		if !w.curr.IsNil() {
			ret = 1
		}
	case r.Struct:
		if cnt := w.curr.NumField(); cnt > 0 {
			if w.curr.Field(cnt-1).Kind() == r.Map {
				cnt--
			}
			ret = cnt
		}
	case r.Slice:
		if sliceType(w.curr.Type().Elem()) == Value {
			panic("doesn't measure the length of a primitive array")
		} else {
			ret = w.curr.Len()
		}
	}
	return
}

// only valid for the members of a flow; panics otherwise
// can return invalid if the markup field doesnt exist in the flow (ex. for tests)
func (w *Walker) Markup() (ret r.Value) {
	if cnt := w.curr.NumField(); cnt > 0 {
		if v := w.curr.Field(cnt - 1); v.Kind() == r.Map {
			ret = v
		}
	}
	return
}

// only valid for the members of a flow; panics otherwise
func (w *Walker) Field() Field {
	if k := w.curr.Kind(); k != r.Struct || w.index == 0 {
		panic("fields only make sense for structs")
	}
	containerType := w.curr.Type()
	field := containerType.Field(w.index - 1)
	tag := tag.ReadTag(field.Tag)
	return Field{field.Name, field.Type, tag}
}

// returns the "focus" of the current iteration.
// falls back to the container itself if next() has yet to be called.
func (w *Walker) Value() (ret r.Value) {
	if !w.focus.IsValid() {
		ret = w.curr
	} else {
		switch typeOf(w.focus.Type()) {
		case Str, Num:
			ret = w.focus.Field(0)
		default:
			ret = w.focus
		}
	}
	return
}

// returns a new walker for the currently focused element;
// panics if the element isn't a valid container.
// the returned iterator points to the container
// and requires a Next() to advance to the first element.
func (w *Walker) Walk() (ret Walker) {
	switch v := w.focus; typeOf(v.Type()) {
	case Flow, Slot, Swap:
		ret = Walker{curr: v}
	default:
		log.Printf("trying to walk a %s a %s", w.focus.Type(), w.focus.Kind())
		panic("can't descend into primitive values")
	}
	return
}

// advance the focus within the current collection.
func (w *Walker) Next() (okay bool) {
	switch curr := w.curr; curr.Kind() {
	case r.Interface:
		okay = w.step(func(int) r.Value {
			// first get the element underlying the interface
			// which is always a pointer, then get the value ( struct ) at the pointer.
			// this callback only happens if there was a valid slot ( index 0 of slotLen 1 )
			return w.curr.Elem().Elem()
		})
	case r.Slice:
		// to be here, we must be unpacking a slice of slots or flow
		// we would have already returned value *as* the slice
		okay = w.step(curr.Index)

	case r.Struct:
		if t := structType(curr.Type()); t == Flow {
			// embedded flow
			// zero index, the state before the first call to Next(), indicates the flow itself
			// on first Next() we want to return the zeroth field; so "index" is used directly as the index.
			okay = w.step(curr.Field)
		} else if t == Swap {
			okay = w.step(func(int) r.Value {
				// first get the element underlying the interface
				// which is always a pointer, then get the value ( struct ) at the pointer.
				// this callback only happens if there was a valid slot ( index 0 of slotLen 1 )
				return w.curr.Field(1).Elem().Elem()
			})

		}

	default:
		// ex. Array, Uintptr, Complex64, Complex128, Chan, Func,
		// Map, Pointer, UnsafePointer, Invalid
		log.Printf("unexpected %s(%s) in generated types", curr.Kind(), curr.Type())
		panic("unexpected generated type")
	}
	return
}

// shouldnt need to be public because callers initiate the traversal.
func typeOf(curr r.Type) (ret Type) {
	switch curr.Kind() {
	default:
		ret = Value // roughly
	case r.Interface:
		ret = Slot
	case r.Struct:
		ret = structType(curr)
	case r.Slice:
		ret = sliceType(curr.Elem())
	}
	return
}

func (w *Walker) step(get func(int) r.Value) (okay bool) {
	if at := w.index; at < w.Len() {
		nextField := get(at)
		w.focus, w.index = nextField, at+1
		okay = true
	}
	return
}
