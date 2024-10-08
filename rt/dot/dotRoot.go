package dot

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
)

type rootDot struct {
	run  rt.Runtime
	name string // either the name of an object, or meta.Variables
}

// print friendly string
func (dot rootDot) String() string {
	var prefix string
	if dot.name == meta.Variables {
		prefix = "@"
	} else {
		prefix = "#"
	}
	return prefix + string(dot.name)
}

func (dot rootDot) Peek(c Cursor) (ret Cursor, err error) {
	if c != nil {
		err = errors.New("unexpected cursor at root")
	} else {
		copy := dot
		ret = rootCursor(copy)
	}
	return
}

func (dot rootDot) Poke(c Cursor, newValue rt.Value) error {
	return errors.New("can't write a value without a dotted path")
}
