// Code generated by "makeops"; edit at your own risk.
package debug

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/rt"
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

func (op *DebugLog) Marshal(n jsn.Marshaler) {
	DebugLog_Marshal(n, op)
}

type DebugLog_Slice []DebugLog

func (op *DebugLog_Slice) GetSize() int    { return len(*op) }
func (op *DebugLog_Slice) SetSize(cnt int) { (*op) = make(DebugLog_Slice, cnt) }

func DebugLog_Repeats_Marshal(n jsn.Marshaler, vals *[]DebugLog) {
	if n.RepeatValues(DebugLog_Type, (*DebugLog_Slice)(vals)) {
		for i := range *vals {
			DebugLog_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

func DebugLog_Optional_Marshal(n jsn.Marshaler, pv **DebugLog) {
	if enc := n.IsEncoding(); enc && *pv != nil {
		DebugLog_Marshal(n, *pv)
	} else if !enc {
		var v DebugLog
		if DebugLog_Marshal(n, &v) {
			*pv = &v
		}
	}
}

func DebugLog_Marshal(n jsn.Marshaler, val *DebugLog) (okay bool) {
	if okay = n.MapValues("log", DebugLog_Type); okay {
		if n.MapKey("", DebugLog_Field_Value) {
			rt.Assignment_Marshal(n, &val.Value)
		}
		if n.MapKey("as", DebugLog_Field_LogLevel) {
			LoggingLevel_Optional_Marshal(n, &val.LogLevel)
		}
		n.EndValues()
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

func (op *DoNothing) Marshal(n jsn.Marshaler) {
	DoNothing_Marshal(n, op)
}

type DoNothing_Slice []DoNothing

func (op *DoNothing_Slice) GetSize() int    { return len(*op) }
func (op *DoNothing_Slice) SetSize(cnt int) { (*op) = make(DoNothing_Slice, cnt) }

func DoNothing_Repeats_Marshal(n jsn.Marshaler, vals *[]DoNothing) {
	if n.RepeatValues(DoNothing_Type, (*DoNothing_Slice)(vals)) {
		for i := range *vals {
			DoNothing_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

func DoNothing_Optional_Marshal(n jsn.Marshaler, pv **DoNothing) {
	if enc := n.IsEncoding(); enc && *pv != nil {
		DoNothing_Marshal(n, *pv)
	} else if !enc {
		var v DoNothing
		if DoNothing_Marshal(n, &v) {
			*pv = &v
		}
	}
}

func DoNothing_Marshal(n jsn.Marshaler, val *DoNothing) (okay bool) {
	if okay = n.MapValues(DoNothing_Type, DoNothing_Type); okay {
		if n.MapKey("why", DoNothing_Field_Reason) {
			value.Text_Unboxed_Optional_Marshal(n, &val.Reason)
		}
		n.EndValues()
	}
	return
}

// LoggingLevel requires a user-specified string.
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

func (op *LoggingLevel) Marshal(n jsn.Marshaler) {
	LoggingLevel_Marshal(n, op)
}

func LoggingLevel_Optional_Marshal(n jsn.Marshaler, val *LoggingLevel) {
	var zero LoggingLevel
	if enc := n.IsEncoding(); !enc || val.Str != zero.Str {
		LoggingLevel_Marshal(n, val)
	}
}

func LoggingLevel_Marshal(n jsn.Marshaler, val *LoggingLevel) {
	n.MarshalValue(LoggingLevel_Type, jsn.MakeEnum(val, &val.Str))
}

type LoggingLevel_Slice []LoggingLevel

func (op *LoggingLevel_Slice) GetSize() int    { return len(*op) }
func (op *LoggingLevel_Slice) SetSize(cnt int) { (*op) = make(LoggingLevel_Slice, cnt) }

func LoggingLevel_Repeats_Marshal(n jsn.Marshaler, vals *[]LoggingLevel) {
	if n.RepeatValues(LoggingLevel_Type, (*LoggingLevel_Slice)(vals)) {
		for i := range *vals {
			LoggingLevel_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
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
