// Code generated by "makeops"; edit at your own risk.
package debug

import (
	"encoding/json"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/export/jsonexp"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

// DebugLog Debug log
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
const DebugLog_Lede = "log"
const DebugLog_Field_Value = "$VALUE"
const DebugLog_Field_LogLevel = "$LOG_LEVEL"

func (op *DebugLog) MarshalCompact(n jsonexp.Context) (ret []byte, err error) {
	return DebugLog_Compact_Marshal(n, op)
}
func (op *DebugLog) UnmarshalCompact(n jsonexp.Context, b []byte) error {
	return DebugLog_Compact_Unmarshal(n, b, op)
}
func (op *DebugLog) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return DebugLog_Detailed_Marshal(n, op)
}
func (op *DebugLog) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return DebugLog_Detailed_Unmarshal(n, b, op)
}

func DebugLog_Compact_Repeats_Marshal(n jsonexp.Context, vals *[]DebugLog) ([]byte, error) {
	return DebugLog_Repeats_Marshal(n, vals, DebugLog_Compact_Marshal)
}
func DebugLog_Detailed_Repeats_Marshal(n jsonexp.Context, vals *[]DebugLog) ([]byte, error) {
	return DebugLog_Repeats_Marshal(n, vals, DebugLog_Detailed_Marshal)
}
func DebugLog_Repeats_Marshal(n jsonexp.Context, vals *[]DebugLog, marshEl func(jsonexp.Context, *DebugLog) ([]byte, error)) (ret []byte, err error) {
	var msgs []json.RawMessage
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		msgs = make([]json.RawMessage, cnt)
		for i, el := range *vals {
			if b, e := marshEl(n, &el); e != nil {
				err = errutil.New(DebugLog_Type, "at", i, "-", e)
				break
			} else {
				msgs[i] = b
			}
		}
	}
	if err == nil {
		ret, err = json.Marshal(msgs)
	}
	return
}

func DebugLog_Compact_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]DebugLog) error {
	return DebugLog_Repeats_Unmarshal(n, b, out, DebugLog_Compact_Unmarshal)
}
func DebugLog_Detailed_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]DebugLog) error {
	return DebugLog_Repeats_Unmarshal(n, b, out, DebugLog_Detailed_Unmarshal)
}
func DebugLog_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]DebugLog, unmarshEl func(jsonexp.Context, []byte, *DebugLog) error) (err error) {
	var vals []DebugLog
	if len(b) > 0 { // generated code collapses optional and empty.
		var msgs []json.RawMessage
		if e := json.Unmarshal(b, &msgs); e != nil {
			err = errutil.New(DebugLog_Type, "-", e)
		} else {
			vals = make([]DebugLog, len(msgs))
			for i, msg := range msgs {
				if e := unmarshEl(n, msg, &vals[i]); e != nil {
					err = errutil.New(DebugLog_Type, "at", i, "-", e)
					break
				}
			}
		}
	}
	if err == nil {
		*out = vals
	}
	return
}

func DebugLog_Compact_Optional_Marshal(n jsonexp.Context, val **DebugLog) (ret []byte, err error) {
	if *val != nil {
		ret, err = DebugLog_Compact_Marshal(n, *val)
	}
	return
}
func DebugLog_Compact_Optional_Unmarshal(n jsonexp.Context, b []byte, out **DebugLog) (err error) {
	if len(b) > 0 {
		var val DebugLog
		if e := DebugLog_Compact_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}

func DebugLog_Compact_Marshal(n jsonexp.Context, val *DebugLog) (ret []byte, err error) {
	var sig jsonexp.CompactFlow
	sig.WriteLede(DebugLog_Lede)
	if b, e := rt.Assignment_Compact_Marshal(n, &val.Value); e != nil {
		err = errutil.Append(err, e)
	} else {
		sig.AddMsg("", b)
	}
	if b, e := LoggingLevel_Compact_Optional_Marshal(n, &val.LogLevel); e != nil {
		err = errutil.Append(err, e)
	} else if len(b) > 0 {
		sig.AddMsg("as", b)
	}
	if err == nil {
		ret, err = sig.MarshalJSON()
	}
	return
}
func DebugLog_Compact_Unmarshal(n jsonexp.Context, b []byte, out *DebugLog) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(DebugLog_Type, "-", e)
	} else if e := rt.Assignment_Compact_Unmarshal(n, msg.Fields[DebugLog_Field_Value], &out.Value); e != nil {
		err = errutil.New(DebugLog_Type+"."+DebugLog_Field_Value, "-", e)
	} else if e := LoggingLevel_Compact_Optional_Unmarshal(n, msg.Fields[DebugLog_Field_LogLevel], &out.LogLevel); e != nil {
		err = errutil.New(DebugLog_Type+"."+DebugLog_Field_LogLevel, "-", e)
	}
	return
}

