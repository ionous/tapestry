// system commands
package game

//
// Code generated by Tapestry; edit at your own risk.
//

import (
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

type QuitGame struct {
	Markup map[string]any
}

// quit_game, a type of flow.
var Zt_QuitGame typeinfo.Flow

// implements typeinfo.Instance
func (*QuitGame) TypeInfo() typeinfo.T {
	return &Zt_QuitGame
}

// implements typeinfo.Markup
func (op *QuitGame) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*QuitGame)(nil)

// holds a slice of type quit_game
type QuitGame_Slice []QuitGame

// implements typeinfo.Instance
func (*QuitGame_Slice) TypeInfo() typeinfo.T {
	return &Zt_QuitGame
}

// implements typeinfo.Repeats
func (op *QuitGame_Slice) Repeats() bool {
	return len(*op) > 0
}

type SaveGame struct {
	Markup map[string]any
}

// save_game, a type of flow.
var Zt_SaveGame typeinfo.Flow

// implements typeinfo.Instance
func (*SaveGame) TypeInfo() typeinfo.T {
	return &Zt_SaveGame
}

// implements typeinfo.Markup
func (op *SaveGame) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*SaveGame)(nil)

// holds a slice of type save_game
type SaveGame_Slice []SaveGame

// implements typeinfo.Instance
func (*SaveGame_Slice) TypeInfo() typeinfo.T {
	return &Zt_SaveGame
}

// implements typeinfo.Repeats
func (op *SaveGame_Slice) Repeats() bool {
	return len(*op) > 0
}

type LoadGame struct {
	Markup map[string]any
}

// load_game, a type of flow.
var Zt_LoadGame typeinfo.Flow

// implements typeinfo.Instance
func (*LoadGame) TypeInfo() typeinfo.T {
	return &Zt_LoadGame
}

// implements typeinfo.Markup
func (op *LoadGame) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*LoadGame)(nil)

// holds a slice of type load_game
type LoadGame_Slice []LoadGame

// implements typeinfo.Instance
func (*LoadGame_Slice) TypeInfo() typeinfo.T {
	return &Zt_LoadGame
}

// implements typeinfo.Repeats
func (op *LoadGame_Slice) Repeats() bool {
	return len(*op) > 0
}

type UndoTurn struct {
	Markup map[string]any
}

// undo_turn, a type of flow.
var Zt_UndoTurn typeinfo.Flow

// implements typeinfo.Instance
func (*UndoTurn) TypeInfo() typeinfo.T {
	return &Zt_UndoTurn
}

// implements typeinfo.Markup
func (op *UndoTurn) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*UndoTurn)(nil)

// holds a slice of type undo_turn
type UndoTurn_Slice []UndoTurn

// implements typeinfo.Instance
func (*UndoTurn_Slice) TypeInfo() typeinfo.T {
	return &Zt_UndoTurn
}

// implements typeinfo.Repeats
func (op *UndoTurn_Slice) Repeats() bool {
	return len(*op) > 0
}

type PrintVersion struct {
	Markup map[string]any
}

// print_version, a type of flow.
var Zt_PrintVersion typeinfo.Flow

// implements typeinfo.Instance
func (*PrintVersion) TypeInfo() typeinfo.T {
	return &Zt_PrintVersion
}

// implements typeinfo.Markup
func (op *PrintVersion) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*PrintVersion)(nil)

// holds a slice of type print_version
type PrintVersion_Slice []PrintVersion

// implements typeinfo.Instance
func (*PrintVersion_Slice) TypeInfo() typeinfo.T {
	return &Zt_PrintVersion
}

// implements typeinfo.Repeats
func (op *PrintVersion_Slice) Repeats() bool {
	return len(*op) > 0
}

// init the terms of all flows in init
// so that they can refer to each other when needed.
func init() {
	Zt_QuitGame = typeinfo.Flow{
		Name:  "quit_game",
		Lede:  "quit_game",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
	}
	Zt_SaveGame = typeinfo.Flow{
		Name:  "save_game",
		Lede:  "save_game",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
	}
	Zt_LoadGame = typeinfo.Flow{
		Name:  "load_game",
		Lede:  "load_game",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
	}
	Zt_UndoTurn = typeinfo.Flow{
		Name:  "undo_turn",
		Lede:  "undo_turn",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
	}
	Zt_PrintVersion = typeinfo.Flow{
		Name:  "print_version",
		Lede:  "print_version",
		Terms: []typeinfo.Term{},
		Slots: []*typeinfo.Slot{
			&rtti.Zt_Execute,
		},
	}
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "game",
	Comment: []string{
		"system commands",
	},

	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_QuitGame,
	&Zt_SaveGame,
	&Zt_LoadGame,
	&Zt_UndoTurn,
	&Zt_PrintVersion,
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]typeinfo.Instance{
	7048942477941640285:  (*LoadGame)(nil),     /* execute=LoadGame */
	16069653899165369674: (*PrintVersion)(nil), /* execute=PrintVersion */
	13962506025236193050: (*QuitGame)(nil),     /* execute=QuitGame */
	12343662000108026632: (*SaveGame)(nil),     /* execute=SaveGame */
	6128819475946940678:  (*UndoTurn)(nil),     /* execute=UndoTurn */
}
