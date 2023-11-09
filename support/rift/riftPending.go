package rift

import "git.sr.ht/~ionous/tapestry/support/charmed"

// used by riftEntry to read values when the entry is finished.
type pendingValue interface {
	FinalizeValue() (any, error)
}

// a final value, ex. from a boolean.
type computedValue struct{ v any }

func (v computedValue) FinalizeValue() (any, error) {
	return v.v, nil
}

// a number --
// note this is a little different than the other types
// because there's no terminal value for it.
type numValue struct{ charmed.NumParser }

// fix? returns float64 because json does
// could also return int64 when its int like
func (v *numValue) FinalizeValue() (ret any, err error) {
	ret, err = v.GetFloat()
	return
}
