package groktest

// implements Registrar to watch incoming calls.
type Mock []string

func (m *Mock) AddKind(kind, ancestor string) (_ error) {
	(*m) = append(*m, "AddKind", kind, ancestor)
	return
}

func (m *Mock) AddKindTrait(kind, trait string) (_ error) {
	(*m) = append(*m, "AddKindTrait", kind, trait)
	return
}
