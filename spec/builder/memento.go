package builder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec"
	"strings"
)

// Memento implements spec.Block. Each chained call targets the surrounding block. For example, in:
//  if c.Cmd("parent").Begin() {
//    c.Cmd("some command", params).Cmds().Val(value).End()
//  }
// "the command, the array, and the val are all considered members of "parent".
type Memento struct {
	factory  *_Factory   // for chaining
	chain    *Memento    // for detecting bad chaining
	pos      Location    // source of the memento
	key      string      // the target of this memento in its parent
	spec     spec.Spec   // cmd data
	specs    spec.Specs  // array data
	val      interface{} // primitive data
	kids     Mementos    // child data, either array elements or command parameters
	cmdBlock bool
}

// Begin starts a new block
func (n *Memento) Begin() (okay bool) {
	if !n.cmdBlock {
		if next, e := n.factory.newCmds(n); e != nil {
			panic(errutil.New(e, Capture(1)))
		} else {
			n = next
		}
	}
	if e := n.factory.newBlock(); e != nil {
		panic(errutil.New(e, Capture(1)))
	} else {
		okay = true
	}
	return
}

// End terminates a block.
func (n *Memento) End() {
	if e := n.factory.endBlock(); e != nil {
		panic(errutil.New(e, Capture(1)))
	}
	return
}

// Cmd adds a new command of name with the passed set of positional args. Args can contain Mementos and literals. Returns a memento which can be passed to arrays or commands, or chained.
// To add data to the new command, pass them via args or follow this call with a call to Begin().
func (n *Memento) Cmd(name string, args ...interface{}) (ret spec.Block) {
	// HACK: because .Val("{cmd}") just looks so odd.
	if strings.Contains(name, "{") && len(args) == 0 {
		ret = n.Val(name)
	} else if next, e := n.factory.newCmd(n, name, args); e != nil {
		panic(errutil.New(e, Capture(1)))
	} else {
		next.cmdBlock = true
		ret = next
	}
	return
}

// Val specifies a single literal value: whether one primitive value or one array of primitive values.
func (n *Memento) Val(val interface{}) (ret spec.Block) {
	if n, e := n.factory.newVal(n, val); e != nil {
		panic(errutil.New(e, Capture(1)))
	} else {
		ret = n
	}
	return
}

// Param adds a key-value parameter to the spec.
// The passed name is the key; the return value is used to specify the value.
func (n *Memento) Param(field string) spec.Slot {
	return Param{n, field}
}
