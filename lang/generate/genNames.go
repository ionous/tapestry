package generate

type Names []string

func (gs *Names) AddName(g string) {
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
