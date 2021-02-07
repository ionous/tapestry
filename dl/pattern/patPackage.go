package pattern

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"github.com/ionous/errutil"
)

var Slats = []composer.Composer{
	(*Determine)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &core.CommandError{Cmd: op})
}

func cmdErrorCtx(op composer.Composer, ctx string, err error) error {
	return errutil.Append(err, &core.CommandError{Cmd: op, Ctx: ctx})
}
