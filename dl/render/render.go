package render

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"github.com/ionous/errutil"
)

var Slats = []composer.Composer{
	(*RenderField)(nil),
	(*RenderName)(nil),
	(*RenderRef)(nil),
	(*RenderTemplate)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &core.CommandError{Cmd: op})
}
