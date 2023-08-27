package assign

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"github.com/ionous/errutil"
)

var cmdError = CmdError       // backwards compat
var cmdErrorCtx = CmdErrorCtx // backwards compat
type CommandError = composer.CommandError

func CmdError(op composer.Composer, err error) error {
	return cmdErrorCtx(op, "", err)
}

func CmdErrorCtx(op composer.Composer, ctx string, err error) error {
	e := &CommandError{Cmd: op, Ctx: ctx}
	return errutil.Append(e, err)
}
