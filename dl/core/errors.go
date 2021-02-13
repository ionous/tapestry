package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"github.com/ionous/errutil"
)

type CommandError struct {
	Cmd composer.Composer
	Ctx string
}

func (e *CommandError) Error() string {
	name := composer.SpecName(e.Cmd)
	var padding string
	if len(e.Ctx) > 0 {
		padding = " "
	}
	return errutil.Sprintf("error in command %q%s%s", name, padding, e.Ctx)
}

func cmdError(op composer.Composer, err error) error {
	return cmdErrorCtx(op, "", err)
}

func cmdErrorCtx(op composer.Composer, ctx string, err error) error {
	// avoid triggering errutil panics for break statements
	if _, ok := err.(DoInterrupt); !ok {
		err = errutil.Append(err, &CommandError{Cmd: op, Ctx: ctx})
	}
	return err
}
