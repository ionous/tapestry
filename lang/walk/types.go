package walk

import r "reflect"

type Type int

//go:generate stringer -type=Type
const (
	Flow Type = iota
	Slot
	Swap
	Str   // an enum
	Num   // rare, tbd: remove. one of a set of fixed set of numbers
	Value // a field of a flow
)

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

// determine if the field is
// a flow, a slot, a swap, a str, a num, generic value;
// or a slice of one of those things.
func fieldType(t r.Type) (ret Type) {
	switch t.Kind() {
	default:
		ret = Value
	case r.Interface:
		ret = Slot
	case r.Struct:
		ret = structType(t)
	case r.Slice:
		ret = sliceType(t.Elem())
	}
	return
}

func sliceType(elType r.Type) (ret Type) {
	switch k := elType.Kind(); k {
	default:
		ret = Value
	case r.Interface:
		ret = Slot
	case r.Struct: // ex. CallPattern { Arguments  []Arg }
		ret = structType(elType)
	}
	return
}

// deduce the command type from the passed struct definition
// relies on the type ( rather than composer, etc. )
// to make using hand created types easier
func structType(t r.Type) (ret Type) {
	ret = Flow // provisionally
	switch t.NumField() {
	case 1:
		switch f := t.Field(0); f.Name {
		case "Num":
			ret = Num
		case "Str":
			ret = Str
		}
	case 2:
		switch f := t.Field(0); f.Name {
		case "Choice":
			ret = Swap
		}
	}
	return
}