func DebugLog_Detailed_Optional_Marshal(n jsonexp.Context, val **DebugLog) (ret []byte, err error) {
	if *val != nil {
		ret, err = DebugLog_Detailed_Marshal(n, *val)
	}
	return
}
func DebugLog_Detailed_Marshal(n jsonexp.Context, val *DebugLog) (ret []byte, err error) {
	fields := make(jsonexp.Fields)
	if b, e := rt.Assignment_Detailed_Marshal(n, &val.Value); e != nil {
		err = errutil.Append(err, e)
	} else {
		fields[DebugLog_Field_Value] = b
	}

	if b, e := LoggingLevel_Detailed_Optional_Marshal(n, &val.LogLevel); e != nil {
		err = errutil.Append(err, e)
	} else if len(b) > 0 {
		fields[DebugLog_Field_LogLevel] = b
	}

	if err == nil {
		ret, err = json.Marshal(jsonexp.Flow{
			Type:   DebugLog_Type,
			Fields: fields,
		})
	}
	return
}

func DebugLog_Detailed_Optional_Unmarshal(n jsonexp.Context, b []byte, out **DebugLog) (err error) {
	if len(b) > 0 {
		var val DebugLog
		if e := DebugLog_Detailed_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}
func DebugLog_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *DebugLog) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(DebugLog_Type, "-", e)
	} else if e := rt.Assignment_Detailed_Unmarshal(n, msg.Fields[DebugLog_Field_Value], &out.Value); e != nil {
		err = errutil.New(DebugLog_Type+"."+DebugLog_Field_Value, "-", e)
	} else if e := LoggingLevel_Detailed_Optional_Unmarshal(n, msg.Fields[DebugLog_Field_LogLevel], &out.LogLevel); e != nil {
		err = errutil.New(DebugLog_Type+"."+DebugLog_Field_LogLevel, "-", e)
	}
	return
}

// DoNothing Statement which does nothing.
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
const DoNothing_Lede = DoNothing_Type
const DoNothing_Field_Reason = "$REASON"

func (op *DoNothing) MarshalCompact(n jsonexp.Context) (ret []byte, err error) {
	return DoNothing_Compact_Marshal(n, op)
}
func (op *DoNothing) UnmarshalCompact(n jsonexp.Context, b []byte) error {
	return DoNothing_Compact_Unmarshal(n, b, op)
}
func (op *DoNothing) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return DoNothing_Detailed_Marshal(n, op)
}
func (op *DoNothing) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return DoNothing_Detailed_Unmarshal(n, b, op)
}

func DoNothing_Compact_Repeats_Marshal(n jsonexp.Context, vals *[]DoNothing) ([]byte, error) {
	return DoNothing_Repeats_Marshal(n, vals, DoNothing_Compact_Marshal)
}
func DoNothing_Detailed_Repeats_Marshal(n jsonexp.Context, vals *[]DoNothing) ([]byte, error) {
	return DoNothing_Repeats_Marshal(n, vals, DoNothing_Detailed_Marshal)
}
func DoNothing_Repeats_Marshal(n jsonexp.Context, vals *[]DoNothing, marshEl func(jsonexp.Context, *DoNothing) ([]byte, error)) (ret []byte, err error) {
	var msgs []json.RawMessage
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		msgs = make([]json.RawMessage, cnt)
		for i, el := range *vals {
			if b, e := marshEl(n, &el); e != nil {
				err = errutil.New(DoNothing_Type, "at", i, "-", e)
				break
			} else {
				msgs[i] = b
			}
		}
	}
	if err == nil {
		ret, err = json.Marshal(msgs)
	}
	return
}

