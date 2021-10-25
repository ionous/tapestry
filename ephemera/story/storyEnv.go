package story

import "git.sr.ht/~ionous/iffy/ephemera/eph"

type StoryEnv struct {
	Recent struct {
		// Scene, Aspect, Test string
		// Nouns[]? Relation, Trait
		// string or eph.Named
		Nouns Nouns
		Test  eph.Named
	}
	Game struct {
		Domain eph.Named
	}
	Current struct {
		// eventually, a stack.
		Domain eph.Named
	}
	Domains []eph.Named
}

// save current domain and set a new one
func (n *StoryEnv) PushDomain(newDomain eph.Named) {
	n.Current.Domain, n.Domains = newDomain, append(n.Domains, n.Current.Domain)
}

// restore current domain
func (n *StoryEnv) PopDomain() {
	end := len(n.Domains) - 1
	n.Current.Domain, n.Domains = n.Domains[end], n.Domains[:end]
}
