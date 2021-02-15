package writer

import (
	"unicode/utf8"
)

// Chunk unifies the four possible output data types: rune, byte, []rune, and string.
type Chunk struct {
	Data interface{}
}

// Reset - forget the chunk's data; empty will return true, valid will return false.
func (c *Chunk) Reset() { c.Data = nil }

// IsClosed - return true if the chunk was explicitly set to writer.Closed
func (c *Chunk) IsClosed() bool { return c.Data == Closed }

// IsEmpty - return true if has no data or if the data is empty.
func (c *Chunk) IsEmpty() (okay bool) {
	switch b := c.Data.(type) {
	case nil:
		okay = true
	case []byte:
		okay = len(b) == 0
	case string:
		okay = len(b) == 0
	}
	return
}

// DecodeRune - return the first rune of the data.
func (c *Chunk) DecodeRune() (ret rune, cnt int) {
	switch b := c.Data.(type) {
	case byte:
		r := rune(b)
		ret, cnt = r, utf8.RuneLen(r)
	case rune:
		ret, cnt = b, utf8.RuneLen(b)
	case []byte:
		ret, cnt = utf8.DecodeRune(b)
	case string:
		ret, cnt = utf8.DecodeRuneInString(b)
	}
	return
}

// DecodeRune - return the last rune in the data.
func (c *Chunk) DecodeLastRune() (ret rune, cnt int) {
	switch b := c.Data.(type) {
	case byte:
		r := rune(b)
		ret, cnt = r, utf8.RuneLen(r)
	case rune:
		ret, cnt = b, utf8.RuneLen(b)
	case []byte:
		ret, cnt = utf8.DecodeLastRune(b)
	case string:
		ret, cnt = utf8.DecodeLastRuneInString(b)
	}
	return
}

// WriteTo - output the chunk's data to w; return the number of bytes written.
func (c *Chunk) WriteTo(w Output) (ret int, err error) {
	switch b := c.Data.(type) {
	case error:
		err = b
	case byte:
		if e := w.WriteByte(b); e != nil {
			err = e
		} else {
			ret = 1
		}
	case []byte:
		ret, err = w.Write(b)
	case rune:
		ret, err = w.WriteRune(b)
	case string:
		ret, err = w.WriteString(b)
	}
	return
}

// String - convert the chunk's data to string.
func (c *Chunk) String() (ret string) {
	switch b := c.Data.(type) {
	case []byte:
		ret = string(b)
	case byte:
		ret = string(b)
	case rune:
		ret = string(b)
	case string:
		ret = b
	}
	return
}
