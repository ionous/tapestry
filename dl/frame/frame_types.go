// Code generated by Tapestry; edit at your own risk.
package frame

import (
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// event, a type of slot.
const Z_Event_Name = "event"

var Z_Event_T = typeinfo.Slot{
	Name: Z_Event_Name,
}

// holds a single slot
// FIX: currently provided by the spec
type FIX_Event_Slot struct{ Value Event }

// implements typeinfo.Inspector for a single slot.
func (*FIX_Event_Slot) Inspect() typeinfo.T {
	return &Z_Event_T
}

// holds a slice of slots
type Event_Slots []Event

// implements typeinfo.Inspector for a series of slots.
func (*Event_Slots) Inspect() typeinfo.T {
	return &Z_Event_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_Frame struct {
	Result string
	Events Event
	Error  string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Frame) Inspect() typeinfo.T {
	return &Z_Frame_T
}

// return a valid markup map, creating it if necessary.
func (op *Frame) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// frame, a type of flow.
const Z_Frame_Name = "frame"

var Z_Frame_T = typeinfo.Flow{
	Name: Z_Frame_Name,
	Lede: "frame",
	Terms: []typeinfo.Term{{
		Name:  "result",
		Label: "result",
		Type:  &prim.Z_Text_T,
	}, {
		Name:    "events",
		Label:   "events",
		Repeats: true,
		Type:    &Z_Event_T,
	}, {
		Name:     "error",
		Label:    "error",
		Optional: true,
		Type:     &prim.Z_Text_T,
	}},
}

// holds a slice of type frame
// FIX: duplicates the spec decl.
type FIX_Frame_Slice []Frame

// implements typeinfo.Inspector
func (*Frame_Slice) Inspect() typeinfo.T {
	return &Z_Frame_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_FrameOutput struct {
	Text   string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FrameOutput) Inspect() typeinfo.T {
	return &Z_FrameOutput_T
}

// return a valid markup map, creating it if necessary.
func (op *FrameOutput) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// frame_output, a type of flow.
const Z_FrameOutput_Name = "frame_output"

// ensure the command implements its specified slots:
var _ Event = (*FrameOutput)(nil)

var Z_FrameOutput_T = typeinfo.Flow{
	Name: Z_FrameOutput_Name,
	Lede: "frame_output",
	Terms: []typeinfo.Term{{
		Name:  "text",
		Label: "_",
		Type:  &prim.Z_Text_T,
	}},
	Slots: []*typeinfo.Slot{
		&Z_Event_T,
	},
}

// holds a slice of type frame_output
// FIX: duplicates the spec decl.
type FIX_FrameOutput_Slice []FrameOutput

// implements typeinfo.Inspector
func (*FrameOutput_Slice) Inspect() typeinfo.T {
	return &Z_FrameOutput_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_SceneStarted struct {
	Domains string
	Markup  map[string]any
}

// implements typeinfo.Inspector
func (*SceneStarted) Inspect() typeinfo.T {
	return &Z_SceneStarted_T
}

// return a valid markup map, creating it if necessary.
func (op *SceneStarted) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// scene_started, a type of flow.
const Z_SceneStarted_Name = "scene_started"

// ensure the command implements its specified slots:
var _ Event = (*SceneStarted)(nil)

var Z_SceneStarted_T = typeinfo.Flow{
	Name: Z_SceneStarted_Name,
	Lede: "scene_started",
	Terms: []typeinfo.Term{{
		Name:    "domains",
		Label:   "_",
		Repeats: true,
		Type:    &prim.Z_Text_T,
	}},
	Slots: []*typeinfo.Slot{
		&Z_Event_T,
	},
}

// holds a slice of type scene_started
// FIX: duplicates the spec decl.
type FIX_SceneStarted_Slice []SceneStarted

// implements typeinfo.Inspector
func (*SceneStarted_Slice) Inspect() typeinfo.T {
	return &Z_SceneStarted_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_SceneEnded struct {
	Domains string
	Markup  map[string]any
}

// implements typeinfo.Inspector
func (*SceneEnded) Inspect() typeinfo.T {
	return &Z_SceneEnded_T
}

// return a valid markup map, creating it if necessary.
func (op *SceneEnded) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// scene_ended, a type of flow.
const Z_SceneEnded_Name = "scene_ended"

// ensure the command implements its specified slots:
var _ Event = (*SceneEnded)(nil)

var Z_SceneEnded_T = typeinfo.Flow{
	Name: Z_SceneEnded_Name,
	Lede: "scene_ended",
	Terms: []typeinfo.Term{{
		Name:    "domains",
		Label:   "_",
		Repeats: true,
		Type:    &prim.Z_Text_T,
	}},
	Slots: []*typeinfo.Slot{
		&Z_Event_T,
	},
}

// holds a slice of type scene_ended
// FIX: duplicates the spec decl.
type FIX_SceneEnded_Slice []SceneEnded

// implements typeinfo.Inspector
func (*SceneEnded_Slice) Inspect() typeinfo.T {
	return &Z_SceneEnded_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_StateChanged struct {
	Noun   string
	Aspect string
	Prev   string
	Trait  string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*StateChanged) Inspect() typeinfo.T {
	return &Z_StateChanged_T
}

// return a valid markup map, creating it if necessary.
func (op *StateChanged) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// state_changed, a type of flow.
const Z_StateChanged_Name = "state_changed"

// ensure the command implements its specified slots:
var _ Event = (*StateChanged)(nil)

var Z_StateChanged_T = typeinfo.Flow{
	Name: Z_StateChanged_Name,
	Lede: "state_changed",
	Terms: []typeinfo.Term{{
		Name:  "noun",
		Label: "noun",
		Type:  &prim.Z_Text_T,
	}, {
		Name:  "aspect",
		Label: "aspect",
		Type:  &prim.Z_Text_T,
	}, {
		Name:  "prev",
		Label: "prev",
		Type:  &prim.Z_Text_T,
	}, {
		Name:  "trait",
		Label: "trait",
		Type:  &prim.Z_Text_T,
	}},
	Slots: []*typeinfo.Slot{
		&Z_Event_T,
	},
}

// holds a slice of type state_changed
// FIX: duplicates the spec decl.
type FIX_StateChanged_Slice []StateChanged

// implements typeinfo.Inspector
func (*StateChanged_Slice) Inspect() typeinfo.T {
	return &Z_StateChanged_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_PairChanged struct {
	A      string
	B      string
	Rel    string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*PairChanged) Inspect() typeinfo.T {
	return &Z_PairChanged_T
}

// return a valid markup map, creating it if necessary.
func (op *PairChanged) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// pair_changed, a type of flow.
const Z_PairChanged_Name = "pair_changed"

// ensure the command implements its specified slots:
var _ Event = (*PairChanged)(nil)

var Z_PairChanged_T = typeinfo.Flow{
	Name: Z_PairChanged_Name,
	Lede: "pair_changed",
	Terms: []typeinfo.Term{{
		Name:  "a",
		Label: "a",
		Type:  &prim.Z_Text_T,
	}, {
		Name:  "b",
		Label: "b",
		Type:  &prim.Z_Text_T,
	}, {
		Name:  "rel",
		Label: "rel",
		Type:  &prim.Z_Text_T,
	}},
	Slots: []*typeinfo.Slot{
		&Z_Event_T,
	},
}

// holds a slice of type pair_changed
// FIX: duplicates the spec decl.
type FIX_PairChanged_Slice []PairChanged

// implements typeinfo.Inspector
func (*PairChanged_Slice) Inspect() typeinfo.T {
	return &Z_PairChanged_T
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "frame",
	Slot: z_slot_list,
	Flow: z_flow_list,
}

// a list of all slots in this this package
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Z_Event_T,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Z_Frame_T,
	&Z_FrameOutput_T,
	&Z_SceneStarted_T,
	&Z_SceneEnded_T,
	&Z_StateChanged_T,
	&Z_PairChanged_T,
}
