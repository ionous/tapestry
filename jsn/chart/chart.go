package chart

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"github.com/ionous/errutil"
)

type Machine struct {
	State
	encoding bool
	stack    chartStack
	cursor   string
	err      error
	custom   Customization
}

// Customization
type Customization map[string]jsn.CustomizedMarshal

// MakeEncoder writes json data
func MakeEncoder(custom Customization) Machine {
	return makeMachine(custom, true)
}

// MakeDecoder reads json data
func MakeDecoder(custom Customization) Machine {
	return makeMachine(custom, false)
}

// newMachine create an empty serializer to produce compact script data.
func makeMachine(custom Customization, encoding bool) Machine {
	return Machine{encoding: encoding, custom: custom}
}

// IsEncoding indicates whether the machine is writing json ( or reading json. )
func (m *Machine) IsEncoding() bool {
	return m.encoding
}

func (m *Machine) CustomizedMarshal(typeName string) (ret jsn.CustomizedMarshal, okay bool) {
	ret, okay = m.custom[typeName]
	return
}

func (m *Machine) SetCursor(id string) {
	m.cursor = id
}

// Data returns the accumulated script tree ready for serialization
func (m *Machine) Errors() error {
	return m.err
}

// FlushCursor - return and reset the most recently recorded cursor position
func (m *Machine) FlushCursor() (ret string) {
	ret, m.cursor = m.cursor, ""
	return
}

func (m *Machine) Error(e error) {
	m.err = errutil.Append(m.err, e)
	return
}

// PushState - enter the passed state saving the current state into history
func (m *Machine) PushState(d State) {
	if d == nil {
		panic("trying to push a nil state")
	}
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
