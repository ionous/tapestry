package dot

import (
	"errors"
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
)

// Dotted - path operations to access the contents of certain targets.
type Dotted interface {
	Peek(Cursor) (Cursor, error)
	Poke(Cursor, rt.Value) error
	writeTo(b *strings.Builder)
}

// Cursor - generic access to objects, lists, and records.
type Cursor interface {
	CurrentValue() rt.Value
	SetAtIndex(int, rt.Value) error
	GetAtIndex(int) (Cursor, error)
	SetAtField(string, rt.Value) error
	GetAtField(string) (Cursor, error)
}

// the final position in a path where we might want to get or put a value.
// returned by FindEndpoint()
type Endpoint struct {
	parent Cursor // the parent of value; needed for writing the value back
	child  Dotted // the final part of the path; holds the current value
}

// this walks up the tree to write back the final value.
func (last Endpoint) SetValue(newValue rt.Value) (err error) {
	return last.child.Poke(last.parent, newValue)

}
func (last Endpoint) GetValue() (ret rt.Value, err error) {
	if at, e := last.child.Peek(last.parent); e != nil {
		err = e
	} else {
		ret = at.CurrentValue()
	}
	return
}

// expects there's at least one path element
func FindEndpoint(run rt.Runtime, name string, path Path) (ret Endpoint, err error) {
	c := NewTarget(run, name)
	return findEndpoint(c, path)
}

// expects there's at least one path element
func findEndpoint(c Cursor, path Path) (ret Endpoint, err error) {
	if end := len(path) - 1; end < 0 {
		err = errors.New("path is empty")
	} else {
		pos, front, last := c, path[:end], path[end]
		for i, dot := range front {
			if next, e := dot.Peek(pos); e != nil {
				err = fmt.Errorf("%s at %d in %s", e, i, path)
				break
			} else {
				pos = next
			}
		}
		if err == nil {
			ret = Endpoint{pos, last}
		}
	}
	return
}
