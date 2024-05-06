package qna

import (
	"encoding/binary"
	"math/rand/v2"
	"time"
	"unsafe"
)

type Randomizer struct {
	src  *rand.PCG
	rand *rand.Rand
	// the appearance of random especially for small numbers
	// means not ever seeing the last number
	lastRandom int64
}

func RandomizedTime(ref *any) Randomizer {
	seed2 := uint64(uintptr(unsafe.Pointer(ref)))
	seed1 := uint64(time.Now().UnixNano())
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

func (r *Randomizer) Unmarshal(b []byte) (err error) {
	v, w := binary.Varint(b)
	if e := r.src.UnmarshalBinary(b[w:]); e != nil {
		err = e
	} else {
		r.lastRandom = v
	}
	return
}

func (r Randomizer) writeRandomizeer(w writeCb) (err error) {
	if b, e := r.src.MarshalBinary(); e != nil {
		err = e
	} else if e := w(nil, "$randomizer", "seed", b); e != nil {
		err = e
	} else if e := w(nil, "$randomizer", "last", r.lastRandom); e != nil {
		err = e
	}
	return
}
