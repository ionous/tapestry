package scope

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

type Stack []rt.Scope

var MaxStack = 25

func (k *Stack) PushScope(scope rt.Scope) {
	if len((*k)) > MaxStack {
		panic("scope overflow")
	}
	(*k) = append((*k), scope)
}

func (k *Stack) PopScope() {
	if cnt := len((*k)); cnt == 0 {
		panic("Stack: popping an empty scopes")
	} else {
		(*k) = (*k)[0 : cnt-1]
	}
}

func (k *Stack) ReplaceScope(scope rt.Scope) (ret rt.Scope) {
	// that's okay; its probably from a replace of an empty stack
	if scope == nil {
		k.PopScope()
	} else if end := len(*k) - 1; end < 0 {
		(*k) = append((*k), scope)
	} else {
		ret, (*k)[end] = (*k)[end], scope
	}
	return
}

// FieldByName returns the value at 'field'
func (k *Stack) FieldByName(field string) (ret g.Value, err error) {
	if cnt := len((*k)); cnt == 0 {
		err = g.UnknownVariable(field)
	} else {
		for i := cnt - 1; i >= 0; i-- {
			ret, err = (*k)[i].FieldByName(field)
			if _, isUnknown := err.(g.Unknown); !isUnknown {
				break // while isUnknown keep going; otherwise done.
			}
		}
	}
	return
}

// SetFieldByName writes the value of 'v' into the value at 'field'.
func (k *Stack) SetFieldByName(field string, v g.Value) (err error) {
	if cnt := len((*k)); cnt == 0 {
		err = g.UnknownVariable(field)
	} else {
		for i := cnt - 1; i >= 0; i-- {
			err = (*k)[i].SetFieldByName(field, v)
			if _, isUnknown := err.(g.Unknown); !isUnknown {
				break // while isUnknown keep going; otherwise done.
			}
		}
	}
	return
}
