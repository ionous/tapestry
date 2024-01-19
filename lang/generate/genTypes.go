package generate

// because references to types arent scoped but the generated code needs to be:
// the generator has to load all possible types before writing them out.
type flowData struct {
	Name  string
	Lede  string
	Slots []string
	Terms []termData
	// Comment []string
	// Markup? - might want comments in there too for shapes
}

type slotData struct {
	Name string
}

type strData struct {
	Name    string
	Options []string
}

type numData struct {
	Name string
}

type termData struct {
	// these are lower-case names
	Name, Label, Type          string
	Private, Optional, Repeats bool
}

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
