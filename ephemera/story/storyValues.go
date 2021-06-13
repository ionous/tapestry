package story

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/value"
	"github.com/ionous/errutil"
)

type Lines = value.Lines

func (op *Certainty) ImportString(k *Importer) (ret string, err error) {
	if str, ok := composer.FindChoice(op, op.Str); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w %q", InvalidValue, op.Str))
	} else {
		ret = str
	}
	return
}
