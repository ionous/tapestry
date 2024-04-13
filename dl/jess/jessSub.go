package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
)

func (op *SubAssignment) Match(input *InputState) (okay bool) {
	return false
}

func (op *SubAssignment) Deduce() (ret rt.Assignment, err error) {
	return nil, nil
}
