package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

func (op *PluralKinds) GetName() string {
	return op.String()
}

func (op *SingularNoun) GetName() string {
	return op.Name
}

// Execute - called by the macro runtime during weave.
func (op *EventBlock) Execute(macro rt.Runtime) error {
	return weave.StoryStatement(macro, op)
}

func (op *EventBlock) Schedule(cat *weave.Catalog) (err error) {
	if opt, ok := op.Target.Value.(interface{ GetName() string }); !ok {
		err = errutil.Fmt("unknown event block target %T", opt)
	} else {
		//
		tgt := opt.GetName()
		// each handler is a rule...
		for _, h := range op.Handlers {
			evt := h.Event.String()
			if flags, e := h.EventPhase.ReadFlags(); e != nil {
				err = errutil.Append(e)
			} else if e := ImportRules(cat, evt, tgt, h.Rules, flags); e != nil {
				err = errutil.Append(e)
			} else {
				// and these are locals used by those rules
				err = declareFields(h.Locals, func(name, class string, aff affine.Affinity, init assign.Assignment) (err error) {
					return cat.AssertField(evt, name, class, aff, init)
				})
			}
		}
	}
	return
}

func (op *EventPhase) ReadFlags() (ret assert.EventTiming, err error) {
	switch str := op.Str; str {
	case EventPhase_Before:
		ret = assert.Before
	case EventPhase_While:
		ret = assert.During
	case EventPhase_After:
		ret = assert.After
	default:
		if len(str) > 0 {
			err = errutil.Fmt("unknown event flags %q", str)
		}
	}
	return
}
