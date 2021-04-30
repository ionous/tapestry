package debug

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"github.com/ionous/errutil"
)

//go:generate capnp compile -I ../../../../../zombiezen.com/go/capnproto2/std -ogo:.. --src-prefix=../../idl ../../idl/debug/debug.capnp
var Slats = []composer.Composer{
	(*DoNothing)(nil),
	(*Log)(nil),
	(*Level)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &composer.CommandError{Cmd: op})
}
