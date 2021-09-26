package detailed

import (
	"git.sr.ht/~ionous/iffy/export/jsn"
	"github.com/ionous/errutil"
)

type DetailedMarshaler struct {
	detailedState
	top    string // for debugging
	stack  detStack
	cursor string
	err    error
}

type detailedState interface {
	jsn.Marshaler
	commit(value interface{})
}

type detRoot struct {
	detBaseState
	out interface{}
}

func (d *detRoot) commit(v interface{}) {
	if d.out != nil {
		d.m.Warning(errutil.New("can only write data once"))
	} else {
		d.out = v
	}
}

func NewDetailedMarshaler() *DetailedMarshaler {
	m := new(DetailedMarshaler)
	m.changeState(&detRoot{detBaseState: detBaseState{m: m}})
	return m
}

// Data returns the accumulated script tree ready for serialization
// FIX, FUTURE: could write a custom json serialization to skip this in memory step.
func (m *DetailedMarshaler) Data() (interface{}, error) {
	var out interface{}
	if det, ok := m.detailedState.(*detRoot); ok {
		out = det.out
	}
	return out, m.err
}

func (m *DetailedMarshaler) flushCursor() (ret string) {
	ret, m.cursor = m.cursor, ""
	return
}

func (m *DetailedMarshaler) pushState(d detailedState) {
	m.stack.push(m.detailedState) // remember the current state
	m.detailedState = d           // new current state
}

// set the current state to the last saved state
func (m *DetailedMarshaler) finishState(data interface{}) {
	m.detailedState = m.stack.pop()
	m.detailedState.commit(data)
}

// replace the top of the stack ( equals a pop and push )
func (m *DetailedMarshaler) changeState(d detailedState) {
	m.detailedState = d // new current state
}
