package grok

// how many sources and how many targets are expected
type MacroType int

const (
	Macro_SourcesOnly MacroType = iota
	Macro_ManySources
	Macro_ManyTargets
	Macro_ManyMany
)
