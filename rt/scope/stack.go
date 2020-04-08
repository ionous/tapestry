package scope

import "github.com/ionous/iffy/rt"

type ScopeStack struct {
	stack []rt.VariableScope
}

func (k *ScopeStack) PushScope(scope rt.VariableScope) {
	if len(k.stack) > 25 {
		panic("stack overflow")
	}
	k.stack = append(k.stack, scope)
}

func (k *ScopeStack) PopScope() {
	if cnt := len(k.stack); cnt == 0 {
		panic("ScopeStack: popping an empty stack")
	} else {
		k.stack = k.stack[0 : cnt-1]
	}
}

// GetVariable returns the value at 'name'
func (k *ScopeStack) GetVariable(name string) (ret interface{}, err error) {
	err = k.visit(name, func(scope rt.VariableScope) (err error) {
		if v, e := scope.GetVariable(name); e != nil {
			err = e
		} else {
			ret = v
		}
		return err
	})
	return
}

// SetVariable writes the value of 'v' into the value at 'name'.
func (k *ScopeStack) SetVariable(name string, v interface{}) (err error) {
	return k.visit(name, func(scope rt.VariableScope) error {
		return scope.SetVariable(name, v)
	})
}

func (k *ScopeStack) visit(name string, visitor func(rt.VariableScope) error) (err error) {
	for i := len(k.stack) - 1; i >= 0; i-- {
		switch e := visitor(k.stack[i]); e.(type) {
		case nil:
			// no error? we're done.
			goto Done
		case UnknownVariable:
			// didn't find? keep looking...
		default:
			// other error? done.
			err = e
			goto Done
		}
	}
	err = UnknownVariable(name)
Done:
	return
}