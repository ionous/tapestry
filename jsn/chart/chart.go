package chart

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"github.com/ionous/errutil"
)

type Machine struct {
	State
	encoding bool
	out      interface{}
	stack    chartStack
	cursor   string
	err      error
	custom   Customization
}

// Customization
type Customization map[string]jsn.CustomizedMarshal

// NewEncoder writes json data
func NewEncoder(custom Customization, init func(*Machine) *StateMix) *Machine {
	return newMachine(custom, true, init)
}

// NewDecoder reads json data
func NewDecoder(custom Customization, init func(*Machine) *StateMix) *Machine {
	return newMachine(custom, false, init)
}

// newMachine create an empty serializer to produce compact script data.
func newMachine(custom Customization, encoding bool, init func(*Machine) *StateMix) *Machine {
	m := &Machine{encoding: encoding, custom: custom}
	if init != nil {
		next := init(m)
		next.OnCommit = func(v interface{}) {
			if m.out != nil {
				m.Error(errutil.New("can only write data once"))
			} else {
				m.out = v
			}
		}
		m.ChangeState(next)
	}
	return m
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
