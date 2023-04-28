package eph

import "git.sr.ht/~ionous/tapestry/imp/assert"

type Context struct {
	c     *Catalog
	d     *Domain
	at    string
	phase assert.Phase
}

func (ctx *Context) BeginDomain(name string, requires []string) error {
	panic("BeginDomain")
}

func (ctx *Context) EndDomain() error {
	panic("EndDomain")
}

// func (ctx *Context) AssertNounPhrase() error {
// 	panic("AssertNounPhrase")
// }
