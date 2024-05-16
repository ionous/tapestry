package shortcut

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
)

func ReadDots(str string) (ret assign.Address, err error) {
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
	dot    []assign.Dot
}

func (a *Shortcut) Finish(str string) (ret assign.Address) {
	name := literal.T(a.name)
	switch a.marker {
	case VarMarker:
		ret = &assign.VariableDot{
			Name: name,
			Dot:  a.dot,
			Markup: map[string]any{
				"shortcut": str,
			},
		}
	case ObjMarker:
		ret = &assign.VariableDot{
			Name: name,
			Dot:  a.dot,
			Markup: map[string]any{
				"shortcut": str,
			},
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
			a.dot = append(a.dot, &assign.AtField{
				Field: literal.T(name),
			})
		}

	case t.IsNumber():
		if idx := v.(int); len(a.name) == 0 {
			err = errors.New("expected an object or variable name")
		} else {
			a.dot = append(a.dot, &assign.AtIndex{
				Index: literal.I(idx),
			})
		}

	default:
		err = fmt.Errorf("unexpected token %s", t)
	}
	return
}
