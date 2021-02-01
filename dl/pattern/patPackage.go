package pattern

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/term"
	"github.com/ionous/errutil"
)

// Rules contained by this package.
// fix: would it be better to list rule sets?
// the rule set elements could be used to find the individual rule types.
var Support = []interface{}{
	(*Rule)(nil),
	//
	//(*term.Preparer)(nil),
	(*term.Number)(nil),
	(*term.Bool)(nil),
	(*term.Text)(nil),
	(*term.Record)(nil),
	(*term.Object)(nil),
	(*term.NumList)(nil),
	(*term.TextList)(nil),
	(*term.RecordList)(nil),
}

var Slats = []composer.Composer{
	(*Determine)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &core.CommandError{Cmd: op})
}

func cmdErrorCtx(op composer.Composer, ctx string, err error) error {
	return errutil.Append(err, &core.CommandError{Cmd: op, Ctx: ctx})
}
