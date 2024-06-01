package shortcut

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
)

// should ReadDots try to keep the shortcut string
// ( ex. for debugging )
var KeepShortcuts bool

func ReadDots(str string) (ret rtti.Address, err error) {
	var a Shortcut
	if e := Tokenize(str, &a); e != nil {
		err = e
	} else {
		ret = a.Finish(str)
	}
	return
}

type Shortcut struct {
	marker Type
	name   string
	dot    []object.Dot
}

func (a *Shortcut) Finish(str string) (ret rtti.Address) {
	name := literal.T(a.name)
	var m map[string]any
	if KeepShortcuts {
		m = map[string]any{
			"shortcut": str,
		}
	}
	switch a.marker {
	case VarMarker:
		ret = &object.VariableDot{
			Name:   name,
			Dot:    a.dot,
			Markup: m,
		}
	case ObjMarker:
		ret = &object.ObjectDot{
			Name:   name,
			Dot:    a.dot,
			Markup: m,
		}
	default:
		panic("unexpected error")
	}
	return
}

func (a *Shortcut) Decoded(t Type, v any) (err error) {
	switch {
	case t.IsMarker():
		if a.marker != 0 {
			err = errors.New("unexpected marker")
		} else {
			a.marker = t
		}

	case t.IsName():
		if a.marker == 0 {
			err = errors.New("expected an object or variable marker")
		} else if name := v.(string); len(a.name) == 0 {
			a.name = name
		} else {
			a.dot = append(a.dot, &object.AtField{
				Field: literal.T(name),
			})
		}

	case t.IsNumber():
		if idx := v.(int); len(a.name) == 0 {
			err = errors.New("expected an object or variable name")
		} else if idx == 0 {
			err = errors.New("zero is an invalid element")
		} else {
			a.dot = append(a.dot, &object.AtIndex{
				Index: literal.I(idx - 1),
			})
		}

	default:
		err = fmt.Errorf("unexpected token %s", t)
	}
	return
}
