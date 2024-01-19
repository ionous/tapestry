package generate

type Groups []string

func (gs *Groups) AddGroup(g string) {
	var dupe bool
	for _, el := range *gs {
		if el == g {
			dupe = true
			break
		}
	}
	if !dupe {
		*gs = append(*gs, g)
	}
}

type groupData struct {
	Name string
	groupContent
}

type groupContent struct {
	Flow []typeData
	Slot []typeData
	Str  []typeData
	Num  []typeData
}

func findType(n string, els []typeData) (ret typeData, okay bool) {
	for _, el := range els {
		if el.GetName() == n {
			ret = el
			okay = true
			break
		}
	}
	return
}
