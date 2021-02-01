package story

import (
	"git.sr.ht/~ionous/iffy/ephemera/decode"
	"github.com/ionous/errutil"
)

func (op *Certainty) ImportString(k *Importer) (ret string, err error) {
	if str, ok := decode.FindChoice(op, op.Str); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w %q", InvalidValue, op.Str))
	} else {
		ret = str
	}
	return
}

// blocks of text might well be a template.
func (op *Lines) ConvertText() (ret interface{}, err error) {
	return convert_text_or_template(op.Str)
}
