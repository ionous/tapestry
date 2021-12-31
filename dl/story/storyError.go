package story

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/reader"
	"github.com/ionous/errutil"
)

type OpError struct {
	Op  composer.Composer
	At  reader.Position
	Err error
}

const UnhandledSwap = errutil.Error("unhandled swap")
const MissingSlot = errutil.Error("missing slot")
const InvalidValue = errutil.Error("invalid value")

func ImportError(op composer.Composer, at reader.Position, e error) error {
	return &OpError{op, at, e}
}

func (e *OpError) Error() string {
	return errutil.Sprintf("%s in %s at %s", e.Err, composer.SpecName(e.Op), e.At.String())
}

func (e *OpError) Unwrap() error {
	return e.Err
}
