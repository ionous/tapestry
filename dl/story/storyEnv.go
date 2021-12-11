package story

type StoryEnv struct {
	Recent struct {
		// Scene, Aspect, Test string
		// Nouns[]? Relation, Trait
		// string or string
		Nouns Nouns
		Test  string
	}
}
