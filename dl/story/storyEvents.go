package story

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *EventBlock) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *EventBlock) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weave.RequireAncestry, func(w *weave.Weaver) (err error) {
		// todo: always assumed to be a kind right now;
		// could auto switch, ex. prefer nouns if a match is found
		if tgt, e := safe.GetText(cat.Runtime(), op.Target); e != nil {
			err = e
		} else {
			pen := w.Pin()
			tgt := lang.Normalize(tgt.String())
			for _, h := range op.Handlers {
				// each handler references a pattern
				pb := mdl.NewPatternBuilder(h.Event.String())
				if e := addFields(pb, mdl.PatternLocals, h.Locals); e != nil {
					err = errutil.Append(e)
				} else if e := addRules(pb, tgt, h.Rules); e != nil {
					err = errutil.Append(e)
				} else if e := pen.ExtendPattern(pb.Pattern); e != nil {
					err = errutil.Append(e)
				}
			}
		}
		return
	})
}
