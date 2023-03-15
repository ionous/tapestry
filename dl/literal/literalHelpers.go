package literal

func B(b bool) *BoolValue        { return &BoolValue{Value: b} }
func I(n int) *NumValue          { return &NumValue{Value: float64(n)} }
func F(n float64) *NumValue      { return &NumValue{Value: n} }
func T(s string) *TextValue      { return &TextValue{Value: s} }
func Ts(s ...string) *TextValues { return &TextValues{Values: s} }
