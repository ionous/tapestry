package jess

type Priority int
type Process func(Query) error

const (
	// these happen immediately after matching
	// primarily so that multi part traits can match correctly.
	// re: The bother is a fixed in place closed container in the kitchen.
	// another option would be to match the whole span and break it up into actual traits later.
	GenerateKinds Priority = iota
	// turn specified nouns into desired nouns
	// waits until after GenerateKinds so that all specified kind names are known
	// ex. `The sapling is a tree. A tree is a kind of thing.` ( though that doesnt work in inform. )
	GenerateNouns
	GenerateDefaultKinds
	GenerateValues // generates implied nouns
	GenerateConnections
	GenerateUnderstanding // awww. love and peas.
	PriorityCount
)

// used internally to generate matched phrases in a good order.
type ProcessingList struct {
	posted [PriorityCount][]Process
}

func (m *ProcessingList) AddToList(i Priority, p Process) {
	m.posted[i] = append(m.posted[i], p)
}

func (m *ProcessingList) ProcessAll(q Query) (err error) {
Loop:
	for i := Priority(0); i < PriorityCount; i++ {
		// fix: can we really add new processes during post?
		// and if so, shouldnt mock panic on misorders?
		for len(m.posted[i]) > 0 {
			posted := m.posted[i]
			next, rest := posted[0], posted[1:]
			if e := next(q); e != nil {
				err = e
				break Loop
			} else {
				m.posted[i] = rest
			}
		}
	}
	return
}
