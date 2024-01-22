package generate

// because references to types arent scoped but the generated code needs to be:
// the generator has to load all possible types before writing them out.
type flowData struct {
	Name   string
	Lede   string
	Slots  []string
	Terms  []termData
	Markup map[string]any
}

type slotData struct {
	Name   string
	Markup map[string]any
}

type strData struct {
	Name    string
	Options []string
	Markup  map[string]any
}

type numData struct {
	Name   string
	Markup map[string]any
}

type termData struct {
	Name, Label, Type          string
	Private, Optional, Repeats bool
	Markup                     map[string]any
}

// shadows the typeinfo.T interface
// these implementations exist to created the data for that interface.
type typeData interface {
	getName() string
}

func (f flowData) getName() string {
	return f.Name
}

func (f slotData) getName() string {
	return f.Name
}

func (f strData) getName() string {
	return f.Name
}

func (f numData) getName() string {
	return f.Name
}
