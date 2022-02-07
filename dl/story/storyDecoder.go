package story

import (
	"encoding/json"
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// Read a story from a story file.
func Decode(dst jsn.Marshalee, msg json.RawMessage, sig cin.Signatures) error {
	return decode(dst, msg, sig)
}

func decode(dst jsn.Marshalee, msg json.RawMessage, reg cin.TypeCreator) error {
	return cin.NewDecoder(reg).
		//SetFlowDecoder(CompactFlowDecoder).
		SetSlotDecoder(CompactSlotDecoder).
		Decode(dst, msg)

	// re: flow decoder
	// right now, we only get here when the flow is a member of a flow;
	// when we know what the type is but havent tried to read from the msg yet.
	// for slots to read a slat, the msg already would have been expanded into its final type.
}

//
func CompactSlotDecoder(m jsn.Marshaler, slot jsn.SlotBlock, msg json.RawMessage) (err error) {
	if err = core.CompactSlotDecoder(m, slot, msg); err != nil {
		// keep this as the provisional error unless we figure out something else
		var unhandled chart.Unhandled
		if errors.As(err, &unhandled) {
			switch typeName := slot.GetType(); typeName {
			case
				rt.Execute_Type,
				rt.BoolEval_Type,
				rt.NumberEval_Type,
				rt.TextEval_Type,
				rt.RecordEval_Type,
				rt.NumListEval_Type,
				rt.TextListEval_Type,
				rt.RecordListEval_Type:
				//
				if op, e := cin.ReadOp(msg); e != nil {
					err = e // if we cant parse it, the regular compact marshal isn't going to be able to either.
				} else if reg := m.(cin.TypeCreator); !reg.HasType(op.Key) {
					if sig, args, e := op.ReadMsg(); e != nil {
						err = e
					} else {
						out := &core.CallPattern{Pattern: core.PatternName{sig.Name}}
						if len(sig.Params) > 0 {
							var call []rt.Arg
							for i, p := range sig.Params {
								var ptr rt.Assignment
								if e := decode(rt.Assignment_Slot{&ptr}, args[i], reg); e != nil {
									err = e
									break
								} else {
									call = append(call, rt.Arg{Name: p.Label, From: ptr})
								}
							}
							out.Arguments.Args = call
						}
						if !slot.SetSlot(out) {
							err = errutil.New("unexpected error setting slot")
						} else {
							err = nil // good to go,
						}
					}
				}
			}
		}
	}
	return
}
