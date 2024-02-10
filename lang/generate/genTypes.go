package generate

import "sort"

type specData struct {
	Name   string
	Markup markup
}

func (d specData) Comment() (ret []string) {
	if c, ok := d.Markup["comment"]; ok {
		switch c := c.(type) {
		case string:
			ret = []string{c}
		case []string:
			ret = c
		}
	}
	return
}

// because references to types arent scoped but the generated code needs to be:
// the generator has to load all possible types before writing them out.
type flowData struct {
	specData
	Lede  string
	Slots []string
	Terms []termData
}

type markup map[string]any

// fix: cache to stop the multiple creation?
func (m markup) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		if k != "pos" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

type slotData struct {
	specData
}

type strData struct {
	specData
	Options []string
}

type numData struct {
	specData
}

type termData struct {
	Name, Label, Type          string
	Private, Optional, Repeats bool
	Markup                     markup
}

// handle transforming _ into a blank string
func (t *termData) SimpleLabel() (ret string) {
	if t.Label != "_" {
		ret = t.Label
	}
	return
}

// shadows the typeinfo.T interface
// these implementations exist to created the data for that interface.
type typeData interface {
	getName() string
	goType() string
	getMarkup() markup
}

func (f specData) getName() string {
	return f.Name
}

func (f specData) getMarkup() markup {
	return f.Markup
}

func (f specData) goType() string {
	return Pascal(f.Name)
}

func (f numData) goType() string {
	return "float64"
}

func (f strData) goType() (ret string) {
	if len(f.Options) == 0 {
		ret = "string"
	} else {
		switch f.Name {
		case "bool":
			ret = "bool"
		default:
			ret = Pascal(f.Name)
		}
	}
	return
}
