package rel

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"github.com/ionous/errutil"
)

var Slats = []composer.Composer{
	(*Relation)(nil),
	(*Relate)(nil),
	(*RelativeOf)(nil),
	(*RelativesOf)(nil),
	(*ReciprocalOf)(nil),
	(*ReciprocalsOf)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &core.CommandError{Cmd: op})
}
