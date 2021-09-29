package chart

import (
	"github.com/ionous/errutil"
)

type Machine struct {
	State
	out    interface{}
	stack  chartStack
	cursor string
	err    error
}

// NewMachine create an empty serializer to produce compact script data.
func NewMachine(init func(*Machine) *StateMix) *Machine {
	m := new(Machine)
	next := init(m)
	next.OnCommit = func(v interface{}) {
		if m.out != nil {
			m.Warning(errutil.New("can only write data once"))
		} else {
			m.out = v
		}
	}
	m.ChangeState(next)
	return m
}

// Data returns the accumulated script tree ready for serialization
func (m *Machine) Data() (interface{}, error) {
	return m.out, m.err
}

// FlushCursor - return and reset the most recently recorded cursor position
func (m *Machine) FlushCursor() (ret string) {
	ret, m.cursor = m.cursor, ""
	return
}

// PushState - enter the passed state saving the current state into history
func (m *Machine) PushState(d State) {
	m.stack.push(m.State) // remember the current state
	m.State = d           // new current state
}

// FinishState - end the current state,
// and send the passed data ( presumably from the current state ) to the most recent prior state.
func (m *Machine) FinishState(data interface{}) {
	m.State = m.stack.pop()
	m.State.Commit(data)
}

// ChangeState - replace the current state without remembering the previous state.
func (m *Machine) ChangeState(d State) {
	m.State = d // new current state
}
