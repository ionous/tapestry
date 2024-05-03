package testpat

import (
	"git.sr.ht/~ionous/tapestry/rt"
)

type Pattern struct {
	Name   string
	Return string     // name of return field; empty if none ( could be an index but slightly safer this way )
	Labels []string   // one label for every parameter
	Fields []rt.Field // flat list of params and locals and an optional return
	Rules  []rt.Rule
}

func (p *Pattern) GetLabels() []string {
	return append(p.Labels, p.Return)
}
