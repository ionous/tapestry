package story

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// Create a story dl from native maps and slices.
func Decode(dst jsn.Marshalee, msg map[string]any, reg cin.Signatures) error {
	return decode(dst, msg, reg)
}

func decode(dst jsn.Marshalee, val any, reg cin.TypeCreator) error {
	return cin.NewDecoder(reg).
		SetSlotDecoder(CompactSlotDecoder).
		DecodeValue(dst, val)
}

func CompactSlotDecoder(m jsn.Marshaler, slot jsn.SlotBlock, body any) (err error) {
	if err = core.CompactSlotDecoder(m, slot, body); err != nil {
		// keep this as the provisional error unless we figure out something else
		// ( important so that callers can see the original unhandled value if any )
		var unhandled chart.Unhandled
		if errors.As(err, &unhandled) {
			switch typeName := slot.GetType(); typeName {
			case StoryStatement_Type:
				if msg, ok := body.(map[string]any); ok {
					if sig, args, e := tryPattern(m, msg, typeName); e != nil {
						err = e
					} else if len(sig) > 0 {
						if !slot.SetSlot(&CallMacro{MacroName: sig, Arguments: args}) {
							err = errutil.New("unexpected error setting pattern slot", sig)
						} else {
							err = nil // clear the unhandled error
						}
					}
				}

			case
				rt.Execute_Type,
				rt.BoolEval_Type,
				rt.NumberEval_Type,
				rt.TextEval_Type,
				rt.RecordEval_Type,
				rt.NumListEval_Type,
				rt.TextListEval_Type,
				rt.RecordListEval_Type:
				if msg, ok := body.(map[string]any); ok {
					if sig, args, e := tryPattern(m, msg, typeName); e != nil {
						err = e
					} else if len(sig) > 0 {
						if !slot.SetSlot(&assign.CallPattern{PatternName: sig, Arguments: args}) {
							err = errutil.New("unexpected error setting pattern slot", sig)
						} else {
							err = nil // clear the unhandled error
							// ( it will error in the weave if no such pattern exists )
						}
					}
				}
			}
		}
	}
	return
}

// read the command as if it were a standard compact encoded golang struct.
// if it's *not* a normal command, treat it as a pattern call.
func tryPattern(m jsn.Marshaler, msg map[string]any, typeName string) (retSig string, retArgs []assign.Arg, err error) {
	// are we in fact parsing with the compact decoder?
	// if so, we can use its registry to figure out what's known and unknown.
	if reg, ok := m.(cin.TypeCreator); ok {
		// separate the signature, markup, and body of the specified op.
		// { "Signature:something:": [msg body], "--": "possible markup" }
		if op, e := cin.ReadOp(msg); e != nil {
			err = e
		} else if !reg.HasType(cin.Hash(op.Sig, typeName)) {
			if sig, args, e := op.ReadMsg(); e != nil {
				err = e
			} else {
				var call []assign.Arg
				for i, p := range sig.Params {
					var val rt.Assignment
					if e := decode(rt.Assignment_Slot{Value: &val}, args[i], reg); e != nil {
						err = e
						break
					}
					call = append(call, assign.Arg{Name: p.Label, Value: val})
				}
				if len(sig.Params) == len(call) {
					retSig, retArgs = sig.Name, call
				}
			}
		}
	}
	return
}
