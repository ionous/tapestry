package assign

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/assign/dot"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// uniform access to objects and variables
// implemented by various Tapestry commands.
type Address interface {
	GetReference(rt.Runtime) (dot.Endpoint, error)
}

// check that op exists and then call GetReference on it.
func GetReference(run rt.Runtime, op Address) (ret dot.Endpoint, err error) {
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
		err = cmdError(op, e)
	} else {
		ret = dot.Index(idx.Int() - 1)
	}
	return
}

// return a dot.Field
func (op *AtField) Resolve(run rt.Runtime) (ret dot.Dotted, err error) {
	if field, e := safe.GetText(run, op.Field); e != nil {
		err = cmdError(op, e)
	} else {
		ret = dot.Field(field.String())
	}
	return
}
