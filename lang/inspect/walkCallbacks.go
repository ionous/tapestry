package inspect

// implements the Events interface
type Callbacks struct {
	OnFlow   func(It) error
	OnField  func(It) error
	OnSlot   func(It) error
	OnRepeat func(It) error
	OnValue  func(It) error
	OnEnd    func(It) error
}

func (c Callbacks) Flow(w It) (err error) {
	if fn := c.OnFlow; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) Slot(w It) (err error) {
	if fn := c.OnSlot; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) Repeat(w It) (err error) {
	if fn := c.OnRepeat; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) Field(w It) (err error) {
	if fn := c.OnField; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) Value(w It) (err error) {
	if fn := c.OnValue; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) End(w It) (err error) {
	if fn := c.OnEnd; fn != nil {
		err = fn(w)
	}
	return
}
