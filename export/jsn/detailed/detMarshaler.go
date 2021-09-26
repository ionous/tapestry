package detailed

import (
	"github.com/ionous/errutil"
)

type DetailedMarshaler struct {
	detailedMarshaler
	out    interface{}
	stack  detStack
	cursor string
	err    error
}

func NewDetailedMarshaler() *DetailedMarshaler {
	m := new(DetailedMarshaler)
	next := newBlock(m)
	next.onCommit = func(v interface{}) {
		if m.out != nil {
			m.Warning(errutil.New("can only write data once"))
		} else {
			m.out = v
		}
	}
	m.changeState(next)
	return m
}

// Data returns the accumulated script tree ready for serialization
// FIX, FUTURE: could write a custom json serialization to skip this in memory step.
func (m *DetailedMarshaler) Data() (interface{}, error) {
	return m.out, m.err
}

func (m *DetailedMarshaler) flushCursor() (ret string) {
	ret, m.cursor = m.cursor, ""
	return
}

func (m *DetailedMarshaler) pushState(d detailedMarshaler) {
	m.stack.push(m.detailedMarshaler) // remember the current state
	m.detailedMarshaler = d           // new current state
}

// set the current state to the last saved state
func (m *DetailedMarshaler) finishState(data interface{}) {
	m.detailedMarshaler = m.stack.pop()
	m.detailedMarshaler.commit(data)
}

// replace the top of the stack ( equals a pop and push )
func (m *DetailedMarshaler) changeState(d detailedMarshaler) {
	m.detailedMarshaler = d // new current state
}
