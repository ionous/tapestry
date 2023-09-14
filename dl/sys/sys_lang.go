// Code generated by "makeops"; edit at your own risk.
package sys

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/rt"
)

// PrintVersion
type PrintVersion struct {
	Markup map[string]any
}

// User implemented slots:
var _ rt.Execute = (*PrintVersion)(nil)

func (*PrintVersion) Compose() composer.Spec {
	return composer.Spec{
		Name: PrintVersion_Type,
		Uses: composer.Type_Flow,
	}
}

const PrintVersion_Type = "print_version"

func (op *PrintVersion) Marshal(m jsn.Marshaler) error {
	return PrintVersion_Marshal(m, op)
}

type PrintVersion_Slice []PrintVersion

func (op *PrintVersion_Slice) GetType() string { return PrintVersion_Type }

func (op *PrintVersion_Slice) Marshal(m jsn.Marshaler) error {
	return PrintVersion_Repeats_Marshal(m, (*[]PrintVersion)(op))
}

func (op *PrintVersion_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *PrintVersion_Slice) SetSize(cnt int) {
	var els []PrintVersion
	if cnt >= 0 {
		els = make(PrintVersion_Slice, cnt)
	}
	(*op) = els
}

func (op *PrintVersion_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return PrintVersion_Marshal(m, &(*op)[i])
}

func PrintVersion_Repeats_Marshal(m jsn.Marshaler, vals *[]PrintVersion) error {
	return jsn.RepeatBlock(m, (*PrintVersion_Slice)(vals))
}

func PrintVersion_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]PrintVersion) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = PrintVersion_Repeats_Marshal(m, pv)
	}
	return
}

type PrintVersion_Flow struct{ ptr *PrintVersion }

func (n PrintVersion_Flow) GetType() string      { return PrintVersion_Type }
func (n PrintVersion_Flow) GetLede() string      { return PrintVersion_Type }
func (n PrintVersion_Flow) GetFlow() interface{} { return n.ptr }
func (n PrintVersion_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*PrintVersion); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func PrintVersion_Optional_Marshal(m jsn.Marshaler, pv **PrintVersion) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = PrintVersion_Marshal(m, *pv)
	} else if !enc {
		var v PrintVersion
		if err = PrintVersion_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func PrintVersion_Marshal(m jsn.Marshaler, val *PrintVersion) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(PrintVersion_Flow{val}); err == nil {
		m.EndBlock()
	}
	return
}

// QuitGame
type QuitGame struct {
	Markup map[string]any
}

// User implemented slots:
var _ rt.Execute = (*QuitGame)(nil)

func (*QuitGame) Compose() composer.Spec {
	return composer.Spec{
		Name: QuitGame_Type,
		Uses: composer.Type_Flow,
	}
}

const QuitGame_Type = "quit_game"

func (op *QuitGame) Marshal(m jsn.Marshaler) error {
	return QuitGame_Marshal(m, op)
}

type QuitGame_Slice []QuitGame

func (op *QuitGame_Slice) GetType() string { return QuitGame_Type }

func (op *QuitGame_Slice) Marshal(m jsn.Marshaler) error {
	return QuitGame_Repeats_Marshal(m, (*[]QuitGame)(op))
}

func (op *QuitGame_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *QuitGame_Slice) SetSize(cnt int) {
	var els []QuitGame
	if cnt >= 0 {
		els = make(QuitGame_Slice, cnt)
	}
	(*op) = els
}

func (op *QuitGame_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return QuitGame_Marshal(m, &(*op)[i])
}

func QuitGame_Repeats_Marshal(m jsn.Marshaler, vals *[]QuitGame) error {
	return jsn.RepeatBlock(m, (*QuitGame_Slice)(vals))
}

func QuitGame_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]QuitGame) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = QuitGame_Repeats_Marshal(m, pv)
	}
	return
}

type QuitGame_Flow struct{ ptr *QuitGame }

