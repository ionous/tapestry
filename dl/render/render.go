package render

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"github.com/ionous/errutil"
)

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &composer.CommandError{Cmd: op})
}
