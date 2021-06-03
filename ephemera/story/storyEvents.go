package story

import (
	"git.sr.ht/~ionous/iffy/ephemera"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

func (op *EventBlock) ImportPhrase(k *Importer) (err error) {
	if opt, ok := op.Target.Opt.(interface {
		NewName(*Importer) (ephemera.Named, error)
	}); !ok {
		err = errutil.Fmt("Unknown interface %T", opt)
	} else if tgt, e := opt.NewName(k); e != nil {
		err = e
	} else {
		if err == nil {
			for _, h := range op.Handlers {
				if evt, e := NewEventName(k, h.Event); e != nil {
					err = e
					break
				} else if flags, e := h.EventPhase.ReadFlags(k); e != nil {
					err = e
					break
				} else if e := h.PatternRules.ImportRules(k, evt, tgt, flags); e != nil {
					err = e
				} else if h.Locals != nil {
					if e := h.Locals.ImportLocals(k, evt); e != nil {
						break
					}
				}
			}
		}
	}
	return
}

func (op *EventPhase) ReadFlags(k *Importer) (ret rt.Flags, err error) {
	switch str := op.Str; str {
	case EventPhase_Before:
		ret = rt.Prefix
	case EventPhase_While:
		ret = rt.Infix
	case EventPhase_After:
		ret = rt.Postfix
	default:
		if len(str) > 0 {
			err = errutil.Fmt("unknown event flags %q", str)
		}
	}
	return
}
