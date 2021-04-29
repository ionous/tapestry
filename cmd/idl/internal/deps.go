package internal

type Deps []string

func (deps Deps) AddDep(s string) Deps {
	if len(s) > 0 {
		i, cnt := 0, len(deps)
		for ; i < cnt; i++ {
			if deps[i] == s {
				break
			}
		}
		if i == cnt { // never found.
			deps = append(deps, s)
		}
	}
	return deps
}