func DoNothing_Compact_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]DoNothing) error {
	return DoNothing_Repeats_Unmarshal(n, b, out, DoNothing_Compact_Unmarshal)
}
func DoNothing_Detailed_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]DoNothing) error {
	return DoNothing_Repeats_Unmarshal(n, b, out, DoNothing_Detailed_Unmarshal)
}
func DoNothing_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]DoNothing, unmarshEl func(jsonexp.Context, []byte, *DoNothing) error) (err error) {
	var vals []DoNothing
	if len(b) > 0 { // generated code collapses optional and empty.
		var msgs []json.RawMessage
		if e := json.Unmarshal(b, &msgs); e != nil {
			err = errutil.New(DoNothing_Type, "-", e)
		} else {
			vals = make([]DoNothing, len(msgs))
			for i, msg := range msgs {
				if e := unmarshEl(n, msg, &vals[i]); e != nil {
					err = errutil.New(DoNothing_Type, "at", i, "-", e)
					break
				}
			}
		}
	}
	if err == nil {
		*out = vals
	}
	return
}

func DoNothing_Compact_Optional_Marshal(n jsonexp.Context, val **DoNothing) (ret []byte, err error) {
	if *val != nil {
		ret, err = DoNothing_Compact_Marshal(n, *val)
	}
	return
}
func DoNothing_Compact_Optional_Unmarshal(n jsonexp.Context, b []byte, out **DoNothing) (err error) {
	if len(b) > 0 {
		var val DoNothing
		if e := DoNothing_Compact_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}

func DoNothing_Compact_Marshal(n jsonexp.Context, val *DoNothing) (ret []byte, err error) {
	var sig jsonexp.CompactFlow
	sig.WriteLede(DoNothing_Lede)
	if b, e := value.Text_Override_Compact_Optional_Marshal(n, &val.Reason); e != nil {
		err = errutil.Append(err, e)
	} else if len(b) > 0 {
		sig.AddMsg("why", b)
	}
	if err == nil {
		ret, err = sig.MarshalJSON()
	}
	return
}
func DoNothing_Compact_Unmarshal(n jsonexp.Context, b []byte, out *DoNothing) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(DoNothing_Type, "-", e)
	} else if e := value.Text_Override_Compact_Optional_Unmarshal(n, msg.Fields[DoNothing_Field_Reason], &out.Reason); e != nil {
		err = errutil.New(DoNothing_Type+"."+DoNothing_Field_Reason, "-", e)
	}
	return
}

func DoNothing_Detailed_Optional_Marshal(n jsonexp.Context, val **DoNothing) (ret []byte, err error) {
	if *val != nil {
		ret, err = DoNothing_Detailed_Marshal(n, *val)
	}
	return
}
func DoNothing_Detailed_Marshal(n jsonexp.Context, val *DoNothing) (ret []byte, err error) {
	fields := make(jsonexp.Fields)
	if b, e := value.Text_Override_Detailed_Optional_Marshal(n, &val.Reason); e != nil {
		err = errutil.Append(err, e)
	} else if len(b) > 0 {
		fields[DoNothing_Field_Reason] = b
	}

	if err == nil {
		ret, err = json.Marshal(jsonexp.Flow{
			Type:   DoNothing_Type,
			Fields: fields,
		})
	}
	return
}

func DoNothing_Detailed_Optional_Unmarshal(n jsonexp.Context, b []byte, out **DoNothing) (err error) {
	if len(b) > 0 {
		var val DoNothing
		if e := DoNothing_Detailed_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}
func DoNothing_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *DoNothing) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(DoNothing_Type, "-", e)
	} else if e := value.Text_Override_Detailed_Optional_Unmarshal(n, msg.Fields[DoNothing_Field_Reason], &out.Reason); e != nil {
		err = errutil.New(DoNothing_Type+"."+DoNothing_Field_Reason, "-", e)
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
const LoggingLevel_Lede = LoggingLevel_Type

func LoggingLevel_Exists(val *LoggingLevel) bool {
	var zero LoggingLevel
	return val.Str != zero.Str
}

func (op *LoggingLevel) MarshalCompact(n jsonexp.Context) (ret []byte, err error) {
	return LoggingLevel_Compact_Marshal(n, op)
}
func (op *LoggingLevel) UnmarshalCompact(n jsonexp.Context, b []byte) error {
	return LoggingLevel_Compact_Unmarshal(n, b, op)
}
func (op *LoggingLevel) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return LoggingLevel_Detailed_Marshal(n, op)
}
func (op *LoggingLevel) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return LoggingLevel_Detailed_Unmarshal(n, b, op)
}

