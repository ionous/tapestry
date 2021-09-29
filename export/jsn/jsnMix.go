package jsn

import "github.com/ionous/errutil"

// MarshalMix implements the Marshaler interface
// providing functions which can be overridden one at a time to customize functionality
// ( ie. for statemachines )
type MarshalMix struct {
	OnMap     func(lede, kind string) bool
	OnKey     func(sig, field string) bool
	OnLiteral func(field string) bool
	OnPick    func(Picker) bool
	OnRepeat  func(hint int) bool
	OnEnd     func()
	OnBool    func(BoolMarshaler)
	OnEnum    func(EnumMarshaler)
	OnNum     func(NumMarshaler)
	OnStr     func(StrMarshaler)
	OnCursor  func(id string)
	OnWarn    func(error)
	OnError   func(error)
}

func (ms *MarshalMix) MapValues(lede, kind string) (ret bool) {
	if call := ms.OnMap; call != nil {
		ret = call(lede, kind)
	} else {
		ms.Error(errutil.New("unexpected map", lede, kind))
	}
	return
}
func (ms *MarshalMix) MapKey(key, field string) (ret bool) {
	if call := ms.OnKey; call != nil {
		ret = call(key, field)
	} else {
		ms.Error(errutil.New("unexpected key", key, field))
	}
	return
}
func (ms *MarshalMix) MapLiteral(field string) (ret bool) {
	if call := ms.OnLiteral; call != nil {
		ret = call(field)
	} else {
		ms.Error(errutil.New("unexpected literal", field))
	}
	return
}
func (ms *MarshalMix) PickValues(val Picker) (ret bool) {
	if call := ms.OnPick; call != nil {
		ret = call(val)
	} else {
		ms.Error(errutil.New("unexpected pick", val))
	}
	return
}
func (ms *MarshalMix) RepeatValues(hint int) (ret bool) {
	if call := ms.OnRepeat; call != nil {
		ret = call(hint)
	} else {
		ms.Error(errutil.New("unexpected repeat", hint))
	}
	return
}
func (ms *MarshalMix) EndValues() {
	if call := ms.OnEnd; call != nil {
		call()
	} else {
		ms.Error(errutil.New("unexpected end"))
	}
}
func (ms *MarshalMix) BoolValue(val BoolMarshaler) {
	if call := ms.OnBool; call != nil {
		call(val)
	} else {
		ms.Error(errutil.New("unexpected value", val))
	}
}
func (ms *MarshalMix) EnumValue(val EnumMarshaler) {
	if call := ms.OnEnum; call != nil {
		call(val)
	} else {
		ms.Error(errutil.New("unexpected value", val))
	}
}
func (ms *MarshalMix) NumValue(val NumMarshaler) {
	if call := ms.OnNum; call != nil {
		call(val)
	} else {
		ms.Error(errutil.New("unexpected value", val))
	}
}
func (ms *MarshalMix) StrValue(val StrMarshaler) {
	if call := ms.OnStr; call != nil {
		call(val)
	} else {
		ms.Error(errutil.New("unexpected value", val))
	}
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
