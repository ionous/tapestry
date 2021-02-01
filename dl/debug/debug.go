package debug

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"github.com/ionous/errutil"
)

var Slats = []composer.Composer{
	(*DoNothing)(nil),
	(*Log)(nil),
	(*Level)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &core.CommandError{Cmd: op})
}
