package jsn

import "github.com/ionous/errutil"

// MarshalMix implements the Marshaler interface
// providing functions which can be overridden one at a time to customize functionality
// ( ie. for statemachines )
type MarshalMix struct {
	OnMap     func(string, string) bool
	OnKey     func(string, string) bool
	OnLiteral func(string, string) bool
	OnSlot    func(string, Spotter) bool
	OnPick    func(string, Picker) bool
	OnRepeat  func(string, Slicer) bool
	OnEnd     func()
	OnValue   func(string, interface{})
	OnWarn    func(error)
}

func (ms *MarshalMix) MapValues(lede, typeName string) (okay bool) {
	if call := ms.OnMap; call != nil {
		okay = call(lede, typeName)
	} else {
		ms.Error(errutil.New("unexpected map", lede, typeName))
	}
	return
}
func (ms *MarshalMix) MapLiteral(lede, typeName string) (okay bool) {
	if call := ms.OnLiteral; call != nil {
		okay = call(lede, typeName)
	} else {
		// fall back to a regular map
		okay = ms.MapValues(lede, typeName)
	}
	return
}
func (ms *MarshalMix) MapKey(key, field string) (okay bool) {
	if call := ms.OnKey; call != nil {
		okay = call(key, field)
	} else {
		ms.Error(errutil.New("unexpected key", key, field))
	}
	return
}
func (ms *MarshalMix) SlotValues(typeName string, val Spotter) (okay bool) {
	if call := ms.OnSlot; call != nil {
		okay = call(typeName, val)
	} else {
		ms.Error(errutil.New("unexpected pick", typeName, val))
	}
	return
}
func (ms *MarshalMix) PickValues(typeName string, val Picker) (okay bool) {
	if call := ms.OnPick; call != nil {
		okay = call(typeName, val)
	} else {
		ms.Error(errutil.New("unexpected pick", typeName, val))
	}
	return
}
func (ms *MarshalMix) RepeatValues(typeName string, val Slicer) (okay bool) {
	if call := ms.OnRepeat; call != nil {
		okay = call(typeName, val)
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
func (ms *MarshalMix) MarshalValue(typeName string, pv interface{}) (okay bool) {
	if call := ms.OnValue; call != nil {
		call(typeName, pv)
		okay = true
	} else {
		ms.Error(errutil.New("unexpected value", typeName, pv))
	}
	return
}
func (ms *MarshalMix) Error(err error) {
	if call := ms.OnWarn; call != nil {
		call(err)
	} else {
		panic(err)
	}
}
