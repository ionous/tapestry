package compact

import "git.sr.ht/~ionous/iffy/export/jsn"

type CompactMarshaler struct {
	compactState
	stack  comStack
	cursor string
	err    error
}

type compactState interface {
	jsn.Marshaler
	// someone writes a fully completed value into us.
	commit(value interface{})
}

type comRoot struct {
	comBlock
	out interface{}
}

func (cr *comRoot) commit(v interface{}) {
	cr.out = v
}

// NewCompactMarshaler create an empty serializer that can produce compact data.
func NewCompactMarshaler() *CompactMarshaler {
	m := new(CompactMarshaler)
	m.changeState(&comRoot{comBlock: comBlock{comValue{m: m}}})
	return m
}

// Data returns the accumulated script tree ready for serialization
// FIX, FUTURE: could write a custom json serialization to skip this in memory step.
func (m *CompactMarshaler) Data() (interface{}, error) {
	var out interface{}
	if cr, ok := m.compactState.(*comRoot); ok {
		out = cr.out
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
}

// set the current state to the last saved state
func (m *CompactMarshaler) finishState(data interface{}) {
	m.compactState = m.stack.pop()
	m.compactState.commit(data)
}

// replace the top of the stack ( equals a pop and push )
func (m *CompactMarshaler) changeState(d compactState) {
	m.compactState = d // new current state
}
