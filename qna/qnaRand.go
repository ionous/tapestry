package qna

import (
	"database/sql"
	"hash/fnv"
	"io"
	"math/rand/v2"
	"runtime/debug"
	"time"
)

const randomizerKey = "$randomizer"

type Randomizer struct {
	src  *rand.PCG
	rand *rand.Rand
	// the appearance of random especially for small numbers
	// means not ever seeing the last number
	lastRandom int64
}

func RandomizedTime() Randomizer {
	seed2 := uint64(0xbadf00d)
	if b, ok := debug.ReadBuildInfo(); ok {
		seed2 = hash(b.Main.Sum)
	}
	seed1 := uint64(time.Now().UnixNano())
	return SeedRandomizer(seed1, seed2)
}

func hash(str string) uint64 {
	w := fnv.New64a()
	io.WriteString(w, str)
	return w.Sum64()
}

func SeedRandomizer(seed1, seed2 uint64) Randomizer {
	src := rand.NewPCG(seed1, seed2)
	return Randomizer{src, rand.New(src), -1}
}

func (r *Randomizer) Random(inclusiveMin, exclusiveMax int) int {
	// say max was 2, and min 0.
	// we should get either 0 or 1.
	// width is 2-0= 2, n = [0..2), ret= 2+0
	width := exclusiveMax - inclusiveMin
	n := r.rand.IntN(width) + inclusiveMin // [0,width)
	if int64(n) == r.lastRandom && width > 1 {
		n = ((n + 1) % width) + inclusiveMin // also, [0,width)
	}
	r.lastRandom = int64(n)
	return n
}

func (r Randomizer) writeRandomizeer(w writeCb) (err error) {
	if b, e := r.src.MarshalBinary(); e != nil {
		err = e
	} else if e := w("", randomizerKey, "seed", b); e != nil {
		err = e
	} else if e := w("", randomizerKey, "last", r.lastRandom); e != nil {
		err = e
	}
	return
}

func (r *Randomizer) readRandomizer(db *sql.DB) (err error) {
	var b []byte
	var last int64
	if e := db.QueryRow(
		`select value from run_value 
		where noun = $1 and field = $2`,
		randomizerKey, "seed").Scan(&b); e != nil {
		err = e
	} else if e := db.QueryRow(
		`select value from run_value 
		where noun = $1 and field = $2`,
		randomizerKey, "last").Scan(&last); e != nil {
		err = e
	} else if e := r.src.UnmarshalBinary(b); e != nil {
		err = e
	} else {
		r.lastRandom = last
	}
	return
}
