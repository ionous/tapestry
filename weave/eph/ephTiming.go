package eph

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/imp/assert"
)

func toTiming(t Timing, always Always) (ret assert.EventTiming) {
	switch t.Str {
	case Timing_Before:
		ret = assert.Before
	case Timing_During:
		ret = assert.During
	case Timing_After:
		ret = assert.After
	case Timing_Later:
		ret = assert.Later
	}
	if always.Str == Always_Always {
		ret |= assert.RunAlways
	}
	return
}

func fromTiming(t assert.EventTiming) (ret Timing, always Always) {
	if t&assert.RunAlways != 0 {
		always.Str = Always_Always
		t ^= assert.RunAlways
	}
	switch t {
	case assert.Before:
		ret.Str = Timing_Before
	case assert.During:
		ret.Str = Timing_During
	case assert.After:
		ret.Str = Timing_After
	case assert.Later:
		ret.Str = Timing_Later
	}
	return
}

// Always requires a predefined string.
type Always struct {
	Str string
}

func (op *Always) String() string {
	return op.Str
}

const Always_Always = "$ALWAYS"

func (*Always) Compose() composer.Spec {
	return composer.Spec{
		Uses: composer.Type_Str,
		Choices: []string{
			Always_Always,
		},
		Strings: []string{
			"always",
		},
	}
}

// Timing requires a predefined string.
type Timing struct {
	Str string
}

func (op *Timing) GetPartition() (ret int, okay bool) {
	// probably shouldnt but just rely on declared order for now.
	spec := op.Compose()
	if _, i := spec.IndexOfChoice(op.Str); i >= 0 {
		ret, okay = i, true
	}
	return
}

func (op *Timing) String() string {
	return op.Str
}

const Timing_Before = "$BEFORE"
const Timing_During = "$DURING"
const Timing_After = "$AFTER"
const Timing_Later = "$LATER"

func (*Timing) Compose() composer.Spec {
	return composer.Spec{
		Uses: composer.Type_Str,
		Choices: []string{
			Timing_Before, Timing_During, Timing_After, Timing_Later,
		},
		Strings: []string{
			"before", "during", "after", "later",
		},
	}
}
