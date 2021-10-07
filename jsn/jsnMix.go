package jsn

import "github.com/ionous/errutil"

// MarshalMix implements the Marshaler interface
// providing functions which can be overridden one at a time to customize functionality
// ( ie. for statemachines )
type MarshalMix struct {
	OnMap     func(lede, typeName string) bool
	OnKey     func(sig, field string) bool
	OnLiteral func(string) bool
	OnPick    func(string, Picker) bool
	OnRepeat  func(string, Slicer) bool
	OnEnd     func()
	OnValue   func(string, interface{})
	OnCursor  func(string)
	OnWarn    func(error)
	OnError   func(error)
}

func (ms *MarshalMix) IsEncoding() bool {
	panic("not implemented")
}
func (ms *MarshalMix) MapValues(lede, typeName string) (ret bool) {
	if call := ms.OnMap; call != nil {
		ret = call(lede, typeName)
	} else {
		ms.Error(errutil.New("unexpected map", lede, typeName))
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
func (ms *MarshalMix) PickValues(typeName string, val Picker) (ret bool) {
	if call := ms.OnPick; call != nil {
		ret = call(typeName, val)
	} else {
		ms.Error(errutil.New("unexpected pick", typeName, val))
	}
	return
}
func (ms *MarshalMix) RepeatValues(typeName string, val Slicer) (ret bool) {
	if call := ms.OnRepeat; call != nil {
		ret = call(typeName, val)
	} else {
		ms.Error(errutil.New("unexpected repeat", typeName, val))
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
func (ms *MarshalMix) GenericValue(typeName string, pv interface{}) {
	if call := ms.OnValue; call != nil {
		call(typeName, pv)
	} else {
		ms.Error(errutil.New("unexpected value", typeName, pv))
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
