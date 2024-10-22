package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/express"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/template"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
	"github.com/ionous/errutil"
)

// transform SayTemplate into a RenderResponse
func (op *SayTemplate) PreImport(cat *weave.Catalog) (ret typeinfo.Instance, err error) {
	if xs, e := template.Parse(op.Template); e != nil {
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

// transform SayResponse into a RenderResponse
// ( post import so it happens after any transforms in its evals have been processed )
func (op *SayResponse) PostImport(cat *weave.Catalog) (ret typeinfo.Instance, err error) {
	// render by lookup if there's no text
	name := inflect.Normalize(op.ResponseName)
	if op.Text == nil {
		ret = &render.RenderResponse{Name: name}
	} else {
		if txt, e := convertEval(op.Text); e != nil {
			err = e
		} else if len(name) == 0 {
			// no name? render by value
			ret = &render.RenderResponse{Text: txt}
		} else {
			// otherwise store the value
			if e := cat.ScheduleCmd(op, weaver.NounPhase, func(w weaver.Weaves, run rt.Runtime) error {
				return w.AddKindFields(kindsOf.Response.String(), []mdl.FieldInfo{{
					Name:     name,
					Affinity: affine.Text,
					Init:     &call.FromText{Value: txt}},
				})
			}); e != nil {
				err = e
			} else {
				// and render by lookup if we stored the text
				ret = &render.RenderResponse{Name: name}
			}
		}
	}
	return
}

func convertEval(txt rt.TextEval) (ret rt.TextEval, err error) {
	if lit, ok := txt.(*literal.TextValue); !ok || len(lit.KindName) > 0 {
		ret = txt
	} else {
		ret, err = jess.ConvertTextTemplate(lit.Value)
	}
	return
}

func convertTextAssignment(str string) (ret rt.Assignment, err error) {
	if txt, e := jess.ConvertTextTemplate(str); e != nil {
		err = e
	} else {
		ret = &call.FromText{Value: txt}
	}
	return
}
