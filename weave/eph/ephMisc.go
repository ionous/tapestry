package eph

import (
	"git.sr.ht/~ionous/tapestry/imp/assert"
)

type Definition struct {
	Path  []string
	Value string
}

func (op *Definition) Weave(k assert.Assertions) (err error) {
	path := append(op.Path, op.Value)
	return k.AssertDefinition(path...)
}

// wrapper for implementing Ephemera with free functions.
// its a shortcut to defining new ephemera structs for every last processing statement;
// which works because its all the same process space.
// could separate these out into commands for inter-process communication, logging, etc. if ever needed.
type PhaseFunction struct {
	OnPhase assert.Phase
	Do      func(assert.World, assert.Assertions) error
}

func (fn PhaseFunction) Weave(k assert.Assertions) (err error) {
	return k.Schedule(fn.OnPhase, fn.Do)
}
