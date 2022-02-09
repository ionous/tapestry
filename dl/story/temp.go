package story

func (op *EventHandler) Rewrite() {
	if xs := op.Provides; xs != nil {
		op.Locals = xs.Locals
		op.Provides = nil
	}
	if xs := op.PatternRules; xs != nil {
		op.Rules = xs.PatternRule
		op.PatternRules = nil
	}
}
func (op *PatternActions) Rewrite() {
	if xs := op.Provides; xs != nil {
		op.Locals = xs.Locals
		op.Provides = nil
	}
	if xs := op.PatternRules; xs != nil {
		op.Rules = xs.PatternRule
		op.PatternRules = nil
	}
}
func (op *PatternDecl) Rewrite() {
	if xs := op.Params; xs != nil {
		op.Props = xs.Props
		op.Params = nil
	}
}
