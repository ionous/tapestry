package debug

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"github.com/ionous/errutil"
)

var Slats = []composer.Composer{
	(*DoNothing)(nil),
	(*Log)(nil),
	(*Level)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &composer.CommandError{Cmd: op})
}
