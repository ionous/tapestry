package dot

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
)

type Cursor = rt.Cursor

func MakeReference(run rt.Runtime, name string) Reference {
	root := rootDot{run, name}
	return Reference{child: root}
}

func MakeReferenceValue(run rt.Runtime, rec rt.Value) Reference {
	root := recordDot{run, rec}
	return Reference{child: root}
}

// the final position in a path where we might want to get or put a value.
// implements rt.Reference
type Reference struct {
	pos   Cursor    // the pos of value; needed for writing the value back
	child rt.Dotted // the final part of the path; holds the current value
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
func (at Reference) Dot(next rt.Dotted) (ret rt.Reference, err error) {
	if pos, e := at.child.Peek(at.pos); e != nil {
		err = e
	} else {
		ret = Reference{pos, next}
	}
	return
}

// step into the current value multiple times
func Path(pos rt.Reference, path []rt.Dotted) (ret rt.Reference, err error) {
	for i, dot := range path {
		if next, e := pos.Dot(dot); e != nil {
			err = fmt.Errorf("%s at %d in %s", e, i, path)
			break
		} else {
			pos = next
		}
	}
	if err == nil {
		ret = pos
	}
	return
}
