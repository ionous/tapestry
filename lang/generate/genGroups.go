package generate

import "sort"

type Group struct {
	Name string
	groupContent
}

type groupContent struct {
	Flow []typeData
	Slot []typeData
	Str  []typeData
	Num  []typeData
}

// dont forget to exclude self
type groupSearch struct {
	list         []Group
	currentGroup string
	refs         map[string]bool
}

// string writer would be faster; oh well.
func (p *groupSearch) getRefs(base string, extra ...string) []string {
	out := make([]string, 0, len(p.refs)+len(extra))
	for k := range p.refs {
		out = append(out, base+k)
	}
	out = append(out, extra...)
	sort.Strings(out)
	return out
}

func (p *groupSearch) setCurrent(at int) Group {
	g := p.list[at]
	p.currentGroup = g.Name
	p.refs = make(map[string]bool)
	return g
}

// return scoped typename
func (p *groupSearch) findType(n string) (retGroup, retName string) {
	for _, g := range p.list {
		if str := g.findType(n); len(str) > 0 {
			if group := g.Name; p.currentGroup != group {
				retGroup = group
				p.refs[group] = true
			}
			retName = str
			break
		}
	}
	return
}

// return scoped typename
func (g *groupContent) findType(n string) (ret string) {
	if _, ok := findType(n, g.Slot); ok {
		ret = Pascal(n)
	} else if _, ok := findType(n, g.Flow); ok {
		ret = Pascal(n)
	} else if _, ok := findType(n, g.Num); ok {
		ret = "float64"
	} else if t, ok := findType(n, g.Str); ok {
		//
		if str := t.(strData); len(str.Options) == 0 {
			ret = "string"
		} else {
			switch n {
			case "bool":
				ret = "bool"
			default:
				ret = "string" // tbd: or maybe dont switch from Str structs just yet
			}
		}
	}
	return

}

func findType(n string, els []typeData) (ret typeData, okay bool) {
	for _, el := range els {
		if el.getName() == n {
			ret = el
			okay = true
			break
		}
	}
	return
}
