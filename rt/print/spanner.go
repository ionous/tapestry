package print

import (
	"bytes"
	"unicode"

	"git.sr.ht/~ionous/tapestry/rt/writer"
)

// Spanner implements ChunkWriter, buffering output with spaces.
// It treats each new Write as a word adding spaces to separate words as necessary.
// FIX: this uses a buffer internally, but does it really need to do so?
type Spanner struct {
	buf bytes.Buffer // note: we cant aggregate buf or io.WriteString will bypasses implementation of Write() in favor of bytes.Buffer.WriteString()
}

func NewSpanner() *Spanner {
	return new(Spanner)
}

func (p *Spanner) Len() int {
	return p.buf.Len()
}
func (p *Spanner) Bytes() []byte {
	return p.buf.Bytes()
}
func (p *Spanner) String() string {
	return p.buf.String()
}

// WriteTo - move the contents of the spanner to the passed output.
func (p *Spanner) WriteTo(w writer.Output) (int, error) {
	i, e := p.buf.WriteTo(w)
	return int(i), e
}

// ChunkOutput - returns an object capable of writing to the spanner.
func (p *Spanner) ChunkOutput() writer.ChunkOutput {
	return p.WriteChunk
}

func (p *Spanner) WriteChunk(c writer.Chunk) (ret int, err error) {
	// writing something?
	if b, cnt := c.DecodeLastRune(); cnt > 0 {
		// and already written something and the thing we are writing is not a space?
		if p.Len() > 0 && !spaceLike(b) {
			p.buf.WriteRune(' ')
		}
		// go (bytes.Buffer) expects us to return how much of the input we wrote,
		// not the actual number of bytes we wrote.
		ret, err = c.WriteTo(&p.buf)
	}
	return
}

func spaceLike(r rune) bool {
	return unicode.IsSpace(r) || unicode.In(r, unicode.Po, unicode.Pi, unicode.Pf)
}

// https://www.compart.com/en/unicode/category/Pi
// Pc     = _Pc // Pc is the set of Unicode characters in category Pc.
// Pd     = _Pd // Pd is the set of Unicode characters in category Pd.
// Pe     = _Pe // Pe is the set of Unicode characters in category Pe.
// Pf     = _Pf // Pf is the set of Unicode characters in category Pf.
// Pi     = _Pi // Pi is the set of Unicode characters in category Pi.
// Po     = _Po // Po is the set of Unicode characters in category Po.
// Ps     = _Ps // Ps is the set of Unicode characters in category Ps.
//
// connector: _
// dash: -|
// end/close: )〞-- closing brackets, some arabic? right ornate quotes
// final: ’” -- angled right quotes
// initial: ‘“ -- angled left quotes
// other: !"'.; -- straight quotes, and full stop
// start/open: ([{ -- opening brackets, some arabic? left ornate quotes
