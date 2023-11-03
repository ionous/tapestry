package rift

import "git.sr.ht/~ionous/tapestry/support/charm"

type Document struct {
	History
	Cursor
}

type Cursor struct {
	Row, Col, Indent int // y, x, tab
	indenting        bool
}

func (doc *Document) Parse(str string, c charm.State) (err error) {
	if e := charm.Parse(str, CountSpaces(&doc.Cursor, c)); e != nil {
		err = e
	} else if e := doc.PopAll(); e != nil {
		err = e
	}
	return
}

func CountSpaces(c *Cursor, next charm.State) charm.State {
	return charm.Self("counting", func(self charm.State, r rune) (ret charm.State) {
		switch r {
		case Newline:
			c.Row++
			c.Col = 0
			c.indenting = true

		default:
			c.Col++
			if r != Space {
				c.indenting = false
			} else if c.indenting {
				c.Indent++
			}
		}
		next = next.NewRune(r)
		switch next.(type) {
		case nil, charm.Terminal:
			ret = next
		default:
			ret = self
		}
		return
	})
}
