package composer

const (
	Type_Slot = "slot"
	Type_Swap = "swap"
	Type_Flow = "flow"
	Type_Str  = "str"
	Type_Num  = "num"
)

type Composer interface {
	Compose() Spec
}

// Spec definition for display in composer.
type Spec struct {
	Name        string
	Uses        string
	Lede        string   // indicates a fluent style command
	OpenStrings bool     // for str types, whether any value is permitted
	Strings     []string // values for str types, generates tokens, labels, and selectors.
	Choices     []string
	Swaps       []interface{}
}

// given a choice like "$NAME" return the friendly value, ex. "name"
func (spec *Spec) IndexOfChoice(key string) (retVal string, retInd int) {
	retInd = -1 // provisionally
	for i, k := range spec.Choices {
		if k == key {
			retVal, retInd = spec.Strings[i], i
			break
		}
	}
	return
}

// given a value like "name" return the more predefined choice, ex. "$NAME"
func (spec *Spec) IndexOfValue(val string) (retKey string, retInd int) {
	retInd = -1 // provisionally
	for i, str := range spec.Strings {
		if str == val {
			retKey, retInd = spec.Choices[i], i
			break
		}
	}
	return
}
