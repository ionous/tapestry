// Code generated by "makeops"; edit at your own risk.
package play

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/jsn"
	"github.com/ionous/errutil"
)

// PlayLog a log message that can optionally be displayed to the client
type PlayLog struct {
	Log         string `if:"label=log,type=text"`
	UserComment string
}

// User implemented slots:
var _ PlayMessage = (*PlayLog)(nil)

func (*PlayLog) Compose() composer.Spec {
	return composer.Spec{
		Name: PlayLog_Type,
		Uses: composer.Type_Flow,
		Lede: "play",
	}
}

const PlayLog_Type = "play_log"
const PlayLog_Field_Log = "$LOG"

func (op *PlayLog) Marshal(m jsn.Marshaler) error {
	return PlayLog_Marshal(m, op)
}

type PlayLog_Slice []PlayLog

func (op *PlayLog_Slice) GetType() string { return PlayLog_Type }

func (op *PlayLog_Slice) Marshal(m jsn.Marshaler) error {
	return PlayLog_Repeats_Marshal(m, (*[]PlayLog)(op))
}

func (op *PlayLog_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *PlayLog_Slice) SetSize(cnt int) {
	var els []PlayLog
	if cnt >= 0 {
		els = make(PlayLog_Slice, cnt)
	}
	(*op) = els
}

func (op *PlayLog_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return PlayLog_Marshal(m, &(*op)[i])
}

func PlayLog_Repeats_Marshal(m jsn.Marshaler, vals *[]PlayLog) error {
	return jsn.RepeatBlock(m, (*PlayLog_Slice)(vals))
}

func PlayLog_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]PlayLog) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = PlayLog_Repeats_Marshal(m, pv)
	}
	return
}

type PlayLog_Flow struct{ ptr *PlayLog }

func (n PlayLog_Flow) GetType() string      { return PlayLog_Type }
func (n PlayLog_Flow) GetLede() string      { return "play" }
func (n PlayLog_Flow) GetFlow() interface{} { return n.ptr }
func (n PlayLog_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*PlayLog); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func PlayLog_Optional_Marshal(m jsn.Marshaler, pv **PlayLog) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = PlayLog_Marshal(m, *pv)
	} else if !enc {
		var v PlayLog
		if err = PlayLog_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func PlayLog_Marshal(m jsn.Marshaler, val *PlayLog) (err error) {
	m.SetComment(&val.UserComment)
	if err = m.MarshalBlock(PlayLog_Flow{val}); err == nil {
		e0 := m.MarshalKey("log", PlayLog_Field_Log)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Marshal(m, &val.Log)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", PlayLog_Field_Log))
		}
		m.EndBlock()
	}
	return
}

const PlayMessage_Type = "play_message"

var PlayMessage_Optional_Marshal = PlayMessage_Marshal

type PlayMessage_Slot struct{ Value *PlayMessage }

func (at PlayMessage_Slot) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(at); err == nil {
		if a, ok := at.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}
func (at PlayMessage_Slot) GetType() string              { return PlayMessage_Type }
func (at PlayMessage_Slot) GetSlot() (interface{}, bool) { return *at.Value, *at.Value != nil }
func (at PlayMessage_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.Value), okay = v.(PlayMessage)
	return
}

func PlayMessage_Marshal(m jsn.Marshaler, ptr *PlayMessage) (err error) {
	slot := PlayMessage_Slot{ptr}
	return slot.Marshal(m)
}

type PlayMessage_Slice []PlayMessage

func (op *PlayMessage_Slice) GetType() string { return PlayMessage_Type }

func (op *PlayMessage_Slice) Marshal(m jsn.Marshaler) error {
	return PlayMessage_Repeats_Marshal(m, (*[]PlayMessage)(op))
}

func (op *PlayMessage_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *PlayMessage_Slice) SetSize(cnt int) {
	var els []PlayMessage
	if cnt >= 0 {
		els = make(PlayMessage_Slice, cnt)
	}
	(*op) = els
}

func (op *PlayMessage_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return PlayMessage_Marshal(m, &(*op)[i])
}

