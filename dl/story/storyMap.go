package story

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
)

// block type to key handler
type BlockMap map[string]KeyMap
type KeyMap map[string]func(jsn.Block, interface{}) error

const BlockEnd = "$end"
const BlockStart = "$start"
const OtherBlocks = "$others"

func (callbacks BlockMap) call(b jsn.Block, key string, val interface{}) (err error) {
	if b != nil {
		if kvm, ok := callbacks[b.GetType()]; ok {
			err = kvm.call(b, key, val)
		}
		if err == nil {
			if kvm, ok := callbacks[OtherBlocks]; ok {
				err = kvm.call(b, key, val)
			}
		}
	}
	return
}

func (kvm KeyMap) call(b jsn.Block, key string, val interface{}) (err error) {
	if fn, ok := kvm[key]; ok {
		err = fn(b, val)
	}
	return
}

func Map(m *chart.Machine, callbacks BlockMap) (scope *BlockScope) {
	return &BlockScope{m: m, callbacks: callbacks}
}

// BlockScope - implements jsn.State keeping track of the current block while marshaling and
// self-terminating when the initial block goes out of scope.
type BlockScope struct {
	// fix: access to machine:
	// 1. pass it with context
	// 2. use return codes
	// 3. pass marshaler or some sort of "machine" interface
	//  as a parameter to end (and,or all calls)
	// 4. always pop on end.
	m         *chart.Machine
	Blocks    BlockStack
	callbacks BlockMap
	atKey     string
}

func (n *BlockScope) Commit(interface{}) {
	// ick
}

func (n *BlockScope) MarshalBlock(newBlock jsn.Block) (err error) {
	// the callbacks scope makes "newBlock" the the new top block *after* this function
	if key := n.atKey; len(key) > 0 {
		n.atKey = ""
		// top here is the *parent* block still
		err = n.callbacks.call(n.Blocks.Top(), key, newBlock)
	}
	if err == nil {
		err = n.callbacks.call(newBlock, BlockStart, nil)
	}
	if err == nil {
		n.Blocks.Push(newBlock)
	}
	return
}
func (n *BlockScope) MarshalKey(lede, key string) (err error) {
	n.atKey = key
	return
}
func (n *BlockScope) MarshalValue(valType string, val interface{}) (err error) {
	if key := n.atKey; len(key) > 0 {
		n.atKey = ""
		err = n.callbacks.call(n.Blocks.Top(), key, val)
	}
	return
}
func (n *BlockScope) EndBlock() {
	n.atKey = ""
	was := n.Blocks.Top()
	var parent jsn.Block
	if prev, ok := n.Blocks.Pop(); ok {
		parent = prev
	} else {
		n.m.FinishState("scope")
	}
	if e := n.callbacks.call(was, BlockEnd, parent); e != nil && e != jsn.Missing {
		n.m.Error(e)
	}
}
