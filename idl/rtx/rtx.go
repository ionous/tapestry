package rtx

import (
	"git.sr.ht/~ionous/iffy/idl/reg"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

// field is a pointer to a BoolEvalImpl
func (op *BoolEval) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if !op.IsValid() {
		err = errutil.New("eval not set")
	} else if ptr, e := op.EvalPtr(); e != nil {
		err = e
	} else if i, e := reg.Unpack(BoolEval_TypeID, ptr); e != nil {
		err = e
	} else if eval, ok := i.(rt.BoolEval); !ok {
		err = errutil.Fmt("unexpected eval %T", i)
	} else {
		ret, err = eval.GetBool(run)
	}
	return
}
