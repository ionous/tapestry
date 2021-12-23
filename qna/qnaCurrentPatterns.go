package qna

import "strings"

// used to answer meta.Running;
type currentPatterns struct {
	calls map[string][]int // depths at which a given pattern was started
	depth int              // a running depth of all patterns
}

// returns the deepest depth of the named pattern
func (cp *currentPatterns) runningPattern(p string) (ret int) {
	p = strings.ToLower(p)
	calls := cp.calls[p]
	if cnt := len(calls); cnt > 0 {
		ret = calls[cnt-1]
	}
	return
}

func (cp *currentPatterns) startedPattern(p string) {
	if cp.calls == nil {
		cp.calls = make(map[string][]int)
	}
	cp.depth++
	cp.calls[p] = append(cp.calls[p], cp.depth)
}

func (cp *currentPatterns) stoppedPattern(p string) {
	was := cp.calls[p]
	cp.calls[p] = was[:len(was)-1]
	cp.depth--
}
