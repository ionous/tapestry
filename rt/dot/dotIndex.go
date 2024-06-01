package dot

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/rt"
)

// zero-based index to pick a value from a list.
type Index int

// raw int
func (dot Index) Index() int {
	return int(dot)
}

// print friendly string
func (dot Index) String() string {
	return "[" + strconv.Itoa(dot.Index()) + "]"
}

func (dot Index) Peek(c Cursor) (Cursor, error) {
	return c.GetAtIndex(int(dot))
}

func (dot Index) Poke(c Cursor, newValue rt.Value) error {
	return c.SetAtIndex(int(dot), newValue)
}
