package scope

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

type Stack []rt.Scope

var MaxStack = 25

func (k *Stack) size() int {
	return len((*k))
}

func (k *Stack) PushScope(scope rt.Scope) {
	if k.size() > MaxStack {
		panic("scope overflow")
	}
	(*k) = append((*k), scope)
}

func (k *Stack) PopScope() {
	if cnt := k.size(); cnt == 0 {
		panic("Stack: popping empty scope stack")
	} else {
		(*k) = (*k)[0 : cnt-1]
	}
}

func (k *Stack) ReplaceScope(scope rt.Scope) (ret rt.Scope) {
	// that's okay; its probably from a replace of an empty stack
	if scope == nil {
		k.PopScope()
	} else if end := k.size() - 1; end < 0 {
		(*k) = append((*k), scope)
	} else {
		ret, (*k)[end] = (*k)[end], scope
	}
	return
}

// FieldByName returns the value at 'field'
func (k *Stack) FieldByName(field string) (ret g.Value, err error) {
	i := k.size() - 1
	for ; i >= 0; i-- {
		if ret, err = (*k)[i].FieldByName(field); !g.IsUnknown(err) {
			break // while isUnknown keep going; otherwise done.
		}
	}
	if i < 0 {
		err = g.UnknownVariable(field)
	}
	return
}

// SetFieldByName writes the value of 'v' into the value at 'field'.
func (k *Stack) SetFieldByName(field string, v g.Value) (err error) {
	i := k.size() - 1
	for ; i >= 0; i-- {
		if err = (*k)[i].SetFieldByName(field, v); !g.IsUnknown(err) {
			break // while isUnknown keep going; otherwise done.
		}
	}
	if i < 0 {
		err = g.UnknownVariable(field)
	}
	return
}
