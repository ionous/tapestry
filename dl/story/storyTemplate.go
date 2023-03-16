package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/express"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/template"
	"git.sr.ht/~ionous/tapestry/template/types"
	"github.com/ionous/errutil"
)

// transform SayTemplate into a RenderResponse
func (op *SayTemplate) PreImport(k *imp.Importer) (interface{}, error) {
	return convertTemplate("", op.Template.Str)
}

// transform SayResponse into a RenderResponse
func (op *SayResponse) PreImport(k *imp.Importer) (interface{}, error) {
	k.WriteEphemera(&eph.EphKinds{
		Kind: kindsOf.Response.String(),
		Contain: []eph.EphParams{{
			Affinity:  eph.Affinity{eph.Affinity_Text},
			Name:      op.Name,
			Initially: &assign.FromText{Value: op.Text},
		}}})
	//
	return &render.RenderResponse{Name: op.Name, Text: op.Text}, nil
}

func convertTemplate(name, tmpl string) (ret *render.RenderResponse, err error) {
	if xs, e := template.Parse(tmpl); e != nil {
		err = e
	} else if got, e := express.Convert(xs); e != nil {
		err = errutil.New(e, xs)
	} else if eval, ok := got.(rt.TextEval); !ok {
		err = errutil.Fmt("render template has unknown expression %T", got)
	} else {
		ret = &render.RenderResponse{Text: eval}
	}
	return
}

// returns a string or a FromText assignment as a slice of bytes
func ConvertText(k *imp.Importer, str string) (ret string, err error) {
	// FIX: it would make more sense if the ephemera stored this as an rt.Assignment
	// upgrade to include Affinity() for literal value, and all would be well.
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
			// ex. storyMarshaller
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
