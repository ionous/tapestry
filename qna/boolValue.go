package qna

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type boolEval struct{ evalValue }

// GetBool, or error if the underlying value isn't a bool
func (q *boolEval) GetBool() (bool, error) {
	return rt.GetBool(q.run, q.eval.(rt.BoolEval))
}

func newBoolValue(run rt.Runtime, v interface{}) (ret rt.Value, err error) {
	switch a := v.(type) {
	case nil: // zero value for unhandled defaults in sqlite
		ret = generic.False // fix? could also return some predefined "constants" for true,false,and... for the other types, nil
	case bool:
		ret = generic.NewBool(a)
	case int64: // sqlite, boolean values can be represented as 1/0
		ret = generic.NewBool(a != 0)
	case []byte:
		var eval rt.BoolEval
		if e := bytesToEval(a, &eval); e != nil {
			err = e
		} else {
			ret = &boolEval{evalValue{run: run, eval: eval}}
		}
	default:
		err = errutil.New("expected boolean value, got %v(%T)", v, v)
	}
	return
}