package debug

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"github.com/ionous/errutil"
)

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &composer.CommandError{Cmd: op})
}
