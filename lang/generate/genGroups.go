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
	Reg  Registry
}

// dont forget to exclude self
type groupSearch struct {
	list         []Group
	currentGroup string
	refs         map[string]bool
	types        typeMap
}

type typeEntry struct {
	group string
	typeData
}

type typeMap map[string]typeEntry

// sorted keys
func (m typeMap) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func newGroupSearch(list []Group) *groupSearch {
	m := make(typeMap)
	for _, g := range list {
		addTypes(m, g, g.Flow)
		addTypes(m, g, g.Slot)
		addTypes(m, g, g.Str)
		addTypes(m, g, g.Num)
	}
	// fix? uses pointer for the sake of sharing currentGroup; use a context?
	return &groupSearch{
		list:  list,
		types: m,
	}
}

func addTypes(m typeMap, g Group, slice []typeData) {
	for _, x := range slice {
		m[x.getName()] = typeEntry{g.Name, x}
	}
}

// a list of groups, sorted.
// string writer (direct to output) would be faster; oh well.
func (p *groupSearch) getImports(base string, extra ...string) []string {
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

// return scoped go type
func (p *groupSearch) findType(n string) (retGroup, retName string) {
	if a, ok := p.types[n]; ok {
		if p.currentGroup != a.group {
			retGroup = a.group
			p.refs[a.group] = true
		}
		retName = a.goType()
	}
	return
}
