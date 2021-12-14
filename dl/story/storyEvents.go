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
		// each handler is a rule...
		for _, h := range op.Handlers {
			evt := h.Event.String()
			if flags, e := h.EventPhase.ReadFlags(k); e != nil {
				err = errutil.Append(e)
			} else if e := h.PatternRules.ImportRules(k, evt, tgt, flags); e != nil {
				err = errutil.Append(e)
			} else {
				// and these are locals used by those rules
				if h.Locals != nil {
					if locals, e := h.Locals.ImportLocals(k, evt); e != nil {
						err = errutil.Append(e)
					} else if len(locals) > 0 {
						k.Write(&eph.EphPatterns{Name: evt, Locals: locals})
					}
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