func (n QuitGame_Flow) GetType() string      { return QuitGame_Type }
func (n QuitGame_Flow) GetLede() string      { return QuitGame_Type }
func (n QuitGame_Flow) GetFlow() interface{} { return n.ptr }
func (n QuitGame_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*QuitGame); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func QuitGame_Optional_Marshal(m jsn.Marshaler, pv **QuitGame) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = QuitGame_Marshal(m, *pv)
	} else if !enc {
		var v QuitGame
		if err = QuitGame_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func QuitGame_Marshal(m jsn.Marshaler, val *QuitGame) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(QuitGame_Flow{val}); err == nil {
		m.EndBlock()
	}
	return
}

// RestoreGame
type RestoreGame struct {
	Markup map[string]any
}

// User implemented slots:
var _ rt.Execute = (*RestoreGame)(nil)

func (*RestoreGame) Compose() composer.Spec {
	return composer.Spec{
		Name: RestoreGame_Type,
		Uses: composer.Type_Flow,
	}
}

const RestoreGame_Type = "restore_game"

func (op *RestoreGame) Marshal(m jsn.Marshaler) error {
	return RestoreGame_Marshal(m, op)
}

type RestoreGame_Slice []RestoreGame

func (op *RestoreGame_Slice) GetType() string { return RestoreGame_Type }

func (op *RestoreGame_Slice) Marshal(m jsn.Marshaler) error {
	return RestoreGame_Repeats_Marshal(m, (*[]RestoreGame)(op))
}

func (op *RestoreGame_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RestoreGame_Slice) SetSize(cnt int) {
	var els []RestoreGame
	if cnt >= 0 {
		els = make(RestoreGame_Slice, cnt)
	}
	(*op) = els
}

func (op *RestoreGame_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RestoreGame_Marshal(m, &(*op)[i])
}

func RestoreGame_Repeats_Marshal(m jsn.Marshaler, vals *[]RestoreGame) error {
	return jsn.RepeatBlock(m, (*RestoreGame_Slice)(vals))
}

func RestoreGame_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RestoreGame) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = RestoreGame_Repeats_Marshal(m, pv)
	}
	return
}

type RestoreGame_Flow struct{ ptr *RestoreGame }

func (n RestoreGame_Flow) GetType() string      { return RestoreGame_Type }
func (n RestoreGame_Flow) GetLede() string      { return RestoreGame_Type }
func (n RestoreGame_Flow) GetFlow() interface{} { return n.ptr }
func (n RestoreGame_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*RestoreGame); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func RestoreGame_Optional_Marshal(m jsn.Marshaler, pv **RestoreGame) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = RestoreGame_Marshal(m, *pv)
	} else if !enc {
		var v RestoreGame
		if err = RestoreGame_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func RestoreGame_Marshal(m jsn.Marshaler, val *RestoreGame) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(RestoreGame_Flow{val}); err == nil {
		m.EndBlock()
	}
	return
}

// SaveGame
type SaveGame struct {
	Markup map[string]any
}

// User implemented slots:
var _ rt.Execute = (*SaveGame)(nil)

func (*SaveGame) Compose() composer.Spec {
	return composer.Spec{
		Name: SaveGame_Type,
		Uses: composer.Type_Flow,
	}
}

const SaveGame_Type = "save_game"

func (op *SaveGame) Marshal(m jsn.Marshaler) error {
	return SaveGame_Marshal(m, op)
}

type SaveGame_Slice []SaveGame

func (op *SaveGame_Slice) GetType() string { return SaveGame_Type }

func (op *SaveGame_Slice) Marshal(m jsn.Marshaler) error {
	return SaveGame_Repeats_Marshal(m, (*[]SaveGame)(op))
}

func (op *SaveGame_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *SaveGame_Slice) SetSize(cnt int) {
	var els []SaveGame
	if cnt >= 0 {
		els = make(SaveGame_Slice, cnt)
	}
	(*op) = els
}

func (op *SaveGame_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return SaveGame_Marshal(m, &(*op)[i])
}

