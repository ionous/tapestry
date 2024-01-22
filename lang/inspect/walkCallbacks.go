package inspect

// implements the Events interface
type Callbacks struct {
	OnFlow   func(Iter) error
	OnField  func(Iter) error
	OnSlot   func(Iter) error
	OnRepeat func(Iter) error
	OnValue  func(Iter) error
	OnEnd    func(Iter) error
}

func (c Callbacks) Flow(w Iter) (err error) {
	if fn := c.OnFlow; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) Slot(w Iter) (err error) {
	if fn := c.OnSlot; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) Repeat(w Iter) (err error) {
	if fn := c.OnRepeat; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) Field(w Iter) (err error) {
	if fn := c.OnField; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) Value(w Iter) (err error) {
	if fn := c.OnValue; fn != nil {
		err = fn(w)
	}
	return
}

func (c Callbacks) End(w Iter) (err error) {
	if fn := c.OnEnd; fn != nil {
		err = fn(w)
	}
	return
}
