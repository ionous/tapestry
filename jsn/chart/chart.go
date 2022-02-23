package chart

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"github.com/ionous/errutil"
)

type Machine struct {
	State
	encoding bool
	stack    chartStack
	Comment  *string
	err      error
}

// MakeEncoder writes data from Tapestry's data-language structures.
// Encoding *skips* empty in-memory structures ( ex. nil slices. )
func MakeEncoder() Machine {
	return makeMachine(true)
}

// MakeDecoder reads data into Tapestry's data-language structures.
// Decoding visits the in-memory structures regardless of whether they currently have data
// ( so that they can be filled with data from whatever source is being decoded. )
func MakeDecoder() Machine {
	return makeMachine(false)
}

// newMachine create an empty serializer to produce compact script data.
func makeMachine(encoding bool) Machine {
	return Machine{encoding: encoding}
}

func (m *Machine) Marshal(tgt jsn.Marshalee, init State) (err error) {
	m.ChangeState(&StateMix{}) // fix. right now, if you try to Finish the initial state (ex. during tests) pop panics
	m.PushState(init)
	if e := tgt.Marshal(m); e != nil {
		err = e
	} else {
		err = m.err
	}
	return
}

// IsEncoding indicates whether the machine is writing json ( or reading json. )
func (m *Machine) IsEncoding() bool {
	return m.encoding
}

func (m *Machine) SetComment(pid *string) {
	m.Comment = pid
}

// Data returns the accumulated script tree ready for serialization
func (m *Machine) Errors() error {
	return m.err
}

func (m *Machine) Error(e error) {
	m.err = errutil.Append(m.err, e)
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
