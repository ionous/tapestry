package compact

import "git.sr.ht/~ionous/iffy/export/jsn"

type CompactMarshaler struct {
	compactState
	top    string // for debugging
	stack  comStack
	cursor string
	err    error
}

type compactState interface {
	jsn.Marshaler
	named() string
	writeData(value interface{})
	readData() interface{}
}

// NewCompactMarshaler create an empty serializer that can produce compact data.
func NewCompactMarshaler() *CompactMarshaler {
	m := &CompactMarshaler{top: "root"}
	m.compactState = &comBlock{comValue{
		m:    m,
		name: "root",
	}}
	return m
}

// Data returns the accumulated script tree ready for serialization
// FIX, FUTURE: could write a custom json serialization to skip this in memory step.
func (m *CompactMarshaler) Data() (interface{}, error) {
	var out interface{}
	if com, ok := m.compactState.(*comBlock); ok {
		out = com.readData()
	}
	return out, m.err
}

func (m *CompactMarshaler) flushCursor() (ret string) {
	ret, m.cursor = m.cursor, ""
	return
}

func (m *CompactMarshaler) pushState(d compactState) {
	m.stack.push(m.compactState) // remember the current state
	m.compactState = d           // new current state
	m.top = m.compactState.named()
}

// set the current state to the last saved state
func (m *CompactMarshaler) popState() (ret compactState) {
	ret, m.compactState = m.compactState, m.stack.pop()
	m.top = m.compactState.named()
	return
}

// replace the top of the stack ( equals a pop and push )
func (m *CompactMarshaler) changeState(d compactState) {
	m.compactState = d // new current state
	m.top = m.compactState.named()
}