func PlayMessage_Repeats_Marshal(m jsn.Marshaler, vals *[]PlayMessage) error {
	return jsn.RepeatBlock(m, (*PlayMessage_Slice)(vals))
}

func PlayMessage_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]PlayMessage) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = PlayMessage_Repeats_Marshal(m, pv)
	}
	return
}

// PlayMode app level change in state
type PlayMode struct {
	Mode        PlayModes `if:"label=mode"`
	UserComment string
}

// User implemented slots:
var _ PlayMessage = (*PlayMode)(nil)

func (*PlayMode) Compose() composer.Spec {
	return composer.Spec{
		Name: PlayMode_Type,
		Uses: composer.Type_Flow,
		Lede: "play",
	}
}

const PlayMode_Type = "play_mode"
const PlayMode_Field_Mode = "$MODE"

func (op *PlayMode) Marshal(m jsn.Marshaler) error {
	return PlayMode_Marshal(m, op)
}

type PlayMode_Slice []PlayMode

func (op *PlayMode_Slice) GetType() string { return PlayMode_Type }

func (op *PlayMode_Slice) Marshal(m jsn.Marshaler) error {
	return PlayMode_Repeats_Marshal(m, (*[]PlayMode)(op))
}

func (op *PlayMode_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *PlayMode_Slice) SetSize(cnt int) {
	var els []PlayMode
	if cnt >= 0 {
		els = make(PlayMode_Slice, cnt)
	}
	(*op) = els
}

func (op *PlayMode_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return PlayMode_Marshal(m, &(*op)[i])
}

func PlayMode_Repeats_Marshal(m jsn.Marshaler, vals *[]PlayMode) error {
	return jsn.RepeatBlock(m, (*PlayMode_Slice)(vals))
}

func PlayMode_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]PlayMode) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = PlayMode_Repeats_Marshal(m, pv)
	}
	return
}

type PlayMode_Flow struct{ ptr *PlayMode }

func (n PlayMode_Flow) GetType() string      { return PlayMode_Type }
func (n PlayMode_Flow) GetLede() string      { return "play" }
func (n PlayMode_Flow) GetFlow() interface{} { return n.ptr }
func (n PlayMode_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*PlayMode); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func PlayMode_Optional_Marshal(m jsn.Marshaler, pv **PlayMode) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = PlayMode_Marshal(m, *pv)
	} else if !enc {
		var v PlayMode
		if err = PlayMode_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func PlayMode_Marshal(m jsn.Marshaler, val *PlayMode) (err error) {
	m.SetComment(&val.UserComment)
	if err = m.MarshalBlock(PlayMode_Flow{val}); err == nil {
		e0 := m.MarshalKey("mode", PlayMode_Field_Mode)
		if e0 == nil {
			e0 = PlayModes_Marshal(m, &val.Mode)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", PlayMode_Field_Mode))
		}
		m.EndBlock()
	}
	return
}

// PlayModes requires a predefined string.
type PlayModes struct {
	Str string
}

func (op *PlayModes) String() string {
	return op.Str
}

const PlayModes_Asm = "$ASM"
const PlayModes_Play = "$PLAY"
const PlayModes_Complete = "$COMPLETE"
const PlayModes_Error = "$ERROR"

func (*PlayModes) Compose() composer.Spec {
	return composer.Spec{
		Name: PlayModes_Type,
		Uses: composer.Type_Str,
		Choices: []string{
			PlayModes_Asm, PlayModes_Play, PlayModes_Complete, PlayModes_Error,
		},
		Strings: []string{
			"asm", "play", "complete", "error",
		},
	}
}

const PlayModes_Type = "play_modes"

func (op *PlayModes) Marshal(m jsn.Marshaler) error {
	return PlayModes_Marshal(m, op)
}

func PlayModes_Optional_Marshal(m jsn.Marshaler, val *PlayModes) (err error) {
	var zero PlayModes
	if enc := m.IsEncoding(); !enc || val.Str != zero.Str {
		err = PlayModes_Marshal(m, val)
	}
	return
}

func PlayModes_Marshal(m jsn.Marshaler, val *PlayModes) (err error) {
	return m.MarshalValue(PlayModes_Type, jsn.MakeEnum(val, &val.Str))
}

type PlayModes_Slice []PlayModes

func (op *PlayModes_Slice) GetType() string { return PlayModes_Type }

func (op *PlayModes_Slice) Marshal(m jsn.Marshaler) error {
	return PlayModes_Repeats_Marshal(m, (*[]PlayModes)(op))
}

func (op *PlayModes_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *PlayModes_Slice) SetSize(cnt int) {
	var els []PlayModes
	if cnt >= 0 {
		els = make(PlayModes_Slice, cnt)
	}
	(*op) = els
}

func (op *PlayModes_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return PlayModes_Marshal(m, &(*op)[i])
}

func PlayModes_Repeats_Marshal(m jsn.Marshaler, vals *[]PlayModes) error {
	return jsn.RepeatBlock(m, (*PlayModes_Slice)(vals))
}

func PlayModes_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]PlayModes) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = PlayModes_Repeats_Marshal(m, pv)
	}
	return
}

