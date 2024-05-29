package assign

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/assign/dot"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// uniform access to objects and variables
// implemented by object dot and variable dot.
type Address interface {
	GetReference(rt.Runtime) (dot.Reference, error)
}

// check that op exists and then call GetPath on it.
func GetReference(run rt.Runtime, op Address) (ret dot.Reference, err error) {
	if op == nil {
		err = errors.New("missing address")
	} else {
		ret, err = op.GetReference(run)
	}
	return
}

// return a dot.Index
func (op *AtIndex) Resolve(run rt.Runtime) (ret dot.Dotted, err error) {
	if idx, e := safe.GetNumber(run, op.Index); e != nil {
		err = CmdError(op, e)
	} else {
		ret = dot.Index(idx.Int() - 1)
	}
	return
}

// return a dot.Field
func (op *AtField) Resolve(run rt.Runtime) (ret dot.Dotted, err error) {
	if field, e := safe.GetText(run, op.Field); e != nil {
		err = CmdError(op, e)
	} else {
		ret = dot.Field(field.String())
	}
	return
}
