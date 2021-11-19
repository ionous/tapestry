package eph

type Conflict struct {
	Reason ReasonForConflict
	Domain string
	Def    Definition
}

func (n *Conflict) Error() string {
	return n.Reason.String()
}

type ReasonForConflict int

//go:generate stringer -type=ReasonForConflict
const (
	Redefined ReasonForConflict = iota
	Duplicated
)

type Definitions map[string]Definition

type Definition struct {
	at, value string
	err       error
}

type DomainConflicts map[string]Definitions

// walks the properly cased named domain's dependencies ( non-recursively ) to find
// whether the new key,value pair contradicts or duplicates any existing value.
func (dc DomainConflicts) CheckConflicts(n string, l DependencyFinder, cat, at, key, value string) (err error) {
	fullKey := cat + " " + key
	if e := dc.checkConflict(n, fullKey, value); e != nil {
		err = e
	} else if deps, e := GetResolvedDependencies(n, l); e != nil {
		err = e
	} else {
		for _, depName := range deps.Ancestors(true) {
			if e := dc.checkConflict(depName, fullKey, value); e != nil {
				err = e
				break
			}
		}
		defines := dc[n]
		if defines == nil {
			defines = make(Definitions)
			dc[n] = defines
		}
		// store the result: rather than checking through all the resolved domains the next time
		// ( if there is a next time )
		defines[fullKey] = Definition{
			at:    at,
			value: value,
			err:   err, // might be null and that's cool.
		}
	}
	return
}

// was anything stored before?
func (dc DomainConflicts) checkConflict(n, key, value string) (err error) {
	if defines, ok := dc[n]; ok {
		if def, ok := defines[key]; ok {
			if def.err != nil {
				err = def.err
			} else {
				var why ReasonForConflict
				if def.value == value {
					why = Duplicated // if its duplicated, the previous entry would have checked for redefined
				} else {
					why = Redefined
				}
				err = &Conflict{why, n, def}
			}
		}
	}
	return
}
