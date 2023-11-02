package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

type History struct {
	els []moment
}

func (h *History) PushIndent(indent int, c charm.State, pop func() error) charm.State {
	if cnt := len(h.els); cnt > 0 {
		if prev := h.els[cnt-1].indent; prev > indent {
			panic("indents should grow")
		}
	}
	h.els = append(h.els, moment{indent, c, pop})
	return c
}

func (h *History) PopAll() (err error) {
	_, err = h.popIndent(-1)
	return
}

// pop to reach the passed indent
func (h *History) PopIndent(indent int) (ret charm.State) {
	if next, e := h.popIndent(indent); e != nil {
		ret = charm.Error(e)
	} else {
		ret = next
	}
	return
}

// pop to reach the passed indent
func (h *History) popIndent(indent int) (ret charm.State, err error) {
	for len(h.els) > 0 {
		top := len(h.els) - 1
		bell, slice := h.els[top], h.els[:top]
		if indent > bell.indent {
			err = errutil.New("unexpected indent")
			break
		} else if bell.indent == indent {
			ret = bell.state
			break
		} else if e := bell.popfn(); e != nil {
			err = e
			break
		} else {
			h.els = slice
		}
	}
	return
}

type moment struct {
	indent int
	state  charm.State
	popfn  func() error
}

func (m *moment) pop() (ret int, err error) {
	return m.indent, m.popfn()
}
