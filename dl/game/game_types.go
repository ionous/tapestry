// Code generated by Tapestry; edit at your own risk.
package game

import (
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_QuitGame struct {
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*QuitGame) Inspect() typeinfo.T {
	return &Z_QuitGame_T
}

// return a valid markup map, creating it if necessary.
func (op *QuitGame) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// quit_game, a type of flow.
const Z_QuitGame_Name = "quit_game"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*QuitGame)(nil)

var Z_QuitGame_T = typeinfo.Flow{
	Name:  Z_QuitGame_Name,
	Lede:  "quit_game",
	Terms: []typeinfo.Term{},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_T,
	},
}

// holds a slice of type quit_game
// FIX: duplicates the spec decl.
type FIX_QuitGame_Slice []QuitGame

// implements typeinfo.Inspector
func (*QuitGame_Slice) Inspect() typeinfo.T {
	return &Z_QuitGame_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_SaveGame struct {
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*SaveGame) Inspect() typeinfo.T {
	return &Z_SaveGame_T
}

// return a valid markup map, creating it if necessary.
func (op *SaveGame) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// save_game, a type of flow.
const Z_SaveGame_Name = "save_game"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*SaveGame)(nil)

var Z_SaveGame_T = typeinfo.Flow{
	Name:  Z_SaveGame_Name,
	Lede:  "save_game",
	Terms: []typeinfo.Term{},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_T,
	},
}

// holds a slice of type save_game
// FIX: duplicates the spec decl.
type FIX_SaveGame_Slice []SaveGame

// implements typeinfo.Inspector
func (*SaveGame_Slice) Inspect() typeinfo.T {
	return &Z_SaveGame_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_RestoreGame struct {
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*RestoreGame) Inspect() typeinfo.T {
	return &Z_RestoreGame_T
}

// return a valid markup map, creating it if necessary.
func (op *RestoreGame) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// restore_game, a type of flow.
const Z_RestoreGame_Name = "restore_game"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*RestoreGame)(nil)

var Z_RestoreGame_T = typeinfo.Flow{
	Name:  Z_RestoreGame_Name,
	Lede:  "restore_game",
	Terms: []typeinfo.Term{},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_T,
	},
}

// holds a slice of type restore_game
// FIX: duplicates the spec decl.
type FIX_RestoreGame_Slice []RestoreGame

// implements typeinfo.Inspector
func (*RestoreGame_Slice) Inspect() typeinfo.T {
	return &Z_RestoreGame_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_UndoTurn struct {
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*UndoTurn) Inspect() typeinfo.T {
	return &Z_UndoTurn_T
}

// return a valid markup map, creating it if necessary.
func (op *UndoTurn) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// undo_turn, a type of flow.
const Z_UndoTurn_Name = "undo_turn"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*UndoTurn)(nil)

var Z_UndoTurn_T = typeinfo.Flow{
	Name:  Z_UndoTurn_Name,
	Lede:  "undo_turn",
	Terms: []typeinfo.Term{},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_T,
	},
}

// holds a slice of type undo_turn
// FIX: duplicates the spec decl.
type FIX_UndoTurn_Slice []UndoTurn

// implements typeinfo.Inspector
func (*UndoTurn_Slice) Inspect() typeinfo.T {
	return &Z_UndoTurn_T
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_PrintVersion struct {
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*PrintVersion) Inspect() typeinfo.T {
	return &Z_PrintVersion_T
}

// return a valid markup map, creating it if necessary.
func (op *PrintVersion) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// print_version, a type of flow.
const Z_PrintVersion_Name = "print_version"

// ensure the command implements its specified slots:
var _ rtti.Execute = (*PrintVersion)(nil)

var Z_PrintVersion_T = typeinfo.Flow{
	Name:  Z_PrintVersion_Name,
	Lede:  "print_version",
	Terms: []typeinfo.Term{},
	Slots: []*typeinfo.Slot{
		&rtti.Z_Execute_T,
	},
}

// holds a slice of type print_version
// FIX: duplicates the spec decl.
type FIX_PrintVersion_Slice []PrintVersion

// implements typeinfo.Inspector
func (*PrintVersion_Slice) Inspect() typeinfo.T {
	return &Z_PrintVersion_T
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "game",
	Flow: z_flow_list,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Z_QuitGame_T,
	&Z_SaveGame_T,
	&Z_RestoreGame_T,
	&Z_UndoTurn_T,
	&Z_PrintVersion_T,
}