// Code generated by Tapestry; edit at your own risk.
package frame

import (
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// event, a type of slot.
var Zt_Event = typeinfo.Slot{
	Name: "event",
}

// holds a single slot
// FIX: currently provided by the spec
type Event_Slot struct{ Value Event }

// implements typeinfo.Inspector for a single slot.
func (*Event_Slot) Inspect() (typeinfo.T, bool) {
	return &Zt_Event, false
}

// holds a slice of slots
type Event_Slots []Event

// implements typeinfo.Inspector for a series of slots.
func (*Event_Slots) Inspect() (typeinfo.T, bool) {
	return &Zt_Event, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type Frame struct {
	Result string
	Events []Event
	Error  string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*Frame) Inspect() (typeinfo.T, bool) {
	return &Zt_Frame, false
}

// return a valid markup map, creating it if necessary.
func (op *Frame) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// frame, a type of flow.
var Zt_Frame = typeinfo.Flow{
	Name: "frame",
	Lede: "frame",
	Terms: []typeinfo.Term{{
		Name:  "result",
		Label: "result",
		Type:  &prim.Zt_Text,
	}, {
		Name:    "events",
		Label:   "events",
		Repeats: true,
		Type:    &Zt_Event,
	}, {
		Name:     "error",
		Label:    "error",
		Optional: true,
		Type:     &prim.Zt_Text,
	}},
}

// holds a slice of type frame
// FIX: duplicates the spec decl.
type Frame_Slice []Frame

// implements typeinfo.Inspector
func (*Frame_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_Frame, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FrameOutput struct {
	Text   string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*FrameOutput) Inspect() (typeinfo.T, bool) {
	return &Zt_FrameOutput, false
}

// return a valid markup map, creating it if necessary.
func (op *FrameOutput) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ Event = (*FrameOutput)(nil)

// frame_output, a type of flow.
var Zt_FrameOutput = typeinfo.Flow{
	Name: "frame_output",
	Lede: "frame_output",
	Terms: []typeinfo.Term{{
		Name: "text",
		Type: &prim.Zt_Text,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_Event,
	},
}

// holds a slice of type frame_output
// FIX: duplicates the spec decl.
type FrameOutput_Slice []FrameOutput

// implements typeinfo.Inspector
func (*FrameOutput_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_FrameOutput, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type SceneStarted struct {
	Domains []string
	Markup  map[string]any
}

// implements typeinfo.Inspector
func (*SceneStarted) Inspect() (typeinfo.T, bool) {
	return &Zt_SceneStarted, false
}

// return a valid markup map, creating it if necessary.
func (op *SceneStarted) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ Event = (*SceneStarted)(nil)

// scene_started, a type of flow.
var Zt_SceneStarted = typeinfo.Flow{
	Name: "scene_started",
	Lede: "scene_started",
	Terms: []typeinfo.Term{{
		Name:    "domains",
		Repeats: true,
		Type:    &prim.Zt_Text,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_Event,
	},
}

// holds a slice of type scene_started
// FIX: duplicates the spec decl.
type SceneStarted_Slice []SceneStarted

// implements typeinfo.Inspector
func (*SceneStarted_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_SceneStarted, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type SceneEnded struct {
	Domains []string
	Markup  map[string]any
}

// implements typeinfo.Inspector
func (*SceneEnded) Inspect() (typeinfo.T, bool) {
	return &Zt_SceneEnded, false
}

// return a valid markup map, creating it if necessary.
func (op *SceneEnded) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ Event = (*SceneEnded)(nil)

// scene_ended, a type of flow.
var Zt_SceneEnded = typeinfo.Flow{
	Name: "scene_ended",
	Lede: "scene_ended",
	Terms: []typeinfo.Term{{
		Name:    "domains",
		Repeats: true,
		Type:    &prim.Zt_Text,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_Event,
	},
}

// holds a slice of type scene_ended
// FIX: duplicates the spec decl.
type SceneEnded_Slice []SceneEnded

// implements typeinfo.Inspector
func (*SceneEnded_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_SceneEnded, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type StateChanged struct {
	Noun   string
	Aspect string
	Prev   string
	Trait  string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*StateChanged) Inspect() (typeinfo.T, bool) {
	return &Zt_StateChanged, false
}

// return a valid markup map, creating it if necessary.
func (op *StateChanged) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ Event = (*StateChanged)(nil)

// state_changed, a type of flow.
var Zt_StateChanged = typeinfo.Flow{
	Name: "state_changed",
	Lede: "state_changed",
	Terms: []typeinfo.Term{{
		Name:  "noun",
		Label: "noun",
		Type:  &prim.Zt_Text,
	}, {
		Name:  "aspect",
		Label: "aspect",
		Type:  &prim.Zt_Text,
	}, {
		Name:  "prev",
		Label: "prev",
		Type:  &prim.Zt_Text,
	}, {
		Name:  "trait",
		Label: "trait",
		Type:  &prim.Zt_Text,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_Event,
	},
}

// holds a slice of type state_changed
// FIX: duplicates the spec decl.
type StateChanged_Slice []StateChanged

// implements typeinfo.Inspector
func (*StateChanged_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_StateChanged, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type PairChanged struct {
	A      string
	B      string
	Rel    string
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*PairChanged) Inspect() (typeinfo.T, bool) {
	return &Zt_PairChanged, false
}

// return a valid markup map, creating it if necessary.
func (op *PairChanged) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ Event = (*PairChanged)(nil)

// pair_changed, a type of flow.
var Zt_PairChanged = typeinfo.Flow{
	Name: "pair_changed",
	Lede: "pair_changed",
	Terms: []typeinfo.Term{{
		Name:  "a",
		Label: "a",
		Type:  &prim.Zt_Text,
	}, {
		Name:  "b",
		Label: "b",
		Type:  &prim.Zt_Text,
	}, {
		Name:  "rel",
		Label: "rel",
		Type:  &prim.Zt_Text,
	}},
	Slots: []*typeinfo.Slot{
		&Zt_Event,
	},
}

// holds a slice of type pair_changed
// FIX: duplicates the spec decl.
type PairChanged_Slice []PairChanged

// implements typeinfo.Inspector
func (*PairChanged_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_PairChanged, true
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name:       "frame",
	Slot:       z_slot_list,
	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// a list of all slots in this this package
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_Event,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_Frame,
	&Zt_FrameOutput,
	&Zt_SceneStarted,
	&Zt_SceneEnded,
	&Zt_StateChanged,
	&Zt_PairChanged,
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]any{
	14657663848717440116: (*Frame)(nil),        /* Frame result:events: */
	2438049115146588168:  (*Frame)(nil),        /* Frame result:events:error: */
	4385780296792938688:  (*FrameOutput)(nil),  /* event=FrameOutput: */
	17021232753503984522: (*PairChanged)(nil),  /* event=PairChanged a:b:rel: */
	14005264853352099464: (*SceneEnded)(nil),   /* event=SceneEnded: */
	12902248384806780167: (*SceneStarted)(nil), /* event=SceneStarted: */
	7027046405509259850:  (*StateChanged)(nil), /* event=StateChanged noun:aspect:prev:trait: */
}
