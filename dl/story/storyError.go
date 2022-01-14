package story

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"github.com/ionous/errutil"
)

type OpError struct {
	Op  composer.Composer
	Err error
}

const UnhandledSwap = errutil.Error("unhandled swap")
const MissingSlot = errutil.Error("missing slot")
const InvalidValue = errutil.Error("invalid value")

func ImportError(op composer.Composer, e error) error {
	return &OpError{op, e}
}

func (e *OpError) Error() string {
	return errutil.Sprintf("%s in %s at %s", e.Err, composer.SpecName(e.Op))
}

func (e *OpError) Unwrap() error {
	return e.Err
}
