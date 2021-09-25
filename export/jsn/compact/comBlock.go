package compact

// supports blocks of data ( or single values ) being written.
type comBlock struct {
	comValue
}

// starts a series of key-values pairs
// the flow is closed ( written ) with a call to EndValues()
func (d *comBlock) MapValues(lede, kind string) {
	next := &comFlow{
		comBlock: comBlock{comValue{
			m:    d.m,
			name: newName("flow"),
		}},
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
		comBlock: comBlock{comValue{
			m:    d.m,
			name: newName("swap"),
		}},
	})
}

func (d *comBlock) RepeatValues(hint int) {
	d.m.pushState(&comSlice{
		comBlock: comBlock{comValue{
			m:    d.m,
			name: newName("slice"),
		}},
		values: make([]interface{}, 0, hint),
	})
}

// EndValues ends the current state and writeDatas its data to the parent state.
func (d *comBlock) EndValues() {
	was := d.m.popState()         // gets rid of us
	d.m.writeData(was.readData()) // write our accumulated data to the new parent
	// again, use "was" not "d" to get to the outermost version of ourself
}
