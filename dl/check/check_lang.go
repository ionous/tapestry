// Code generated by "makeops"; edit at your own risk.
package check

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/jsn"
	"github.com/ionous/errutil"
)

// CheckOutput Run an activity and verify that it produces the expected output.
type CheckOutput struct {
	Name   string        `if:"label=_,type=text"`
	Expect string        `if:"label=expect,type=text"`
	Test   core.Activity `if:"label=do"`
}

func (*CheckOutput) Compose() composer.Spec {
	return composer.Spec{
		Name: CheckOutput_Type,
		Uses: composer.Type_Flow,
		Lede: "check",
	}
}

const CheckOutput_Type = "check_output"

const CheckOutput_Field_Name = "$NAME"
const CheckOutput_Field_Expect = "$EXPECT"
const CheckOutput_Field_Test = "$TEST"

func (op *CheckOutput) Marshal(m jsn.Marshaler) error {
	return CheckOutput_Marshal(m, op)
}

type CheckOutput_Slice []CheckOutput

func (op *CheckOutput_Slice) GetType() string { return CheckOutput_Type }

func (op *CheckOutput_Slice) Marshal(m jsn.Marshaler) error {
	return CheckOutput_Repeats_Marshal(m, (*[]CheckOutput)(op))
}

func (op *CheckOutput_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *CheckOutput_Slice) SetSize(cnt int) {
	var els []CheckOutput
	if cnt >= 0 {
		els = make(CheckOutput_Slice, cnt)
	}
	(*op) = els
}

func (op *CheckOutput_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return CheckOutput_Marshal(m, &(*op)[i])
}

func CheckOutput_Repeats_Marshal(m jsn.Marshaler, vals *[]CheckOutput) error {
	return jsn.RepeatBlock(m, (*CheckOutput_Slice)(vals))
}

func CheckOutput_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]CheckOutput) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = CheckOutput_Repeats_Marshal(m, pv)
	}
	return
}

type CheckOutput_Flow struct{ ptr *CheckOutput }

func (n CheckOutput_Flow) GetType() string      { return CheckOutput_Type }
func (n CheckOutput_Flow) GetLede() string      { return "check" }
func (n CheckOutput_Flow) GetFlow() interface{} { return n.ptr }
func (n CheckOutput_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*CheckOutput); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func CheckOutput_Optional_Marshal(m jsn.Marshaler, pv **CheckOutput) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = CheckOutput_Marshal(m, *pv)
	} else if !enc {
		var v CheckOutput
		if err = CheckOutput_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func CheckOutput_Marshal(m jsn.Marshaler, val *CheckOutput) (err error) {
	if err = m.MarshalBlock(CheckOutput_Flow{val}); err == nil {
		e0 := m.MarshalKey("", CheckOutput_Field_Name)
		if e0 == nil {
			e0 = literal.Text_Unboxed_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", CheckOutput_Field_Name))
		}
		e1 := m.MarshalKey("expect", CheckOutput_Field_Expect)
		if e1 == nil {
			e1 = literal.Text_Unboxed_Marshal(m, &val.Expect)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", CheckOutput_Field_Expect))
		}
		e2 := m.MarshalKey("do", CheckOutput_Field_Test)
		if e2 == nil {
			e2 = core.Activity_Marshal(m, &val.Test)
		}
		if e2 != nil && e2 != jsn.Missing {
			m.Error(errutil.New(e2, "in flow at", CheckOutput_Field_Test))
		}
		m.EndBlock()
	}
	return
}

var Slats = []composer.Composer{
	(*CheckOutput)(nil),
}

var Signatures = map[uint64]interface{}{
	17982686721133675429: (*CheckOutput)(nil), /* Check:expect:do: */
}
