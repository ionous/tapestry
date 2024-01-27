package assign

import (
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"github.com/ionous/errutil"
)

var cmdError = CmdError       // backwards compat
var cmdErrorCtx = CmdErrorCtx // backwards compat

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
	return errutil.Sprintf("# %s%s%s", name, padding, e.Ctx)
}

func CmdError(op typeinfo.Instance, err error) error {
	return cmdErrorCtx(op, "", err)
}

func CmdErrorCtx(op typeinfo.Instance, ctx string, err error) error {
	e := &CommandError{Cmd: op, Ctx: ctx}
	return errutil.Append(e, err)
}
