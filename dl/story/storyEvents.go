package story

import (
	"git.sr.ht/~ionous/iffy/dl/eph"
	"github.com/ionous/errutil"
)

func (op *PluralKinds) GetName(k *Importer) (ret string, err error) {
	ret = op.String()
	return
}

func (op *NamedNoun) GetName(k *Importer) (ret string, err error) {
	err = errutil.New("named noun for event block not implemented")
	return
}

func (op *EventBlock) ImportPhrase(k *Importer) (err error) {
	if opt, ok := op.Target.Value.(interface {
		GetName(*Importer) (string, error)
	}); !ok {
		err = errutil.Fmt("unknown event block target %T at %s", opt, op.At)
	} else if tgt, e := opt.GetName(k); e != nil {
		err = e
	} else {
		if err == nil {
			for _, h := range op.Handlers {
				evt := h.Event.String()
				if flags, e := h.EventPhase.ReadFlags(k); e != nil {
					err = e
					break
				} else if e := h.PatternRules.ImportRules(k, evt, tgt, flags); e != nil {
					err = e
				} else {
					println("EventBlock not implemented")
					// if h.Locals != nil {
					// 	if e := h.Locals.ImportLocals(k, evt); e != nil {
					// 		break
					// 	}
					// }
				}
			}
		}
	}
	return
}

func (op *EventPhase) ReadFlags(k *Importer) (ret eph.EphTiming, err error) {
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