func SaveGame_Repeats_Marshal(m jsn.Marshaler, vals *[]SaveGame) error {
	return jsn.RepeatBlock(m, (*SaveGame_Slice)(vals))
}

func SaveGame_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]SaveGame) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = SaveGame_Repeats_Marshal(m, pv)
	}
	return
}

type SaveGame_Flow struct{ ptr *SaveGame }

func (n SaveGame_Flow) GetType() string      { return SaveGame_Type }
func (n SaveGame_Flow) GetLede() string      { return SaveGame_Type }
func (n SaveGame_Flow) GetFlow() interface{} { return n.ptr }
func (n SaveGame_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*SaveGame); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func SaveGame_Optional_Marshal(m jsn.Marshaler, pv **SaveGame) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = SaveGame_Marshal(m, *pv)
	} else if !enc {
		var v SaveGame
		if err = SaveGame_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func SaveGame_Marshal(m jsn.Marshaler, val *SaveGame) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(SaveGame_Flow{val}); err == nil {
		m.EndBlock()
	}
	return
}

// UndoTurn
type UndoTurn struct {
	Markup map[string]any
}

// User implemented slots:
var _ rt.Execute = (*UndoTurn)(nil)

func (*UndoTurn) Compose() composer.Spec {
	return composer.Spec{
		Name: UndoTurn_Type,
		Uses: composer.Type_Flow,
	}
}

const UndoTurn_Type = "undo_turn"

func (op *UndoTurn) Marshal(m jsn.Marshaler) error {
	return UndoTurn_Marshal(m, op)
}

type UndoTurn_Slice []UndoTurn

func (op *UndoTurn_Slice) GetType() string { return UndoTurn_Type }

func (op *UndoTurn_Slice) Marshal(m jsn.Marshaler) error {
	return UndoTurn_Repeats_Marshal(m, (*[]UndoTurn)(op))
}

func (op *UndoTurn_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *UndoTurn_Slice) SetSize(cnt int) {
	var els []UndoTurn
	if cnt >= 0 {
		els = make(UndoTurn_Slice, cnt)
	}
	(*op) = els
}

func (op *UndoTurn_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return UndoTurn_Marshal(m, &(*op)[i])
}

func UndoTurn_Repeats_Marshal(m jsn.Marshaler, vals *[]UndoTurn) error {
	return jsn.RepeatBlock(m, (*UndoTurn_Slice)(vals))
}

func UndoTurn_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]UndoTurn) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = UndoTurn_Repeats_Marshal(m, pv)
	}
	return
}

type UndoTurn_Flow struct{ ptr *UndoTurn }

func (n UndoTurn_Flow) GetType() string      { return UndoTurn_Type }
func (n UndoTurn_Flow) GetLede() string      { return UndoTurn_Type }
func (n UndoTurn_Flow) GetFlow() interface{} { return n.ptr }
func (n UndoTurn_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*UndoTurn); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func UndoTurn_Optional_Marshal(m jsn.Marshaler, pv **UndoTurn) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = UndoTurn_Marshal(m, *pv)
	} else if !enc {
		var v UndoTurn
		if err = UndoTurn_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func UndoTurn_Marshal(m jsn.Marshaler, val *UndoTurn) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(UndoTurn_Flow{val}); err == nil {
		m.EndBlock()
	}
	return
}

var Slats = []composer.Composer{
	(*PrintVersion)(nil),
	(*QuitGame)(nil),
	(*RestoreGame)(nil),
	(*SaveGame)(nil),
	(*UndoTurn)(nil),
}

var Signatures = map[uint64]interface{}{
	16069653899165369674: (*PrintVersion)(nil), /* execute=PrintVersion */
	13962506025236193050: (*QuitGame)(nil),     /* execute=QuitGame */
	8293164373151279469:  (*RestoreGame)(nil),  /* execute=RestoreGame */
	12343662000108026632: (*SaveGame)(nil),     /* execute=SaveGame */
	6128819475946940678:  (*UndoTurn)(nil),     /* execute=UndoTurn */
}
