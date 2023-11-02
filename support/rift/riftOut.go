package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

type History struct {
	els []moment
}

// indent of the most recent history, or -1 if history is empty
func (h *History) CurrentIndent() (ret int) {
	if cnt := len(h.els); cnt == 0 {
		ret = -1
	} else {
		top := h.els[cnt-1]
		ret = top.indent
	}
	return
}

func (h *History) PushIndent(indent int, c charm.State, pop func() error) (ret charm.State) {
	// fix? there are some +1s that probably need adding for colon and dash
	// so this accepts equal indents right now
	if curr := h.CurrentIndent(); indent < curr {
		e := errutil.New("indents should grow")
		ret = charm.Error(e)
	} else {
		h.els = append(h.els, moment{indent, c, pop})
		ret = c
	}
	return
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
