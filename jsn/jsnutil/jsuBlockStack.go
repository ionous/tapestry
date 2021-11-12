package jsnutil

import (
	"git.sr.ht/~ionous/iffy/jsn"
)

// BlockStack - keep track of the current block while marshaling.
// ( so that end block can be called. )
type BlockStack []StackedBlock

type StackedBlock struct {
	jsn.Block
	depth int
}

func (k *BlockStack) Push(b jsn.Block) {
	(*k) = append(*k, StackedBlock{b, 0})
}

func (k *BlockStack) Top() (ret *StackedBlock, okay bool) {
	if end := len(*k) - 1; end >= 0 {
		ret, okay = &(*k)[end], true
	}
	return
}

// return false if empty
func (k *BlockStack) Pop() (ret StackedBlock, okay bool) {
	if end := len(*k) - 1; end >= 0 {
		ret, (*k) = (*k)[end], (*k)[:end]
		okay = true
	}
	return
}
