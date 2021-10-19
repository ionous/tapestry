package cout

import "encoding/json"

// the compact format supports writing swaps in one of two ways:
// the expanded format: { "swapName choice:": value, flow, etc. }
// or, when embedded in a flow: { "flow param choice:": value, flow, etc. }
// the embedded format saves an unnecessary level of depth for the most common situation.
// while the expanded format should be allowed in flows, it's for use in repeats and slots.
type comSwap struct {
	typeName, choice string // ex."noun_phrase" "$KIND_OF_NOUN"
	value            interface{}
}

// the "unpacked" value.
func (c *comSwap) SetValue(v interface{}) *comSwap {
	c.value = v
	return c
}

func (c *comSwap) MarshalJSON() ([]byte, error) {
	var sig Sig
	sig.WriteLede(c.typeName)
	sig.WriteLabel(c.choice)
	return json.Marshal(map[string]interface{}{
		sig.String(): c.value,
	})
}
