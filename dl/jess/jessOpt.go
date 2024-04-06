package jess

// a generic method to read an optional element.
// optional elements are implemented by pointer types
// the pointers are required to implement Interpreter.
// returns true if matched.
func Optional[M any,
	IM interface {
		// the syntax for this feels very strange
		// the method takes an 'out' pointer ( so go can determine the type by inference )
		// "interface" here is a reused keyword: signifying a type constraint.
		// it means we want pointers to M to implement Interpreter.
		// *phew*
		*M
		Interpreter
	}](q Query, input *InputState, out **M) (okay bool) {
	var v M
	if next := *input; IM(&v).Match(q, &next) {
		*out, *input, okay = &v, next, true
	}
	return
}

// go is unable to reliably test typed interfaces for nil without reflection.
// however, look how easy this is. *sigh*
// https://go.dev/doc/faq#nil_error
func ref[M any,
	IM interface {
		*M
		NounBuilder
	}](in *M) nounBuilderRef {
	return nounBuilderRef{IM(in), in == nil}
}

type nounBuilderRef struct {
	NounBuilder
	IsNil bool
}
