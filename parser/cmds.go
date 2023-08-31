package parser

// Commands - reference for all of objects implementing Scanner
type Commands struct {
	*Action
	*AllOf
	*AnyOf
	*Focus
	*Multi
	*Noun
	*Refine
	*Word
}