func LoggingLevel_Compact_Optional_Marshal(n jsonexp.Context, val *LoggingLevel) (ret []byte, err error) {
	var zero LoggingLevel
	if val.Str != zero.Str {
		ret, err = LoggingLevel_Compact_Marshal(n, val)
	}
	return
}
func LoggingLevel_Compact_Marshal(n jsonexp.Context, val *LoggingLevel) ([]byte, error) {
	var out string
	if str, ok := composer.FindChoice(val, val.Str); !ok {
		out = val.Str
	} else {
		out = str
	}
	return json.Marshal(out)
}

var LoggingLevel_Compact_Optional_Unmarshal = LoggingLevel_Compact_Unmarshal

func LoggingLevel_Compact_Unmarshal(n jsonexp.Context, b []byte, out *LoggingLevel) (err error) {
	var msg jsonexp.Str
	if len(b) > 0 {
		if e := json.Unmarshal(b, &msg); e != nil {
			err = errutil.New(LoggingLevel_Type, "-", e)
		}
	}
	if err == nil {
		out.Str = msg.Value
	}
	return
}

func LoggingLevel_Detailed_Optional_Marshal(n jsonexp.Context, val *LoggingLevel) (ret []byte, err error) {
	var zero LoggingLevel
	if val.Str != zero.Str {
		ret, err = LoggingLevel_Detailed_Marshal(n, val)
	}
	return
}
func LoggingLevel_Detailed_Marshal(n jsonexp.Context, val *LoggingLevel) ([]byte, error) {
	return json.Marshal(jsonexp.Str{
		Type:  LoggingLevel_Type,
		Value: val.Str,
	})
}

var LoggingLevel_Detailed_Optional_Unmarshal = LoggingLevel_Detailed_Unmarshal

func LoggingLevel_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *LoggingLevel) (err error) {
	var msg jsonexp.Str
	if len(b) > 0 {
		if e := json.Unmarshal(b, &msg); e != nil {
			err = errutil.New(LoggingLevel_Type, "-", e)
		}
	}
	if err == nil {
		out.Str = msg.Value
	}
	return
}

func LoggingLevel_Compact_Repeats_Marshal(n jsonexp.Context, vals *[]LoggingLevel) ([]byte, error) {
	return LoggingLevel_Repeats_Marshal(n, vals, LoggingLevel_Compact_Marshal)
}
func LoggingLevel_Detailed_Repeats_Marshal(n jsonexp.Context, vals *[]LoggingLevel) ([]byte, error) {
	return LoggingLevel_Repeats_Marshal(n, vals, LoggingLevel_Detailed_Marshal)
}
func LoggingLevel_Repeats_Marshal(n jsonexp.Context, vals *[]LoggingLevel, marshEl func(jsonexp.Context, *LoggingLevel) ([]byte, error)) (ret []byte, err error) {
	var msgs []json.RawMessage
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		msgs = make([]json.RawMessage, cnt)
		for i, el := range *vals {
			if b, e := marshEl(n, &el); e != nil {
				err = errutil.New(LoggingLevel_Type, "at", i, "-", e)
				break
			} else {
				msgs[i] = b
			}
		}
	}
	if err == nil {
		ret, err = json.Marshal(msgs)
	}
	return
}

func LoggingLevel_Compact_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]LoggingLevel) error {
	return LoggingLevel_Repeats_Unmarshal(n, b, out, LoggingLevel_Compact_Unmarshal)
}
func LoggingLevel_Detailed_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]LoggingLevel) error {
	return LoggingLevel_Repeats_Unmarshal(n, b, out, LoggingLevel_Detailed_Unmarshal)
}
func LoggingLevel_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]LoggingLevel, unmarshEl func(jsonexp.Context, []byte, *LoggingLevel) error) (err error) {
	var vals []LoggingLevel
	if len(b) > 0 { // generated code collapses optional and empty.
		var msgs []json.RawMessage
		if e := json.Unmarshal(b, &msgs); e != nil {
			err = errutil.New(LoggingLevel_Type, "-", e)
		} else {
			vals = make([]LoggingLevel, len(msgs))
			for i, msg := range msgs {
				if e := unmarshEl(n, msg, &vals[i]); e != nil {
					err = errutil.New(LoggingLevel_Type, "at", i, "-", e)
					break
				}
			}
		}
	}
	if err == nil {
		*out = vals
	}
	return
}

var Slats = []composer.Composer{
	(*DebugLog)(nil),
	(*DoNothing)(nil),
	(*LoggingLevel)(nil),
}
