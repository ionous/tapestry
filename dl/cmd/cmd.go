package cmd

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"github.com/ionous/errutil"
)

type cmdError struct {
	Cmd typeinfo.Instance
	Ctx string
}

func (e cmdError) Error() string {
	name := "<nil>"
	if e.Cmd != nil {
		t := e.Cmd.TypeInfo()
		name = t.TypeName()
	}
	var padding string
	if len(e.Ctx) > 0 {
		padding = " "
	}
	return fmt.Sprintf("# %s%s%s", name, padding, e.Ctx)
}

func Error(op typeinfo.Instance, err error) error {
	return ErrorCtx(op, "", err)
}

func ErrorCtx(op typeinfo.Instance, ctx string, err error) error {
	// fix: implement As and Is
	e := &cmdError{Cmd: op, Ctx: ctx}
	return errutil.Append(e, err)
}
