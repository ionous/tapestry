package composer

const (
	Type_Slot = "slot"
	Type_Swap = "swap"
	Type_Flow = "flow"
	Type_Str  = "str"
	Type_Num  = "num"
)

// Slot definition for display in the composer.
type Slot struct {
	Name  string
	Type  interface{} // nil instance, ex. (*core.Comparator)(nil)
	Desc  string
	Group string // display group(s)
}

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

func (spec *Spec) UsesStr() bool {
	return spec.OpenStrings || len(spec.Strings) > 0
}

func (spec *Spec) FindChoice(choice string) (ret string, okay bool) {
	if len(choice) > 0 {
		if choice[0] != '$' {
			if spec.OpenStrings {
				ret = choice
				okay = true
			}
		} else if s, i := spec.IndexOfChoice(choice); i >= 0 {
			ret = s
			okay = true
		}
	}
	return
}

func (spec *Spec) IndexOfChoice(choice string) (retStr string, retInd int) {
	retInd = -1 // provisionally
	for i, c := range spec.Choices {
		if c == choice {
			retStr, retInd = spec.Strings[i], i
			break
		}
	}
	return
}
