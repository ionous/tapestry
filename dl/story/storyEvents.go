package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
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
	return imp.StoryStatement(macro, op)
}

func (op *EventBlock) PostImport(k *imp.Importer) (err error) {
	if opt, ok := op.Target.Value.(interface{ GetName() string }); !ok {
		err = errutil.Fmt("unknown event block target %T", opt)
	} else {
		tgt := opt.GetName()
		// each handler is a rule...
		for _, h := range op.Handlers {
			evt := h.Event.String()
			if flags, e := h.EventPhase.ReadFlags(); e != nil {
				err = errutil.Append(e)
			} else if e := ImportRules(k, evt, tgt, h.Rules, flags); e != nil {
				err = errutil.Append(e)
			} else if locals := ImportLocals(h.Locals); len(locals) > 0 {
				// and these are locals used by those rules
				k.WriteEphemera(&eph.EphPatterns{PatternName: evt, Locals: locals})
			}
		}
	}
	return
}

func (op *EventPhase) ReadFlags() (ret eph.EphTiming, err error) {
	switch str := op.Str; str {
	case EventPhase_Before:
		ret = eph.EphTiming{eph.EphTiming_Before}
	case EventPhase_While:
		ret = eph.EphTiming{eph.EphTiming_During}
	case EventPhase_After:
		ret = eph.EphTiming{eph.EphTiming_After}
	default:
		if len(str) > 0 {
			err = errutil.Fmt("unknown event flags %q", str)
		}
	}
	return
}
