package dot

import (
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

func MakeReference(run rt.Runtime, name string) Reference {
	root := rootDot{run, name}
	return Reference{child: root}
}

// the final position in a path where we might want to get or put a value.
type Reference struct {
	pos   Cursor // the pos of value; needed for writing the value back
	child Dotted // the final part of the path; holds the current value
}

// write a value
func (at Reference) SetValue(newValue rt.Value) (err error) {
	return at.child.Poke(at.pos, newValue)
}

// read a value
func (at Reference) GetValue() (ret rt.Value, err error) {
	if at, e := at.child.Peek(at.pos); e != nil {
		err = e
	} else {
		ret = at.CurrentValue()
	}
	return
}

// step into the current value
func (at Reference) Dot(next Dotted) (ret Reference, err error) {
	if pos, e := at.child.Peek(at.pos); e != nil {
		err = e
	} else {
		ret = Reference{pos, next}
	}
	return
}

// step into the current value multiple times
func (at Reference) DotPath(path Path) (ret Reference, err error) {
	for i, dot := range path {
		if next, e := at.Dot(dot); e != nil {
			err = fmt.Errorf("%s at %d in %s", e, i, path)
			break
		} else {
			at = next
		}
	}
	if err == nil {
		ret = at
	}
	return
}
