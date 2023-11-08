package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
)

type History struct {
	els []moment
}

type moment struct {
	indent int
	state  charm.State
	popfn  func() error
}

func (h *History) Push(indent int, c charm.State) charm.State {
	return h.PushCallback(indent, c, func() (_ error) { return })
}

func (h *History) PushCallback(indent int, c charm.State, pop func() error) charm.State {
	h.els = append(h.els, moment{indent, c, pop})
	return c
}

func (h *History) PopAll() (err error) {
	_, err = h.pop(-1)
	return
}

// remove elements off the stack until a matching indent is found
func (h *History) Pop(i int) (ret charm.State) {
	if n, e := h.pop(i); e != nil {
		ret = charm.Error(e)
	} else {
		ret = n
	}
	return
}

func (h *History) pop(i int) (ret charm.State, err error) {
	for len(h.els) > 0 {
		end := len(h.els) - 1
		if top := h.els[end]; top.indent == i {
			ret = top.state
			break
		} else if e := top.popfn(); e != nil {
			err = e
			break
		} else {
			h.els = h.els[:end]
		}
	}
	return
}
