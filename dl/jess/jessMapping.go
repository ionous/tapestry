package jess

// ----
func (op *MapLocations) Match(q Query, input *InputState) (okay bool) {
	panic("xxx")
}

func (op *MapLocations) Generate(Registrar) (err error) {
	panic("yyy")
}

// ----
func (op *MapDirections) Match(q Query, input *InputState) (okay bool) {
	panic("xxx")
}

func (op *MapDirections) Generate(Registrar) (err error) {
	panic("yyy")
}

// ----
func (op *MapConnections) Match(q Query, input *InputState) (okay bool) {
	panic("xxx")
}

func (op *MapConnections) Generate(Registrar) (err error) {
	panic("yyy")
}

// ----
func (op *DirectionFromLinks) Match(q Query, input *InputState) (okay bool) {
	panic("xxx")
}

// ----
func (op *Direction) Match(q Query, input *InputState) (okay bool) {
	panic("xxx")
}

// ----
func (op *Links) Match(q Query, input *InputState) (okay bool) {
	panic("xxx")
}

// ----
func (op *AdditionalLinks) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Links.Match(q, &next) {
		Optional(q, &next, &op.AdditionalLinks)
		*input, okay = next, true
	}
	return
}
