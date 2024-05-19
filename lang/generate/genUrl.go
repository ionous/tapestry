package generate

import "path"

var apiDocs = struct {
	site, types, slots string
}{
	"https://tapestry.ionous.net/api/",
	"idl",
	"slot",
}

// vaguely like GlobalData for documentation generation
// fix: maybe documentation generation should have used this package
// instead of the go type data. leave go data uses for runtime info
var hackForLinks *groupSearch

func (p *groupSearch) findGroup(name string) (ret string, okay bool) {
	if n, ok := p.findType(name); ok { // typeEntry
		ret, okay = n.group, true
	}
	return
}

func (p *groupSearch) linkByName(name string) (ret string, okay bool) {
	if n, ok := p.findType(name); ok { // typeEntry
		okay = true
		switch t := n.typeData.(type) {
		case slotData:
			ret = linkToSlot(name, "")
		case flowData:
			if len(t.Slots) > 0 {
				ret = linkToSlot(t.Slots[0], name)
			} else {
				ret = linkToType(n.group, name)
			}
		default:
			ret = linkToType(n.group, name)
		}
	}
	return
}

// slot name is something like "story_statement"
func linkToSlot(slotName, typeName string) string {
	out := apiDocs.site + path.Join(apiDocs.slots, slotName)
	if len(typeName) > 0 {
		out += "#" + typeName
	}
	return out
}

// idl name is a .tells file without the extension.
func linkToType(idlName, typeName string) string {
	out := apiDocs.site + path.Join(apiDocs.types, idlName)
	if len(typeName) > 0 {
		out += "#" + typeName
	}
	return out
}
