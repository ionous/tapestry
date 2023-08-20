package rules

type eventPrefix int

//go:generate stringer -type=eventPrefix -linecomment
const (
	instead eventPrefix = iota // instead of
	before
	after
	report
	//
	numPrefixes = iota
)

// with the theory that sqlite sorts asc by default
// ( spaces numbers out so some theoretical extra 'first' or 'last' prefix could add or subtract to the numbers )
func (p eventPrefix) rank() (ret int) {
	return ranks[p]
}

var ranks = []int{-20, -10, 10, 20, 0}
