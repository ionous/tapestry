package grok

type MacroType int

const (
	ManyToOne MacroType = iota
	OneToMany
	ManyToMany
)
