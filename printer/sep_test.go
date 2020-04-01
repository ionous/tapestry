package printer

import (
	"io"
	"testing"

	testify "github.com/stretchr/testify/assert"
)

func TestPrintSep(t *testing.T) {
	assert := testify.New(t)
	//
	if s, e := write(AndSeparator, "pizza"); assert.NoError(e) {
		assert.Equal("pizza", s)
	}
	if s, e := write(AndSeparator, "apple", "hedgehog", "washington", "mushroom"); assert.NoError(e) {
		assert.Equal("apple, hedgehog, washington, and mushroom", s)
	}
	if s, e := write(AndSeparator, "apple", "hedgehog"); assert.NoError(e) {
		assert.Equal("apple and hedgehog", s, "serial comma only after two items")
	}
	//
	if s, e := write(OrSeparator, "pistachio"); assert.NoError(e) {
		assert.Equal("pistachio", s)
	}
	if s, e := write(OrSeparator, "apple", "hedgehog", "washington", "mushroom"); assert.NoError(e) {
		assert.Equal("apple, hedgehog, washington, or mushroom", s)
	}
	if s, e := write(OrSeparator, "washington", "mushroom"); assert.NoError(e) {
		assert.Equal("washington or mushroom", s, "serial comma only after two items")
	}
}

func write(sep func(w io.Writer) io.WriteCloser, names ...string) (ret string, err error) {
	var buffer Spanner
	w := sep(&buffer)
	for _, n := range names {
		if _, e := io.WriteString(w, n); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		// normally PopWriter would call close, but we arent using the runtime here.
		if e := w.Close(); e != nil {
			err = e
		} else {
			ret = buffer.String()
		}
	}
	return
}