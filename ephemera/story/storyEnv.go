package story

import "git.sr.ht/~ionous/iffy/ephemera"

type StoryEnv struct {
	Recent struct {
		// Scene, Aspect, Test string
		// Nouns[]? Relation, Trait
		// string or ephemera.Named
		Nouns Nouns
		Test  ephemera.Named
	}
	Game struct {
		Domain ephemera.Named
	}
	Current struct {
		// eventually, a stack.
		Domain ephemera.Named
	}
	Domains []ephemera.Named
}

// save current domain and set a new one
func (n *StoryEnv) PushDomain(newDomain ephemera.Named) {
	n.Current.Domain, n.Domains = newDomain, append(n.Domains, n.Current.Domain)
}

// restore current domain
func (n *StoryEnv) PopDomain() {
	end := len(n.Domains) - 1
	n.Current.Domain, n.Domains = n.Domains[end], n.Domains[:end]
}
