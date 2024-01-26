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
func (*QuitGame) Inspect() (typeinfo.T, bool) {
	return &Zt_QuitGame, false
}

// return a valid markup map, creating it if necessary.
func (op *QuitGame) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*QuitGame)(nil)

// quit_game, a type of flow.
var Zt_QuitGame = typeinfo.Flow{
	Name:  "quit_game",
	Lede:  "quit_game",
	Terms: []typeinfo.Term{},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Execute,
	},
}

// holds a slice of type quit_game
// FIX: duplicates the spec decl.
type FIX_QuitGame_Slice []QuitGame

// implements typeinfo.Inspector
func (*QuitGame_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_QuitGame, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_SaveGame struct {
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*SaveGame) Inspect() (typeinfo.T, bool) {
	return &Zt_SaveGame, false
}

// return a valid markup map, creating it if necessary.
func (op *SaveGame) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*SaveGame)(nil)

// save_game, a type of flow.
var Zt_SaveGame = typeinfo.Flow{
	Name:  "save_game",
	Lede:  "save_game",
	Terms: []typeinfo.Term{},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Execute,
	},
}

// holds a slice of type save_game
// FIX: duplicates the spec decl.
type FIX_SaveGame_Slice []SaveGame

// implements typeinfo.Inspector
func (*SaveGame_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_SaveGame, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_RestoreGame struct {
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*RestoreGame) Inspect() (typeinfo.T, bool) {
	return &Zt_RestoreGame, false
}

// return a valid markup map, creating it if necessary.
func (op *RestoreGame) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*RestoreGame)(nil)

// restore_game, a type of flow.
var Zt_RestoreGame = typeinfo.Flow{
	Name:  "restore_game",
	Lede:  "restore_game",
	Terms: []typeinfo.Term{},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Execute,
	},
}

// holds a slice of type restore_game
// FIX: duplicates the spec decl.
type FIX_RestoreGame_Slice []RestoreGame

// implements typeinfo.Inspector
func (*RestoreGame_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_RestoreGame, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_UndoTurn struct {
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*UndoTurn) Inspect() (typeinfo.T, bool) {
	return &Zt_UndoTurn, false
}

// return a valid markup map, creating it if necessary.
func (op *UndoTurn) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*UndoTurn)(nil)

// undo_turn, a type of flow.
var Zt_UndoTurn = typeinfo.Flow{
	Name:  "undo_turn",
	Lede:  "undo_turn",
	Terms: []typeinfo.Term{},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Execute,
	},
}

// holds a slice of type undo_turn
// FIX: duplicates the spec decl.
type FIX_UndoTurn_Slice []UndoTurn

// implements typeinfo.Inspector
func (*UndoTurn_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_UndoTurn, true
}

// FIX: for now we are generating side by side with the old definitions
// also should have user comment here
type FIX_PrintVersion struct {
	Markup map[string]any
}

// implements typeinfo.Inspector
func (*PrintVersion) Inspect() (typeinfo.T, bool) {
	return &Zt_PrintVersion, false
}

// return a valid markup map, creating it if necessary.
func (op *PrintVersion) GetMarkup(ensure bool) map[string]any {
	if ensure && op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	return op.Markup
}

// ensure the command implements its specified slots:
var _ rtti.Execute = (*PrintVersion)(nil)

// print_version, a type of flow.
var Zt_PrintVersion = typeinfo.Flow{
	Name:  "print_version",
	Lede:  "print_version",
	Terms: []typeinfo.Term{},
	Slots: []*typeinfo.Slot{
		&rtti.Zt_Execute,
	},
}

// holds a slice of type print_version
// FIX: duplicates the spec decl.
type FIX_PrintVersion_Slice []PrintVersion

// implements typeinfo.Inspector
func (*PrintVersion_Slice) Inspect() (typeinfo.T, bool) {
	return &Zt_PrintVersion, true
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name:       "game",
	Flow:       z_flow_list,
	Signatures: z_signatures,
}

// a list of all flows in this this package
// ( ex. for reading blockly blocks )
var z_flow_list = []*typeinfo.Flow{
	&Zt_QuitGame,
	&Zt_SaveGame,
	&Zt_RestoreGame,
	&Zt_UndoTurn,
	&Zt_PrintVersion,
}

// a list of all command signatures
// ( for processing and verifying story files )
var z_signatures = map[uint64]any{
	16069653899165369674: (*PrintVersion)(nil), /* execute=PrintVersion */
	13962506025236193050: (*QuitGame)(nil),     /* execute=QuitGame */
	8293164373151279469:  (*RestoreGame)(nil),  /* execute=RestoreGame */
	12343662000108026632: (*SaveGame)(nil),     /* execute=SaveGame */
	6128819475946940678:  (*UndoTurn)(nil),     /* execute=UndoTurn */
}
