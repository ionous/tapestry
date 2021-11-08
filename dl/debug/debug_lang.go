// Code generated by "makeops"; edit at your own risk.
package debug

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

// DebugLog Debug log
// User implements: Execute.
type DebugLog struct {
	Value    rt.Assignment `if:"label=_"`
	LogLevel LoggingLevel  `if:"label=as,optional"`
}

func (*DebugLog) Compose() composer.Spec {
	return composer.Spec{
		Name: DebugLog_Type,
		Uses: composer.Type_Flow,
		Lede: "log",
	}
}

const DebugLog_Type = "debug_log"

const DebugLog_Field_Value = "$VALUE"
const DebugLog_Field_LogLevel = "$LOG_LEVEL"

func (op *DebugLog) Marshal(m jsn.Marshaler) error {
	return DebugLog_Marshal(m, op)
}

type DebugLog_Slice []DebugLog

func (op *DebugLog_Slice) GetType() string {
	return DebugLog_Type
}

func (op *DebugLog_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *DebugLog_Slice) SetSize(cnt int) {
	var els []DebugLog
	if cnt >= 0 {
		els = make(DebugLog_Slice, cnt)
	}
	(*op) = els
}

func (op *DebugLog_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return DebugLog_Marshal(m, &(*op)[i])
}

func DebugLog_Repeats_Marshal(m jsn.Marshaler, vals *[]DebugLog) error {
	return jsn.RepeatBlock(m, (*DebugLog_Slice)(vals))
}

func DebugLog_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]DebugLog) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = DebugLog_Repeats_Marshal(m, pv)
	}
	return
}

func DebugLog_Optional_Marshal(m jsn.Marshaler, pv **DebugLog) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = DebugLog_Marshal(m, *pv)
	} else if !enc {
		var v DebugLog
		if err = DebugLog_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func DebugLog_Marshal(m jsn.Marshaler, val *DebugLog) (err error) {
	if err = m.MarshalBlock(jsn.MakeFlow("log", DebugLog_Type, val)); err == nil {
		e0 := m.MarshalKey("", DebugLog_Field_Value)
		if e0 == nil {
			e0 = rt.Assignment_Marshal(m, &val.Value)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", DebugLog_Field_Value))
		}
		e1 := m.MarshalKey("as", DebugLog_Field_LogLevel)
		if e1 == nil {
			e1 = LoggingLevel_Optional_Marshal(m, &val.LogLevel)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", DebugLog_Field_LogLevel))
		}
		m.EndBlock()
	}
	return
}

// DoNothing Statement which does nothing.
// User implements: Execute.
type DoNothing struct {
	Reason string `if:"label=why,optional,type=text"`
}

func (*DoNothing) Compose() composer.Spec {
	return composer.Spec{
		Name: DoNothing_Type,
		Uses: composer.Type_Flow,
	}
}

const DoNothing_Type = "do_nothing"

const DoNothing_Field_Reason = "$REASON"

func (op *DoNothing) Marshal(m jsn.Marshaler) error {
	return DoNothing_Marshal(m, op)
}

type DoNothing_Slice []DoNothing

func (op *DoNothing_Slice) GetType() string {
	return DoNothing_Type
}

func (op *DoNothing_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *DoNothing_Slice) SetSize(cnt int) {
	var els []DoNothing
	if cnt >= 0 {
		els = make(DoNothing_Slice, cnt)
	}
	(*op) = els
}

func (op *DoNothing_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return DoNothing_Marshal(m, &(*op)[i])
}

func DoNothing_Repeats_Marshal(m jsn.Marshaler, vals *[]DoNothing) error {
	return jsn.RepeatBlock(m, (*DoNothing_Slice)(vals))
}

func DoNothing_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]DoNothing) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = DoNothing_Repeats_Marshal(m, pv)
	}
	return
}

func DoNothing_Optional_Marshal(m jsn.Marshaler, pv **DoNothing) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = DoNothing_Marshal(m, *pv)
	} else if !enc {
		var v DoNothing
		if err = DoNothing_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func DoNothing_Marshal(m jsn.Marshaler, val *DoNothing) (err error) {
	if err = m.MarshalBlock(jsn.MakeFlow(DoNothing_Type, DoNothing_Type, val)); err == nil {
		e0 := m.MarshalKey("why", DoNothing_Field_Reason)
		if e0 == nil {
			e0 = value.Text_Unboxed_Optional_Marshal(m, &val.Reason)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", DoNothing_Field_Reason))
		}
		m.EndBlock()
	}
	return
}

// LoggingLevel requires a predefined string.
type LoggingLevel struct {
	Str string
}

func (op *LoggingLevel) String() string {
	return op.Str
}

const LoggingLevel_Note = "$NOTE"
const LoggingLevel_ToDo = "$TO_DO"
const LoggingLevel_Fix = "$FIX"
const LoggingLevel_Info = "$INFO"
const LoggingLevel_Warning = "$WARNING"
const LoggingLevel_Error = "$ERROR"

func (*LoggingLevel) Compose() composer.Spec {
	return composer.Spec{
		Name: LoggingLevel_Type,
		Uses: composer.Type_Str,
		Choices: []string{
			LoggingLevel_Note, LoggingLevel_ToDo, LoggingLevel_Fix, LoggingLevel_Info, LoggingLevel_Warning, LoggingLevel_Error,
		},
		Strings: []string{
			"note", "to_do", "fix", "info", "warning", "error",
		},
	}
}

const LoggingLevel_Type = "logging_level"

func (op *LoggingLevel) Marshal(m jsn.Marshaler) error {
	return LoggingLevel_Marshal(m, op)
}

func LoggingLevel_Optional_Marshal(m jsn.Marshaler, val *LoggingLevel) (err error) {
	var zero LoggingLevel
	if enc := m.IsEncoding(); !enc || val.Str != zero.Str {
		err = LoggingLevel_Marshal(m, val)
	}
	return
}

func LoggingLevel_Marshal(m jsn.Marshaler, val *LoggingLevel) (err error) {
	return m.MarshalValue(LoggingLevel_Type, jsn.MakeEnum(val, &val.Str))
}

type LoggingLevel_Slice []LoggingLevel

func (op *LoggingLevel_Slice) GetType() string {
	return LoggingLevel_Type
}

func (op *LoggingLevel_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *LoggingLevel_Slice) SetSize(cnt int) {
	var els []LoggingLevel
	if cnt >= 0 {
		els = make(LoggingLevel_Slice, cnt)
	}
	(*op) = els
}

func (op *LoggingLevel_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return LoggingLevel_Marshal(m, &(*op)[i])
}

func LoggingLevel_Repeats_Marshal(m jsn.Marshaler, vals *[]LoggingLevel) error {
	return jsn.RepeatBlock(m, (*LoggingLevel_Slice)(vals))
}

func LoggingLevel_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]LoggingLevel) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = LoggingLevel_Repeats_Marshal(m, pv)
	}
	return
}

var Slats = []composer.Composer{
	(*DebugLog)(nil),
	(*DoNothing)(nil),
	(*LoggingLevel)(nil),
}

var Signatures = map[uint64]interface{}{
	5700043876155103121:  (*DebugLog)(nil),  /* Log: */
	17593113710683116377: (*DebugLog)(nil),  /* Log:as: */
	5234640093503358177:  (*DoNothing)(nil), /* DoNothing */
	15838679201884235887: (*DoNothing)(nil), /* DoNothing why: */
}
