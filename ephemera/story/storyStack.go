package story

import (
	"git.sr.ht/~ionous/iffy/jsn"
)

type BlockStack []jsn.Block

func (k *BlockStack) Push(b jsn.Block) {
	(*k) = append(*k, b)
}

func (k *BlockStack) Top() (ret jsn.Block) {
	if end := len(*k) - 1; end >= 0 {
		ret = (*k)[end]
	}
	return
}

// return false if empty
func (k *BlockStack) Pop() (ret jsn.Block, okay bool) {
	if end := len(*k) - 1; end >= 0 {
		ret, (*k) = (*k)[end], (*k)[:end]
		okay = true
	}
	return
}
