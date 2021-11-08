package story

import (
	"bytes"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/render"
	"git.sr.ht/~ionous/iffy/ephemera/express"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/template"
	"git.sr.ht/~ionous/iffy/template/types"
	"github.com/ionous/errutil"
)

func (op *RenderTemplate) ImportStub(k *Importer) (ret interface{}, err error) {
	if xs, e := template.Parse(op.Template.Str); e != nil {
		err = e
	} else if got, e := express.Convert(xs); e != nil {
		err = errutil.New(e, xs)
	} else if eval, ok := got.(rt.TextEval); !ok {
		err = errutil.Fmt("render template has unknown expression %T", got)
	} else {
		ret = &render.RenderExp{eval}
		// pretty.Println(eval)
	}
	return
}

// returns a string or a FromText assignment as a slice of bytes
func ConvertText(str string) (ret interface{}, err error) {
	if xs, e := template.Parse(str); e != nil {
		err = e
	} else if str, ok := getSimpleString(xs); ok {
		ret = str // okay; return the string.
	} else {
		if got, e := express.Convert(xs); e != nil {
			err = errutil.New(e, xs)
		} else if eval, ok := got.(rt.TextEval); !ok {
			err = errutil.Fmt("render template has unknown expression %T", got)
		} else {
			var buf bytes.Buffer
			if e := cout.Marshal(&buf, &core.FromText{eval}); e != nil {
				err = e
			} else {
				ret = buf.Bytes() // okay; return bytes.
			}
		}
	}
	return
}

// see if the parsed expression contained anything other than text
// if true, return that text
func getSimpleString(xs template.Expression) (ret string, okay bool) {
	switch len(xs) {
	case 0:
		okay = true
	case 1:
		if quote, ok := xs[0].(types.Quote); ok {
			ret, okay = quote.Value(), true
		}
	}
	return
}
