package compact

import "github.com/ionous/errutil"

type CompactMarshaler struct {
	compactMarshaler
	out    interface{}
	stack  comStack
	cursor string
	err    error
}

// NewCompactMarshaler create an empty serializer to produce compact script data.
func NewCompactMarshaler() *CompactMarshaler {
	m := new(CompactMarshaler)
	next := newBlock(m)
	next.onCommit = func(v interface{}) {
		if m.out != nil {
			m.Warning(errutil.New("can only write data once"))
		} else {
			m.out = v
		}
	}
	m.changeState(next)
	// m.changeState(&comRoot{comBlock: comBlock{comValue{m: m}}})
	return m
}

// Data returns the accumulated script tree ready for serialization
func (m *CompactMarshaler) Data() (interface{}, error) {
	return m.out, m.err
}

func (m *CompactMarshaler) flushCursor() (ret string) {
	ret, m.cursor = m.cursor, ""
	return
}

func (m *CompactMarshaler) pushState(d compactMarshaler) {
	m.stack.push(m.compactMarshaler) // remember the current state
	m.compactMarshaler = d           // new current state
}

// set the current state to the last saved state
func (m *CompactMarshaler) finishState(data interface{}) {
	m.compactMarshaler = m.stack.pop()
	m.compactMarshaler.commit(data)
}

// replace the top of the stack ( equals a pop and push )
func (m *CompactMarshaler) changeState(d compactMarshaler) {
	m.compactMarshaler = d // new current state
}
