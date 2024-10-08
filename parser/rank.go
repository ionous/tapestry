package parser

// RankNoun implementations accumulate targets for actions during calls to RankNouns.
type RankNoun interface {
	// returning false indicates some critical error that should cancel ranking.
	RankNoun(Cursor, NounInstance) bool
}

// RankNouns visits each noun in a bounds, calling RankNoun.
// it returns false if RankNoun returns false.
func RankNouns(bounds Bounds, cs Cursor, r RankNoun) bool {
	return !bounds(func(n NounInstance) bool {
		return !r.RankNoun(cs, n)
	})
}

type RankOne struct {
	Filters
	Ranking
}

// Ranking accumulates Nouns at a given Rank.
// Rank counts the number of words that match a given NounInstance
// It's possible for different nouns to share the same rank for some given set of words.
// For example, the "real eiffel tower" and the "toy eiffel tower"
// would share a rank of two for the words: "tower eiffel"
type Ranking struct {
	Rank  int
	Nouns []NounInstance
}

func (r *Ranking) Empty() bool {
	return r.Rank == 0
}

func (r *Ranking) AddRanking(n NounInstance, rank int) {
	switch {
	case rank > r.Rank:
		r.Rank, r.Nouns = rank, []NounInstance{n}
	case rank == r.Rank:
		r.Nouns = append(r.Nouns, n)
	}
}

// RankNoun - never returns false
func (m *RankOne) RankNoun(cs Cursor, n NounInstance) bool {
	if m.MatchesNoun(n) {
		var rank int
		for ; ; rank++ {
			if name := cs.CurrentWord(); len(name) > 0 && n.HasName(name) {
				cs = cs.Skip(1)
			} else {
				break
			}
		}
		if rank > 0 {
			m.Ranking.AddRanking(n, rank)
		}
	}
	return true // always keep going
}

type RankAll struct {
	Filters
	Context Context
	// we dont know what follows the keyword "all"
	// if it turns out that its a word which identifies one or more objects
	// then we really dont want "all" anymore, we simply want those objects.
	// in the meantime, accumulate all "unmentioned" objects
	Implied   []NounInstance
	Plurals   []string
	WordCount int
	Ranking
}

// RankNoun - never returns fals.
func (m *RankAll) RankNoun(cs Cursor, n NounInstance) bool {
	if m.MatchesNoun(n) {
		var rank, cnt int
		for {
			var matches bool
			if name := cs.CurrentWord(); len(name) > 0 {
				if n.HasName(name) {
					matches = true
					rank++
				}
				// note: we dont test whether the noun applies to this plural
				// we will have to test all accumulated nouns anyway.
				if name := cs.CurrentWord(); m.Context.IsPlural(name) {
					m.Plurals = append(m.Plurals, name)
					matches = true
				}
			}
			//
			if matches {
				cs = cs.Skip(1)
				cnt++
			} else {
				break
			}
		}
		if rank > 0 {
			m.Ranking.AddRanking(n, rank)
		} else if m.Ranking.Empty() {
			m.Implied = append(m.Implied, n)
		}
		if cnt > m.WordCount {
			m.WordCount = cnt
		}
	}
	return true // always keep going
}
