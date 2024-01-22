package inspect

import "errors"

// implements Events by delegating calls to whatever is on the top of the stack.
type Stack struct {
	Events
	events []Events
}

func (s *Stack) Replace(next Events) {
	s.Events = next
}

// always returns nil
func (s *Stack) Push(next Events) (_ error) {
	if s.Events != nil {
		s.events = append(s.events, s.Events)
	}
	s.Events = next
	return
}

func (s *Stack) Pop() (err error) {
	// various extra steps to avoid having to
	// implement all of the event functions manually
	if s.Events == nil {
		err = errors.New("stack is empty")
	} else {
		var top Events
		if end := len(s.events) - 1; end >= 0 {
			top = s.events[end]
			s.events = s.events[:end]
		}
		s.Events = top // nil if its the last one
	}
	return
}
