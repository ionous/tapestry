package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"github.com/ionous/errutil"
)

func cmdError(op composer.Composer, err error) error {
	return cmdErrorCtx(op, "", err)
}

func cmdErrorCtx(op composer.Composer, ctx string, err error) error {
	// avoid triggering errutil panics for break statements
	if _, ok := err.(DoInterrupt); !ok {
		e := &composer.CommandError{Cmd: op, Ctx: ctx}
		err = errutil.Append(err, e)
	}
	return err
}
