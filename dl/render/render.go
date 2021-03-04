package render

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"github.com/ionous/errutil"
)

var Slats = []composer.Composer{
	(*RenderField)(nil),
	(*RenderName)(nil),
	(*RenderPattern)(nil),
	(*RenderRef)(nil),
	(*RenderTemplate)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &composer.CommandError{Cmd: op})
}
