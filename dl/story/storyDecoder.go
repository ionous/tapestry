package story

import (
	"encoding/json"
	"errors"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// Read a story from a story file.
func DecodeJson(dst jsn.Marshalee, msg json.RawMessage, reg cin.Signatures) error {
	return cin.NewDecoder(reg).
		SetSlotDecoder(CompactSlotDecoder).
		DecodeJson(dst, msg)
}

func decode(dst jsn.Marshalee, obj r.Value, reg cin.TypeCreator) error {
	return cin.NewDecoder(reg).
		SetSlotDecoder(CompactSlotDecoder).
		DecodeMsg(dst, obj)
}

func CompactSlotDecoder(m jsn.Marshaler, slot jsn.SlotBlock, msg r.Value) (err error) {
	if err = core.CompactSlotDecoder(m, slot, msg); err != nil {
		// keep this as the provisional error unless we figure out something else
		// ( important so that callers can see the original unhandled value if any )
		var unhandled chart.Unhandled
		if errors.As(err, &unhandled) {
			switch typeName := slot.GetType(); typeName {
			case StoryStatement_Type:
				if sig, args, e := tryPattern(m, msg, typeName); e != nil {
					err = e
				} else if len(sig) > 0 {
					if !slot.SetSlot(&CallMacro{MacroName: sig, Arguments: args}) {
						err = errutil.New("unexpected error setting pattern slot", sig)
					} else {
						err = nil // clear the unhandled error
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
	return
}

// read the command as if it were a standard compact encoded golang struct.
// if we don't find it, then we'll treat it as a pattern call.
func tryPattern(m jsn.Marshaler, msg r.Value, typeName string) (retSig string, retArgs []assign.Arg, err error) {
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
					el := args.Index(i).Elem()
					if e := decode(rt.Assignment_Slot{Value: &val}, el, reg); e != nil {
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
