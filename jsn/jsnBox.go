package jsn

func BoxBool(v *bool) (ret BoxedBool) {
	return BoxedBool{v}
}
func BoxFloat64(v *float64) BoxedFloat {
	return BoxedFloat{v}
}
func BoxString(v *string) BoxedString {
	return BoxedString{v}
}

type BoxedBool struct {
	v *bool
}

type BoxedFloat struct {
	v *float64
}

type BoxedString struct {
	v *string
}

func (box BoxedBool) GetValue() (ret interface{}) {
	if *box.v {
		ret = "$TRUE"
	} else {
		ret = "$FALSE"
	}
	return
}

func (box BoxedBool) GetCompactValue() (ret interface{}) {
	if *box.v {
		ret = "true"
	} else {
		ret = "false"
	}
	return
}

func (box BoxedBool) SetValue(v interface{}) (okay bool) {
	switch v.(string) {
	case "$TRUE", "true":
		box.setValue(true)
		okay = true
	case "$FALSE", "false":
		box.setValue(false)
		okay = true
	}
	return
}

func (box BoxedBool) setValue(b bool) {
	*box.v = b
}

func (box BoxedFloat) GetValue() interface{} {
	return *box.v
}
func (box BoxedFloat) SetValue(v interface{}) (okay bool) {
	if n, ok := v.(float64); ok {
		*box.v = n
		okay = true
	}
	return
}

func (box BoxedString) GetValue() interface{} {
	return *box.v
}
func (box BoxedString) SetValue(v interface{}) (okay bool) {
	if n, ok := v.(string); ok {
		*box.v = n
		okay = true
	}
	return
}
