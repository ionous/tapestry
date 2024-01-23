// Code generated by Tapestry; edit at your own risk.
package play

import (
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// play_message, a type of slot.
const Z_PlayMessage_Type = "play_message"

var Z_PlayMessage_Info = typeinfo.Slot{
	Name: Z_PlayMessage_Type,
	Markup: map[string]any{
		"comment": "a client-server message for the play app",
	},
}

// holds a single slot
// FIX: currently provided by the spec
type FIX_PlayMessage_Slot struct{ Value PlayMessage }

// implements typeinfo.Inspector for a single slot.
func (*FIX_PlayMessage_Slot) Inspect() typeinfo.T {
	return &Z_PlayMessage_Info
}

// holds a slice of slots
type PlayMessage_Slots []PlayMessage

// implements typeinfo.Inspector for a series of slots.
func (*PlayMessage_Slots) Inspect() typeinfo.T {
	return &Z_PlayMessage_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_PlayLog struct {
	Log    string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*PlayLog) Inspect() typeinfo.T {
	return &Z_PlayLog_Info
}

// return a valid markup map, creating it if necessary.
func (op *PlayLog) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// play_log, a type of flow.
const Z_PlayLog_Type = "play_log"

// ensure the command implements its specified slots:
var _ PlayMessage = (*PlayLog)(nil)

var Z_PlayLog_Info = typeinfo.Flow{
	Name: Z_PlayLog_Type,
	Lede: "play",
	Terms: []typeinfo.Term{{
		Name:  "log",
		Label: "log",
		Type:  &prim.Z_Text_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_PlayMessage_Info,
	},
	Markup: map[string]any{
		"comment": "a log message that can optionally be displayed to the client.",
	},
}

// holds a slice of type play_log
// FIX: duplicates the spec decl.
type FIX_PlayLog_Slice []PlayLog

// implements typeinfo.Inspector
func (*PlayLog_Slice) Inspect() typeinfo.T {
	return &Z_PlayLog_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_PlayMode struct {
	Mode   string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*PlayMode) Inspect() typeinfo.T {
	return &Z_PlayMode_Info
}

// return a valid markup map, creating it if necessary.
func (op *PlayMode) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// play_mode, a type of flow.
const Z_PlayMode_Type = "play_mode"

// ensure the command implements its specified slots:
var _ PlayMessage = (*PlayMode)(nil)

var Z_PlayMode_Info = typeinfo.Flow{
	Name: Z_PlayMode_Type,
	Lede: "play",
	Terms: []typeinfo.Term{{
		Name:  "mode",
		Label: "mode",
		Type:  &Z_PlayModes_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_PlayMessage_Info,
	},
	Markup: map[string]any{
		"comment": "app level change in state.",
	},
}

// holds a slice of type play_mode
// FIX: duplicates the spec decl.
type FIX_PlayMode_Slice []PlayMode

// implements typeinfo.Inspector
func (*PlayMode_Slice) Inspect() typeinfo.T {
	return &Z_PlayMode_Info
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_PlayOut struct {
	Out    string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*PlayOut) Inspect() typeinfo.T {
	return &Z_PlayOut_Info
}

// return a valid markup map, creating it if necessary.
func (op *PlayOut) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// play_out, a type of flow.
const Z_PlayOut_Type = "play_out"

// ensure the command implements its specified slots:
var _ PlayMessage = (*PlayOut)(nil)

var Z_PlayOut_Info = typeinfo.Flow{
	Name: Z_PlayOut_Type,
	Lede: "play",
	Terms: []typeinfo.Term{{
		Name:  "out",
		Label: "out",
		Type:  &prim.Z_Text_Info,
	}},
	Slots: []*typeinfo.Slot{
		&Z_PlayMessage_Info,
	},
	Markup: map[string]any{
		"comment": "output from the game itself.",
	},
}

// holds a slice of type play_out
// FIX: duplicates the spec decl.
type FIX_PlayOut_Slice []PlayOut

// implements typeinfo.Inspector
func (*PlayOut_Slice) Inspect() typeinfo.T {
	return &Z_PlayOut_Info
}

// play_modes, a type of str enum.
const Z_PlayModes_Type = "play_modes"

const (
	W_PlayModes_Asm      = "$ASM"
	W_PlayModes_Play     = "$PLAY"
	W_PlayModes_Complete = "$COMPLETE"
	W_PlayModes_Error    = "$ERROR"
)

var Z_PlayModes_Info = typeinfo.Str{
	Name: Z_PlayModes_Type,
	Options: []string{
		W_PlayModes_Asm,
		W_PlayModes_Play,
		W_PlayModes_Complete,
		W_PlayModes_Error,
	},
	Markup: map[string]any{
		"comment": "enum for play play_mode",
	},
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "play",
	Slot: z_slot_list,
	Flow: z_flow_list,
}

// a list of all slots in this this package
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Z_PlayMessage_Info,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Z_PlayLog_Info,
	&Z_PlayMode_Info,
	&Z_PlayOut_Info,
}
