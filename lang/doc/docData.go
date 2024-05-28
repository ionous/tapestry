package doc

import (
	"cmp"
	"fmt"
	"path"
	"slices"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

type GlobalData struct {
	// idl name -> typeset
	types map[string]typeinfo.TypeSet
	// slots name -> slots info
	slots map[string]SlotInfo
	// cmds name -> cmds
	cmds map[string]FlowInfo
	// str/num name -> flow
	prims map[string]PrimInfo
	//
	allCommands []FlowInfo
}

type SlotInfo struct {
	Idl string
	*typeinfo.Slot
}

type FlowInfo struct {
	Idl string // ex. "story"
	*typeinfo.Flow
	Spec string // Define something:
}

type PrimInfo struct {
	Idl string
	typeinfo.T
}

// link to go package documentation
func (c FlowInfo) SourceLink() string {
	name := inflect.Pascal(c.Name)
	pkgPath := "dl/" + c.Idl
	return path.Join(SourceUrl, pkgPath+"#"+name)
}

func (c FlowInfo) Terms() (ret []typeinfo.Term) {
	// filter out hidden terms;
	for i, t := range c.Flow.Terms {
		if t.Private && ret == nil {
			ret = c.Flow.Terms[:i]
		}
		if !t.Private && ret != nil {
			ret = append(ret, t)
		}
	}
	if ret == nil {
		ret = c.Flow.Terms
	}
	return
}

func makeGlobalData(idl []typeinfo.TypeSet) GlobalData {
	types := make(map[string]typeinfo.TypeSet)
	flow := make(map[string]FlowInfo)
	slot := make(map[string]SlotInfo)
	prim := make(map[string]PrimInfo)
	var flat []FlowInfo
	for _, ts := range idl {
		types[ts.Name] = ts
		//
		for _, t := range ts.Flow {
			if _, hidden := t.Markup["internal"]; !hidden {
				spec := BuildSpec(t)
				i := FlowInfo{ts.Name, t, spec}
				flow[t.Name] = i
				flat = append(flat, i)
			}
		}
		for _, t := range ts.Slot {
			slot[t.Name] = SlotInfo{ts.Name, t}
		}
		for _, t := range ts.Str {
			prim[t.Name] = PrimInfo{ts.Name, t}
		}
		for _, t := range ts.Num {
			prim[t.Name] = PrimInfo{ts.Name, t}
		}
	}
	// sort all the commands by spec.
	slices.SortFunc(flat, func(a, b FlowInfo) int {
		// we want specs ending in colons to be listed before ones with a new word.
		// ex. Numeral: before Numeral words:
		// colon is less than underscore ( but greater than space )
		x, y := strings.Replace(a.Spec, " ", "_", -1), strings.Replace(b.Spec, " ", "_", -1)
		return cmp.Compare(x, y)
	})

	return GlobalData{
		types, slot, flow, prim, flat,
	}
}

func (g *GlobalData) linkByIdl(name string) (ret string, err error) {
	if _, ok := g.types[name]; !ok {
		err = fmt.Errorf("couldnt find idl %q", name)
	} else {
		ret = linkToType(name, "")
	}
	return
}

func (g *GlobalData) linkByType(t typeinfo.T) (ret string, err error) {
	return g.linkByName(t.TypeName())
}

// ugh. part of the issue is that go-templates dont take multiple parametrs
func (g *GlobalData) linkByName(name string) (ret string, err error) {
	if _, ok := g.slots[name]; ok {
		ret = linkToSlot(name, "")
	} else if flow, ok := g.cmds[name]; ok {
		ret = linkToType(flow.Idl, name)
	} else if prim, ok := g.prims[name]; ok {
		ret = linkToType(prim.Idl, name)
	} else {
		err = fmt.Errorf("couldnt find type %q", name)
	}
	return
}

// Build the document style signature for this flow
// it's different than the actual signature because,
// among other things, it includes markers for optional elements.
func BuildSpec(t *typeinfo.Flow) string {
	var str strings.Builder
	str.WriteString(inflect.Pascal(t.Lede))
	if len(t.Terms) == 0 {
		str.WriteString(":")
	} else {
		for i, t := range t.Terms {
			if t.Optional {
				str.WriteRune('[')
			}
			if len(t.Label) > 0 {
				if i == 0 {
					str.WriteRune(' ')
				}
				str.WriteString(inflect.Camelize(t.Label))
			}
			str.WriteRune(':')
			if t.Optional {
				str.WriteRune(']')
			}
		}
	}
	return str.String()
}

// idl name is a .tells file without the extension.
func linkToType(idlName, typeName string) string {
	out := path.Join(baseUrl, typesFolder, idlName)
	if len(typeName) > 0 {
		out += "#" + typeName
	}
	return out
}

// slot name is something like "story_statement"
func linkToSlot(slotName, typeName string) string {
	out := path.Join(baseUrl, slotFolder, slotName)
	if len(typeName) > 0 {
		out += "#" + typeName
	}
	return out
}
