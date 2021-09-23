package jsn

type DetailedMarshaler struct {
	detailedState
	top    string
	stack  jsnDetails
	cursor string
}

type detailedState interface {
	Marshaler
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

func (m *DetailedMarshaler) Data() (ret interface{}) {
	if det, ok := m.detailedState.(*detBaseState); ok {
		ret = det.readData()
	}
	return
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

// doesnt change state
func (m *DetailedMarshaler) makeValue(kind string, value interface{}) detValueData {
	return detValueData{
		Id:    m.flushCursor(),
		Type:  kind,
		Value: value,
	}
}
