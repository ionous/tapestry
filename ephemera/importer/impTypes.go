package importer

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
)

type BlockStack []jsn.Block

func (k *BlockStack) Push(b jsn.Block) {
	(*k) = append(*k, b)
}

func (k *BlockStack) At(typeName string) (ret jsn.Block, okay bool) {
	if end := len(*k) - 1; end >= 0 {
		if at := (*k)[end]; at.GetType() == typeName {
			ret, okay = at, true
		}
	}
	return
}

// return false if empty
func (k *BlockStack) Pop() (okay bool) {
	if end := len(*k) - 1; end >= 0 {
		(*k) = (*k)[:end]
		okay = true
	} else {
		(*k) = nil
	}
	return
}

// BlockScope - implements jsn.State keeping track of the current block while marshaling and
// self-terminating when the initial block goes out of scope.
type BlockScope struct {
	m     *chart.Machine
	wrap  chart.StateMix
	Stack BlockStack
}

func NewBlockScope(m *chart.Machine, wrap chart.StateMix) *BlockScope {
	// fix: access to machine:
	// 1. pass it with context
	// 2. use return codes
	// 3. pass marshaler or some sort of "machine" interface
	//  as a parameter to end (and,or all calls)
	// 4. always pop on end.
	return &BlockScope{m: m, wrap: wrap}
}

func (n *BlockScope) At(typeName string) (ret jsn.Block, okay bool) {
	return n.Stack.At(typeName)
}
func (n *BlockScope) MarshalBlock(b jsn.Block) (err error) {
	if n.wrap.OnBlock != nil {
		err = n.wrap.OnBlock(b)
	}
	if err == nil {
		n.Stack.Push(b)
	}
	return
}
func (n *BlockScope) MarshalKey(lede, key string) (err error) {
	if n.wrap.OnKey != nil {
		err = n.wrap.OnKey(lede, key)
	}
	return
}
func (n *BlockScope) MarshalValue(valType string, val interface{}) (err error) {
	if n.wrap.OnValue != nil {
		err = n.wrap.OnValue(valType, val)
	}
	return
}
func (n *BlockScope) EndBlock() {
	if n.wrap.OnEnd != nil {
		n.wrap.OnEnd()
	}
	if !n.Stack.Pop() {
		n.m.FinishState("scope")
	}
}
func (n *BlockScope) Commit(v interface{}) {
	if n.wrap.OnCommit != nil {
		n.wrap.OnCommit(v)
	}
}
