package detailed

import "git.sr.ht/~ionous/iffy/export/jsn"

type DetailedMarshaler struct {
	detailedState
	top    string // for debugging
	stack  detStack
	cursor string
	err    error
}

type detailedState interface {
	jsn.Marshaler
	named() string
	writeData(value interface{})
	readData() interface{}
}

func NewDetailedMarshaler() *DetailedMarshaler {
	m := &DetailedMarshaler{top: "root"}
	m.detailedState = &detBaseState{
		m:    m,
		name: "root",
	}
	return m
}

// Data returns the accumulated script tree ready for serialization
// FIX, FUTURE: could write a custom json serialization to skip this in memory step.
func (m *DetailedMarshaler) Data() (interface{}, error) {
	var out interface{}
	if det, ok := m.detailedState.(*detBaseState); ok {
		out = det.readData()
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
	m.top = m.detailedState.named()
}

// set the current state to the last saved state
func (m *DetailedMarshaler) popState() (ret detailedState) {
	ret, m.detailedState = m.detailedState, m.stack.pop()
	m.top = m.detailedState.named()
	return
}

// replace the top of the stack ( equals a pop and push )
func (m *DetailedMarshaler) changeState(d detailedState) {
	m.detailedState = d // new current state
	m.top = m.detailedState.named()
}
