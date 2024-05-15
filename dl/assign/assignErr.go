package assign

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"github.com/ionous/errutil"
)

type CommandError struct {
	Cmd typeinfo.Instance
	Ctx string
}

func (e CommandError) Error() string {
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

func CmdError(op typeinfo.Instance, err error) error {
	return CmdErrorCtx(op, "", err)
}

func CmdErrorCtx(op typeinfo.Instance, ctx string, err error) error {
	// fix: implement As and Is
	e := &CommandError{Cmd: op, Ctx: ctx}
	return errutil.Append(e, err)
}
