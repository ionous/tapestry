package rel

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
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
	return errutil.Append(err, &composer.CommandError{Cmd: op})
}
