package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
)

type Cursor struct {
	Col, Row int // x,y
}

// update the cursor; errors on all control characters except Newline.
func (c *Cursor) NewRune(r rune) charm.State {
	switch {
	case r == Newline:
		c.Row++
		c.Col = 0
	default:
		c.Col++
	}
	return c
}
