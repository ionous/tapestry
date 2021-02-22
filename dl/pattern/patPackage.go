package pattern

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"github.com/ionous/errutil"
)

var Slats = []composer.Composer{
	(*Determine)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &composer.CommandError{Cmd: op})
}

func cmdErrorCtx(op composer.Composer, ctx string, err error) error {
	return errutil.Append(err, &composer.CommandError{Cmd: op, Ctx: ctx})
}
