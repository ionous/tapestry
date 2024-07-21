package dot

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/rt"
)

type recordDot struct {
	ks  rt.Kinds
	val rt.Value
}

func (dot recordDot) Peek(c Cursor) (ret Cursor, err error) {
	if c != nil {
		err = errors.New("unexpected cursor at root")
	} else {
		ret = MakeValueCursor(dot.ks, dot.val)
	}
	return
}

func (dot recordDot) Poke(c Cursor, newValue rt.Value) error {
	return errors.New("can't write a value without a dotted path")
}
