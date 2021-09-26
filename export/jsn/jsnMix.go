package jsn

import "github.com/ionous/errutil"

// MarshalMix implements the Marshaler interface
// providing functions which can be overridden one at a time to customize functionality
// ( ie. for statemachines )
type MarshalMix struct {
	OnMap     func(lede, kind string)
	OnKey     func(sig, field string)
	OnLiteral func(field string)
	OnPick    func(kind, choice string)
	OnRepeat  func(hint int)
	OnEnd     func()
	OnValue   func(kind string, value interface{})
	OnCursor  func(id string)
	OnWarn    func(error)
	OnError   func(error)
}

func (ms *MarshalMix) MapValues(lede, kind string) {
	if call := ms.OnMap; call != nil {
		call(lede, kind)
	} else {
		ms.Error(errutil.New("unexpected map", lede, kind))
	}
}
func (ms *MarshalMix) MapKey(key, field string) {
	if call := ms.OnKey; call != nil {
		call(key, field)
	} else {
		ms.Error(errutil.New("unexpected key", key, field))
	}
}
func (ms *MarshalMix) MapLiteral(field string) {
	if call := ms.OnLiteral; call != nil {
		call(field)
	} else {
		ms.Error(errutil.New("unexpected literal", field))
	}
}
func (ms *MarshalMix) PickValues(kind, choice string) {
	if call := ms.OnPick; call != nil {
		call(kind, choice)
	} else {
		ms.Error(errutil.New("unexpected pick", kind, choice))
	}
}
func (ms *MarshalMix) RepeatValues(hint int) {
	if call := ms.OnRepeat; call != nil {
		call(hint)
	} else {
		ms.Error(errutil.New("unexpected repeat", hint))
	}
}
func (ms *MarshalMix) EndValues() {
	if call := ms.OnEnd; call != nil {
		call()
	} else {
		ms.Error(errutil.New("unexpected end"))
	}
}
func (ms *MarshalMix) SpecifyValue(kind string, value interface{}) {
	if call := ms.OnValue; call != nil {
		call(kind, value)
	} else {
		ms.Error(errutil.New("unexpected value", kind, value))
	}
}
func (ms *MarshalMix) SpecifyEnum(kind string, val Enumeration) {
	var out string
	if str, ok := val.FindChoice(); !ok {
		out = val.String()
	} else {
		out = str
	}
	ms.SpecifyValue(kind, out)
}
func (ms *MarshalMix) SetCursor(id string) {
	if call := ms.OnCursor; call != nil {
		call(id)
	} else {
		ms.Error(errutil.New("unexpected cursor", id))
	}
}
func (ms *MarshalMix) Warning(err error) {
	if call := ms.OnWarn; call != nil {
		call(err)
	} else {
		ms.Error(err)
	}
}
func (ms *MarshalMix) Error(err error) {
	if call := ms.OnError; call != nil {
		call(err)
	} else {
		panic(err)
	}
}
