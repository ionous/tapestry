package scope

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// a linked list of scopes:
// if a variable isnt found in the first scope, the next is searched, and so on.
type Chain struct {
	rt.Scope
	next  *Chain
	depth int // counts the number of pushed scopes; not accurate if the scopes contain other scopes.
}

const MaxDepth = 25

func MakeChain(top rt.Scope) Chain {
	return Chain{Scope: top}
}

// rewrite the current scope as the passed scope
// returns the previous scope
// ( used for function calls; they "hide" all variables from the caller stack )
func (s *Chain) ReplaceScope(top rt.Scope) Chain {
	prev, next := *s, MakeChain(top)
	*s = next
	return prev
}

func (s *Chain) RestoreScope(c Chain) {
	*s = c
}

func (s *Chain) PushScope(top rt.Scope) {
	if s.depth > MaxDepth {
		panic("scope overflow")
	}
	clone := *s // save the values before they are overwritten
	*s = Chain{Scope: top, next: &clone, depth: s.depth + 1}
}

func (s *Chain) PopScope() {
	if s.next == nil {
		panic("scope underflow")
	}
	*s = *s.next
}

// FieldByName returns the value at 'field'
func (s *Chain) FieldByName(field string) (ret g.Value, err error) {
	var found bool
	for top := s; top != nil; top = top.next {
		if v, e := top.Scope.FieldByName(field); !g.IsUnknown(e) {
			ret, err, found = v, e, true
			break
		}
	}
	if !found {
		err = g.UnknownVariable(field)
	}
	return
}

// SetFieldByName writes the value of 'v' into the value at 'field'.
func (s *Chain) SetFieldByName(field string, v g.Value) (err error) {
	var found bool
	for top := s; top != nil; top = top.next {
		if e := top.Scope.SetFieldByName(field, v); !g.IsUnknown(e) {
			err, found = e, true
			break
		}
	}
	if !found {
		err = g.UnknownVariable(field)
	}
	return
}

// SetFieldDirty - tell the current scope the named value has changed.
func (s *Chain) SetFieldDirty(field string) (err error) {
	var found bool
	for top := s; top != nil; top = top.next {
		if e := top.Scope.SetFieldDirty(field); !g.IsUnknown(e) {
			err, found = e, true
			break
		}
	}
	if !found {
		err = g.UnknownVariable(field)
	}
	return
}
