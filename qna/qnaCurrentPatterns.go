package qna

import "strings"

type currentPatterns struct {
	calls map[string][]int
	depth int
}

// returns the deepest depth of the named pattern
func (cp *currentPatterns) isRunning(p string) (ret int) {
	p = strings.ToLower(p)
	calls := cp.calls[p]
	if cnt := len(calls); cnt > 0 {
		ret = calls[cnt-1]
	}
	return
}

func (cp *currentPatterns) startedPattern(p string) {
	p = strings.ToLower(p)
	if cp.calls == nil {
		cp.calls = make(map[string][]int)
	}
	cp.depth++
	cp.calls[p] = append(cp.calls[p], cp.depth)
}

func (cp *currentPatterns) stoppedPattern(p string) {
	p = strings.ToLower(p)
	was := cp.calls[p]
	cp.calls[p] = was[:len(was)-1]
	cp.depth--
}
