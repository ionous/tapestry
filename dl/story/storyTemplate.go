package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/express"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/template"
	"git.sr.ht/~ionous/tapestry/template/types"
	"github.com/ionous/errutil"
)

func (op *RenderTemplate) PreImport(k *imp.Importer) (ret interface{}, err error) {
	if xs, e := template.Parse(op.Template.Str); e != nil {
		err = e
	} else if got, e := express.Convert(xs); e != nil {
		err = errutil.New(e, xs)
	} else if eval, ok := got.(rt.TextEval); !ok {
		err = errutil.Fmt("render template has unknown expression %T", got)
	} else {
		// note: this used to wrap this in a "RenderExp" which called the eval
		// that was just a marker; not needed so removed.
		//  &render.RenderExp{Expression: eval}
		ret = eval
		// pretty.Println(eval)
	}
	return
}

// returns a string or a FromText assignment as a slice of bytes
func ConvertText(k *imp.Importer, str string) (ret string, err error) {
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
			ret, err = k.Marshal(&assign.FromText{Value: eval})
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