// PlayOut output from the game itself
type PlayOut struct {
	Out         string `if:"label=out,type=text"`
	UserComment string
}

// User implemented slots:
var _ PlayMessage = (*PlayOut)(nil)

func (*PlayOut) Compose() composer.Spec {
	return composer.Spec{
		Name: PlayOut_Type,
		Uses: composer.Type_Flow,
		Lede: "play",
	}
}

const PlayOut_Type = "play_out"
const PlayOut_Field_Out = "$OUT"

func (op *PlayOut) Marshal(m jsn.Marshaler) error {
	return PlayOut_Marshal(m, op)
}

type PlayOut_Slice []PlayOut

func (op *PlayOut_Slice) GetType() string { return PlayOut_Type }

func (op *PlayOut_Slice) Marshal(m jsn.Marshaler) error {
	return PlayOut_Repeats_Marshal(m, (*[]PlayOut)(op))
}

func (op *PlayOut_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *PlayOut_Slice) SetSize(cnt int) {
	var els []PlayOut
	if cnt >= 0 {
		els = make(PlayOut_Slice, cnt)
	}
	(*op) = els
}

func (op *PlayOut_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return PlayOut_Marshal(m, &(*op)[i])
}

func PlayOut_Repeats_Marshal(m jsn.Marshaler, vals *[]PlayOut) error {
	return jsn.RepeatBlock(m, (*PlayOut_Slice)(vals))
}

func PlayOut_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]PlayOut) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = PlayOut_Repeats_Marshal(m, pv)
	}
	return
}

type PlayOut_Flow struct{ ptr *PlayOut }

func (n PlayOut_Flow) GetType() string      { return PlayOut_Type }
func (n PlayOut_Flow) GetLede() string      { return "play" }
func (n PlayOut_Flow) GetFlow() interface{} { return n.ptr }
func (n PlayOut_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*PlayOut); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func PlayOut_Optional_Marshal(m jsn.Marshaler, pv **PlayOut) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = PlayOut_Marshal(m, *pv)
	} else if !enc {
		var v PlayOut
		if err = PlayOut_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func PlayOut_Marshal(m jsn.Marshaler, val *PlayOut) (err error) {
	m.SetComment(&val.UserComment)
	if err = m.MarshalBlock(PlayOut_Flow{val}); err == nil {
		e0 := m.MarshalKey("out", PlayOut_Field_Out)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Marshal(m, &val.Out)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", PlayOut_Field_Out))
		}
		m.EndBlock()
	}
	return
}

var Slots = []interface{}{
	(*PlayMessage)(nil),
}

var Slats = []composer.Composer{
	(*PlayLog)(nil),
	(*PlayMode)(nil),
	(*PlayModes)(nil),
	(*PlayOut)(nil),
}

var Signatures = map[uint64]interface{}{
	7172632378119087629:  (*PlayLog)(nil),  /* Play log: */
	18089831122107554052: (*PlayMode)(nil), /* Play mode: */
	10234677752751325483: (*PlayOut)(nil),  /* Play out: */
}
