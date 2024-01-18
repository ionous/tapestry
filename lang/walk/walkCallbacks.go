package walk

// implements the Events interface
type Callbacks struct {
	OnFlow   func(Walker) error
	OnField  func(Walker) error
	OnSlot   func(Walker) error
	OnRepeat func(Walker) error
	OnValue  func(Walker) error
	OnEnd    func(Walker) error
}

func (c Callbacks) Flow(w Walker) (err error) {
	if fn := c.OnFlow; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) Slot(w Walker) (err error) {
	if fn := c.OnSlot; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) Repeat(w Walker) (err error) {
	if fn := c.OnRepeat; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) Field(w Walker) (err error) {
	if fn := c.OnField; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) Value(w Walker) (err error) {
	if fn := c.OnValue; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) End(w Walker) (err error) {
	if fn := c.OnEnd; fn != nil {
		err = fn(w)
	}
	return
}
