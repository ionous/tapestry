package compact

// supports blocks of data ( or single values ) being written.
type comBlock struct {
	comValue
}

// starts a series of key-values pairs
// the flow is closed ( written ) with a call to EndValues()
func (d *comBlock) MapValues(lede, kind string) {
	next := &comFlow{
		comBlock: comBlock{comValue{m: d.m}},
	}
	next.sig.WriteLede(lede)
	d.m.pushState(next)
}

func (d *comBlock) SetCursor(id string) {
	d.m.cursor = id //  for now, overwrite without error checking.
}

// ex."noun_phrase" "$KIND_OF_NOUN"
func (d *comBlock) PickValues(kind, choice string) {
	d.m.pushState(&comSwap{
		comBlock: comBlock{comValue{m: d.m}},
	})
}

func (d *comBlock) RepeatValues(hint int) {
	d.m.pushState(&comSlice{
		comBlock: comBlock{comValue{m: d.m}},
		values:   make([]interface{}, 0, hint),
	})
}
