package generate

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
