package qna

import "math/rand"

// the nil value initialize itself
type Randomizer struct {
	rand       *rand.Rand
	lastRandom int // the appearance of random for players usually means not ever seeing the last value
}

func (r *Randomizer) Reset(seed int64) {
	r.rand = rand.New(rand.NewSource(seed))
	r.lastRandom = -1
}

func (r *Randomizer) Random(inclusiveMin, exclusiveMax int) int {
	if r.rand == nil {
		r.Reset(1)
	}
	// say max was 2, and min 0.
	// we should get either 0 or 1.
	// width is 2-0= 2, n = [0..2), ret= 2+0
	width := exclusiveMax - inclusiveMin
	n := r.rand.Intn(width) + inclusiveMin // [0,width)
	if n == r.lastRandom && width > 1 {
		n = ((n + 1) % width) + inclusiveMin // also, [0,width)
	}
	r.lastRandom = n
	return n
}
