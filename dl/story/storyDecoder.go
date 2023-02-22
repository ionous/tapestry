package story

import (
	"encoding/json"
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/assign"
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
		SetFlowDecoder(core.CompactFlowDecoder).
		SetSlotDecoder(CompactSlotDecoder).
		Decode(dst, msg)

	// re: flow decoder
	// right now, we only get here when the flow is a member of a flow;
	// when we know what the type is but havent tried to read from the msg yet.
	// for slots to read a slat, the msg already would have been expanded into its final type.
}

func CompactSlotDecoder(m jsn.Marshaler, slot jsn.SlotBlock, msg json.RawMessage) (err error) {
	if err = core.CompactSlotDecoder(m, slot, msg); err != nil {
		// keep this as the provisional error unless we figure out something else
		// ( important so that callers can see the original unhandled value if any )
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
				// read the command as if it were a standard compact encoded golang struct.
				// ( at this point, we can overwrite the unhandled error since we are trying to handle it. )
				if op, e := cin.ReadOp(msg); e != nil {
					err = e
				} else if reg := m.(cin.TypeCreator); !reg.HasType(cin.Hash(op.Sig, typeName)) {
					// if we didn't find it, then we'll treat it as a pattern call.
					// ( it will error out later in assembly if no such pattern exists )
					if sig, args, e := op.ReadMsg(); e != nil {
						err = e
					} else {
						out := &core.CallPattern{Pattern: core.PatternName{sig.Name}}
						var call []core.Arg
						for i, p := range sig.Params {
							arg := args[i]
							// fix: temp: backwards compat:
							var str string
							var flag bool
							var num float64
							var val assign.Assignment
							if e := json.Unmarshal(arg, &str); e == nil {
								val = &assign.FromText{Value: T(str)}
							} else if e := json.Unmarshal(arg, &flag); e == nil {
								val = &assign.FromBool{Value: B(flag)}
							} else if e := json.Unmarshal(arg, &num); e == nil {
								val = &assign.FromNumber{Value: F(num)}
							} else if e := decode(assign.Assignment_Slot{&val}, arg, reg); e != nil {
								err = e
								break
							}
							call = append(call, core.Arg{Name: p.Label, Value: val})
						}
						if len(sig.Params) == len(call) {
							out.Arguments = call
							if slot.SetSlot(out) {
								err = nil // good to go
							} else {
								// we tried handle the command, parsed it even, but failed to assign it:
								// report that back as an error.
								err = errutil.New("unexpected error setting slot")
							}
						}
					}
				}
			}
		}
	}
	return
}
