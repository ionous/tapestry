package parser_test

import (
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPhraseMaker(t *testing.T) {
	assert := testify.New(t)
	assert.EqualValues(sliceOf.String("look"), Phrases("look"))
	assert.EqualValues(sliceOf.String("look", "l"), Phrases("look/l"))

	assert.EqualValues([]string{
		"look inside something", "look inside wicked",
		"look in something", "look in wicked",
		"l inside something", "l inside wicked",
		"l in something", "l in wicked"},
		Phrases("look/l inside/in something/wicked"))
}

// generate permutations from inform-like slash phrases
func Phrases(phrase string) (ps []string) {
	// step 1 split the phrase into space chunks
	var multi [][]string
	for _, f := range strings.Fields(phrase) {
		multi = append(multi, strings.Split(f, "/"))
	}
	for _, m := range multi {
		ps = permute(ps, m)
	}
	return
}

func permute(a, b []string) (ps []string) {
	if len(a) == 0 {
		ps = b
	} else {
		for _, a := range a {
			for _, b := range b {
				n := strings.Join([]string{a, b}, " ")
				ps = append(ps, n)
			}
		}
	}
	return
}
