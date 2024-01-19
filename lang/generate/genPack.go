package generate

// dont forget to exclude self
type packageHelper struct {
	groups         []groupData
	currentPackage string
}

func (p *packageHelper) FindScope(n string) (ret string, okay bool) {
	for _, g := range p.groups {
		_, slot := findType(n, g.Slot)
		_, flow := findType(n, g.Flow)
		if slot || flow {
			ret = p.scope(g.Name)
			break
		}
	}
	return
}
func (p *packageHelper) scope(group string) (ret string) {
	if p.currentPackage != group {
		ret = group + "."
	}
	return
}

// return scoped typename
func (p *packageHelper) TypeName(n string) (ret string) {
	for _, g := range p.groups {
		if _, ok := findType(n, g.Slot); ok {
			ret = p.scope(g.Name) + Pascal(n)
			break

		} else if _, ok := findType(n, g.Flow); ok {
			ret = p.scope(g.Name) + Pascal(n)
			break

		} else if _, ok := findType(n, g.Num); ok {
			ret = "float64"
			break

		} else {
			if t, ok := findType(n, g.Str); ok {
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
			break
		}
	}
	return
}
