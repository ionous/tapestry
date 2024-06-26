// Low level communication with a game server.
// None of these commands are available for use in game scripts.
package frame

//
// Code generated by Tapestry; edit at your own risk.
//

import (
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// notification, a type of slot.
var Zt_Notification = typeinfo.Slot{
	Name: "notification",
	Markup: map[string]any{
		"comment": "Marker interface used by all frame events.",
	},
}

// Holds a single slot.
type Notification_Slot struct{ Value Notification }

// Implements [typeinfo.Instance] for a single slot.
func (*Notification_Slot) TypeInfo() typeinfo.T {
	return &Zt_Notification
}

// Holds a slice of slots.
type Notification_Slots []Notification

// Implements [typeinfo.Instance] for a slice of slots.
func (*Notification_Slots) TypeInfo() typeinfo.T {
	return &Zt_Notification
}

// Implements [typeinfo.Repeats] for a slice of slots.
func (op *Notification_Slots) Repeats() bool {
	return len(*op) > 0
}

// The results of a a player initiated turn, or other client to server query.
type Frame struct {
	Result string
	Events []Notification
	Error  string
	Markup map[string]any
}

// frame, a type of flow.
var Zt_Frame typeinfo.Flow

// Implements [typeinfo.Instance]
func (*Frame) TypeInfo() typeinfo.T {
	return &Zt_Frame
}

// Implements [typeinfo.Markup]
func (op *Frame) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Holds a slice of type Frame.
type Frame_Slice []Frame

// Implements [typeinfo.Instance] for a slice of Frame.
func (*Frame_Slice) TypeInfo() typeinfo.T {
	return &Zt_Frame
}

// Implements [typeinfo.Repeats] for a slice of Frame.
func (op *Frame_Slice) Repeats() bool {
	return len(*op) > 0
}

// Printed text that should be visible the player.
type FrameOutput struct {
	Text   string
	Markup map[string]any
}

// frame_output, a type of flow.
var Zt_FrameOutput typeinfo.Flow

// Implements [typeinfo.Instance]
func (*FrameOutput) TypeInfo() typeinfo.T {
	return &Zt_FrameOutput
}

// Implements [typeinfo.Markup]
func (op *FrameOutput) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Notification = (*FrameOutput)(nil)

// Holds a slice of type FrameOutput.
type FrameOutput_Slice []FrameOutput

// Implements [typeinfo.Instance] for a slice of FrameOutput.
func (*FrameOutput_Slice) TypeInfo() typeinfo.T {
	return &Zt_FrameOutput
}

// Implements [typeinfo.Repeats] for a slice of FrameOutput.
func (op *FrameOutput_Slice) Repeats() bool {
	return len(*op) > 0
}

// One or more scenes ( aka domain ) have started.
type SceneStarted struct {
	Domains []string
	Markup  map[string]any
}

// scene_started, a type of flow.
var Zt_SceneStarted typeinfo.Flow

// Implements [typeinfo.Instance]
func (*SceneStarted) TypeInfo() typeinfo.T {
	return &Zt_SceneStarted
}

// Implements [typeinfo.Markup]
func (op *SceneStarted) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Notification = (*SceneStarted)(nil)

// Holds a slice of type SceneStarted.
type SceneStarted_Slice []SceneStarted

// Implements [typeinfo.Instance] for a slice of SceneStarted.
func (*SceneStarted_Slice) TypeInfo() typeinfo.T {
	return &Zt_SceneStarted
}

// Implements [typeinfo.Repeats] for a slice of SceneStarted.
func (op *SceneStarted_Slice) Repeats() bool {
	return len(*op) > 0
}

// One or more scenes ( aka domain ) have finished.
type SceneEnded struct {
	Domains []string
	Markup  map[string]any
}

// scene_ended, a type of flow.
var Zt_SceneEnded typeinfo.Flow

// Implements [typeinfo.Instance]
func (*SceneEnded) TypeInfo() typeinfo.T {
	return &Zt_SceneEnded
}

// Implements [typeinfo.Markup]
func (op *SceneEnded) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Notification = (*SceneEnded)(nil)

// Holds a slice of type SceneEnded.
type SceneEnded_Slice []SceneEnded

// Implements [typeinfo.Instance] for a slice of SceneEnded.
func (*SceneEnded_Slice) TypeInfo() typeinfo.T {
	return &Zt_SceneEnded
}

// Implements [typeinfo.Repeats] for a slice of SceneEnded.
func (op *SceneEnded_Slice) Repeats() bool {
	return len(*op) > 0
}

// Some object in the game has changed from one state to another.
type StateChanged struct {
	Noun   string
	Aspect string
	Prev   string
	Trait  string
	Markup map[string]any
}

// state_changed, a type of flow.
var Zt_StateChanged typeinfo.Flow

// Implements [typeinfo.Instance]
func (*StateChanged) TypeInfo() typeinfo.T {
	return &Zt_StateChanged
}

// Implements [typeinfo.Markup]
func (op *StateChanged) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Notification = (*StateChanged)(nil)

// Holds a slice of type StateChanged.
type StateChanged_Slice []StateChanged

// Implements [typeinfo.Instance] for a slice of StateChanged.
func (*StateChanged_Slice) TypeInfo() typeinfo.T {
	return &Zt_StateChanged
}

// Implements [typeinfo.Repeats] for a slice of StateChanged.
func (op *StateChanged_Slice) Repeats() bool {
	return len(*op) > 0
}

// The relationship between two objects has changed. For instance, an object's whereabouts as indicated by a parent-child type relation.
type PairChanged struct {
	A      string
	B      string
	Rel    string
	Markup map[string]any
}

// pair_changed, a type of flow.
var Zt_PairChanged typeinfo.Flow

// Implements [typeinfo.Instance]
func (*PairChanged) TypeInfo() typeinfo.T {
	return &Zt_PairChanged
}

// Implements [typeinfo.Markup]
func (op *PairChanged) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// Ensures the command implements its specified slots.
var _ Notification = (*PairChanged)(nil)

// Holds a slice of type PairChanged.
type PairChanged_Slice []PairChanged

// Implements [typeinfo.Instance] for a slice of PairChanged.
func (*PairChanged_Slice) TypeInfo() typeinfo.T {
	return &Zt_PairChanged
}

// Implements [typeinfo.Repeats] for a slice of PairChanged.
func (op *PairChanged_Slice) Repeats() bool {
	return len(*op) > 0
}

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_Frame = typeinfo.Flow{
		Name: "frame",
		Lede: "frame",
		Terms: []typeinfo.Term{{
			Name:  "result",
			Label: "result",
			Markup: map[string]any{
				"comment": "The result of a query assignment sent to the server.",
			},
			Type: &prim.Zt_Text,
		}, {
			Name:    "events",
			Label:   "events",
			Repeats: true,
			Markup: map[string]any{
				"comment": "Changes as a result of a player's turn or a client query.",
			},
			Type: &Zt_Notification,
		}, {
			Name:     "error",
			Label:    "error",
			Optional: true,
			Markup: map[string]any{
				"comment": []string{"Any critical server errors.", "Returned to the client as a string", "so that it can be displayed to the player."},
			},
			Type: &prim.Zt_Text,
		}},
		Markup: map[string]any{
			"comment": "The results of a a player initiated turn, or other client to server query.",
		},
	}
	Zt_FrameOutput = typeinfo.Flow{
		Name: "frame_output",
		Lede: "frame_output",
		Terms: []typeinfo.Term{{
			Name: "text",
			Markup: map[string]any{
				"comment": "The text to show to the player.",
			},
			Type: &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_Notification,
		},
		Markup: map[string]any{
			"comment": "Printed text that should be visible the player.",
		},
	}
	Zt_SceneStarted = typeinfo.Flow{
		Name: "scene_started",
		Lede: "scene_started",
		Terms: []typeinfo.Term{{
			Name:    "domains",
			Repeats: true,
			Markup: map[string]any{
				"comment": "The names of the scenes.",
			},
			Type: &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_Notification,
		},
		Markup: map[string]any{
			"comment": "One or more scenes ( aka domain ) have started.",
		},
	}
	Zt_SceneEnded = typeinfo.Flow{
		Name: "scene_ended",
		Lede: "scene_ended",
		Terms: []typeinfo.Term{{
			Name:    "domains",
			Repeats: true,
			Markup: map[string]any{
				"comment": "The names of the scenes.",
			},
			Type: &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_Notification,
		},
		Markup: map[string]any{
			"comment": "One or more scenes ( aka domain ) have finished.",
		},
	}
	Zt_StateChanged = typeinfo.Flow{
		Name: "state_changed",
		Lede: "state_changed",
		Terms: []typeinfo.Term{{
			Name:  "noun",
			Label: "noun",
			Markup: map[string]any{
				"comment": "The object who's state has changed.",
			},
			Type: &prim.Zt_Text,
		}, {
			Name:  "aspect",
			Label: "aspect",
			Markup: map[string]any{
				"comment": "The name of the set of states involved.",
			},
			Type: &prim.Zt_Text,
		}, {
			Name:  "prev",
			Label: "prev",
			Markup: map[string]any{
				"comment": "The name of the old state.",
			},
			Type: &prim.Zt_Text,
		}, {
			Name:  "trait",
			Label: "trait",
			Markup: map[string]any{
				"comment": "The name of the new state.",
			},
			Type: &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_Notification,
		},
		Markup: map[string]any{
			"comment": "Some object in the game has changed from one state to another.",
		},
	}
	Zt_PairChanged = typeinfo.Flow{
		Name: "pair_changed",
		Lede: "pair_changed",
		Terms: []typeinfo.Term{{
			Name:  "a",
			Label: "a",
			Markup: map[string]any{
				"comment": "The id of the primary object.",
			},
			Type: &prim.Zt_Text,
		}, {
			Name:  "b",
			Label: "b",
			Markup: map[string]any{
				"comment": "The id of the secondary object.",
			},
			Type: &prim.Zt_Text,
		}, {
			Name:  "rel",
			Label: "rel",
			Markup: map[string]any{
				"comment": "The name of the relation in question.",
			},
			Type: &prim.Zt_Text,
		}},
		Slots: []*typeinfo.Slot{
			&Zt_Notification,
		},
		Markup: map[string]any{
			"comment": "The relationship between two objects has changed. For instance, an object's whereabouts as indicated by a parent-child type relation.",
		},
	}
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "frame",
	Comment: []string{
		"Low level communication with a game server.",
		"None of these commands are available for use in game scripts.",
	},

	Slot:       z_slot_list,
	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// A list of all slots in this this package.
// ( ex. for generating blockly shapes )
var z_slot_list = []*typeinfo.Slot{
	&Zt_Notification,
}

// A list of all flows in this this package.
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_Frame,
	&Zt_FrameOutput,
	&Zt_SceneStarted,
	&Zt_SceneEnded,
	&Zt_StateChanged,
	&Zt_PairChanged,
}

// gob like registration
func Register(reg func(any)) {
	reg((*Frame)(nil))
	reg((*FrameOutput)(nil))
	reg((*SceneStarted)(nil))
	reg((*SceneEnded)(nil))
	reg((*StateChanged)(nil))
	reg((*PairChanged)(nil))
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]typeinfo.Instance{
	14657663848717440116: (*Frame)(nil),        /* Frame result:events: */
	2438049115146588168:  (*Frame)(nil),        /* Frame result:events:error: */
	1471826827262877163:  (*FrameOutput)(nil),  /* notification=FrameOutput: */
	7717838589161657001:  (*PairChanged)(nil),  /* notification=PairChanged a:b:rel: */
	5707328743025875669:  (*SceneEnded)(nil),   /* notification=SceneEnded: */
	6789405419350772834:  (*SceneStarted)(nil), /* notification=SceneStarted: */
	4672682947926144893:  (*StateChanged)(nil), /* notification=StateChanged noun:aspect:prev:trait: */
}
