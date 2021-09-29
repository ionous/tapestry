// Code generated by "makeops"; edit at your own risk.
package debug

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/export/jsn"
	"git.sr.ht/~ionous/iffy/rt"
)

// DebugLog Debug log
type DebugLog struct {
	Value    rt.Assignment `if:"label=_"`
	LogLevel LoggingLevel  `if:"label=as,optional"`
}

func (*DebugLog) GetType() string {
	return DebugLog_Type
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

func DebugLog_Repeats_Marshal(n jsn.Marshaler, vals *[]DebugLog) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		if n.RepeatValues(cnt) {
			for _, el := range *vals {
				DebugLog_Marshal(n, &el)
			}
			n.EndValues()
		}
	}
	return
}

func DebugLog_Optional_Marshal(n jsn.Marshaler, val **DebugLog) {
	if *val != nil {
		DebugLog_Marshal(n, *val)
	}
}

func DebugLog_Marshal(n jsn.Marshaler, val *DebugLog) {
	if n.MapValues("log", DebugLog_Type) {
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
type DoNothing struct {
	Reason string `if:"label=why,optional,type=text"`
}

func (*DoNothing) GetType() string {
	return DoNothing_Type
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

func DoNothing_Repeats_Marshal(n jsn.Marshaler, vals *[]DoNothing) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		if n.RepeatValues(cnt) {
			for _, el := range *vals {
				DoNothing_Marshal(n, &el)
			}
			n.EndValues()
		}
	}
	return
}

func DoNothing_Optional_Marshal(n jsn.Marshaler, val **DoNothing) {
	if *val != nil {
		DoNothing_Marshal(n, *val)
	}
}

func DoNothing_Marshal(n jsn.Marshaler, val *DoNothing) {
	if n.MapValues(DoNothing_Type, DoNothing_Type) {
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

func (*LoggingLevel) GetType() string {
	return LoggingLevel_Type
}

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
	if val.Str != zero.Str {
		LoggingLevel_Marshal(n, val)
	}
}

func LoggingLevel_Marshal(n jsn.Marshaler, val *LoggingLevel) {
	n.SpecifyEnum(jsn.MakeEnum(val, &val.Str))
}

func LoggingLevel_Repeats_Marshal(n jsn.Marshaler, vals *[]LoggingLevel) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		if n.RepeatValues(cnt) {
			for _, el := range *vals {
				LoggingLevel_Marshal(n, &el)
			}
			n.EndValues()
		}
	}
	return
}

var Slats = []composer.Composer{
	(*DebugLog)(nil),
	(*DoNothing)(nil),
	(*LoggingLevel)(nil),
}
